package session

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type (
	TokenManager interface {
		GenerateToken(ctx context.Context, user user.User, sessionID shared.UUID[Session]) (string, error)
		ParseToken(ctx context.Context, token string) (*Claim, error)
	}

	Claim struct {
		Sub         string `json:"sub"` // subject (user)
		Iat         int64  `json:"iat"` // issued at (seconds)
		Exp         int64  `json:"exp"` // expiration time (seconds)
		Jti         string `json:"jti"` // JWT ID（SessionID）
		Picture     string `json:"picture"`
		DisplayName string `json:"display_name"`
		DisplayID   string `json:"display_id"`
	}
)

func NewClaim(user user.User, sessionID shared.UUID[Session]) Claim {
	picture := ""
	if user.Picture() != nil {
		picture = *user.Picture()
	}

	return Claim{
		Sub:         user.UserID().String(),
		Iat:         time.Now().Unix(),
		Exp:         time.Now().Add(24 * time.Hour).Unix(),
		Jti:         sessionID.String(),
		Picture:     picture,
		DisplayName: user.DisplayName(),
		DisplayID:   user.DisplayID(),
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

const (
	Issuer   = "neko-dream"
	Audience = "neko-dream"
)

func (c *Claim) GenMapClaim() *jwt.MapClaims {
	return &jwt.MapClaims{
		"exp":          c.Exp,
		"iat":          c.Iat,
		"jti":          c.Jti,
		"sub":          c.Sub,
		"iss":          Issuer,
		"aud":          Audience,
		"picture":      c.Picture,
		"display_name": c.DisplayName,
		"display_id":   c.DisplayID,
	}
}
