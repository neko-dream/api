package admin

import (
	"context"
	"time"

	"github.com/neko-dream/api/internal/domain/model/shared"
)

// AdminSession admin のリフレッシュトークン(revocable)。
// TokenHash は HMAC(refresh_token, pepper)。生トークンは保持しない。
type AdminSession struct {
	AdminSessionID shared.UUID[AdminSession]
	AdminUserID    shared.UUID[AdminUser]
	TokenHash      string
	ExpiresAt      time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (s *AdminSession) IsExpired(now time.Time) bool {
	return !now.Before(s.ExpiresAt)
}

type AdminSessionRepository interface {
	Create(ctx context.Context, session AdminSession) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*AdminSession, error)
	Rotate(ctx context.Context, id shared.UUID[AdminSession], oldTokenHash, newTokenHash string, expiresAt time.Time) (bool, error)
	DeleteByTokenHash(ctx context.Context, tokenHash string) error
}
