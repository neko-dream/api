package session

import (
	"context"
	"time"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type (
	TokenManager interface {
		GenerateToken(ctx context.Context, userID shared.UUID[user.User], sessionID shared.UUID[Session]) (string, error)
		ParseToken(ctx context.Context, token string) (*Claim, error)
	}

	Claim struct {
		Sub string `json:"sub"` // subject (user)
		Iat int64  `json:"iat"` // issued at (seconds)
		Exp int64  `json:"exp"` // expiration time (seconds)
		Jti string `json:"jti"` // JWT ID（SessionID）
	}
)

func NewClaim(userID shared.UUID[user.User], sessionID shared.UUID[Session]) Claim {
	return Claim{
		Sub: userID.String(),
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(24 * time.Hour).Unix(),
		Jti: sessionID.String(),
	}
}

func (c *Claim) IsExpired() bool {
	return time.Now().Unix() > c.Exp
}

func (c *Claim) UserID() string {
	return c.Sub
}

func (c *Claim) SessionID() string {
	return c.Jti
}
