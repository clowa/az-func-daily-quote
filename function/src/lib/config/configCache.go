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
	MongodbConnectionString string
	MongodbDatabase         string
	MongodbCollection       string
	ApiPort                 string
}

func (c *Config) LoadConfig() {
	c.ApiPort = ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		c.ApiPort = ":" + val
	}

	c.MongodbConnectionString = os.Getenv("MONGODB_CONNECTION_STRING")
	if c.MongodbConnectionString == "" {
		log.Fatal("Environment variable \"MONGODB_CONNECTION_STRING\" not set")
	}

	c.MongodbDatabase = os.Getenv("MONGODB_DATABASE")
	if c.MongodbDatabase == "" {
		log.Fatal("Environment variable \"MONGODB_DATABASE\" not set")
	}

	c.MongodbCollection = os.Getenv("MONGODB_COLLECTION")
	if c.MongodbCollection == "" {
		log.Fatal("Environment variable \"MONGODB_COLLECTION\" not set")
	}
}

func (c *Config) PrintConfig() {
	log.Printf("MONGODB_CONNECTION_STRING: %s\n", c.MongodbConnectionString)
	log.Printf("MONGODB_DATABASE: %s\n", c.MongodbDatabase)
	log.Printf("MONGODB_COLLECTION: %s\n", c.MongodbCollection)
}
