package handler

import (
	"context"
	"net/http"

	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/presentation/oas"
	http_utils "github.com/neko-dream/api/pkg/http"
	"go.opentelemetry.io/otel"
)

func (h *userHandler) GetCurrentUser(ctx context.Context) (*oas.CurrentUser, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "userHandler.GetCurrentUser")
	defer span.End()
	if w, ok := ctx.Value(http_utils.HTTPResponseContextKey).(http.ResponseWriter); ok {
		w.Header().Set("Cache-Control", "no-store")
	}
	authCtx, err := h.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}
	u, err := h.currentUser.FindByID(ctx, authCtx.UserID)
	if err != nil {
		return nil, err
	}
	if u == nil || u.User.WithdrawalDate != nil {
		return nil, messages.UserNotFoundError
	}
	return &oas.CurrentUser{
		UserID:      authCtx.UserID.String(),
		User:        u.User.ToResponse(),
		Permissions: oas.UserPermissions{CanCreateTalkSession: u.CanCreateTalkSession},
	}, nil
}
