package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Address   string   `json:"address"`
	Upstreams []string `json:"upstreams"`
}

func NewConfig() *Config {
	content, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal("Failed to open config file")
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatalf("Error during unmarshal config: %s\n", err)
	}

	return &config
}
