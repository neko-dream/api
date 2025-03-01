package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

// GCMEncryptor AES-GCMを使用した暗号化用の構造体
type GCMEncryptor struct {
	key []byte
}

func NewGCMEncryptor(key []byte) *GCMEncryptor {
	// 鍵長の検証
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		panic(fmt.Errorf("%w: %d バイト (16, 24, または 32 バイトである必要があります)", ErrInvalidKeyLength, len(key)))
	}
	return &GCMEncryptor{
		key: key,
	}
}

// EncryptBytes
func (e *GCMEncryptor) EncryptBytes(plaintext []byte) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("%w: 暗号化ブロックの作成に失敗しました: %v", ErrEncryption, err)
	}

	// GCMモードを初期化
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("%w: GCMモードの初期化に失敗しました: %v", ErrEncryption, err)
	}

	// Nonceを生成 rand.Readerを使用しているため衝突確率はあるが、一旦考慮しない。本当はカウントアップすべき？
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("%w: Nonceの生成に失敗しました: %v", ErrEncryption, err)
	}

	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)

	// version.暗号文.nonceの形式で結合
	result := fmt.Sprintf("%s.%s.%s",
		Version1,
		base64.StdEncoding.EncodeToString(ciphertext),
		base64.StdEncoding.EncodeToString(nonce),
	)

	return result, nil
}

// DecryptBytes
func (e *GCMEncryptor) DecryptBytes(ciphertext string) ([]byte, error) {
	parts := strings.Split(ciphertext, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidFormat
	}

	version := parts[0]
	if version != Version1 {
		return nil, fmt.Errorf("%w: バージョン: %s", ErrUnsupportedVersion, version)
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("%w: 暗号化ブロックの作成に失敗しました: %v", ErrDecryption, err)
	}

	// GCMモードを初期化
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("%w: GCMモードの初期化に失敗しました: %v", ErrDecryption, err)
	}
	// Base64デコード
	encryptedData, nonceStr := parts[1], parts[2]
	encrypted, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("%w: 暗号文のデコードに失敗しました: %v", ErrDecryption, err)
	}
	nonce, err := base64.StdEncoding.DecodeString(nonceStr)
	if err != nil {
		return nil, fmt.Errorf("%w: Nonceのデコードに失敗しました: %v", ErrDecryption, err)
	}
	// Nonceのサイズを検証
	if len(nonce) != aesGCM.NonceSize() {
		return nil, fmt.Errorf("%w: 無効なNonceサイズ: %d バイト (12バイトである必要があります)", ErrDecryption, len(nonce))
	}

	// 復号
	plaintext, err := aesGCM.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: 復号化に失敗しました: %v", ErrDecryption, err)
	}

	return plaintext, nil
}

// EncryptString 文字列を暗号化
func (e *GCMEncryptor) EncryptString(value string) (string, error) {
	return e.EncryptBytes([]byte(value))
}

// DecryptString 文字列を復号化
func (e *GCMEncryptor) DecryptString(ciphertext string) (string, error) {
	plaintext, err := e.DecryptBytes(ciphertext)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// EncryptInt 整数を暗号化
func (e *GCMEncryptor) EncryptInt(value int64) (string, error) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(value))
	return e.EncryptBytes(buf)
}

// DecryptInt 整数を復号化
func (e *GCMEncryptor) DecryptInt(ciphertext string) (int64, error) {
	plaintext, err := e.DecryptBytes(ciphertext)
	if err != nil {
		return 0, err
	}
	if len(plaintext) != 8 {
		return 0, ErrInvalidInteger
	}
	return int64(binary.BigEndian.Uint64(plaintext)), nil
}
