package talksession

import (
	"encoding/json"
	"fmt"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type RestrictionAttribute struct {
	Key         RestrictionAttributeKey
	Description string
	Fn          func(user user.User) bool
}

type RestrictionAttributeKey string

const (
	DemographicsCity          RestrictionAttributeKey = "demographics.city"
	DemographicsPrefecture    RestrictionAttributeKey = "demographics.prefecture"
	DemographicsGender        RestrictionAttributeKey = "demographics.gender"
	DemographicsHouseholdSize RestrictionAttributeKey = "demographics.household_size"
	DemographicsOccupation    RestrictionAttributeKey = "demographics.occupation"
	DemographicsBirth         RestrictionAttributeKey = "demographics.birth"
)

var (
	RestrictionAttributeKeyMap = map[RestrictionAttributeKey]RestrictionAttribute{
		DemographicsCity: {Key: DemographicsCity, Description: "市区町村", Fn: func(user user.User) bool {
			if user.Demographics() == nil {
				return false
			}
			return user.Demographics().City() != nil
		}},
		DemographicsPrefecture: {Key: DemographicsPrefecture, Description: "都道府県", Fn: func(user user.User) bool {
			if user.Demographics() == nil {
				return false
			}
			return user.Demographics().Prefecture() != nil
		}},
		DemographicsGender: {Key: DemographicsGender, Description: "性別", Fn: func(user user.User) bool {
			if user.Demographics() == nil {
				return false
			}
			return user.Demographics().Gender() != nil
		}},
		DemographicsHouseholdSize: {Key: DemographicsHouseholdSize, Description: "世帯人数", Fn: func(user user.User) bool {
			if user.Demographics() == nil {
				return false
			}
			return user.Demographics().HouseholdSize() != nil
		}},
		DemographicsOccupation: {Key: DemographicsOccupation, Description: "職業", Fn: func(user user.User) bool {
			if user.Demographics() == nil {
				return false
			}
			return user.Demographics().Occupation() != nil
		}},
		DemographicsBirth: {Key: DemographicsBirth, Description: "誕生年", Fn: func(user user.User) bool {
			if user.Demographics() == nil {
				return false
			}
			return user.Demographics().YearOfBirth() != nil
		}},
	}
)

func (k *RestrictionAttributeKey) RestrictionAttribute() RestrictionAttribute {
	return RestrictionAttributeKeyMap[*k]
}

func (k *RestrictionAttributeKey) IsValid() bool {
	_, ok := RestrictionAttributeKeyMap[*k]
	return ok
}

var (
	ErrInvalidRestrictionAttribute = messages.APIError{
		Code:       "restriction_attribute_invalid",
		StatusCode: 400,
		Message:    "不正な値です。",
	}
)

type Restrictions []string

func (s *Restrictions) Scan(src interface{}) error {
	if src == nil {
		*s = nil
		return nil
	}

	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	default:
		return fmt.Errorf("unsupported type for StringSlice: %T", src)
	}
}

func (s Restrictions) Value() (interface{}, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}
