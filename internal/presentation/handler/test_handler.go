package handler

import (
	"context"

	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"github.com/neko-dream/api/internal/presentation/oas"
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

// DummyInit implements oas.TestHandler.
func (t *testHandler) DummyInit(ctx context.Context, req *oas.DummyInitReq) (oas.DummyInitRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "testHandler.DummyInit")
	defer span.End()

	_ = ctx

	t.DummyInitializer.Initialize()

	return &oas.DummyInitOK{}, nil
}

// Test implements oas.TestHandler.
func (t *testHandler) Test(ctx context.Context) (oas.TestRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "testHandler.Test")
	defer span.End()

	_ = ctx

	panic("unimplemented")
}
