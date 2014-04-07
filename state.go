package main

import (
	"encoding/json"
	"log"
	"net"
)

type SBState struct {
	con			net.Conn
	chQuit		chan bool
}

func InitState(con net.Conn) *SBState {
	s := SBState{con: con}
	return &s
}

func (s *SBState) HandleReply(reply []byte) bool {
	var m Reply
	err := json.Unmarshal(reply, &m)
	deny(err)
	log.Println(m)
	// log.Println(string(reply))

	return true
}