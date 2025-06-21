package utils

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	// Try to load .env from current directory first
	if err := godotenv.Load(); err != nil {
		// If that fails, try to find .env by walking up the directory tree
		cwd, _ := os.Getwd()
		for dir := cwd; dir != "/" && dir != "."; dir = filepath.Dir(dir) {
			envPath := filepath.Join(dir, ".env")
			if _, err := os.Stat(envPath); err == nil {
				if err := godotenv.Load(envPath); err == nil {
					log.Printf("Loaded .env from: %s", envPath)
					return nil
				}
			}
		}
		log.Println("Error loading .env file, using default environment variables")
	}

	return nil
}
