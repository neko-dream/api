package talksession

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type RestrictionAttribute struct {
	Key         RestrictionAttributeKey
	Description string
	Order       int
	DependsOn   []RestrictionAttributeKey
	// IsSatisfied ユーザーが条件を満たしているかを判定する
	IsSatisfied func(user user.User) bool
}

type RestrictionAttributeKey string

const (
	DemographicsCity       RestrictionAttributeKey = "demographics.city"
	DemographicsPrefecture RestrictionAttributeKey = "demographics.prefecture"
	DemographicsGender     RestrictionAttributeKey = "demographics.gender"
	DemographicsBirth      RestrictionAttributeKey = "demographics.birth"
)

var (
	RestrictionAttributeKeyMap = map[RestrictionAttributeKey]RestrictionAttribute{
		DemographicsGender: {
			Key:         DemographicsGender,
			Description: "性別",
			Order:       0,
			IsSatisfied: func(user user.User) bool {
				if user.Demographics() == nil {
					return false
				}
				return user.Demographics().Gender() != nil
			}},
		DemographicsBirth: {
			Key:         DemographicsBirth,
			Description: "生年月日",
			Order:       1,
			IsSatisfied: func(user user.User) bool {
				if user.Demographics() == nil {
					return false
				}
				return user.Demographics().DateOfBirth() != nil
			}},
		DemographicsCity: {
			Key:         DemographicsCity,
			Description: "市区町村",
			Order:       2,
			IsSatisfied: func(user user.User) bool {
				if user.Demographics() == nil {
					return false
				}
				return user.Demographics().City() != nil
			},
			DependsOn: []RestrictionAttributeKey{DemographicsPrefecture},
		},
		DemographicsPrefecture: {
			Key:         DemographicsPrefecture,
			Description: "都道府県",
			Order:       3,
			IsSatisfied: func(user user.User) bool {
				if user.Demographics() == nil {
					return false
				}
				return user.Demographics().Prefecture() != nil
			}},
	}
)

func (k *RestrictionAttributeKey) RestrictionAttribute() RestrictionAttribute {
	return RestrictionAttributeKeyMap[*k]
}

func (k *RestrictionAttributeKey) IsValid() error {
	attr, ok := RestrictionAttributeKeyMap[*k]
	if !ok {
		return errors.New(attr.Description + "が不正な値です")
	}
	return nil
}

var (
	ErrInvalidRestrictionAttribute = messages.APIError{
		Code:       "restriction_attribute_invalid",
		StatusCode: 400,
		Message:    "不正な値です。",
	}
)

type Restrictions []string

func (s *Restrictions) Scan(src any) error {
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

func (s Restrictions) Value() (any, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}
