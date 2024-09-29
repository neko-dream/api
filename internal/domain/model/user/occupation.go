package user

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
	OccupationOther                                  // その他
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
		OccupationOther:            "その他",
	}
)

func (o Occupation) String() string {
	str, ok := OccupationMap[o]
	if !ok {
		return ""
	}
	return str
}

func NewOccupation(occupation string) (Occupation, error) {
	for key, val := range OccupationMap {
		if val == occupation {
			return key, nil
		}
	}
	return OccupationOther, nil
}
