package main

import (
	"encoding/json"
	"log"
	"net"
)

type SBState struct {
	con			net.Conn
	chMessages	chan Message
	chRequests	chan Request
	chQuit		chan bool

	Player
	Owner			Player
	ProfileId		string
	Stock			map[Card][3]int
	Gold			int
}

type Message struct {
	Text 		string
	From		Player
	Channel		Channel
}

var (
	CardTypes = make(map[CardId]Card)
	CardRarities = make(map[Card]int)
)

func InitState(con net.Conn, owner string) *SBState {
	s := SBState{con: con}
	s.Owner = Player(owner)

	s.chMessages = make(chan Message, 1)
	s.chRequests = make(chan Request, 10)
	s.chQuit = make(chan bool, 5)

	s.SendRequest(Request{"msg": "JoinLobby"})

	return &s
}

func (s *SBState) HandleReply(reply []byte) bool {
	if len(reply) < 2 {
		log.Println("HandleReply: Reply too short")
		return false
	}

	var m Reply
	err := json.Unmarshal(reply, &m)
	deny(err)
	// log.Println(string(reply))

	switch m.Msg {
	case "CardTypes":
		var v MCardTypes
		json.Unmarshal(reply, &v)
		for _, cardType := range v.CardTypes {
			CardTypes[cardType.Id] = cardType.Name
			CardRarities[cardType.Name] = cardType.Rarity
		}
		log.Println(m)

	case "LibraryView":
		var v MLibraryView
		json.Unmarshal(reply, &v)

		if v.ProfileId == s.ProfileId {

			stock := make(map[Card][3]int)
			for _, card := range CardTypes {
				stock[card] = [3]int{0, 0, 0}
			}

			for _, card := range v.Cards {
				if card.Tradable {
					name := CardTypes[card.TypeId]
					st := stock[name]
					st[card.Level]++
					stock[name] = st
				}
			}

			s.Stock = stock

			log.Println("Read out stock")
		}

	case "ProfileDataInfo":
		var v MProfileDataInfo
		json.Unmarshal(reply, &v)
		s.Gold = v.ProfileData.Gold

	case "ProfileInfo":
		var v MProfileInfo
		json.Unmarshal(reply, &v)

		s.Player = v.Profile.Name
		s.ProfileId = v.Profile.Id

	case "RoomChatMessage":
		var v MRoomChatMessage
		json.Unmarshal(reply, &v)
		if v.From != s.Player {
			s.chMessages <- Message{v.Text, v.From, Channel(v.RoomName)}
		}

		log.Println("Chat message:", v.Text, "From:", v.From, "In:", v.RoomName)

	case "TradeInviteForward":
		var v MTradeInviteForward
		json.Unmarshal(reply, &v)
		log.Println(string(reply))
		s.SendRequest(Request{"msg": "TradeInvite", "profile": v.Inviter.Id})

	case "Whisper":
		var v MWhisper
		json.Unmarshal(reply, &v)
		if v.From != s.Player {
			s.chMessages <- Message{v.Text, v.From, Channel("WHISPER")}
		}

		log.Println("Whisper message:", v.Text, "From:", v.From)

	default:
		log.Println(m)
	}

	return true
}

func (s *SBState) SendRequest(req Request) {
	log.Println("Send request:", req)
	s.chRequests <- req
}

func (s *SBState) JoinRoom(room Channel) {
	s.SendRequest(Request{"msg": "RoomEnter", "roomName": room})
}

func (s *SBState) LeaveRoom(room Channel) {
	s.SendRequest(Request{"msg": "RoomExit", "roomName": room})
}