package config

import (
	"log"
	"os"
)

var GlobalConfig Config

func GetConfig() *Config {
	return &GlobalConfig
}

type Config struct {
	CosmosConnectionString string
	CosmosDatabase         string
	CosmosContainer        string
	ApiPort                string
}

func (c *Config) LoadConfig() {
	c.ApiPort = ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		c.ApiPort = ":" + val
	}

	c.CosmosConnectionString = os.Getenv("COSMOS_CONNECTION_STRING")
	if c.CosmosConnectionString == "" {
		log.Fatal("Environment variable \"COSMOS_CONNECTION_STRING\" not set")
	}

	c.CosmosDatabase = os.Getenv("COSMOS_DATABASE")
	if c.CosmosDatabase == "" {
		log.Fatal("Environment variable \"COSMOS_DATABASE\" not set")
	}

	c.CosmosContainer = os.Getenv("COSMOS_CONTAINER")
	if c.CosmosContainer == "" {
		log.Fatal("Environment variable \"COSMOS_CONTAINER\" not set")
	}
}

func (c *Config) PrintConfig() {
	log.Printf("COSMOS_HOST: %s\n", c.CosmosConnectionString)
	log.Printf("COSMOS_DATABASE: %s\n", c.CosmosDatabase)
	log.Printf("COSMOS_CONTAINER: %s\n", c.CosmosContainer)
}
