package talksession

import (
	"fmt"

	"github.com/neko-dream/api/internal/domain/model/shared"
)

type (
	Location struct {
		talkSessionID shared.UUID[TalkSession]
		latitude      float64
		longitude     float64
	}
)

func NewLocation(
	talkSessionID shared.UUID[TalkSession],
	latitude float64,
	longitude float64,
) *Location {
	return &Location{
		talkSessionID: talkSessionID,
		latitude:      latitude,
		longitude:     longitude,
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

func (l *Location) ToGeographyText() string {
	return fmt.Sprintf("POINT(%f %f)", l.longitude, l.latitude)
}
