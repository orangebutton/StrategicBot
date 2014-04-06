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
		Value	float64
	}
}

// Determines the price of a specified card
func DeterminePrice(card Card, num int, buy bool) int {
	return 0
}

// Checks if player gets a discount
func DiscountInPercent() int {
	return 0
}

// Load scroll values from local JSON file
func LoadScrollValues() *SBScrollValues {
	file, err := os.Open("Applications/StrategicBot/values.json")
	deny(err)

	decoder := json.NewDecoder(file)
	scrollValues := new(SBScrollValues)
	decoder.Decode(scrollValues)

	// Todo: This for loop should also set a default value for new scrolls
	for key, scroll := range scrollValues.Data {
		log.Println("Key is", key, "and Value is", scroll)
		log.Println(scroll.Name, "is valued", scroll.Value, "gold.")
	}

	return scrollValues
}

// Updates scroll values
func UpdateScrollValues(scrollValues *SBScrollValues, card Card, num int, buy bool) *SBScrollValues {
	return scrollValues
}

// Store scroll values in local JSON file
func StoreScrollValues(scrollValues *SBScrollValues) {
	file, err := os.OpenFile("Applications/StrategicBot/values.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	deny(err)

	encoder := json.NewEncoder(file)
	encoder.Encode(scrollValues)

	log.Println("Stored new values")
}