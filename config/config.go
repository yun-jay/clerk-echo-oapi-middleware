package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ClerkSecretKey string
}

// LoadConfig loads configuration from environment variables or .env file
func LoadConfig() *Config {
	// Load variables from .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, loading environment variables directly")
	}

	key := os.Getenv("CLERK_SECRET_KEY")
	if key == "" {
		log.Fatal("CLERK_SECRET_KEY environment variable is not set")
	}

	return &Config{
		ClerkSecretKey: key,
	}
}
