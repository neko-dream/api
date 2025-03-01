package crypto

type Encryptor interface {
	EncryptBytes(plaintext []byte) (string, error)
	DecryptBytes(ciphertext string) ([]byte, error)
	EncryptString(plaintext string) (string, error)
	DecryptString(ciphertext string) (string, error)
	EncryptInt(plaintext int64) (string, error)
	DecryptInt(ciphertext string) (int64, error)
}

type EncryptedValue struct {
	value string
}

func NewEncryptedValue(value string) EncryptedValue {
	return EncryptedValue{value: value}
}

func (e EncryptedValue) Value() string {
	return e.value
}
