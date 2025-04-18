package user

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"go.opentelemetry.io/otel"
)

type (
	UserDemographic struct {
		UserDemographicID shared.UUID[UserDemographic] // ユーザーのデモグラフィックスID
		dateOfBirth       *DateOfBirth                 // ユーザーの生年
		gender            *Gender                      // ユーザーの性別
		city              *City                        // ユーザーの居住地
		prefecture        *string                      // ユーザーの居住地の都道府県
	}
)

func (u *UserDemographic) ID() shared.UUID[UserDemographic] {
	return u.UserDemographicID
}

func (u *UserDemographic) DateOfBirth() *DateOfBirth {
	return u.dateOfBirth
}

func (u *UserDemographic) Prefecture() *string {
	return u.prefecture
}

// ユーザーの年齢を返す
func (u *UserDemographic) Age(ctx context.Context) int {
	ctx, span := otel.Tracer("user").Start(ctx, "UserDemographic.Age")
	defer span.End()

	return u.dateOfBirth.Age(ctx)
}

func (u *UserDemographic) Gender() *Gender {
	return u.gender
}

func (u *UserDemographic) City() *City {
	return u.city
}

func (u *UserDemographic) ChangeDateOfBirth(dateOfBirth *DateOfBirth) {
	u.dateOfBirth = dateOfBirth
}

func NewUserDemographic(
	ctx context.Context,
	UserDemographicID shared.UUID[UserDemographic],
	dateOfBirth *int,
	gender *string,
	city *string,
	prefecture *string,
) UserDemographic {
	ctx, span := otel.Tracer("user").Start(ctx, "NewUserDemographic")
	defer span.End()

	var (
		dateOfBirthOut *DateOfBirth
		genderOut      *Gender
		cityOut        *City
	)

	// 誕生日のバリデーション
	if dateOfBirth != nil {
		dateOfBirthOut = NewDateOfBirth(dateOfBirth)
	}

	// 性別のバリデーション
	if gender != nil && *gender != "" {
		genderOut = NewGender(gender)
	}

	// 居住地のバリデーション
	if city != nil && *city != "" {
		cityOut = NewCity(city)
	}
	return UserDemographic{
		UserDemographicID: UserDemographicID,
		dateOfBirth:       dateOfBirthOut,
		gender:            genderOut,
		city:              cityOut,
		prefecture:        prefecture,
	}
}
