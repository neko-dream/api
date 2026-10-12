package crypto

import (
	"errors"
	"fmt"
	"strings"

	"github.com/neko-dream/api/internal/domain/model/crypto"
	"github.com/neko-dream/api/internal/infrastructure/config"
)

var (
	ErrUnsupportedVersion = errors.New("対応していない暗号化バージョンです")
	ErrInvalidFormat      = errors.New("不正な暗号文フォーマットです")
	ErrInvalidInteger     = errors.New("不正な整数フォーマットです")
	ErrEncryption         = errors.New("暗号化エラー")
	ErrDecryption         = errors.New("復号化エラー")
	ErrInvalidKeyLength   = errors.New("キーの長さが不正です")
)

const (
	Version1 = "v1"
)

func NewEncryptor(config *config.Config) (crypto.Encryptor, error) {
	switch config.ENCRYPTION_VERSION {
	case Version1:
		return NewGCMEncryptor([]byte(config.ENCRYPTION_SECRET))
	default:
		return nil, fmt.Errorf("%w: バージョン: %s", ErrUnsupportedVersion, config.ENCRYPTION_VERSION)
	}
}

func GetEncryptorFromCiphertext(ciphertext string, key []byte) (crypto.Encryptor, error) {
	parts := strings.Split(ciphertext, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidFormat
	}

	return NewEncryptor(&config.Config{
		ENCRYPTION_VERSION: parts[0],
		ENCRYPTION_SECRET:  string(key),
	})
}
