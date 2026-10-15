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
	// PreviousTokenHash / RotatedAt は直前のローテーションで無効になったトークンと、その時刻。
	PreviousTokenHash *string
	RotatedAt         *time.Time
}

// RotationGracePeriod はローテーション後も直前のリフレッシュトークンを受け付ける時間。
// 複数タブがほぼ同時にリフレッシュしても、後から来た側が失敗しないようにする。
const RotationGracePeriod = 30 * time.Second

func (s *AdminSession) IsExpired(now time.Time) bool {
	return !now.Before(s.ExpiresAt)
}

// IsPreviousTokenInGracePeriod は tokenHash が直前のトークンで、ローテーションから猶予期間内かを返す。
func (s *AdminSession) IsPreviousTokenInGracePeriod(tokenHash string, now time.Time) bool {
	return s.PreviousTokenHash != nil && *s.PreviousTokenHash == tokenHash &&
		s.RotatedAt != nil && now.Before(s.RotatedAt.Add(RotationGracePeriod))
}

type AdminSessionRepository interface {
	Create(ctx context.Context, session AdminSession) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*AdminSession, error)
	// Rotate は oldTokenHash を直前のトークンとして残し、newTokenHash に置き換える。
	Rotate(ctx context.Context, id shared.UUID[AdminSession], oldTokenHash, newTokenHash string, expiresAt, rotatedAt time.Time) (bool, error)
	DeleteByTokenHash(ctx context.Context, tokenHash string) error
}
