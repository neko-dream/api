package handler

import (
	"context"

	"github.com/neko-dream/server/internal/presentation/oas"
	"go.opentelemetry.io/otel"
)

type healthHandler struct{}

func NewHealthHandler() oas.HealthHandler {
	return &healthHandler{}
}

// Health ヘルスチェック
func (h *healthHandler) Health(ctx context.Context) (oas.HealthRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "healthHandler.Health")
	defer span.End()

	_ = ctx

	return &oas.HealthOK{}, nil
}
