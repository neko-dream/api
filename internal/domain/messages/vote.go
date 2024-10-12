package messages

import "net/http"

var (
	// Unvoteを投票することはできない
	VoteUnvoteNotAllowed = &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "VOTE-001",
		Message:    "投票は必須です",
	}
)
