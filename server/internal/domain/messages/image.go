package messages

var (
	ImageSizeTooLargeError = &APIError{
		StatusCode: 400,
		Code:       "IMG-0001",
		Message:    "画像サイズが大きすぎます。",
	}
	ImageUnsupportedExtError = &APIError{
		StatusCode: 400,
		Code:       "IMG-0002",
		Message:    "対応していない画像形式です。",
	}
	ImageOpenFailedError = &APIError{
		StatusCode: 400,
		Code:       "IMG-0003",
		Message:    "画像の読み込みに失敗しました。画像が壊れている可能性があります。",
	}
	ImageDecodeFailedError = &APIError{
		StatusCode: 400,
		Code:       "IMG-0004",
		Message:    "画像のデコードに失敗しました。画像が壊れている可能性があります。",
	}

	ImageFilePathEmptyError = &APIError{
		StatusCode: 500,
		Code:       "IMG-0005",
		Message:    "画像のファイルパスが空です。",
	}
	ImageUploadFailedError = &APIError{
		StatusCode: 500,
		Code:       "IMG-0006",
		Message:    "画像のアップロードに失敗しました。",
	}
)
