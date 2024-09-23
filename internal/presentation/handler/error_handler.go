package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/neko-dream/server/internal/domain/messages"
)

func CustomErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	apiErr := &messages.APIError{}
	statusCode := http.StatusInternalServerError
	if err == nil {
		return
	}

	// エラーを独自型にキャスト
	if errors.As(err, &apiErr) {
		statusCode = apiErr.StatusCode
	} else {
		// キャストできない場合はデフォルトエラーを作成
		apiErr = &messages.APIError{
			Code:    strconv.Itoa(statusCode),
			Message: err.Error(),
		}
	}

	// JSONレスポンスを作成
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(apiErr)

	// // JSONレスポンスを作成
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(statusCode)
	// json.NewEncoder(w).Encode(myErr)
}
