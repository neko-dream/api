package messages

import "net/http"

var (
	OpinionContentBadLength = &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "OPINION-001",
		Message:    "意見は5~140文字で入力してください",
	}
	OpinionParentOpinionIDIsSame = &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "OPINION-002",
		Message:    "親意見IDと意見IDが同じです",
	}
	OpinionCreateFailed = &APIError{
		StatusCode: http.StatusInternalServerError,
		Code:       "OPINION-003",
		Message:    "意見の投稿に失敗しました。時間をおいて再度お試しください",
	}
	OpinionAlreadyVoted = &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "OPINION-004",
		Message:    "この意見へはすでに投票しています",
	}
	OpinionTitleBadLength = &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "OPINION-005",
		Message:    "タイトルは5~50文字で入力してください",
	}
	OpinionContentBadLengthForUpdate = &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "OPINION-006",
		Message:    "意見は5~140文字で入力してください",
	}
	OpinionContentFailedToFetch = &APIError{
		StatusCode: http.StatusInternalServerError,
		Code:       "OPINION-007",
		Message:    "意見の取得に失敗しました。時間をおいて再度お試しください",
	}
	OpinionReferenceImageUploadFailed = &APIError{
		StatusCode: http.StatusInternalServerError,
		Code:       "OPINION-008",
		Message:    "参考画像のアップロードに失敗しました。時間をおいて再度お試しください",
	}
	OpinionNotFound = &APIError{
		StatusCode: http.StatusNotFound,
		Code:       "OPINION-009",
		Message:    "意見が見つかりません",
	}
	OpinionReportFailed = &APIError{
		StatusCode: http.StatusInternalServerError,
		Code:       "OPINION-010",
		Message:    "意見の通報に失敗しました。時間をおいて再度お試しください",
	}
	OpinionSeedIsOwnerOnly = &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "OPINION-011",
		Message:    "シード意見はセッション成者のみが投票できます",
	}
)
