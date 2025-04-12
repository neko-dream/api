package jwt

import (
	"context"
	"errors"
	"net/http"

	"braces.dev/errtrace"
	"github.com/golang-jwt/jwt/v5"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/config"
	http_utils "github.com/neko-dream/server/pkg/http"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type tokenManager struct {
	secret string
	session.SessionRepository
}

// SetSession implements session.TokenManager.
func (j *tokenManager) SetSession(ctx context.Context) context.Context {
	ctx, span := otel.Tracer("jwt").Start(ctx, "tokenManager.SetSession")
	defer span.End()

	r := http_utils.GetHTTPRequest(ctx)
	if r == nil {
		return ctx
	}

	sessionCookie, err := r.Cookie("SessionId")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return ctx
		}
		// NoCookie以外の場合、エラーをハンドリング
		utils.HandleError(ctx, err, "CookieError")
		return ctx
	}
	// セッションIDを取得
	claim, err := j.Parse(ctx, sessionCookie.Value)
	if err != nil {
		return ctx
	}
	// トークンの有効性を確認 || // スキップするOperationの場合以外は、ユーザー登録済みか確認
	if claim.IsExpired(ctx) || !claim.IsRegistered {
		return ctx
	}

	sessID, err := claim.SessionID()
	if err != nil {
		return ctx
	}

	// サーバー側でセッションの有効性を確認
	sess, err := j.SessionRepository.FindBySessionID(ctx, sessID)
	if err != nil || sess == nil || !sess.IsActive(ctx) {
		return ctx
	}

	return session.SetSession(ctx, claim)
}

func (j *tokenManager) Generate(ctx context.Context, user user.User, sessionID shared.UUID[session.Session]) (string, error) {
	ctx, span := otel.Tracer("jwt").Start(ctx, "tokenManager.Generate")
	defer span.End()

	claim := session.NewClaim(ctx, user, sessionID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim.GenMapClaim())
	return errtrace.Wrap2(token.SignedString([]byte(j.secret)))
}

func (j *tokenManager) Parse(ctx context.Context, token string) (*session.Claim, error) {
	ctx, span := otel.Tracer("jwt").Start(ctx, "tokenManager.Parse")
	defer span.End()

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
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
