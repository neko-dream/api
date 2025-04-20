package talksession_consent

import (
	"time"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type TalkSessionConsent struct {
	TalkSessionID shared.UUID[talksession.TalkSession]
	UserID        shared.UUID[user.User]
	ConsentedAt   time.Time
	Restrictions  []string
}

func NewTalkSessionConsent(
	talkSessionID shared.UUID[talksession.TalkSession],
	userID shared.UUID[user.User],
	consentedAt time.Time,
	restrictions []string,
) TalkSessionConsent {
	return TalkSessionConsent{
		TalkSessionID: talkSessionID,
		UserID:        userID,
		ConsentedAt:   consentedAt,
		Restrictions:  restrictions,
	}
}

