package talksession

import "github.com/neko-dream/server/internal/domain/messages"

type RestrictionAttribute struct {
	Key         RestrictionAttributeKey
	Description string
}

type RestrictionAttributeKey string

const (
	DemographicsCity          RestrictionAttributeKey = "demographics.city"
	DemographicsPrefecture    RestrictionAttributeKey = "demographics.prefecture"
	DemographicsGender        RestrictionAttributeKey = "demographics.gender"
	DemographicsHouseholdSize RestrictionAttributeKey = "demographics.household_size"
	DemographicsOccupation    RestrictionAttributeKey = "demographics.occupation"
	DemographicsBirth         RestrictionAttributeKey = "demographics.birth"
	AuthRegister              RestrictionAttributeKey = "auth.register"
)

var (
	RestrictionAttributeKeyMap = map[RestrictionAttributeKey]RestrictionAttribute{
		DemographicsCity:          {Key: DemographicsCity, Description: "市区町村"},
		DemographicsPrefecture:    {Key: DemographicsPrefecture, Description: "都道府県"},
		DemographicsGender:        {Key: DemographicsGender, Description: "性別"},
		DemographicsHouseholdSize: {Key: DemographicsHouseholdSize, Description: "世帯人数"},
		DemographicsOccupation:    {Key: DemographicsOccupation, Description: "職業"},
		DemographicsBirth:         {Key: DemographicsBirth, Description: "誕生年"},
		AuthRegister:              {Key: AuthRegister, Description: "ユーザー登録"},
	}
)

func (k RestrictionAttributeKey) RestrictionAttribute() RestrictionAttribute {
	return RestrictionAttributeKeyMap[k]
}

func (k RestrictionAttributeKey) IsValid() bool {
	_, ok := RestrictionAttributeKeyMap[k]
	return ok
}

var (
	ErrInvalidRestrictionAttribute = messages.APIError{
		Code:       "restriction_attribute_invalid",
		StatusCode: 400,
		Message:    "不正な値です。",
	}
)
