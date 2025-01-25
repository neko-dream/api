package user

import "github.com/samber/lo"

type Gender int

const (
	GenderMale Gender = iota + 1
	GenderFemale
	GenderOther
)

var (
	GenderMap = map[Gender]string{
		GenderMale:   "男性",
		GenderFemale: "女性",
		GenderOther:  "その他",
	}
)

func (g Gender) String() string {
	str, ok := GenderMap[g]
	if !ok {
		return ""
	}
	return str
}

func NewGender(s *string) *Gender {
	if s == nil {
		return nil
	}
	if *s == "" {
		return nil
	}
	switch *s {
	case "男性":
		return lo.ToPtr(GenderMale)
	case "女性":
		return lo.ToPtr(GenderFemale)
	case "その他":
		return lo.ToPtr(GenderOther)
	default:
		return nil
	}
}
