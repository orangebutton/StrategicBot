package main

import (
	"strings"
)

const(
	SBHelpMessage = "I am an automated trading unit of the strategic angels"
)

func (s *SBState) StartMessageHandling() {
	go func() {
		for {
			m := <-s.chMessages
			s.HandleMessage(m)
		}
	}()
}

func (s *SBState) HandleMessage(m Message) {
	command, _ := ParseCommandAndArgs(m.Text)

	switch command {
	case "!help":
		s.Reply(SBHelpMessage, false, m)
	case "!price":
		s.Reply("Cannot tell the price yet", true, m)
	case "!stock":
		s.Reply("Don't know how to tell the stock", false, m)
	}
}

func ParseCommandAndArgs(text string) (command, args string) {
	text = strings.TrimSpace(strings.ToLower(text))
	strs := strings.SplitN(text, " ", 2)

	command = strings.TrimSpace(strs[0])

	if len(strs) > 1 {
		args = strings.TrimSpace(strs[1])
	}

	if !strings.HasPrefix(command, "!") {
		command = "!" + command
	}

	return command, args
}

func (s *SBState) Say(room Channel, text string) {
	s.SendRequest(Request{"msg": "RoomChatMessage", "text": text, "roomName": room})
}

func (s *SBState) Reply(text string, whisper bool, m Message) {
	if m.Channel == "WHISPER" || whisper {
		s.Whisper(m.From, text)
	} else {
		s.Say(m.Channel, text)
	}
}

func (s *SBState) Whisper(player Player, text string) {
	s.SendRequest(Request{"msg": "Whisper", "text": text, "toProfileName": player})
}