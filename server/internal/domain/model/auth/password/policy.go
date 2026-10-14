package password_auth

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/neko-dream/api/internal/domain/messages"
)

const minPasswordLength = 12
const bcryptMaxInputBytes = 72

// PasswordPolicy applies to newly registered or changed passwords, not login.
// MaxBytes accounts for the existing pepper+password bcrypt input format.
type PasswordPolicy struct {
	MinLength                   int
	MaxBytes                    int
	ForbidControlCharacters     bool
	ForbidSurroundingWhitespace bool
}

func NewPasswordPolicy(pepperBytes int) (PasswordPolicy, error) {
	maxBytes := bcryptMaxInputBytes - pepperBytes
	if pepperBytes < 0 || maxBytes < minPasswordLength {
		return PasswordPolicy{}, fmt.Errorf("invalid password hashing configuration")
	}
	return PasswordPolicy{
		MinLength:                   minPasswordLength,
		MaxBytes:                    maxBytes,
		ForbidControlCharacters:     true,
		ForbidSurroundingWhitespace: true,
	}, nil
}

func (p PasswordPolicy) Validate(password string) error {
	if utf8.RuneCountInString(password) < p.MinLength {
		return messages.PasswordTooShort
	}
	if len(password) > p.MaxBytes {
		return messages.PasswordTooLong
	}
	if p.ForbidControlCharacters && strings.ContainsFunc(password, unicode.IsControl) {
		return messages.PasswordInvalidCharacter
	}
	// コピー＆ペーストで混入しやすい先頭・末尾の空白（全角空白を含む）を拒否する。
	if p.ForbidSurroundingWhitespace && strings.TrimSpace(password) != password {
		return messages.PasswordInvalidCharacter
	}
	return nil
}
