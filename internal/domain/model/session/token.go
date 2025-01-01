package session

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/samber/lo"
)

type (
	sessionContextKey string

	TokenManager interface {
		Generate(ctx context.Context, user user.User, sessionID shared.UUID[Session]) (string, error)
		Parse(ctx context.Context, token string) (*Claim, error)
		SetSession(ctx context.Context) context.Context
	}

	Claim struct {
		Sub         string  `json:"sub"` // subject (user)
		Iat         int64   `json:"iat"` // issued at (seconds)
		Exp         int64   `json:"exp"` // expiration time (seconds)
		Jti         string  `json:"jti"` // JWT ID（SessionID）
		IconURL     *string `json:"iconURL,omitempty"`
		DisplayName *string `json:"displayName,omitempty"`
		DisplayID   *string `json:"displayId,omitempty"`
		IsVerify    bool    `json:"isVerify"`
	}
)

func NewClaim(ctx context.Context, user user.User, sessionID shared.UUID[Session]) Claim {
	now := clock.Now(ctx)
	return Claim{
		Sub:         user.UserID().String(),
		Iat:         now.Unix(),
		Exp:         now.Add(time.Second * 60 * 60 * 24 * 7).Unix(),
		Jti:         sessionID.String(),
		IconURL:     user.ProfileIconURL(),
		DisplayID:   user.DisplayID(),
		DisplayName: user.DisplayName(),
		IsVerify:    user.Verify(),
	}
}

func NewClaimFromMap(claims jwt.MapClaims) Claim {
	var picture, displayName, displayID *string

	if claims["iconURL"] != nil {
		picture = lo.ToPtr(claims["iconURL"].(string))
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
		IconURL:     picture,
		DisplayName: displayName,
		DisplayID:   displayID,
		IsVerify:    claims["isVerify"].(bool),
	}
}

func (c *Claim) UserID() (shared.UUID[user.User], error) {
	return shared.ParseUUID[user.User](c.Sub)
}
func (c *Claim) SessionID() (shared.UUID[Session], error) {
	return shared.ParseUUID[Session](c.Jti)
}

func (c *Claim) Audience() string {
	return Audience
}
func (c *Claim) Issuer() string {
	return Issuer
}

func (c *Claim) IsExpired(ctx context.Context) bool {
	return clock.Now(ctx).Unix() > c.Exp
}
func (c *Claim) IssueAt() time.Time {
	return time.Unix(c.Iat, 0)
}
func (c *Claim) ExpiresAt() time.Time {
	return time.Unix(c.Exp, 0)
}

const (
	Issuer   = "https://api.kotohiro.com"
	Audience = "https://api.kotohiro.com"
)

func (c *Claim) GenMapClaim() *jwt.MapClaims {
	return &jwt.MapClaims{
		"exp":         c.Exp,
		"iat":         c.Iat,
		"jti":         c.Jti,
		"sub":         c.Sub,
		"iss":         Issuer,
		"aud":         Audience,
		"iconURL":     c.IconURL,
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
	if ctx == nil {
		return nil
	}

	value := ctx.Value(sessKey)
	if value == nil {
		return nil
	}

	claim, ok := value.(*Claim)
	if !ok {
		return nil
	}

	return claim
}
