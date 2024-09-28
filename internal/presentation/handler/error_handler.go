package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/ogen-go/ogen/ogenerrors"
)

func CustomErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	apiErr := &messages.APIError{}

	switch {
	case errors.Is(err, ogenerrors.ErrSecurityRequirementIsNotSatisfied):
		apiErr = messages.ForbiddenError
	case errors.As(err, &apiErr):
	default:
		apiErr = messages.InternalServerError
	}

	// JSONレスポンスを作成
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiErr.StatusCode)
	if err := json.NewEncoder(w).Encode(apiErr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
