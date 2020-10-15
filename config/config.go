package config

import (
	"encoding/json"
	"log"
	"os"
)

// Cfg for config
var Cfg *Config = &Config{}

// Config for config struct
type Config struct {
	ServerConfigurations ServerConfigurations
	Logger               Logger
	DBNAME               string
	DBUSER               string
	DBPASS               string
	DBHOST               string
	DBPORT               string
	ViralLaunchEmail     string
	ViralLaunchID        string
}

// ServerConfigurations for server config
type ServerConfigurations struct {
	Port         string
	InstanceName string
}

// Logger for logger type
type Logger struct {
	FileLocation string
	Level        string
}

// LoadConfig loads config info
func LoadConfig(path string) {
	file, errOpenFile := os.Open(path)
	if errOpenFile != nil {
		log.Fatal(errOpenFile)
	}

	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)

	if err != nil {
		log.Fatal(err)
	}

	Cfg = &configuration
}
