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
	saltedPepperedPassword := password + salt + pepper
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(saltedPepperedPassword))
	return err == nil
}
func getBinaryBySHA256(s string) []byte {
	r := sha256.Sum256([]byte(s))
	return r[:]
}

func HashEmail(email, pepper string) (string, error) {
	mac := hmac.New(sha256.New, getBinaryBySHA256(pepper))
	_, err := mac.Write([]byte(email))
	return string(mac.Sum(nil)), err
}

func VerifyEmail(email, pepper, hashedEmail string) bool {
	mac := hmac.New(sha256.New, getBinaryBySHA256(pepper))
	_, err := mac.Write([]byte(email))
	if err != nil {
		return false
	}
	expectedMAC := mac.Sum(nil)
	return hmac.Equal([]byte(hashedEmail), expectedMAC)
}
