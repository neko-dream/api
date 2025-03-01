package crypto

type Encryptor interface {
	EncryptBytes(plaintext []byte) (string, error)
	DecryptBytes(ciphertext string) ([]byte, error)
	EncryptString(plaintext string) (string, error)
	DecryptString(ciphertext string) (string, error)
	EncryptInt(plaintext int64) (string, error)
	DecryptInt(ciphertext string) (int64, error)
}
