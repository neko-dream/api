package user

import "github.com/samber/lo"

type City string

func (m City) String() string {
	return string(m)
}

func NewCity(city *string) *City {
	if city == nil {
		return nil
	}
	if *city == "" {
		return nil
	}
	return lo.ToPtr(City(*city))
}
