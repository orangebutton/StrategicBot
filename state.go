package main

import (
	"encoding/json"
	"log"
	"net"
)

type SBState struct {
	con			net.Conn
}

func (s *SBState) HandleReply(reply []byte) {
	var m Reply
	err := json.Unmarshal(reply, &m)
	deny(err)
	log.Println(m)
	log.Println(string(reply))
}