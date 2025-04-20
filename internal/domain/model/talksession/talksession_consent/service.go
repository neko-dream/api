package talksession_consent

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type TalkSessionConsentService interface {
	TakeConsent(context.Context, shared.UUID[talksession.TalkSession], shared.UUID[user.User], []string) error
	HasConsented(context.Context, shared.UUID[talksession.TalkSession], shared.UUID[user.User]) (bool, error)
}
type talkSessionConsentService struct {
	talkSessionConsentRepository TalkSessionConsentRepository
}

func NewTalkSessionConsentService(
	talkSessionConsentRepository TalkSessionConsentRepository,
) TalkSessionConsentService {
	return &talkSessionConsentService{
		talkSessionConsentRepository: talkSessionConsentRepository,
	}
}

func (s *talkSessionConsentService) TakeConsent(
	ctx context.Context,
	talkSessionID shared.UUID[talksession.TalkSession],
	userID shared.UUID[user.User],
	restrictions []string,
) error {
	ctx, span := otel.Tracer("talksession_consent").Start(ctx, "talkSessionConsentService.TakeConsent")
	defer span.End()

	hasConsented, err := s.HasConsented(ctx, talkSessionID, userID)
    if err != nil {
        utils.HandleError(ctx, err, "TakeConsentできなかった。")
        return messages.TalkSessionGetConsentFailed
    }
	if hasConsented {
		return messages.TalkSessionAlreadyConsented
	}
	consent := NewTalkSessionConsent(
		talkSessionID,
		userID,
		time.Now(),
		restrictions,
	)
	return s.talkSessionConsentRepository.Store(ctx, consent)
}

func (s *talkSessionConsentService) HasConsented(
	ctx context.Context,
	talkSessionID shared.UUID[talksession.TalkSession],
	userID shared.UUID[user.User],
) (bool, error) {
	ctx, span := otel.Tracer("talksession_consent").Start(ctx, "talkSessionConsentService.HasConsented")
	defer span.End()

	consents, err := s.talkSessionConsentRepository.FindByTalkSessionIDAndUserID(ctx, talkSessionID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
            return false, nil
        }
	}
	if consents != nil {
        return true, nil
    }

	return false, nil
}
