package user

type (
	UserDemographics struct {
		yearOfBirth   YearOfBirth  // ユーザーの生年
		occupation    Occupation   // ユーザーの職業
		gender        Gender       // ユーザーの性別
		municipality  Municipality // ユーザーの居住地
		HouseholdSize int          // ユーザーの世帯人数
	}
)

func NewUserDemographics(
	yearOfBirth YearOfBirth,
	occupation Occupation,
	gender Gender,
	municipality Municipality,
	householdSize int,
) UserDemographics {
	return UserDemographics{
		yearOfBirth:   yearOfBirth,
		occupation:    occupation,
		gender:        gender,
		municipality:  municipality,
		HouseholdSize: householdSize,
	}
}
