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
	ES_index           string  `json:"es_index"`
	ES_user            string  `json:"es_user"`
	ES_password        string  `json:"es_password"`
	ES_ip              string  `json:"es_ip"`
	ES_port            string  `json:"es_port"`
	ES_distance        float64 `json:"es_distance"`
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
