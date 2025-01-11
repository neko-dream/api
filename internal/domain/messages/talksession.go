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
	TalkSessionDescriptionTooLong = &APIError{
		StatusCode: 400,
		Code:       "TALKSESSION-0006",
		Message:    "セッションの説明が長すぎます。400文字以内で入力してください。",
	}
	TalkSessionThemeTooLong = &APIError{
		StatusCode: 400,
		Code:       "TALKSESSION-0007",
		Message:    "セッションのテーマが長すぎます。20文字以内で入力してください。",
	}
	InvalidScheduledEndTime = &APIError{
		StatusCode: 400,
		Code:       "TALKSESSION-0008",
		Message:    "終了予定時刻が現在時刻より前です。",
	}
	TalkSessionValidationFailed = &APIError{
		StatusCode: 400,
		Code:       "TALKSESSION-0009",
		Message:    "セッションのバリデーションに失敗しました。",
	}
)
