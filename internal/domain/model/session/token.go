package session

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/samber/lo"
)

type (
	sessionContextKey string

	TokenManager interface {
		Generate(ctx context.Context, user user.User, sessionID shared.UUID[Session]) (string, error)
		Parse(ctx context.Context, token string) (*Claim, error)
	}

	Claim struct {
		Sub         string  `json:"sub"` // subject (user)
		Iat         int64   `json:"iat"` // issued at (seconds)
		Exp         int64   `json:"exp"` // expiration time (seconds)
		Jti         string  `json:"jti"` // JWT ID（SessionID）
		Picture     *string `json:"picture,omitempty"`
		DisplayName *string `json:"displayName,omitempty"`
		DisplayID   *string `json:"displayId,omitempty"`
		IsVerify    bool    `json:"is_verify"`
	}
)

func NewClaim(user user.User, sessionID shared.UUID[Session]) Claim {
	return Claim{
		Sub:       user.UserID().String(),
		Iat:       time.Now().Unix(),
		Exp:       time.Now().Add(24 * time.Hour).Unix(),
		Jti:       sessionID.String(),
		Picture:   user.Picture(),
		DisplayID: user.DisplayID(),
		IsVerify:  user.Verify(),
	}
}

func NewClaimFromMap(claims jwt.MapClaims) Claim {
	var picture, displayName, displayID *string

	if claims["picture"] != nil {
		picture = lo.ToPtr(claims["picture"].(string))
	}
	if claims["displayName"] != nil {
		displayName = lo.ToPtr(claims["displayName"].(string))
	}
	if claims["displayId"] != nil {
		displayID = lo.ToPtr(claims["displayId"].(string))
	}

	return Claim{
		Sub:         claims["sub"].(string),
		Iat:         int64(claims["iat"].(float64)),
		Exp:         int64(claims["exp"].(float64)),
		Jti:         claims["jti"].(string),
		Picture:     picture,
		DisplayName: displayName,
		DisplayID:   displayID,
		IsVerify:    claims["isVerify"].(bool),
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
	Issuer   = "kotohiro.com"
	Audience = "kotohiro.com"
)

func (c *Claim) GenMapClaim() *jwt.MapClaims {
	return &jwt.MapClaims{
		"exp":         c.Exp,
		"iat":         c.Iat,
		"jti":         c.Jti,
		"sub":         c.Sub,
		"iss":         Issuer,
		"aud":         Audience,
		"picture":     c.Picture,
		"displayName": c.DisplayName,
		"displayId":   c.DisplayID,
		"isVerify":    c.IsVerify,
	}
}

// SessionContextKey
var (
	sessKey sessionContextKey = "session"
)

func SetSession(ctx context.Context, claim *Claim) context.Context {
	return context.WithValue(ctx, sessKey, claim)
}

func GetSession(ctx context.Context) *Claim {
	claim, _ := ctx.Value(sessKey).(*Claim)
	return claim
}
