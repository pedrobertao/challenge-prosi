package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	MongoURI string
	DBName   string
}

func Load() *Config {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return &Config{
		Port:     getEnv("PORT", "3030"),
		MongoURI: getEnv("MONGODB_URI", "mongodb://mongodb:27017"),
		DBName:   getEnv("MONGODB_NAME", "blog"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
