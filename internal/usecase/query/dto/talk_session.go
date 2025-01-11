package dto

import (
	"os/user"
	"time"

	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
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
	OpinionCount int
	User         User
	Latitude     *float64
	Longitude    *float64
}
