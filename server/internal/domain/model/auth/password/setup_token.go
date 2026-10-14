package password_auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
)

const SetupTokenTTL = 24 * time.Hour
const SetupTokenResendInterval = time.Minute

// SetupToken is only for initial password setup, never password reset or login.
type SetupToken struct {
	UserID               shared.UUID[user.User]
	TokenHash            []byte
	EmailSubject         string
	CreatedAt, ExpiresAt time.Time
}

type SetupTokenRepository interface {
	// Serialize operator provisioning for an email without rewriting legacy auth rows.
	LockEmailSubject(context.Context, string) error
	// All writers lock the user before reading or changing its setup token.
	LockUser(context.Context, shared.UUID[user.User]) error
	FindUserIDByTokenHash(context.Context, []byte) (shared.UUID[user.User], error)
	FindByUserID(context.Context, shared.UUID[user.User]) (SetupToken, error)
	Save(context.Context, SetupToken) error
	Delete(context.Context, shared.UUID[user.User]) error
}

func SetupTokenHash(token string) ([]byte, bool) {
	if len(token) != 44 {
		return nil, false
	}
	decoded, err := base64.StdEncoding.Strict().DecodeString(token)
	if err != nil || len(decoded) != 32 {
		return nil, false
	}
	sum := sha256.Sum256([]byte(token))
	return sum[:], true
}
