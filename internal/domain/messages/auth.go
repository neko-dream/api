package messages

var (
	InvalidStateError = &APIError{
		StatusCode: 400,
		Code:       "AUTH-0001",
		Message:    "リクエストが不正です。",
	}
	InvalidProviderError = &APIError{
		StatusCode: 400,
		Code:       "AUTH-0002",
		Message:    "認証プロバイダが不正です。",
	}
)
