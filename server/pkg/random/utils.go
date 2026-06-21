package random

import (
	"crypto/rand"
	"fmt"
)

func GenerateRandom() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}
