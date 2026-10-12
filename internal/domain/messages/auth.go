package messages

var (
	InvalidEmailLoginPINError  = &APIError{StatusCode: 400, Code: "AUTH-0017", Message: "確認コードが正しくないか、期限切れです。メールの6桁の確認コードを確認してください。"}
	EmailLoginPINAttemptsError = &APIError{StatusCode: 429, Code: "AUTH-0018", Message: "確認コードを5回間違えたため、このメールでは確認できません。新しい確認メールを送信してください。"}
	InvalidEmailLoginLinkError = &APIError{StatusCode: 400, Code: "AUTH-0015", Message: "確認リンクが無効です。ログインを開始したブラウザで開くか、再度ログインしてください。"}
	EmailLoginRateLimitError   = &APIError{StatusCode: 429, Code: "AUTH-0016", Message: "確認メールの再送は60秒以上待ってから行ってください。"}
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
	// 認証情報の解析に失敗した場合のエラー
	AuthenticationFailedError = &APIError{
		StatusCode: 401,
		Code:       "AUTH-0011",
		Message:    "認証情報の解析に失敗しました。再度ログインしてください。",
	}
	// セッション情報の解析に失敗した場合のエラー
	SessionParseError = &APIError{
		StatusCode: 401,
		Code:       "AUTH-0012",
		Message:    "セッション情報の解析に失敗しました。再度ログインしてください。",
	}
	// ゲストユーザーがセッション作成を試みた場合のエラー
	GuestCannotCreateSessionError = &APIError{
		StatusCode: 403,
		Code:       "AUTH-0013",
		Message:    "ゲストユーザーはセッションを作成できません。LINEまたはGoogleアカウントでログインしてください。",
	}
	// ゲストユーザーが制限付きセッションに参加を試みた場合のエラー
	GuestCannotJoinRestrictedSessionError = &APIError{
		StatusCode: 403,
		Code:       "AUTH-0014",
		Message:    "ゲストユーザーは制限付きセッションに参加できません。LINEまたはGoogleアカウントでログインしてください。",
	}
)
