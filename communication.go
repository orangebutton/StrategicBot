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
	command, args := ParseCommandAndArgs(m.Text)

	switch command {
	case "!help":
		reply := SBHelpMessage + args
		s.Say(m.Channel, reply)
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