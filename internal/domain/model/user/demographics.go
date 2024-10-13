package user

import (
	"github.com/neko-dream/server/internal/domain/model/shared"
)

type (
	UserDemographics struct {
		userDemographicsID shared.UUID[UserDemographics] // ユーザーのデモグラフィックスID
		yearOfBirth        *YearOfBirth                  // ユーザーの生年
		occupation         *Occupation                   // ユーザーの職業
		gender             *Gender                       // ユーザーの性別
		municipality       *Municipality                 // ユーザーの居住地
		householdSize      *HouseholdSize                // ユーザーの世帯人数
	}
)

func (u *UserDemographics) UserDemographicsID() shared.UUID[UserDemographics] {
	return u.userDemographicsID
}

func (u *UserDemographics) YearOfBirth() *YearOfBirth {
	return u.yearOfBirth
}

// ユーザーの年齢を返す
func (u *UserDemographics) Age() int {
	return u.yearOfBirth.Age()
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

func (u *UserDemographics) Municipality() *Municipality {
	return u.municipality
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
	municipality *Municipality,
	householdSize *HouseholdSize,
) UserDemographics {
	return UserDemographics{
		userDemographicsID: userDemographicsID,
		yearOfBirth:        yearOfBirth,
		occupation:         occupation,
		gender:             gender,
		municipality:       municipality,
		householdSize:      householdSize,
	}
}
