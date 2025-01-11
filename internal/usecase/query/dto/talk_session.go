package dto

import (
	"time"

	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type TalkSession struct {
	TalkSessionID    shared.UUID[talksession.TalkSession]
	OwnerID          shared.UUID[user.User]
	Theme            string
	ScheduledEndTime time.Time
	CreatedAt        time.Time
	Description      *string
	City             *string
	Prefecture       *string
}

type TalkSessionWithDetail struct {
	TalkSession
	User         User
	OpinionCount int
	Latitude     *float64
	Longitude    *float64
}

// Latitude, Longitudeはnull, または0の場合はfalseを返す
func (t *TalkSessionWithDetail) HasLocation() bool {
	return t.Latitude != nil && t.Longitude != nil && *t.Latitude != 0 && *t.Longitude != 0
}
