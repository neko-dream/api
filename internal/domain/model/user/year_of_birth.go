package user

import "github.com/neko-dream/server/internal/domain/messages"

var (
	ErrInvalidYearOfBirth = &messages.APIError{
		StatusCode: 400,
		Code:       "invalid_year_of_birth",
		Message:    "生年月日が不正です",
	}
)

type YearOfBirth int

func NewYearOfBirth(year int) (YearOfBirth, error) {
	if year < 1900 {
		return YearOfBirth(year), ErrInvalidYearOfBirth
	}

	return YearOfBirth(year), nil
}
