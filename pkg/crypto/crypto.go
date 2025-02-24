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

	"github.com/pkg/errors"
)

type Version string

const (
	Version1 Version = "v1"
)

var (
	ErrUnsupportedVersion = errors.New("対応していない暗号化バージョンです")
	ErrInvalidFormat      = errors.New("不正な暗号文フォーマットです")
	ErrInvalidInteger     = errors.New("不正な整数フォーマットです")
)

func (v Version) Validate() error {
	switch v {
	case Version1:
		return nil
	default:
		return errors.Wrapf(ErrUnsupportedVersion, "バージョン: %s", v)
	}
}

type Encrypter interface {
	EncryptBytes(plaintext []byte) (string, error)
	DecryptBytes(ciphertext string) ([]byte, error)
	EncryptString(value string) (string, error)
	DecryptString(ciphertext string) (string, error)
	EncryptInt(value int64) (string, error)
	DecryptInt(ciphertext string) (int64, error)
}

func NewEncrypter(version Version, key []byte) (Encrypter, error) {
	if err := version.Validate(); err != nil {
		return nil, err
	}

	switch version {
	case Version1:
		return NewCBCEncrypter(key), nil
	default:
		return nil, errors.Wrapf(ErrUnsupportedVersion, "version: %s", version)
	}
}

func GetEncrypterFromCiphertext(ciphertext string, key []byte) (Encrypter, error) {
	parts := strings.Split(ciphertext, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidFormat
	}

	version := Version(parts[0])
	return NewEncrypter(version, key)
}

type CBCEncrypter struct {
	key []byte
}

func NewCBCEncrypter(key []byte) *CBCEncrypter {
	return &CBCEncrypter{
		key: key,
	}
}

// EncryptBytes バイト列を暗号化する基本関数
func (e *CBCEncrypter) EncryptBytes(plaintext []byte) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", errors.Wrap(err, "暗号化ブロックの作成に失敗しました")
	}

	// IVを生成
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", errors.Wrap(err, "初期化ベクトルの生成に失敗しました")
	}

	// パディング
	plaintext = pkcs7Pad(plaintext, aes.BlockSize)

	// 暗号化
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
func (e *CBCEncrypter) DecryptBytes(ciphertext string) ([]byte, error) {
	parts := strings.Split(ciphertext, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidFormat
	}

	version := Version(parts[0])
	if version != Version1 {
		return nil, errors.Wrapf(ErrUnsupportedVersion, "version: %s", version)
	}

	// Base64デコード
	encryptedData, iv := parts[1], parts[2]
	encrypted, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, errors.Wrap(err, "暗号文のデコードに失敗しました")
	}
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, errors.Wrap(err, "初期化ベクトルのデコードに失敗しました")
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, errors.Wrap(err, "暗号化ブロックの作成に失敗しました")
	}

	// 復号化
	plaintext := make([]byte, len(encrypted))
	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(plaintext, encrypted)

	// パディングを除去
	return pkcs7Unpad(plaintext), nil
}

// EncryptString 文字列を暗号化
func (e *CBCEncrypter) EncryptString(value string) (string, error) {
	return e.EncryptBytes([]byte(value))
}

// DecryptString 文字列を復号化
func (e *CBCEncrypter) DecryptString(ciphertext string) (string, error) {
	plaintext, err := e.DecryptBytes(ciphertext)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// EncryptInt 整数を暗号化
func (e *CBCEncrypter) EncryptInt(value int64) (string, error) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(value))
	return e.EncryptBytes(buf)
}

// DecryptInt 整数を復号化
func (e *CBCEncrypter) DecryptInt(ciphertext string) (int64, error) {
	plaintext, err := e.DecryptBytes(ciphertext)
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
