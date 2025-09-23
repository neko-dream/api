package handler

import (
	"context"
	"slices"

	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/session"
	"github.com/neko-dream/api/internal/presentation/oas"
	"go.opentelemetry.io/otel"
)

type securityHandler struct {
	session.TokenManager
	session.SessionRepository
}

var skipOperations = []string{
	"EstablishUser",
	"GetTokenInfo",
	"RevokeToken",
}

var skipOperationsForWithdrawal = []string{
	"ReactivateUser",
	"RevokeToken",
	"DevAuthorize",
	"PasswordLogin",
}

func (s *securityHandler) HandleCookieAuth(ctx context.Context, operationName string, t oas.CookieAuth) (context.Context, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "securityHandler.HandleSessionId")
	defer span.End()
	// セッションIDを取得
	claim, err := s.TokenManager.Parse(ctx, t.GetAPIKey())
	if err != nil {
		return ctx, messages.ForbiddenError
	}
	// トークンの有効性を確認
	if claim.IsExpired(ctx) {
		return ctx, messages.TokenExpiredError
	}

	if claim.IsWithdrawn &&
		!slices.Contains(skipOperationsForWithdrawal, operationName) {
		// 31日以降の場合はTokenExpiredError
		if claim.WithdrawalDate.AddDate(0, 0, 31).Before(clock.Now(ctx)) {
			return ctx, messages.TokenExpiredError
		}

		// 退会済みの場合は30日以内なら復活可能なエラーを返す
		return ctx, messages.UserWithdrawnRecoverableError
	}

	// スキップするOperationの場合以外は、ユーザー登録済みか確認
	if !claim.IsRegistered &&
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

	if !sess.IsActive(ctx) {
		return ctx, messages.TokenExpiredError
	}

	return session.SetSession(ctx, claim), nil
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
