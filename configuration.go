package main

import (
	"encoding/json"
	"log"
	"os"
)

type SBConfiguration struct {
	Email string
	Password string

	LogSystem bool
}

func LoadConfiguration() *SBConfiguration {
	file, err := os.Open("Applications/StrategicBot/configuration.json")
	deny(err)

	decoder := json.NewDecoder(file)
	configuration := new(SBConfiguration)
	decoder.Decode(configuration)

	log.Println("Configuration file loaded...")
	log.Println("Email:", configuration.Email)
	log.Println("LogSystem:", configuration.LogSystem)

	return configuration
}