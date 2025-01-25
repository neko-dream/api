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
		occupation        *Occupation                  // ユーザーの職業
		gender            *Gender                      // ユーザーの性別
		city              *City                        // ユーザーの居住地
		householdSize     *HouseholdSize               // ユーザーの世帯人数
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

func (u *UserDemographic) Occupation() *Occupation {
	return u.occupation
}

func (u *UserDemographic) Gender() *Gender {
	return u.gender
}

func (u *UserDemographic) City() *City {
	return u.city
}

func (u *UserDemographic) HouseholdSize() *HouseholdSize {
	return u.householdSize
}

func (u *UserDemographic) ChangeYearOfBirth(yearOfBirth *YearOfBirth) {
	u.yearOfBirth = yearOfBirth
}

func NewUserDemographic(
	ctx context.Context,
	UserDemographicID shared.UUID[UserDemographic],
	yearOfBirth *int,
	occupation *string,
	gender *string,
	city *string,
	householdSize *int,
	prefecture *string,
) UserDemographic {
	ctx, span := otel.Tracer("user").Start(ctx, "NewUserDemographic")
	defer span.End()

	var (
		yearOfBirthOut   *YearOfBirth
		occupationOut    *Occupation
		genderOut        *Gender
		cityOut          *City
		householdSizeOut *HouseholdSize
	)

	// 誕生日のバリデーション
	if yearOfBirth != nil &&
		*yearOfBirth >= 1900 &&
		*yearOfBirth < clock.Now(ctx).Year() {
		yearOfBirthOut = NewYearOfBirth(yearOfBirth)
	}

	// 職業のバリデーション
	if occupation != nil &&
		*occupation != "" {
		occupationOut = NewOccupation(occupation)
	}

	// 性別のバリデーション
	if gender != nil && *gender != "" {
		genderOut = NewGender(gender)
	}

	// 居住地のバリデーション
	if city != nil && *city != "" {
		cityOut = NewCity(city)
	}

	// 世帯人数のバリデーション
	if householdSize != nil && *householdSize > 0 {
		householdSizeOut = NewHouseholdSize(householdSize)
	}

	return UserDemographic{
		UserDemographicID: UserDemographicID,
		yearOfBirth:       yearOfBirthOut,
		occupation:        occupationOut,
		gender:            genderOut,
		city:              cityOut,
		householdSize:     householdSizeOut,
		prefecture:        prefecture,
	}
}
