package talk_session_usecase

import (
	"context"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/shared/time"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/db"
)

type (
	CreateTalkSessionUseCase interface {
		Execute(context.Context, CreateTalkSessionInput) (CreateTalkSessionOutput, error)
	}

	CreateTalkSessionInput struct {
		OwnerID          shared.UUID[user.User]
		Theme            string
		ScheduledEndTime time.Time
		Latitude         *float64
		Longitude        *float64
	}

	CreateTalkSessionOutput struct {
		TalkSession talksession.TalkSession
		OwnerUser   *user.User
		Location    *talksession.Location
	}

	createTalkSessionInteractor struct {
		talksession.TalkSessionRepository
		user.UserRepository
		*db.DBManager
	}
)

func (i *createTalkSessionInteractor) Execute(ctx context.Context, input CreateTalkSessionInput) (CreateTalkSessionOutput, error) {
	var output CreateTalkSessionOutput

	if err := i.ExecTx(ctx, func(ctx context.Context) error {
		talkSession := talksession.NewTalkSession(
			shared.NewUUID[talksession.TalkSession](),
			input.Theme,
			input.OwnerID,
			nil,
			time.Now(ctx),
			input.ScheduledEndTime,
			nil,
		)
		if err := i.TalkSessionRepository.Create(ctx, talkSession); err != nil {
			return messages.TalkSessionCreateFailed
		}
		output.TalkSession = *talkSession

		// オーナーのユーザー情報を取得
		ownerUser, err := i.UserRepository.FindByID(ctx, input.OwnerID)
		if err != nil {
			return messages.ForbiddenError
		}
		output.OwnerUser = ownerUser

		return nil
	}); err != nil {
		return output, errtrace.Wrap(err)
	}

	return output, nil
}

func NewCreateTalkSessionUseCase(
	talkSessionRepository talksession.TalkSessionRepository,
	userRepository user.UserRepository,
	DBManager *db.DBManager,
) CreateTalkSessionUseCase {
	return &createTalkSessionInteractor{
		TalkSessionRepository: talkSessionRepository,
		UserRepository:        userRepository,
		DBManager:             DBManager,
	}
}
