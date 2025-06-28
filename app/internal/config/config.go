// Package config provides configuration management for the application.
// It handles loading environment variables from .env files and provides
// default values for essential application settings.
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration values needed by the application.
// These values are typically loaded from environment variables or .env files.
type Config struct {
	Port     string // Server port number (e.g., "3030", "8080")
	MongoURI string // MongoDB connection URI (e.g., "mongodb://localhost:27017")
	DBName   string // MongoDB database name to use
}

// Load reads configuration from environment variables and .env file.
// It attempts to load a .env file first, then reads environment variables
// with fallback to sensible default values.
// Returns a pointer to a Config struct with all values populated.
func Load() *Config {
	// Try to load .env file - this will fail silently in production
	// environments where .env files might not exist
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create and return config with environment variables or defaults
	return &Config{
		Port:     getEnv("PORT", "3030"),                           // Default to port 3030
		MongoURI: getEnv("MONGODB_URI", "mongodb://mongodb:27017"), // Default to Docker MongoDB service
		DBName:   getEnv("MONGODB_NAME", "blog"),                   // Default database name
	}
}

// getEnv retrieves an environment variable value with a fallback default.
// If the environment variable exists and is not empty, it returns that value.
// Otherwise, it returns the provided default value.
//
// Parameters:
//   - key: the environment variable name to look up
//   - defaultValue: the value to return if the environment variable is not set
//
// Returns the environment variable value or the default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
