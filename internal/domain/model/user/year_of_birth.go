package user

import "time"

type YearOfBirth int

func NewYearOfBirth(year int) YearOfBirth {
	return YearOfBirth(year)
}

func (y YearOfBirth) Age() int {
	return time.Now().Year() - int(y)
}
