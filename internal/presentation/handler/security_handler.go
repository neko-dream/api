package handler

import (
	"context"

	"github.com/neko-dream/server/internal/presentation/oas"
)

type securityHandler struct {
}

func (s *securityHandler) HandleSessionId(ctx context.Context, operationName string, t oas.SessionId) (context.Context, error) {
	panic("unimplemented")
}

func NewSecurityHandler() oas.SecurityHandler {
	return &securityHandler{}
}
