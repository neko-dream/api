package talksession_usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/talksession/talksession_consent"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type TakeConsentUseCase interface {
	Execute(ctx context.Context, input TakeConsentUseCaseInput) error
}

type TakeConsentUseCaseInput struct {
	TalkSessionID shared.UUID[talksession.TalkSession]
	UserID        shared.UUID[user.User]
}

type takeConsentUseCase struct {
	talkSessionConsentService talksession_consent.TalkSessionConsentService
	talkSessionRep            talksession.TalkSessionRepository
}

func NewTakeConsentUseCase(
	talkSessionConsentService talksession_consent.TalkSessionConsentService,
	talkSessionRep talksession.TalkSessionRepository,
) TakeConsentUseCase {
	return &takeConsentUseCase{
		talkSessionConsentService: talkSessionConsentService,
		talkSessionRep:            talkSessionRep,
	}
}

func (uc *takeConsentUseCase) Execute(ctx context.Context, input TakeConsentUseCaseInput) error {
	ctx, span := otel.Tracer("talksession_command").Start(ctx, "takeConsentUseCase.Execute")
	defer span.End()

	// セッションが存在するか確認
	talkSession, err := uc.talkSessionRep.FindByID(ctx, input.TalkSessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return messages.TalkSessionNotFound
		}
		utils.HandleError(ctx, err, "セッション取得に失敗。")
		return messages.InternalServerError
	}
	if err := uc.talkSessionConsentService.TakeConsent(
		ctx,
		input.TalkSessionID,
		input.UserID,
		talkSession.RestrictionList(),
	); err != nil {
		utils.HandleError(ctx, err, "Consentの保存に失敗しました。")
		return err
	}

	return nil
}
