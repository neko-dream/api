package dto

import (
	"context"

	"github.com/google/uuid"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type User struct {
	DisplayID   string
	DisplayName string
	IconURL     *string
}

type UserAuth struct {
	UserAuthID    uuid.UUID
	Provider      string
	IsVerified    bool
	Email         *string
	EmailVerified bool
}

type UserDemographic struct {
	UserDemographicID uuid.UUID
	UserID            shared.UUID[user.User]
	DateOfBirth       *int
	Gender            *int
	City              *string
	Prefecture        *string
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

func (u *UserDemographic) Age(ctx context.Context) *int {
	ctx, span := otel.Tracer("dto").Start(ctx, "UserDemographic.Age")
	defer span.End()

	if u.DateOfBirth == nil {
		return nil
	}
	return lo.ToPtr(user.NewDateOfBirth(u.DateOfBirth).Age(ctx))
}

type UserDetail struct {
	User
	UserAuth
	*UserDemographic
}

func (u *User) ToResponse() oas.User {
	return oas.User{
		DisplayID:   u.DisplayID,
		DisplayName: u.DisplayName,
		IconURL:     utils.ToOptNil[oas.OptNilString](u.IconURL),
	}
}
