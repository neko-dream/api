package user

import (
	"github.com/samber/lo"
)

type Occupation int

const (
	OccupationFullTimeEmployee Occupation = iota + 1 // 正社員
	OccupationContractEmployee                       // 契約社員
	OccupationPublicServant                          // 公務員
	OccupationSelfEmployed                           // 自営業
	OccupationExecutive                              // 会社役員
	OccupationPartTimeEmployee                       // パート・アルバイト
	OccupationHomemaker                              // 専業主婦
	OccupationStudent                                // 学生
	OccupationUnemployed                             // 無職
	OccupationOther                                  // 無回答
)

var (
	OccupationMap = map[Occupation]string{
		OccupationFullTimeEmployee: "正社員",
		OccupationContractEmployee: "契約社員",
		OccupationPublicServant:    "公務員",
		OccupationSelfEmployed:     "自営業",
		OccupationExecutive:        "会社役員",
		OccupationPartTimeEmployee: "パート・アルバイト",
		OccupationHomemaker:        "家事従事者",
		OccupationStudent:          "学生",
		OccupationUnemployed:       "無職",
		OccupationOther:            "無回答",
	}
)

func (o Occupation) String() string {
	str, ok := OccupationMap[o]
	if !ok {
		return ""
	}
	return str
}

func NewOccupation(occupation *string) *Occupation {
	if occupation == nil {
		return nil
	}
	if *occupation == "" {
		return nil
	}
	for key, val := range OccupationMap {
		if val == *occupation {
			return lo.ToPtr(key)
		}
	}
	return nil
}
