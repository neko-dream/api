package auth

import (
	"context"
	"os"
	"sync"

	"braces.dev/errtrace"
	"github.com/golang-jwt/jwt/v5"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type tokenManager struct {
	secret string
}

func (j *tokenManager) Generate(ctx context.Context, user user.User, sessionID shared.UUID[session.Session]) (string, error) {
	claim := session.NewClaim(user, sessionID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim.GenMapClaim())
	return errtrace.Wrap2(token.SignedString([]byte(j.secret)))
}

func (j *tokenManager) Parse(ctx context.Context, token string) (*session.Claim, error) {

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	mapClaim, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errtrace.Wrap(jwt.ErrInvalidKeyType)
	}
	claim := session.NewClaimFromMap(mapClaim)
	return &claim, nil
}

var (
	secretReadOnce  sync.Once
	secretSingleton string
)

func initSecret() {
	secretReadOnce.Do(func() {
		secretSingleton = os.Getenv("TOKEN_SECRET")
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
