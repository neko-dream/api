package user

import "github.com/samber/lo"

type Municipality string

func (m Municipality) String() string {
	return string(m)
}

func NewMunicipality(municipality *string) *Municipality {
	if municipality == nil {
		return nil
	}
	if *municipality == "" {
		return nil
	}
	return lo.ToPtr(Municipality(*municipality))
}
