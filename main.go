package main

import (
	"log"
	"os"
	"time"
)

func main() {
	config := LoadConfiguration()

	if config.LogSystem {
		file, err := os.OpenFile("Applications/StrategicBot/bot.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		deny(err)
		log.SetOutput(file)
	}

	for {
		StartBot(config.Email, config.Password)
	}
}

func StartBot(email, password string) {
	defer func() {
		log.Println("Shut bot down.")
	}()

	s, chAlive := Connect(email, password)
	log.Println(s, chAlive)

	s.JoinRoom("strategic trading")
	s.Say("strategic trading", "hello")

	for {
		timeout := time.After(time.Minute * 1)
		InnerLoop:
		for {
			select {
				case <-chAlive:
					break InnerLoop
				case <-s.chQuit:
					log.Println("Bot Quit")
					s.chQuit <- true
					return
				case <-timeout:
					log.Println("Time out")
					return
			}
		}
	}
}

// Some error handling, could be improved
func deny(err error) {
	if err != nil {
		panic(err)
	}
}