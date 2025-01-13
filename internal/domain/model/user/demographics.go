package user

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"go.opentelemetry.io/otel"
)

type (
	UserDemographics struct {
		userDemographicsID shared.UUID[UserDemographics] // ユーザーのデモグラフィックスID
		yearOfBirth        *YearOfBirth                  // ユーザーの生年
		occupation         *Occupation                   // ユーザーの職業
		gender             *Gender                       // ユーザーの性別
		city               *City                         // ユーザーの居住地
		householdSize      *HouseholdSize                // ユーザーの世帯人数
		prefecture         *string                       // ユーザーの居住地の都道府県
	}
)

func (u *UserDemographics) UserDemographicsID() shared.UUID[UserDemographics] {
	return u.userDemographicsID
}

func (u *UserDemographics) YearOfBirth() *YearOfBirth {
	return u.yearOfBirth
}

func (u *UserDemographics) Prefecture() *string {
	return u.prefecture
}

// ユーザーの年齢を返す
func (u *UserDemographics) Age(ctx context.Context) int {
	ctx, span := otel.Tracer("user").Start(ctx, "UserDemographics.Age")
	defer span.End()

	return u.yearOfBirth.Age(ctx)
}

func (u *UserDemographics) Occupation() Occupation {
	if u.occupation == nil {
		return OccupationOther
	}
	return *u.occupation
}

func (u *UserDemographics) Gender() Gender {
	if u.gender == nil {
		return GenderPreferNotToSay
	}
	return *u.gender
}

func (u *UserDemographics) City() *City {
	return u.city
}

func (u *UserDemographics) HouseholdSize() *HouseholdSize {
	return u.householdSize
}

func (u *UserDemographics) ChangeYearOfBirth(yearOfBirth *YearOfBirth) {
	u.yearOfBirth = yearOfBirth
}

func NewUserDemographics(
	userDemographicsID shared.UUID[UserDemographics],
	yearOfBirth *YearOfBirth,
	occupation *Occupation,
	gender *Gender,
	city *City,
	householdSize *HouseholdSize,
	prefecture *string,
) UserDemographics {
	return UserDemographics{
		userDemographicsID: userDemographicsID,
		yearOfBirth:        yearOfBirth,
		occupation:         occupation,
		gender:             gender,
		city:               city,
		householdSize:      householdSize,
		prefecture:         prefecture,
	}
}
