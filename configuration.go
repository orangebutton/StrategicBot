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
	Owner string
}

func LoadConfiguration() *SBConfiguration {
	file, err := os.Open("Applications/StrategicBot/configuration.json")
	deny(err)

	decoder := json.NewDecoder(file)
	config := new(SBConfiguration)
	decoder.Decode(config)

	log.Println("===== Configuration =====")
	log.Println("Email:", config.Email)
	log.Println("LogSystem:", config.LogSystem)
	log.Println("Owner:", config.Owner, "\n")

	return config
}