package talksession_consent

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type TalkSessionConsentRepository interface {
	Store(context.Context, TalkSessionConsent) error
	FindByTalkSessionIDAndUserID(context.Context, shared.UUID[talksession.TalkSession], shared.UUID[user.User]) (*TalkSessionConsent, error)
}
