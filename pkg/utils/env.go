package utils

import (
	"github.com/joho/godotenv"
)

func LoadEnv() error {
	if err := godotenv.Load(".env"); err != nil {
		return nil
	}
	return nil
}
