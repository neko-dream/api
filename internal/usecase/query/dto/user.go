package dto

import (
	"context"

	"github.com/google/uuid"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/samber/lo"
)

type User struct {
	DisplayID   string
	DisplayName string
	IconURL     *string
}

type UserAuth struct {
	UserAuthID uuid.UUID
	Provider   string
	IsVerified bool
}

type UserDemographics struct {
	UserDemographicsID uuid.UUID
	UserID             shared.UUID[user.User]
	YearOfBirth        *int
	Occupation         *int
	Gender             int
	City               *string
	Prefecture         *string
	HouseholdSize      *int
}

func (u *UserDemographics) GenderString() string {
	return user.Gender(u.Gender).String()
}

func (u *UserDemographics) OccupationString() string {
	if u.Occupation == nil {
		return user.OccupationOther.String()
	}
	return user.Occupation(*u.Occupation).String()
}

func (u *UserDemographics) Age(ctx context.Context) *int {
	if u.YearOfBirth == nil {
		return nil
	}
	return lo.ToPtr(user.NewYearOfBirth(u.YearOfBirth).Age(ctx))
}

type UserDetail struct {
	User
	UserAuth
	*UserDemographics
}
