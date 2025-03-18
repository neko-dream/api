package messages

import "net/http"

var (
	PolicyAlreadyConsented = &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "POLICY-001",
		Message:    "すでに同意済みです",
	}
	PolicyNotFound = &APIError{
		StatusCode: http.StatusNotFound,
		Code:       "POLICY-002",
		Message:    "ポリシーが見つかりません",
	}
	PolicyFetchFailed = &APIError{
		StatusCode: http.StatusInternalServerError,
		Code:       "POLICY-003",
		Message:    "ポリシーを取得できませんでした。運営までお問い合わせください。",
	}
)
