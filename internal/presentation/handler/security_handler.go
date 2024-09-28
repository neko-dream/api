package handler

import (
	"context"
	"log"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/presentation/oas"
)

type securityHandler struct {
	session.TokenManager
}

func (s *securityHandler) HandleSessionId(ctx context.Context, operationName string, t oas.SessionId) (context.Context, error) {
	log.Println("HandleSessionId")
	token, err := s.TokenManager.Parse(ctx, t.GetAPIKey())
	if err != nil {
		return ctx, messages.ForbiddenError
	}
	return session.SetSession(ctx, token), nil
}

func NewSecurityHandler(
	tokenManager session.TokenManager,
) oas.SecurityHandler {
	return &securityHandler{
		TokenManager: tokenManager,
	}
}
