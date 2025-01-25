package dto

import (
	"context"

	"github.com/google/uuid"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
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

type UserDemographic struct {
	UserDemographicID uuid.UUID
	UserID            shared.UUID[user.User]
	YearOfBirth       *int
	Occupation        *int
	Gender            *int
	City              *string
	Prefecture        *string
	HouseholdSize     *int
}

func (u *UserDemographic) GenderString() *string {
	if u.Gender == nil {
		return nil
	}

	str := user.Gender(*u.Gender).String()
	if str == "" {
		return nil
	}
	return &str
}

func (u *UserDemographic) OccupationString() string {
	if u.Occupation == nil {
		return user.OccupationOther.String()
	}
	return user.Occupation(*u.Occupation).String()
}

func (u *UserDemographic) Age(ctx context.Context) *int {
	ctx, span := otel.Tracer("dto").Start(ctx, "UserDemographic.Age")
	defer span.End()

	if u.YearOfBirth == nil {
		return nil
	}
	return lo.ToPtr(user.NewYearOfBirth(u.YearOfBirth).Age(ctx))
}

type UserDetail struct {
	User
	UserAuth
	*UserDemographic
}
