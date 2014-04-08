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
	log.Println(m)
	// log.Println(string(reply))

	return true
}

func (s *SBState) SendRequest(req Request) {
	log.Println("Send request:", req)
	s.chRequests <- req
}