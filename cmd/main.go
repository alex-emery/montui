package main

import (
	"log"
	"os"

	"github.com/alex-emery/montui/internal/ui"
	"github.com/alex-emery/montui/pkg/montui"
	"github.com/dotenv-org/godotenvvault"
)

func main() {
	err := godotenvvault.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretID := os.Getenv("SECRET_ID")
	secretKey := os.Getenv("SECRET_KEY")
	montui, err := montui.New(secretID, secretKey, "./sqlite.db")

	if err != nil {
		log.Fatal("failed to create montui")
	}

	tui, err := ui.New(montui)
	if err != nil {
		log.Fatal("failed to start", err)
	}

	err = ui.Run(tui)
	if err != nil {
		log.Fatal("failed to run", err)
	}
}
