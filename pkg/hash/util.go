package hash

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// GenerateRandomToken 指定バイト長のランダムなトークンをbase64で返す
func GenerateRandomToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// BcryptHash pepperを付与してbcryptでハッシュ化する
func BcryptHash(plaintext, pepper string, cost int) (string, error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return "", fmt.Errorf("cost must be between %d and %d", bcrypt.MinCost, bcrypt.MaxCost)
	}
	h, err := bcrypt.GenerateFromPassword([]byte(pepper+plaintext), cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash: %w", err)
	}
	return string(h), nil
}

// BcryptVerify pepperを付与してbcryptで検証する
func BcryptVerify(plaintext, pepper, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pepper+plaintext)) == nil
}

// BcryptHashHMAC pepperをキーとしてHMAC-SHA256で事前ハッシュしてからbcryptでハッシュ化する。
// HMAC-SHA256の出力は常に32バイトのため、bcryptの72バイト上限に依存しない。
// BcryptHash（pepper+plaintextの単純連結）とは互換性がない。
func BcryptHashHMAC(plaintext, pepper string, cost int) (string, error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return "", fmt.Errorf("cost must be between %d and %d", bcrypt.MinCost, bcrypt.MaxCost)
	}
	pre := hmacSHA256(plaintext, pepper)
	h, err := bcrypt.GenerateFromPassword(pre, cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash: %w", err)
	}
	return string(h), nil
}

// BcryptVerifyHMAC BcryptHashHMACで生成したハッシュを検証する
func BcryptVerifyHMAC(plaintext, pepper, hashed string) bool {
	pre := hmacSHA256(plaintext, pepper)
	return bcrypt.CompareHashAndPassword([]byte(hashed), pre) == nil
}

// hmacSHA256 pepperをキーとしてHMAC-SHA256を計算する（32バイト固定）
func hmacSHA256(plaintext, pepper string) []byte {
	mac := hmac.New(sha256.New, sha256sum(pepper))
	mac.Write([]byte(plaintext))
	return mac.Sum(nil)
}

// HMACHash pepperをキーとしてHMAC-SHA256でハッシュ化する
func HMACHash(plaintext, pepper string) (string, error) {
	mac := hmac.New(sha256.New, sha256sum(pepper))
	_, err := mac.Write([]byte(plaintext))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

// HMACVerify pepperをキーとしてHMAC-SHA256で検証する
func HMACVerify(plaintext, pepper, hashed string) bool {
	mac := hmac.New(sha256.New, sha256sum(pepper))
	_, err := mac.Write([]byte(plaintext))
	if err != nil {
		return false
	}
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(hashed), []byte(expected))
}

func sha256sum(s string) []byte {
	r := sha256.Sum256([]byte(s))
	return r[:]
}
