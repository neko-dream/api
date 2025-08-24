package messages

var (
	UserNotFoundError = &APIError{
		StatusCode: 404,
		Code:       "USER-0000",
		Message:    "ユーザが見つかりません。再度ログインしてください。",
	}

	UserDisplayIDAlreadyExistsError = &APIError{
		StatusCode: 400,
		Code:       "USER-0001",
		Message:    "そのIDは既に使用されています。",
	}
	UserDisplayIDInvalidError = &APIError{
		StatusCode: 400,
		Code:       "USER-0002",
		Message:    "IDは半角英数字で入力してください。",
	}
	UserDisplayIDTooLong = &APIError{
		StatusCode: 400,
		Code:       "USER-0003",
		Message:    "IDは30文字以内で入力してください。",
	}
	UserDisplayIDTooShort = &APIError{
		StatusCode: 400,
		Code:       "USER-0004",
		Message:    "IDは4文字以上で入力してください。",
	}
	UserUpdateError = &APIError{
		StatusCode: 500,
		Code:       "USER-0005",
		Message:    "ユーザ情報の更新に失敗しました。",
	}
	UserNotFound = &APIError{
		StatusCode: 404,
		Code:       "USER-0006",
		Message:    "ユーザが見つかりません。",
	}
	UserDisplayNameTooShort = &APIError{
		StatusCode: 400,
		Code:       "USER-0007",
		Message:    "ユーザ名は1文字以上で入力してください。",
	}

	// 退会関連のエラー
	UserAlreadyWithdrawn = &APIError{
		StatusCode: 400,
		Code:       "USER-0011",
		Message:    "既に退会済みです",
	}
	UserNotWithdrawn = &APIError{
		StatusCode: 400,
		Code:       "USER-0012",
		Message:    "このアカウントは退会していません",
	}
	UserReactivationPeriodExpired = &APIError{
		StatusCode: 403,
		Code:       "USER-0013",
		Message:    "復活可能期間（30日）を過ぎています",
	}
)
