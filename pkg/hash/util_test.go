package hash_test

import (
	"testing"

	"github.com/neko-dream/server/pkg/hash"
)

func Test_HashEmail(t *testing.T) {
    salt, err := hash.GenerateSalt(16)
    if err != nil {
        t.Fatalf("failed to generate salt: %v", err)
    }

    password := "password123"
    pepper := "pepper"
    cost := 10
    hashedPassword, err := hash.HashPassword(password, salt, pepper, cost)
    if err != nil {
        t.Fatalf("failed to hash password: %v", err)
    }

    isValid := hash.VerifyPassword(password, hashedPassword)
    if !isValid {
        t.Fatalf("password verification failed")
    }
}

