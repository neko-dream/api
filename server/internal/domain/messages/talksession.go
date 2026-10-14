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
	TalkSessionUpdateFailed = &APIError{
		StatusCode: 500,
		Code:       "TALKSESSION-0010",
		Message:    "セッションの更新に失敗しました。",
	}
	TalkSessionAlreadyConsented = &APIError{
		StatusCode: 400,
		Code:       "TALKSESSION-0011",
		Message:    "セッションに既に同意しています。",
	}
	TalkSessionGetConsentFailed = &APIError{
		StatusCode: 500,
		Code:       "TALKSESSION-0012",
		Message:    "セッションの同意取得に失敗しました。",
	}
	InvalidConsentTime = &APIError{
		StatusCode: 400,
		Code:       "TALKSESSION-0013",
		Message:    "同意の時間が不正です。",
	}
	RestrictionIsZero = &APIError{
		StatusCode: 400,
		Code:       "TALKSESSION-0014",
		Message:    "制限が空です。運営にお問い合わせください。",
	}
	InvalidTalkSessionID = &APIError{
		StatusCode: 400,
		Code:       "TALKSESSION-0015",
		Message:    "トークセッションIDが不正です。",
	}
	InvalidUserID = &APIError{
		StatusCode: 400,
		Code:       "TALKSESSION-0016",
		Message:    "ユーザーIDが不正です。",
	}
	TalkSessionIsFinished = &APIError{
		StatusCode: 400,
		Code:       "TALKSESSION-0017",
		Message:    "セッションは終了しています。",
	}
	TalkSessionPassphraseRequired = &APIError{
		StatusCode: 403,
		Code:       "TALKSESSION-0018",
		Message:    "このセッションは合言葉が必要です。",
	}
	TalkSessionPassphraseInvalid = &APIError{
		StatusCode: 403,
		Code:       "TALKSESSION-0019",
		Message:    "合言葉が正しくありません。",
	}
	TalkSessionSurveyAlreadyExists = &APIError{
		StatusCode: 409,
		Code:       "TALKSESSION-0020",
		Message:    "このセッションには既にアンケートが設定されています。",
	}
	TalkSessionSurveyNotFound = &APIError{
		StatusCode: 404,
		Code:       "TALKSESSION-0021",
		Message:    "アンケートが見つかりません。",
	}
	TalkSessionSurveyAlreadySubmitted = &APIError{
		StatusCode: 409,
		Code:       "TALKSESSION-0022",
		Message:    "既にアンケートに回答済みです。",
	}
	TalkSessionSurveyAnswerRequired = &APIError{
		StatusCode: 400,
		Code:       "TALKSESSION-0023",
		Message:    "このセッションに参加するにはアンケートの回答が必要です。",
	}
	TalkSessionSurveyAnswerInvalid = &APIError{
		StatusCode: 400,
		Code:       "TALKSESSION-0024",
		Message:    "アンケートの回答が不正です。",
	}
	TalkSessionSurveyCannotHardDelete = &APIError{
		StatusCode: 409,
		Code:       "TALKSESSION-0025",
		Message:    "既に回答がある項目は削除できません。",
	}
	TalkSessionSurveyQuestionTypeImmutable = &APIError{
		StatusCode: 409,
		Code:       "TALKSESSION-0026",
		Message:    "既に回答がある質問の種類は変更できません。",
	}
	TalkSessionCreationNotAllowed = &APIError{
		StatusCode: 403,
		Code:       "TALKSESSION-0027",
		Message:    "セッションの作成は許可されていません。運営にお問い合わせください。",
	}
	TalkSessionOrganizationMembershipRequired = &APIError{
		StatusCode: 403,
		Code:       "TALKSESSION-0028",
		Message:    "この組織のメンバーではありません。組織を指定せずにログインし直してください。",
	}
)
