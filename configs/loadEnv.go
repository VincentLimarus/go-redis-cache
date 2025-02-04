package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvironment() {
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	} else if os.IsNotExist(err) {
		log.Println(".env file not found, using environment variables from Docker")
	} else {
		log.Fatal("Error checking .env file:", err)
	}
}