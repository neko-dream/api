package talksession_command

import (
	"context"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/conclusion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"go.opentelemetry.io/otel"
)

type (
	AddConclusionCommand interface {
		Execute(context.Context, AddConclusionCommandInput) error
	}

	AddConclusionCommandInput struct {
		TalkSessionID shared.UUID[talksession.TalkSession]
		UserID        shared.UUID[user.User]
		Conclusion    string
	}

	addConclusionCommandHandler struct {
		talksession.TalkSessionRepository
		conclusion.ConclusionRepository
	}
)

func NewAddConclusionCommandHandler(
	TalkSessionRepo talksession.TalkSessionRepository,
	concRepo conclusion.ConclusionRepository,
) AddConclusionCommand {
	return &addConclusionCommandHandler{
		TalkSessionRepository: TalkSessionRepo,
		ConclusionRepository:  concRepo,
	}
}

func (i *addConclusionCommandHandler) Execute(ctx context.Context, input AddConclusionCommandInput) error {
	ctx, span := otel.Tracer("talksession_command").Start(ctx, "addConclusionCommandHandler.Execute")
	defer span.End()

	// TalkSessionのドメインロジックのような気もする。
	res, err := i.TalkSessionRepository.FindByID(ctx, input.TalkSessionID)
	if err != nil {
		return err
	}

	// オーナーでなければ結論を作成できない
	if res.OwnerUserID().UUID() != input.UserID.UUID() {
		return messages.TalkSessionNotOwner
	}

	// まだ終了していないトークセッションに対しては結論を作成できない
	if !res.IsFinished(ctx) {
		return messages.TalkSessionNotFinished
	}

	conc, err := i.ConclusionRepository.FindByTalkSessionID(ctx, input.TalkSessionID)
	if err != nil {
		return err
	}
	if conc != nil {
		return messages.TalkSessionConclusionAlreadySet
	}

	conclusion := conclusion.NewConclusion(
		input.TalkSessionID,
		input.Conclusion,
		input.UserID,
	)
	if err := i.ConclusionRepository.Create(ctx, *conclusion); err != nil {
		return err
	}

	return nil
}
