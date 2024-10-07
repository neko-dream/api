package talksession

import (
	"fmt"

	"github.com/neko-dream/server/internal/domain/model/shared"
)

type (
	Location struct {
		talkSessionID shared.UUID[TalkSession]
		latitude      float64
		longitude     float64
		city          string
		prefecture    string
	}
)

func NewLocation(
	talkSessionID shared.UUID[TalkSession],
	latitude float64,
	longitude float64,
	city string,
	prefecture string,
) *Location {
	return &Location{
		talkSessionID: talkSessionID,
		latitude:      latitude,
		longitude:     longitude,
		city:          city,
		prefecture:    prefecture,
	}
}

func (l *Location) TalkSessionID() shared.UUID[TalkSession] {
	return l.talkSessionID
}

func (l *Location) Latitude() float64 {
	return l.latitude
}

func (l *Location) Longitude() float64 {
	return l.longitude
}

func (l *Location) City() string {
	return l.city
}

func (l *Location) Prefecture() string {
	return l.prefecture
}

func (l *Location) ToGeographyText() string {
	return fmt.Sprintf("POINT(%f %f)", l.longitude, l.latitude)
}
