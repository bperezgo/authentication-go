package config

import (
	"log"

	"github.com/tkanos/gonfig"
)

type Config struct {
	Port          int
	AuthJwtSecret string
}

var config *Config

// Singleton to get the configuration
func GetConfig() *Config {
	if config != nil {
		return config
	}
	config = &Config{}
	err := gonfig.GetConf("configuration.json", config)
	if err != nil {
		log.Panic("[ERROR] Failed getting the configuration envs:", err)
	}
	return config
}
