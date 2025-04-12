package password_auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GenerateSalt(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func HashPassword(password, salt, pepper string, cost int) (string, error) {
	saltedPepperedPassword := password + salt + pepper
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return "", fmt.Errorf("cost must be between %d and %d", bcrypt.MinCost, bcrypt.MaxCost)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(saltedPepperedPassword), cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hash), nil
}

func VerifyPassword(password string, salt, pepper, hashedPassword string) bool {
	// 検証時もペッパーを使うのじゃ
	saltedPepperedPassword := password + salt + pepper
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(saltedPepperedPassword))
	return err == nil
}
