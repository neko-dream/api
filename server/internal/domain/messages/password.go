package messages

var (
	PasswordTooShort = &APIError{StatusCode: 400, Code: "PASSWORD-TOO-SHORT", Message: "パスワードは12文字以上で入力してください。"}
	PasswordTooLong  = &APIError{StatusCode: 400, Code: "PASSWORD-TOO-LONG", Message: "パスワードが長すぎます。短いものにしてください。"}
)
