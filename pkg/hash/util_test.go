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

func Test_BcryptHashHMAC(t *testing.T) {
	hashed, err := hash.BcryptHashHMAC("passphrase", "pepper", 10)
	if err != nil {
		t.Fatalf("failed to hash: %v", err)
	}
	if !hash.BcryptVerifyHMAC("passphrase", "pepper", hashed) {
		t.Fatal("verification failed")
	}
}

func Test_BcryptHashHMAC_LongInput(t *testing.T) {
	// pepper(32バイト) + plaintext(72バイト) = 104バイト → BcryptHashだと上限超過するケース
	longPepper := "12345678901234567890123456789012" // 32バイト
	longPlain := "1234567890123456789012345678901234567890123456789012345678901234567890ab" // 72バイト
	hashed, err := hash.BcryptHashHMAC(longPlain, longPepper, 10)
	if err != nil {
		t.Fatalf("long input should not fail: %v", err)
	}
	if !hash.BcryptVerifyHMAC(longPlain, longPepper, hashed) {
		t.Fatal("verification failed for long input")
	}
}

func Test_BcryptHashHMAC_WrongPassword(t *testing.T) {
	hashed, err := hash.BcryptHashHMAC("correct", "pepper", 10)
	if err != nil {
		t.Fatalf("failed to hash: %v", err)
	}
	if hash.BcryptVerifyHMAC("wrong", "pepper", hashed) {
		t.Fatal("should not verify with wrong password")
	}
}

func Test_BcryptHashHMAC_WrongPepper(t *testing.T) {
	hashed, err := hash.BcryptHashHMAC("passphrase", "pepper-a", 10)
	if err != nil {
		t.Fatalf("failed to hash: %v", err)
	}
	if hash.BcryptVerifyHMAC("passphrase", "pepper-b", hashed) {
		t.Fatal("should not verify with wrong pepper")
	}
}

func Test_BcryptHashHMAC_IncompatibleWithBcryptHash(t *testing.T) {
	// BcryptHashとBcryptHashHMACは互換性がない
	hashed, err := hash.BcryptHash("passphrase", "pepper", 10)
	if err != nil {
		t.Fatalf("failed to hash: %v", err)
	}
	if hash.BcryptVerifyHMAC("passphrase", "pepper", hashed) {
		t.Fatal("BcryptHashHMAC should not verify hash created by BcryptHash")
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
