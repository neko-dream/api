package talksession

import (
	"encoding/json"
	"fmt"

	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/user"
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
	// SurveyAnswered はアンケート回答済みであることを要求する。
	// このキーの IsSatisfied は常に false を返す（実際の判定は TalkSessionAccessControl が
	// SurveyRepository を参照して行う）。
	SurveyAnswered RestrictionAttributeKey = "survey.answered"
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
		SurveyAnswered: {
			Key:         SurveyAnswered,
			Description: "アンケート回答",
			Order:       4,
			// セッション情報が必要なため、ここでは判定できない。
			// 実際の判定は TalkSessionAccessControl 側で行う。
			IsSatisfied: func(_ user.User) bool { return false },
		},
	}
)

func (k *RestrictionAttributeKey) RestrictionAttribute() RestrictionAttribute {
	return RestrictionAttributeKeyMap[*k]
}

func (k *RestrictionAttributeKey) IsValid() error {
	if _, ok := RestrictionAttributeKeyMap[*k]; !ok {
		// 見つからなかったキーの説明文は空なので、キーそのものを文言に入れる
		return fmt.Errorf("%sは不正な値です", *k)
	}
	return nil
}

var (
	// 共有のエラー。文言は書き換えず WithMessage で複製して返す
	ErrInvalidRestrictionAttribute = &messages.APIError{
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
