package main

import (
	"log"
	"os"
)

func main() {
	log.Println("Starting Up...")

	config := LoadConfiguration()

	if config.LogSystem {
		file, err := os.OpenFile("Applications/StrategicBot/bot.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		deny(err)
		log.SetOutput(file)
	}

	log.Println("Successfully started!")

	Connect(config.Email, config.Password)

	ScrollValues := LoadScrollValues()

	log.Println(ScrollValues)

	StoreScrollValues(ScrollValues)
}

// Some error handling, could be improved
func deny(err error) {
	if err != nil {
		panic(err)
	}
}