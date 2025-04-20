package talksession_command

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/talksession/talksession_consent"
	"github.com/neko-dream/server/internal/domain/model/user"
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
}

func NewTakeConsentUseCase(
	talkSessionConsentService talksession_consent.TalkSessionConsentService,
) TakeConsentUseCase {
	return &takeConsentUseCase{
		talkSessionConsentService: talkSessionConsentService,
	}
}

func (uc *takeConsentUseCase) Execute(ctx context.Context, input TakeConsentUseCaseInput) error {
	ctx, span := otel.Tracer("talksession_command").Start(ctx, "takeConsentUseCase.Execute")
	defer span.End()

	err := uc.talkSessionConsentService.TakeConsent(
		ctx,
		input.TalkSessionID,
		input.UserID,
		[]string{},
	)
	if err != nil {
		return err
	}

	return nil
}
