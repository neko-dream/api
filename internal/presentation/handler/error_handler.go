package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/pkg/utils"
	onerror "github.com/ogen-go/ogen/ogenerrors"
)

func CustomErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	apiErr := &messages.APIError{}

	switch {
	case errors.Is(err, onerror.ErrSecurityRequirementIsNotSatisfied):
		apiErr = messages.ForbiddenError
	case errors.Is(err, onerror.ErrSecurityRequirementIsNotSatisfied):
		apiErr = messages.ForbiddenError
	case errors.As(err, &apiErr):
	default:
		switch err.(type) {
		case *onerror.DecodeBodyError:
			apiErr = messages.BadRequestError
			apiErr.Message = err.(*onerror.DecodeBodyError).Err.Error()
		case *onerror.DecodeParamError:
			apiErr = messages.BadRequestError
			apiErr.Message = err.(*onerror.DecodeParamError).Err.Error()
		case *onerror.DecodeParamsError:
			apiErr = messages.BadRequestError
			apiErr.Message = err.(*onerror.DecodeParamsError).Err.Error()
		case *onerror.DecodeRequestError:
			apiErr = messages.BadRequestError
			apiErr.Message = err.(*onerror.DecodeRequestError).Err.Error()
		default:
			utils.HandleErrorWithCaller(ctx, err, "failed to handle error", 3)
			apiErr = messages.InternalServerError
		}
	}

	// JSONレスポンスを作成
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiErr.StatusCode)
	if err := json.NewEncoder(w).Encode(apiErr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
