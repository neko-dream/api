package password_auth

import (
	"context"
	"time"

	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
)

// Only SHA-256 hashes of the two independent random secrets are persisted.
type LoginChallenge struct {
	UserID                                          shared.UUID[user.User]
	TokenHash, BrowserHash, CredentialHash, PINHash []byte
	CreatedAt, ExpiresAt                            time.Time
	PINAttempts                                     int
	Consumed                                        bool
}

type LoginChallengeRepository interface {
	Save(context.Context, LoginChallenge, time.Time) error
	// Consume must atomically reject expired, consumed and browser-mismatched links.
	Consume(context.Context, []byte, []byte, time.Time) (LoginChallenge, error)
	// FindForPIN locks the browser-bound challenge until its transaction ends.
	FindForPIN(context.Context, []byte) (LoginChallenge, error)
	RecordPINFailure(context.Context, []byte) (int, error)
}
