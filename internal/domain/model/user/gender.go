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
		GenderPreferNotToSay: "回答しない",
	}
)

func (g *Gender) String() string {
	str, ok := GenderMap[*g]
	if !ok {
		return ""
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
	for key, val := range GenderMap {
		if val == *s {
			return key
		}
	}
	return GenderPreferNotToSay
}
