package messages

var (
	// 組織に所属していない場合のエラー
	ErrNotInOrganization = &APIError{
		StatusCode: 403,
		Code:       "AUTH-NOT-IN-ORGANIZATION",
		Message:    "組織への所属が必要です",
	}

	ForbiddenError = &APIError{
		StatusCode: 403,
		Code:       "AUTH-0000",
		Message:    "ログインしてください。",
	}
	InvalidStateError = &APIError{
		StatusCode: 400,
		Code:       "AUTH-0001",
		Message:    "無効なstateです",
	}
	ExpiredStateError = &APIError{
		StatusCode: 400,
		Code:       "AUTH-0002",
		Message:    "stateの有効期限が切れています",
	}
	InvalidProviderError = &APIError{
		StatusCode: 400,
		Code:       "AUTH-0003",
		Message:    "認証プロバイダーが不正です。",
	}
	InvalidUserIDError = &APIError{
		StatusCode: 400,
		Code:       "AUTH-0004",
		Message:    "ユーザIDが不正です。",
	}
	TokenExpiredError = &APIError{
		StatusCode: 401,
		Code:       "AUTH-0005",
		Message:    "トークンが期限切れです。再ログインしてください。",
	}
	TokenGenerateError = &APIError{
		StatusCode: 500,
		Code:       "AUTH-0006",
		Message:    "トークンの生成に失敗しました。",
	}
	TokenNotUserRegisteredError = &APIError{
		StatusCode: 401,
		Code:       "AUTH-0007",
		Message:    "ユーザー登録が完了していません。登録を完了してください。",
	}
	InvalidPasswordOrEmailError = &APIError{
		StatusCode: 401,
		Code:       "AUTH-0008",
		Message:    "メールアドレス,IDまたはパスワードが不正です。",
	}
	InvalidPasswordError = &APIError{
		StatusCode: 401,
		Code:       "AUTH-0009",
		Message:    "パスワードが不正です。",
	}
	UserWithdrawnRecoverableError = &APIError{
		StatusCode: 403,
		Code:       "AUTH-0010",
		Message:    "このアカウントは退会済みです。30日以内であれば復活可能です。",
	}
)
