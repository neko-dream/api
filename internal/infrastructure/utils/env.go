package utils

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}
	return nil
}
