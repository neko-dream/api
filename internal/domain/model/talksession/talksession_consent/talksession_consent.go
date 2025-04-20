package talksession_consent

import (
	"time"

	"github.com/neko-dream/server/internal/domain/messages"
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
) (TalkSessionConsent, error) {
    if consentedAt.IsZero() {
        return TalkSessionConsent{}, messages.InvalidConsentTime
    }
	if len(restrictions) == 0 {
		return TalkSessionConsent{}, messages.RestrictionIsZero
	}

    if talkSessionID.IsZero() {
        return TalkSessionConsent{}, messages.InvalidTalkSessionID
    }
    if userID.IsZero() {
        return TalkSessionConsent{}, messages.InvalidUserID
    }

	return TalkSessionConsent{
		TalkSessionID: talkSessionID,
		UserID:        userID,
		ConsentedAt:   consentedAt,
		Restrictions:  restrictions,
	}, nil
}

