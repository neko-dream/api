package messages

import "net/http"

var (
	// Unvoteを投票することはできない
	VoteUnvoteNotAllowed = &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "VOTE-001",
		Message:    "投票は必須です",
	}
	VoteFailed = &APIError{
		StatusCode: http.StatusInternalServerError,
		Code:       "VOTE-003",
		Message:    "投票に失敗しました。時間をおいて再度お試しください",
	}
)
