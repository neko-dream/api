package user

type (
	UserDemographics struct {
		yearOfBirth   *YearOfBirth   // ユーザーの生年
		occupation    *Occupation    // ユーザーの職業
		gender        *Gender        // ユーザーの性別
		municipality  *Municipality  // ユーザーの居住地
		householdSize *HouseholdSize // ユーザーの世帯人数
	}
)

func (u *UserDemographics) YearOfBirth() *YearOfBirth {
	return u.yearOfBirth
}

// ユーザーの年齢を返す
func (u *UserDemographics) Age() int {
	return u.yearOfBirth.Age()
}

func (u *UserDemographics) Occupation() *Occupation {
	return u.occupation
}

func (u *UserDemographics) Gender() *Gender {
	return u.gender
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
	yearOfBirth *YearOfBirth,
	occupation *Occupation,
	gender *Gender,
	municipality *Municipality,
	householdSize *HouseholdSize,
) UserDemographics {
	return UserDemographics{
		yearOfBirth:   yearOfBirth,
		occupation:    occupation,
		gender:        gender,
		municipality:  municipality,
		householdSize: householdSize,
	}
}
