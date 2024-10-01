package user

import "github.com/samber/lo"

type HouseholdSize int

func NewHouseholdSize(size *int) *HouseholdSize {
	if size == nil {
		return nil
	}
	if *size == 0 {
		return nil
	}
	return lo.ToPtr(HouseholdSize(*size))
}
