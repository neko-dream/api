package handler

import (
	"context"
	"net/http"
	"slices"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/presentation/oas"
)

type securityHandler struct {
	session.TokenManager
	session.SessionRepository
}

var skipOperations = []string{
	"RegisterUser",
	"OAuthTokenInfo",
}

func (s *securityHandler) HandleSessionId(ctx context.Context, operationName string, t oas.SessionId) (context.Context, error) {
	// セッションIDを取得
	claim, err := s.TokenManager.Parse(ctx, t.GetAPIKey())
	if err != nil {
		return ctx, messages.ForbiddenError
	}
	// トークンの有効性を確認
	if claim.IsExpired() {
		return ctx, messages.TokenExpiredError
	}

	// スキップするOperationの場合以外は、ユーザー登録済みか確認
	if !claim.IsVerify &&
		!slices.Contains(skipOperations, operationName) {
		return ctx, messages.TokenNotUserRegisteredError
	}

	sessID, err := claim.SessionID()
	if err != nil {
		return ctx, messages.InternalServerError
	}

	// サーバー側でセッションの有効性を確認
	sess, err := s.SessionRepository.FindBySessionID(ctx, sessID)
	if err != nil {
		return ctx, messages.TokenExpiredError
	}
	if sess == nil {
		return ctx, messages.ForbiddenError
	}

	if !sess.IsActive() {
		return ctx, messages.TokenExpiredError
	}

	return session.SetSession(ctx, claim), nil
}

func GetSessionFromHTTPRequest(ctx context.Context, r *http.Request) context.Context {
	// r.Cookies()

	return ctx
}

func NewSecurityHandler(
	tokenManager session.TokenManager,
	sessRepository session.SessionRepository,
) oas.SecurityHandler {
	return &securityHandler{
		TokenManager:      tokenManager,
		SessionRepository: sessRepository,
	}
}
