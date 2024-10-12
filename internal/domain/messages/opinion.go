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
)
