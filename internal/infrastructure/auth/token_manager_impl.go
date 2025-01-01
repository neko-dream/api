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
		utils.HandleError(ctx, err, "CookieError")
		return nil
	}
	// セッションIDを取得
	claim, err := j.Parse(ctx, sessionCookie.Value)
	if err != nil {
		return ctx
	}
	// トークンの有効性を確認
	if claim.IsExpired(ctx) {
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

	if !sess.IsActive(ctx) {
		return ctx
	}

	return session.SetSession(ctx, claim)
}

func (j *tokenManager) Generate(ctx context.Context, user user.User, sessionID shared.UUID[session.Session]) (string, error) {
	claim := session.NewClaim(ctx, user, sessionID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim.GenMapClaim())
	return errtrace.Wrap2(token.SignedString([]byte(j.secret)))
}

func (j *tokenManager) Parse(ctx context.Context, token string) (*session.Claim, error) {

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// アルゴリズムの確認
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			utils.HandleError(ctx, jwt.ErrInvalidKeyType, "InvalidKeyType")
			return nil, errtrace.Wrap(jwt.ErrInvalidKeyType)
		}

		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	mapClaim, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errtrace.Wrap(jwt.ErrInvalidKeyType)
	}

	// issuerの検証
	if mapClaim["iss"].(string) != session.Issuer {
		return nil, errtrace.Wrap(errors.New("invalid issuer"))
	}

	claim := session.NewClaimFromMap(mapClaim)
	return &claim, nil
}

func NewTokenManager(
	sessRepo session.SessionRepository,
	conf *config.Config,
) session.TokenManager {
	return &tokenManager{
		secret:            conf.TokenSecret,
		SessionRepository: sessRepo,
	}
}

func NewTokenManagerWithSecret(secret string, sessRepo session.SessionRepository) session.TokenManager {
	return &tokenManager{
		secret:            secret,
		SessionRepository: sessRepo,
	}
}
