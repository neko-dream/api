package messages

var (
	ActionItemInvalidStatus = &APIError{
		StatusCode: 400,
		Code:       "ACTION_ITEM-0000",
		Message:    "'未着手', '進行中', '完了', '保留', '中止'のみ入力可能です。",
	}
	ActionItemInvalidContent = &APIError{
		StatusCode: 400,
		Code:       "ACTION_ITEM-0001",
		Message:    "タイムラインの内容は1文字以上40文字以下で入力してください。",
	}
	ActionItemInvalidSequence = &APIError{
		StatusCode: 400,
		Code:       "ACTION_ITEM-0002",
		Message:    "シーケンスは0以上の整数で入力してください。",
	}
	ActionItemNotFound = &APIError{
		StatusCode: 404,
		Code:       "ACTION_ITEM-0003",
		Message:    "アクションアイテムが見つかりません。",
	}
)
