package messages

var (
	TalkSessionCreateFailed = &APIError{
		StatusCode: 500,
		Code:       "TALKSESSION-0000",
		Message:    "セッションの作成に失敗しました。",
	}
	TalkSessionNotFinished = &APIError{
		StatusCode: 400,
		Code:       "TALKSESSION-0001",
		Message:    "セッションが終了していません。",
	}
	TalkSessionNotOwner = &APIError{
		StatusCode: 400,
		Code:       "TALKSESSION-0002",
		Message:    "セッションのオーナーではありません。",
	}
	TalkSessionConclusionNotSet = &APIError{
		StatusCode: 404,
		Code:       "TALKSESSION-0003",
		Message:    "結論はまだありません。",
	}
	TalkSessionConclusionAlreadySet = &APIError{
		StatusCode: 400,
		Code:       "TALKSESSION-0004",
		Message:    "結論は既に設定されています。",
	}
	TalkSessionNotFound = &APIError{
		StatusCode: 404,
		Code:       "TALKSESSION-0005",
		Message:    "セッションが見つかりません。",
	}
)
