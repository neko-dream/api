package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/neko-dream/server/internal/domain/messages"
	onerror "github.com/ogen-go/ogen/ogenerrors"
	"github.com/stretchr/testify/assert"
)

func TestCustomErrorHandler(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedCode   string
		expectedBody   bool
	}{
		{
			name: "APIError",
			err: &messages.APIError{
				StatusCode: http.StatusBadRequest,
				Code:       "BAD_REQUEST",
				Message:    "Bad request",
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "BAD_REQUEST",
			expectedBody:   true,
		},
		{
			name:           "SecurityRequirementNotSatisfied",
			err:            onerror.ErrSecurityRequirementIsNotSatisfied,
			expectedStatus: http.StatusForbidden,
			expectedCode:   messages.ForbiddenError.Code,
			expectedBody:   true,
		},
		{
			name:           "ContextDeadlineExceeded",
			err:            context.DeadlineExceeded,
			expectedStatus: http.StatusRequestTimeout,
			expectedCode:   "REQUEST_TIMEOUT",
			expectedBody:   true,
		},
		{
			name:           "ContextCanceled",
			err:            context.Canceled,
			expectedStatus: 499,
			expectedCode:   "CLIENT_CLOSED_REQUEST",
			expectedBody:   true,
		},
		{
			name: "DecodeBodyError",
			err: &onerror.DecodeBodyError{
				Err: errors.New("invalid JSON"),
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   messages.BadRequestError.Code,
			expectedBody:   true,
		},
		{
			name: "WrappedDecodeBodyError",
			err: fmt.Errorf("wrapped: %w", &onerror.DecodeBodyError{
				Err: errors.New("invalid JSON"),
			}),
			expectedStatus: http.StatusBadRequest,
			expectedCode:   messages.BadRequestError.Code,
			expectedBody:   true,
		},
		{
			name: "NoContent",
			err: &messages.APIError{
				StatusCode: http.StatusNoContent,
				Code:       "NO_CONTENT",
				Message:    "No content",
			},
			expectedStatus: http.StatusNoContent,
			expectedCode:   "",
			expectedBody:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()

			// Act
			CustomErrorHandler(context.Background(), rec, req, tt.err)

			// Assert
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedBody {
				assert.NotEmpty(t, rec.Body.String())
				assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
			} else {
				// No Content等の場合はボディが空であることを確認
				assert.Empty(t, rec.Body.String())
			}
		})
	}
}

func TestShouldSkipBody(t *testing.T) {
	tests := []struct {
		statusCode int
		shouldSkip bool
	}{
		{http.StatusOK, false},
		{http.StatusNoContent, true},
		{http.StatusNotModified, true},
		{http.StatusContinue, true},
		{http.StatusSwitchingProtocols, true},
		{http.StatusProcessing, true},
		{http.StatusBadRequest, false},
		{http.StatusInternalServerError, false},
	}

	for _, tt := range tests {
		t.Run(http.StatusText(tt.statusCode), func(t *testing.T) {
			assert.Equal(t, tt.shouldSkip, shouldSkipBody(tt.statusCode))
		})
	}
}
