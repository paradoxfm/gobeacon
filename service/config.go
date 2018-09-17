package service

import (
	"encoding/json"
	"os"
	"log"
)

var config *Configuration

type Configuration struct {
	Server_port        string  `json:"server_port"`
	Cassandra_user     string  `json:"cassandra_user"`
	Cassandra_password string  `json:"cassandra_password"`
	Cassandra_ip       string  `json:"cassandra_ip"`
	Cassandra_keyspace string  `json:"cassandra_keyspace"`
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
