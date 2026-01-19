package user

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	guestPrefix = "ゲスト"
	randomChars = "abcdefghijklmnopqrstuvwxyz0123456789"
)

// GenerateGuestDisplayID 8桁のランダム数字を生成
// 例: 12345678
func GenerateGuestDisplayID() string {
	// 10000000〜99999999の範囲
	num, _ := rand.Int(rand.Reader, big.NewInt(90000000))
	return fmt.Sprintf("%08d", num.Int64()+10000000)
}

// GenerateGuestDisplayName 「ゲスト」+6文字のランダム英数字を生成
// 例: ゲストa1b2c3
func GenerateGuestDisplayName() string {
	b := make([]byte, 6)
	for i := range b {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(randomChars))))
		b[i] = randomChars[idx.Int64()]
	}
	return guestPrefix + string(b)
}
