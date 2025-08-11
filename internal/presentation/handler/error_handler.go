package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/pkg/utils"
	onerror "github.com/ogen-go/ogen/ogenerrors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// CustomErrorHandler アプリケーション全体のエラーハンドリングを行う
func CustomErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "CustomErrorHandler")
	defer span.End()

	span.RecordError(err)
	span.SetAttributes(
		attribute.String("error.type", fmt.Sprintf("%T", err)),
		attribute.String("request.path", r.URL.Path),
		attribute.String("request.method", r.Method),
	)

	apiErr := convertToAPIError(ctx, err)

	writeErrorResponse(ctx, w, apiErr, span)
}

// convertToAPIError エラーを適切なAPIエラーに変換する
func convertToAPIError(ctx context.Context, err error) *messages.APIError {
	apiErr := &messages.APIError{}

	switch {
	case errors.As(err, &apiErr):
		return apiErr
	case errors.Is(err, onerror.ErrSecurityRequirementIsNotSatisfied):
		return messages.ForbiddenError
	case errors.Is(err, context.DeadlineExceeded):
		return &messages.APIError{
			StatusCode: http.StatusRequestTimeout,
			Code:       "REQUEST_TIMEOUT",
			Message:    "リクエストがタイムアウトしました。",
		}
	case errors.Is(err, context.Canceled):
		return &messages.APIError{
			StatusCode: 499, // Client Closed Request
			Code:       "CLIENT_CLOSED_REQUEST",
			Message:    "クライアントがリクエストを閉じました。",
		}
	default:
		// 特定のエラータイプの処理
		return handleSpecificErrors(ctx, err)
	}
}

// handleSpecificErrors 特定のエラータイプを処理する
func handleSpecificErrors(ctx context.Context, err error) *messages.APIError {
	var (
		decodeBodyErr    *onerror.DecodeBodyError
		decodeParamErr   *onerror.DecodeParamError
		decodeParamsErr  *onerror.DecodeParamsError
		decodeRequestErr *onerror.DecodeRequestError
	)

	switch {
	case errors.As(err, &decodeBodyErr):
		utils.HandleErrorWithCaller(ctx, err, "DecodeBodyError", 3)
		return messages.BadRequestError
	case errors.As(err, &decodeParamErr):
		utils.HandleErrorWithCaller(ctx, err, "DecodeParamError", 3)
		return messages.BadRequestError
	case errors.As(err, &decodeParamsErr):
		utils.HandleErrorWithCaller(ctx, err, "DecodeParamsError", 3)
		return messages.BadRequestError
	case errors.As(err, &decodeRequestErr):
		utils.HandleErrorWithCaller(ctx, err, "DecodeRequestError", 3)
		return messages.BadRequestError
	default:
		// 予期しないエラーの場合はログに記録
		utils.HandleErrorWithCaller(ctx, err, "unexpected error in handler", 3)
		return messages.InternalServerError
	}
}

// writeErrorResponse エラーレスポンスを書き込む
func writeErrorResponse(ctx context.Context, w http.ResponseWriter, apiErr *messages.APIError, span trace.Span) {
	span.SetAttributes(
		attribute.Int("response.status_code", apiErr.StatusCode),
		attribute.String("response.error_code", apiErr.Code),
	)

	// ボディを持たないステータスコードの場合は早期リターン
	if shouldSkipBody(apiErr.StatusCode) {
		w.WriteHeader(apiErr.StatusCode)
		return
	}

	// ヘッダーを設定してからWriteHeaderを呼ぶ
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(apiErr.StatusCode)

	if err := json.NewEncoder(w).Encode(apiErr); err != nil {
		// レスポンスの書き込みに失敗した場合のログ
		utils.HandleErrorWithCaller(ctx, err, "failed to encode error response", 1)
	}
}

// shouldSkipBody ステータスコードがレスポンスボディを持つべきでないかを判定
func shouldSkipBody(statusCode int) bool {
	return statusCode == http.StatusNoContent ||
		statusCode == http.StatusNotModified ||
		statusCode == http.StatusContinue ||
		statusCode == http.StatusSwitchingProtocols ||
		statusCode == http.StatusProcessing ||
		(statusCode >= 100 && statusCode < 200)
}
