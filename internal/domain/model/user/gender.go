package user

type Gender int

const (
	GenderMale Gender = iota + 1
	GenderFemale
	GenderOther
	GenderPreferNotToSay
)

var (
	GenderMap = map[Gender]string{
		GenderMale:           "男性",
		GenderFemale:         "女性",
		GenderOther:          "その他",
		GenderPreferNotToSay: "未回答",
	}
)

func (g Gender) String() string {
	str, ok := GenderMap[g]
	if !ok {
		return "未回答"
	}
	return str
}

func NewGender(s *string) Gender {
	if s == nil {
		return GenderPreferNotToSay
	}
	if *s == "" {
		return GenderPreferNotToSay
	}
	switch *s {
	case "男性":
		return GenderMale
	case "女性":
		return GenderFemale
	case "その他":
		return GenderOther
	case "未回答":
		return GenderPreferNotToSay
	default:
		return GenderPreferNotToSay
	}
}
