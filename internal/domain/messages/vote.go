package messages

import "net/http"

var (
	// Unvoteを投票することはできない
	VoteUnvoteNotAllowed = &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "VOTE-001",
		Message:    "投票は必須です",
	}
	VoteAlreadyVoted = &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "VOTE-002",
		Message:    "この意見へはすでに投票しています",
	}
	VoteFailed = &APIError{
		StatusCode: http.StatusInternalServerError,
		Code:       "VOTE-003",
		Message:    "投票に失敗しました。時間をおいて再度お試しください",
	}
)
