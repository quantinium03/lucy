package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Config(key string) string {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load the .env file")
	}

	return os.Getenv(key)
}
