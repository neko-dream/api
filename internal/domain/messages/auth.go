package messages

var (
	ForbiddenError = &APIError{
		StatusCode: 403,
		Code:       "AUTH-0000",
		Message:    "ログインしてください。",
	}
	InvalidStateError = &APIError{
		StatusCode: 400,
		Code:       "AUTH-0001",
		Message:    "リクエストが不正です。",
	}
	InvalidProviderError = &APIError{
		StatusCode: 400,
		Code:       "AUTH-0002",
		Message:    "認証プロバイダーが不正です。",
	}
	InvalidUserIDError = &APIError{
		StatusCode: 400,
		Code:       "AUTH-0003",
		Message:    "ユーザIDが不正です。",
	}
	TokenExpiredError = &APIError{
		StatusCode: 401,
		Code:       "AUTH-0004",
		Message:    "トークンが期限切れです。再ログインしてください。",
	}
	TokenGenerateError = &APIError{
		StatusCode: 500,
		Code:       "AUTH-0005",
		Message:    "トークンの生成に失敗しました。",
	}
)
