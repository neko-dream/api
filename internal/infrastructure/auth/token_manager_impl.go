package auth

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type tokenManager struct {
	secret string
}

const (
	Issuer   = "neko-dream"
	Audience = "neko-dream"
)

func (j *tokenManager) GenerateToken(ctx context.Context, userID shared.UUID[user.User], sessionID shared.UUID[session.Session]) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID.String(),
		"iss": Issuer,
		"aud": Audience,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
		"jti": sessionID.String(),
	})
	return token.SignedString([]byte(j.secret))
}

func (j *tokenManager) ParseToken(ctx context.Context, token string) (*session.Claim, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrInvalidKeyType
	}
	return &session.Claim{
		Sub: claims["sub"].(string),
		Iat: int64(claims["iat"].(float64)),
		Exp: int64(claims["exp"].(float64)),
		Jti: claims["jti"].(string),
	}, nil
}

var (
	secretReadOnce  sync.Once
	secretSingleton string
)

func initSecret() {
	secretReadOnce.Do(func() {
		secretSingleton = os.Getenv("JWT_SECRET")
	})
}

func NewTokenManager(secret string) session.TokenManager {
	if secret == "" {
		initSecret()
		secret = secretSingleton
	}
	return &tokenManager{
		secret: secret,
	}
}
