package main

import (
	"encoding/json"
	"log"
	"net"
)

type SBState struct {
	con			net.Conn
	chQuit		chan bool
	chRequests	chan Request
}

var (
	CardTypes = make(map[CardId]Card)
	CardRarities = make(map[Card]int)
)

func InitState(con net.Conn) *SBState {
	s := SBState{con: con}
	s.chQuit = make(chan bool, 5)
	s.chRequests = make(chan Request, 10)

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
		log.Println("Read out card types and rarities:", CardTypes, CardRarities)

	case "RoomChatMessage":
		var v MRoomChatMessage
		json.Unmarshal(reply, &v)
		log.Println(m)
		log.Println("Chat message:", v.Text)
		
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

func (s *SBState) Say(room Channel, text string) {
	s.SendRequest(Request{"msg": "RoomChatMessage", "text": text, "roomName": room})
}

func (s *SBState) Whisper(player Player, text string) {
	s.SendRequest(Request{"msg": "Whisper", "text": text, "toProfileName": player})
}