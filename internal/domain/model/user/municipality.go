package user

type Municipality string

func (m Municipality) String() string {
	return string(m)
}

func NewMunicipality(municipality string) Municipality {
	return Municipality(municipality)
}
