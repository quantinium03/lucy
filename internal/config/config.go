package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Config(key string) string {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Failed to load the environment variables")
	}
	return os.Getenv(key)
}
