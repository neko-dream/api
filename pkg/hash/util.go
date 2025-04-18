package hash

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
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
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return "", fmt.Errorf("cost must be between %d and %d", bcrypt.MinCost, bcrypt.MaxCost)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hash), nil
}

func VerifyPassword(password string, hashedPassword string) bool {
	res := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return res == nil
}

func getBinaryBySHA256(s string) []byte {
	r := sha256.Sum256([]byte(s))
	return r[:]
}

func HashEmail(email, pepper string) (string, error) {
	mac := hmac.New(sha256.New, getBinaryBySHA256(pepper))
	_, err := mac.Write([]byte(email))
	if err != nil {
		return "", err
	}

	hashedEmail := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(hashedEmail), err
}

func VerifyEmail(email, pepper, hashedEmail string) bool {
	mac := hmac.New(sha256.New, getBinaryBySHA256(pepper))
	_, err := mac.Write([]byte(email))
	if err != nil {
		return false
	}
	expectedMAC := mac.Sum(nil)
	expected:= base64.StdEncoding.EncodeToString(expectedMAC)
	return hmac.Equal([]byte(hashedEmail), []byte(expected))
}
