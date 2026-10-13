package messages

var (
	ManagedUserEmailExists   = &APIError{StatusCode: 409, Code: "USER-EMAIL-EXISTS", Message: "このメールアドレスのメールログイン用ユーザーは登録済みです。"}
	InvalidPasswordSetupLink = &APIError{StatusCode: 400, Code: "PASSWORD-SETUP-LINK", Message: "設定リンクが無効か、期限切れです。管理者に設定メールの再送を依頼してください。"}
	PasswordSetupNotPending  = &APIError{StatusCode: 409, Code: "PASSWORD-SETUP-NOT-PENDING", Message: "このユーザーには初回パスワード設定メールを再送できません。"}
	PasswordSetupRateLimit   = &APIError{StatusCode: 429, Code: "PASSWORD-SETUP-RATE-LIMIT", Message: "設定メールの再送は60秒以上待ってから行ってください。"}
)
