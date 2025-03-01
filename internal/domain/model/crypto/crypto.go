package crypto

import "context"

type Encryptor interface {
	EncryptBytes(ctx context.Context, plaintext []byte) (string, error)
	DecryptBytes(ctx context.Context, ciphertext string) ([]byte, error)
	EncryptString(ctx context.Context, plaintext string) (string, error)
	DecryptString(ctx context.Context, ciphertext string) (string, error)
	EncryptInt(ctx context.Context, plaintext int64) (string, error)
	DecryptInt(ctx context.Context, ciphertext string) (int64, error)
}
