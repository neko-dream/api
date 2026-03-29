package hash_test

import (
	"testing"

	"github.com/neko-dream/api/pkg/hash"
)

func Test_BcryptHash(t *testing.T) {
	hashed, err := hash.BcryptHash("password", "pepper", 10)
	if err != nil {
		t.Fatalf("failed to hash: %v", err)
	}
	if !hash.BcryptVerify("password", "pepper", hashed) {
		t.Fatal("verification failed")
	}
}

func Test_BcryptVerify_WrongPassword(t *testing.T) {
	hashed, err := hash.BcryptHash("correct", "pepper", 10)
	if err != nil {
		t.Fatalf("failed to hash: %v", err)
	}
	if hash.BcryptVerify("wrong", "pepper", hashed) {
		t.Fatal("should not verify with wrong password")
	}
}

func Test_BcryptVerify_WrongPepper(t *testing.T) {
	hashed, err := hash.BcryptHash("password", "pepper-a", 10)
	if err != nil {
		t.Fatalf("failed to hash: %v", err)
	}
	if hash.BcryptVerify("password", "pepper-b", hashed) {
		t.Fatal("should not verify with wrong pepper")
	}
}

func Test_HMACHash(t *testing.T) {
	hashed, err := hash.HMACHash("email@example.com", "pepper")
	if err != nil {
		t.Fatalf("failed to hash: %v", err)
	}
	if !hash.HMACVerify("email@example.com", "pepper", hashed) {
		t.Fatal("verification failed")
	}
}

func Test_HMACVerify_WrongValue(t *testing.T) {
	hashed, err := hash.HMACHash("correct@example.com", "pepper")
	if err != nil {
		t.Fatalf("failed to hash: %v", err)
	}
	if hash.HMACVerify("wrong@example.com", "pepper", hashed) {
		t.Fatal("should not verify with wrong value")
	}
}

func Test_HMACVerify_WrongPepper(t *testing.T) {
	hashed, err := hash.HMACHash("email@example.com", "pepper-a")
	if err != nil {
		t.Fatalf("failed to hash: %v", err)
	}
	if hash.HMACVerify("email@example.com", "pepper-b", hashed) {
		t.Fatal("should not verify with wrong pepper")
	}
}
