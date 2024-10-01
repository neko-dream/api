package user

import (
	"time"

	"github.com/samber/lo"
)

type YearOfBirth int

func NewYearOfBirth(year *int) *YearOfBirth {
	if year == nil {
		return nil
	}
	if *year == 0 {
		return nil
	}
	if *year < 1900 {
		return nil
	}
	return lo.ToPtr(YearOfBirth(*year))

}

func (y YearOfBirth) Age() int {
	return time.Now().Year() - int(y)
}
