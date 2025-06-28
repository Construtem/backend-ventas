package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from the .env file
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading the .env file: %v", err)
	}
}

// GetEnv gets the value of an environment variable
func GetEnv(key string) string {
	return os.Getenv(key)
}
