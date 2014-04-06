package main

import (
	"encoding/json"
	"log"
	"os"
)

type SBScrollValues struct {
	Data []struct {
		Id		CardId
		Name	Card
		Value	int
	}
}

func LoadScrollValues() *SBScrollValues {
	file, err := os.Open("Applications/StrategicBot/values.json")
	deny(err)

	decoder := json.NewDecoder(file)
	scrollValues := new(SBScrollValues)
	decoder.Decode(scrollValues)

	for key, scroll := range scrollValues.Data {
		log.Println("Key is", key, "and Value is", scroll)
		log.Println(scroll.Name, "is valued", scroll.Value, "gold.")
	}

	return scrollValues
}