package user

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"go.opentelemetry.io/otel"
)

type (
	UserDemographic struct {
		UserDemographicID shared.UUID[UserDemographic] // ユーザーのデモグラフィックスID
		yearOfBirth       *YearOfBirth                 // ユーザーの生年
		gender            *Gender                      // ユーザーの性別
		city              *City                        // ユーザーの居住地
		prefecture        *string                      // ユーザーの居住地の都道府県
	}
)

func (u *UserDemographic) ID() shared.UUID[UserDemographic] {
	return u.UserDemographicID
}

func (u *UserDemographic) YearOfBirth() *YearOfBirth {
	return u.yearOfBirth
}

func (u *UserDemographic) Prefecture() *string {
	return u.prefecture
}

// ユーザーの年齢を返す
func (u *UserDemographic) Age(ctx context.Context) int {
	ctx, span := otel.Tracer("user").Start(ctx, "UserDemographic.Age")
	defer span.End()

	return u.yearOfBirth.Age(ctx)
}

func (u *UserDemographic) Gender() *Gender {
	return u.gender
}

func (u *UserDemographic) City() *City {
	return u.city
}

func (u *UserDemographic) ChangeYearOfBirth(yearOfBirth *YearOfBirth) {
	u.yearOfBirth = yearOfBirth
}

func NewUserDemographic(
	ctx context.Context,
	UserDemographicID shared.UUID[UserDemographic],
	yearOfBirth *int,
	gender *string,
	city *string,
	prefecture *string,
) UserDemographic {
	ctx, span := otel.Tracer("user").Start(ctx, "NewUserDemographic")
	defer span.End()

	var (
		yearOfBirthOut *YearOfBirth
		genderOut      *Gender
		cityOut        *City
	)

	// 誕生日のバリデーション
	if yearOfBirth != nil &&
		*yearOfBirth >= 1900 &&
		*yearOfBirth < clock.Now(ctx).Year() {
		yearOfBirthOut = NewYearOfBirth(yearOfBirth)
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
		yearOfBirth:       yearOfBirthOut,
		gender:            genderOut,
		city:              cityOut,
		prefecture:        prefecture,
	}
}
