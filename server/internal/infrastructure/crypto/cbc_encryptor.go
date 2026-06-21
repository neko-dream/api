package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"go.opentelemetry.io/otel"
	"io"
	"strings"
)

// Deprecated: CBCEncryptorは非推奨です。代わりにGCMEncryptorを使用してください。
type CBCEncryptor struct {
	key []byte
}

func NewCBCEncryptor(key []byte) *CBCEncryptor {
	return &CBCEncryptor{
		key: key,
	}
}

// EncryptBytes
func (e *CBCEncryptor) EncryptBytes(ctx context.Context, plaintext []byte) (string, error) {
	ctx, span := otel.Tracer("crypto").Start(ctx, "CBCEncryptor.EncryptBytes")
	defer span.End()

	_ = ctx

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("%w: 暗号化ブロックの作成に失敗しました: %v", ErrEncryption, err)
	}

	// IVを生成
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("%w: 初期化ベクトルの生成に失敗しました: %v", ErrEncryption, err)
	}

	plaintext = pkcs7Pad(plaintext, aes.BlockSize)

	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	// version.暗号文.ivの形式で結合
	result := fmt.Sprintf("%s.%s.%s",
		Version1, // バージョンを直接指定
		base64.StdEncoding.EncodeToString(ciphertext),
		base64.StdEncoding.EncodeToString(iv),
	)

	return result, nil
}

// DecryptBytes 暗号文をバイト列に復号化する基本関数
func (e *CBCEncryptor) DecryptBytes(ctx context.Context, ciphertext string) ([]byte, error) {
	ctx, span := otel.Tracer("crypto").Start(ctx, "CBCEncryptor.DecryptBytes")
	defer span.End()

	_ = ctx

	parts := strings.Split(ciphertext, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidFormat
	}

	version := parts[0]
	if version != Version1 {
		return nil, fmt.Errorf("%w: バージョン: %s", ErrUnsupportedVersion, version)
	}

	// Base64デコード
	encryptedData, iv := parts[1], parts[2]
	encrypted, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("%w: 暗号文のデコードに失敗しました: %v", ErrDecryption, err)
	}
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, fmt.Errorf("%w: 初期化ベクトルのデコードに失敗しました: %v", ErrDecryption, err)
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("%w: 暗号化ブロックの作成に失敗しました: %v", ErrDecryption, err)
	}

	// 復号化
	plaintext := make([]byte, len(encrypted))
	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(plaintext, encrypted)

	// パディングを除去
	return pkcs7Unpad(plaintext), nil
}

// EncryptString 文字列を暗号化
func (e *CBCEncryptor) EncryptString(ctx context.Context, value string) (string, error) {
	ctx, span := otel.Tracer("crypto").Start(ctx, "CBCEncryptor.EncryptString")
	defer span.End()

	return e.EncryptBytes(ctx, []byte(value))
}

// DecryptString 文字列を復号化
func (e *CBCEncryptor) DecryptString(ctx context.Context, ciphertext string) (string, error) {
	ctx, span := otel.Tracer("crypto").Start(ctx, "CBCEncryptor.DecryptString")
	defer span.End()

	plaintext, err := e.DecryptBytes(ctx, ciphertext)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// EncryptInt 整数を暗号化
func (e *CBCEncryptor) EncryptInt(ctx context.Context, value int64) (string, error) {
	ctx, span := otel.Tracer("crypto").Start(ctx, "CBCEncryptor.EncryptInt")
	defer span.End()

	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(value))
	return e.EncryptBytes(ctx, buf)
}

// DecryptInt 整数を復号化
func (e *CBCEncryptor) DecryptInt(ctx context.Context, ciphertext string) (int64, error) {
	ctx, span := otel.Tracer("crypto").Start(ctx, "CBCEncryptor.DecryptInt")
	defer span.End()

	plaintext, err := e.DecryptBytes(ctx, ciphertext)
	if err != nil {
		return 0, err
	}
	if len(plaintext) != 8 {
		return 0, ErrInvalidInteger
	}
	return int64(binary.BigEndian.Uint64(plaintext)), nil
}

// pkcs7Pad PKCS#7パディングを追加
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

// pkcs7Unpad PKCS#7パディングを除去
func pkcs7Unpad(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
