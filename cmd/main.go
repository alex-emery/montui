package main

import (
	"log"
	"os"

	"github.com/alex-emery/montui/internal/ui"
	"github.com/alex-emery/montui/pkg/montui"
	"github.com/dotenv-org/godotenvvault"
	"go.uber.org/zap"
)

func NewLogger(outpath string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()

	if outpath != "" {
		cfg.OutputPaths = []string{
			outpath,
		}
	}

	return cfg.Build()
}

func main() {
	logger, err := NewLogger("./logs.txt")
	if err != nil {
		log.Fatal(err)
	}

	err = godotenvvault.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	secretID := os.Getenv("SECRET_ID")
	secretKey := os.Getenv("SECRET_KEY")

	if secretID == "" {
		logger.Fatal("SECRET_ID env not set")
	}

	if secretKey == "" {
		logger.Fatal("SECRET_KEY env not set")
	}

	montui, err := montui.New(secretID, secretKey, "./sqlite.db", logger)
	if err != nil {
		logger.Fatal("failed to create montui")
	}

	tui, err := ui.New(montui)
	if err != nil {
		logger.Fatal("failed to start", zap.Error(err))
	}

	err = ui.Run(tui)
	if err != nil {
		logger.Fatal("failed to run", zap.Error(err))
	}
}
