package auth

import (
	"context"
	"errors"
	"net/http"
	"os"
	"sync"

	"braces.dev/errtrace"
	"github.com/golang-jwt/jwt/v5"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	http_utils "github.com/neko-dream/server/pkg/http"
)

type tokenManager struct {
	secret string
	session.SessionRepository
}

// SetSession implements session.TokenManager.
func (j *tokenManager) SetSession(ctx context.Context) context.Context {
	r := http_utils.GetHTTPRequest(ctx)
	if r == nil {
		return nil
	}
	sessionCookie, err := r.Cookie("SessionId")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil
		}
	}
	// セッションIDを取得
	claim, err := j.Parse(ctx, sessionCookie.Value)
	if err != nil {
		return ctx
	}
	// トークンの有効性を確認
	if claim.IsExpired() {
		return ctx
	}

	// スキップするOperationの場合以外は、ユーザー登録済みか確認
	if !claim.IsVerify {
		return ctx
	}

	sessID, err := claim.SessionID()
	if err != nil {
		return ctx
	}

	// サーバー側でセッションの有効性を確認
	sess, err := j.SessionRepository.FindBySessionID(ctx, sessID)
	if err != nil {
		return ctx
	}
	if sess == nil {
		return ctx
	}

	if !sess.IsActive() {
		return ctx
	}

	return session.SetSession(ctx, claim)
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

func NewTokenManager(
	sessRepo session.SessionRepository,
) session.TokenManager {
	initSecret()
	return &tokenManager{
		secret:            secretSingleton,
		SessionRepository: sessRepo,
	}
}

func NewTokenManagerWithSecret(secret string, sessRepo session.SessionRepository) session.TokenManager {
	return &tokenManager{
		secret:            secret,
		SessionRepository: sessRepo,
	}
}
