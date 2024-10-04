package talk_session_usecase

import (
	"context"
	"os/user"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/infrastructure/db"
)

type (
	CreateTalkSessionUseCase interface {
		Execute(context.Context, CreateTalkSessionInput) (AuthLoginOutput, error)
	}

	CreateTalkSessionInput struct {
		Theme   string
		OwnerID string
	}

	AuthLoginOutput struct {
		talksession.TalkSession
	}

	createTalkSessionInteractor struct {
		talksession.TalkSessionRepository
		*db.DBManager
	}
)

func (i *createTalkSessionInteractor) Execute(ctx context.Context, input CreateTalkSessionInput) (AuthLoginOutput, error) {
	var output AuthLoginOutput

	if err := i.ExecTx(ctx, func(ctx context.Context) error {
		talkSession := talksession.NewTalkSession(
			shared.NewUUID[talksession.TalkSession](),
			input.Theme,
			shared.MustParseUUID[user.User](input.OwnerID),
			nil,
		)
		if err := i.TalkSessionRepository.Create(ctx, talkSession); err != nil {
			return messages.TalkSessionCreateFailed
		}

		output.TalkSession = *talkSession
		return nil
	}); err != nil {
		return output, errtrace.Wrap(err)
	}

	return output, nil
}

func NewCreateTalkSessionUseCase(
	talkSessionRepository talksession.TalkSessionRepository,
	DBManager *db.DBManager,
) CreateTalkSessionUseCase {
	return &createTalkSessionInteractor{
		TalkSessionRepository: talkSessionRepository,
		DBManager:             DBManager,
	}
}
