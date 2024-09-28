package utils

import (
	"fmt"

	"braces.dev/errtrace"
	"github.com/joho/godotenv"
)

func LoadEnv() error {
	if err := godotenv.Load(".env"); err != nil {
		return errtrace.Wrap(fmt.Errorf("failed to load .env file: %w", err))
	}
	return nil
}
