package password_auth

import (
	"crypto/rand"
	"math/big"
)

func GeneratePassword(length int) string {
	var letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var specialBytes = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	var numBytes = "0123456789"

	var allBytes string

	allBytes += letterBytes

	allBytes += specialBytes

	allBytes += numBytes

	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(allBytes))))
		b[i] = allBytes[n.Int64()]
	}

	return string(b)
}
