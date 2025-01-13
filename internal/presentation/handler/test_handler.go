package handler

import (
	"context"

	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/internal/presentation/oas"
	"go.opentelemetry.io/otel"
)

type testHandler struct {
	*db.DummyInitializer
}

func NewTestHandler(
	dummyInitializer *db.DummyInitializer,
) oas.TestHandler {
	return &testHandler{
		DummyInitializer: dummyInitializer,
	}
}

// DummiInit implements oas.TestHandler.
func (t *testHandler) DummiInit(ctx context.Context) (oas.DummiInitRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "testHandler.DummiInit")
	defer span.End()

	_ = ctx

	t.DummyInitializer.Initialize()

	return &oas.DummiInitOK{}, nil
}

// Test implements oas.TestHandler.
func (t *testHandler) Test(ctx context.Context) (oas.TestRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "testHandler.Test")
	defer span.End()

	_ = ctx

	panic("unimplemented")
}
