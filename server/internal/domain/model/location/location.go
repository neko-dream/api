package location

type (
	Location struct {
		ID         string
		Prefecture string
		City       string
		Town       string
		Postal     string
		Latitude   float64
		Longitude  float64
	}

	Coordinates struct {
		Latitude  float64
		Longitude float64
	}

	LocationRepository interface {
		FindByPostal(postal string) ([]Location, error)
		FindByCoordinates(coords Coordinates) ([]Location, error)
	}
)
