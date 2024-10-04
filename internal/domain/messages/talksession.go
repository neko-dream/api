package messages

var (
	TalkSessionCreateFailed = &APIError{
		StatusCode: 500,
		Code:       "TALKSESSION-0000",
		Message:    "セッションの作成に失敗しました。",
	}
)
