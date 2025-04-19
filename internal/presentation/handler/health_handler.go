package handler

import (
	"context"

	"github.com/neko-dream/server/internal/presentation/oas"
)

type healthHandler struct{}


func NewHealthHandler() oas.HealthHandler {
	return &healthHandler{}
}


// Health ヘルスチェック
func (h *healthHandler) Health(ctx context.Context) (oas.HealthRes, error) {
	return &oas.HealthOK{}, nil
}
