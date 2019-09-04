package service

import (
	"encoding/json"
	"log"
	"os"
)

var config *Configuration

type Configuration struct {
	ServerKey string `json:"server_key"`
	AppleValidationKey string `json:"apple_validation_key"`
	AppleValidationUrl string `json:"apple_validation_url"`
}

func Config() *Configuration {
	if config == nil {
		file, _ := os.Open("config.json")
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&config)
		if err != nil {
			log.Fatal("json config error: " + err.Error())
		}
	}
	return config
}
