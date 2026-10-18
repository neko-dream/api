package messages

import (
	"errors"
)

// APIError は API のエラー。StatusCode・Code・Message を応答に使い、Cause に元のエラーを持つ。
// 定義済みのエラー（パッケージ変数）は共有されるので、書き換えずに WithMessage / WithCause で複製して使う。
type APIError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	// Cause は元のエラー。応答には含めず、ログや errors.Is / errors.As で辿るために持つ。
	Cause error `json:"-"`
}

func (e *APIError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

// Unwrap は errors.Is / errors.As で Cause を辿れるようにする。
func (e *APIError) Unwrap() error {
	return e.Cause
}

// Is は Code が同じなら同じ種類のエラーとみなす。
// WithMessage / WithCause は複製を返すので、ポインターではなく Code で比較する。
func (e *APIError) Is(target error) bool {
	t, ok := target.(*APIError)
	return ok && e.Code == t.Code
}

// WithMessage は文言を差し替えた複製を返す。元のエラーは変更しない。
func (e *APIError) WithMessage(message string) *APIError {
	return &APIError{StatusCode: e.StatusCode, Code: e.Code, Message: message, Cause: e.Cause}
}

// WithCause は元のエラーを Cause に持つ複製を返す。
// cause が既に APIError を含む場合は、そちらを返す（二重に包まない）。
func (e *APIError) WithCause(cause error) *APIError {
	if cause == nil {
		return e
	}
	var existing *APIError
	if errors.As(cause, &existing) {
		return existing
	}
	return &APIError{StatusCode: e.StatusCode, Code: e.Code, Message: e.Message, Cause: cause}
}
