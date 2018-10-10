package db

import (
	"encoding/json"
	"log"
	"os"
)

var config *Configuration

type Configuration struct {
	CassandraUser     string `json:"cassandra_user"`
	CassandraPassword string `json:"cassandra_password"`
	CassandraIp       string `json:"cassandra_ip"`
	CassandraKey      string `json:"cassandra_keyspace"`
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
