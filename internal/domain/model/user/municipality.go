package user

type Municipality string

func NewMunicipality(municipality string) Municipality {
	return Municipality(municipality)
}
