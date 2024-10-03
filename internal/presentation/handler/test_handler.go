package handler

import (
	"context"

	"github.com/neko-dream/server/internal/presentation/oas"
)

type testHandler struct {
}

func NewTestHandler() oas.TestHandler {
	return &testHandler{}
}

func (t *testHandler) Test(ctx context.Context, req oas.OptTestReq) (oas.TestRes, error) {

	return &oas.TestOK{
		URL: "https://example.com",
	}, nil
}
