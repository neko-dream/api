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
	TakeConsent(context.Context, shared.UUID[talksession.TalkSession], shared.UUID[user.User], talksession.Restrictions) error
	HasConsented(context.Context, shared.UUID[talksession.TalkSession], shared.UUID[user.User]) (bool, error)
}
type talkSessionConsentService struct {
	talkSessionConsentRepository TalkSessionConsentRepository
	talkSessionRep               talksession.TalkSessionRepository
}

func NewTalkSessionConsentService(
	talkSessionConsentRepository TalkSessionConsentRepository,
	talkSessionRep talksession.TalkSessionRepository,
) TalkSessionConsentService {
	return &talkSessionConsentService{
		talkSessionConsentRepository: talkSessionConsentRepository,
		talkSessionRep:               talkSessionRep,
	}
}

func (s *talkSessionConsentService) TakeConsent(
	ctx context.Context,
	talkSessionID shared.UUID[talksession.TalkSession],
	userID shared.UUID[user.User],
	restrictions talksession.Restrictions,
) error {
	ctx, span := otel.Tracer("talksession_consent").Start(ctx, "talkSessionConsentService.TakeConsent")
	defer span.End()

	hasConsented, err := s.HasConsented(ctx, talkSessionID, userID)
	if err != nil {
		utils.HandleError(ctx, err, "TakeConsentできなかった。")
		return err
	}
	if hasConsented {
		return messages.TalkSessionAlreadyConsented
	}
	consent, err := NewTalkSessionConsent(
		talkSessionID,
		userID,
		time.Now(),
		restrictions,
	)
	if err != nil {
		return err
	}
	return s.talkSessionConsentRepository.Store(ctx, consent)
}

func (s *talkSessionConsentService) HasConsented(
	ctx context.Context,
	talkSessionID shared.UUID[talksession.TalkSession],
	userID shared.UUID[user.User],
) (bool, error) {
	ctx, span := otel.Tracer("talksession_consent").Start(ctx, "talkSessionConsentService.HasConsented")
	defer span.End()

	// トークセッションが存在するか確認
	talkSession, err := s.talkSessionRep.FindByID(ctx, talkSessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, messages.TalkSessionNotFound
		}
		utils.HandleError(ctx, err, "セッション取得に失敗。")
		return false, messages.InternalServerError
	}

	// restrictionsがない場合は常にtrueを返す
	if talkSession.Restrictions() == nil || len(talkSession.Restrictions()) == 0 {
		return true, nil
	}
	consents, err := s.talkSessionConsentRepository.FindByTalkSessionIDAndUserID(ctx, talkSessionID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		utils.HandleError(ctx, err, "Consentの取得に失敗しました。")
		return false, messages.TalkSessionGetConsentFailed
	}
	if consents != nil {
		return true, nil
	}

	return false, nil
}
