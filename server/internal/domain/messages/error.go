package messages

// APIError は API のエラー。StatusCode・Code・Message を応答に使い、Cause に元のエラーを持つ。
// 定義済みのエラー（パッケージ変数）は共有されるので、書き換えずに WithMessage / WithCause で複製して使う。
type APIError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	// Cause は元のエラー。応答には含めず、ログや errors.Is / errors.As で辿るために持つ。
	Cause error `json:"-"`

	// origin は複製元の定義済みのエラー。errors.Is はこれで比較する
	// （Code は重複しているものや空のものがあるため、比較には使わない）。
	origin *APIError
}

func (e *APIError) Error() string {
	if e.Cause != nil {
		if cause := e.Cause.Error(); cause != e.Message {
			return e.Message + ": " + cause
		}
	}
	return e.Message
}

// Unwrap は errors.Is / errors.As で Cause を辿れるようにする。
func (e *APIError) Unwrap() error {
	return e.Cause
}

// Is は同じ定義済みのエラー（またはその複製）なら同じ種類のエラーとみなす。
func (e *APIError) Is(target error) bool {
	t, ok := target.(*APIError)
	if !ok || t == nil || e == nil {
		return false
	}
	return e.root() == t.root()
}

func (e *APIError) root() *APIError {
	if e.origin != nil {
		return e.origin
	}
	return e
}

// WithMessage は文言を差し替えた複製を返す。元のエラーは変更しない。
func (e *APIError) WithMessage(message string) *APIError {
	c := e.clone()
	c.Message = message
	return c
}

// WithCause は元のエラーを Cause に持つ複製を返す。元のエラーは変更しない。
// cause に別の APIError が含まれていても、このエラー（コード・文言）を返す。中のエラーは errors.Is / errors.As で辿れる。
func (e *APIError) WithCause(cause error) *APIError {
	c := e.clone()
	if cause != nil {
		c.Cause = cause
	}
	return c
}

func (e *APIError) clone() *APIError {
	return &APIError{StatusCode: e.StatusCode, Code: e.Code, Message: e.Message, Cause: e.Cause, origin: e.root()}
}
