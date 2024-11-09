package talk_session_usecase

import (
	"context"
	"time"
	"unicode/utf8"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/db"
	"github.com/neko-dream/server/pkg/utils"
)

type (
	CreateTalkSessionUseCase interface {
		Execute(context.Context, CreateTalkSessionInput) (CreateTalkSessionOutput, error)
	}

	CreateTalkSessionInput struct {
		OwnerID          shared.UUID[user.User]
		Theme            string
		Description      *string
		ScheduledEndTime time.Time
		Latitude         *float64
		Longitude        *float64
		City             *string
		Prefecture       *string
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

	// Themeは20文字
	if utf8.RuneCountInString(input.Theme) > 20 {
		return output, messages.TalkSessionThemeTooLong
	}
	// Descriptionは400文字
	if input.Description != nil && utf8.RuneCountInString(*input.Description) > 400 {
		return output, messages.TalkSessionDescriptionTooLong
	}

	if err := i.ExecTx(ctx, func(ctx context.Context) error {
		talkSessionID := shared.NewUUID[talksession.TalkSession]()
		var location *talksession.Location
		if input.Latitude != nil && input.Longitude != nil {
			location = talksession.NewLocation(
				talkSessionID,
				*input.Latitude,
				*input.Longitude,
			)
		}
		talkSession := talksession.NewTalkSession(
			talkSessionID,
			input.Theme,
			input.Description,
			input.OwnerID,
			clock.Now(ctx),
			input.ScheduledEndTime,
			location,
			input.City,
			input.Prefecture,
		)

		if err := i.TalkSessionRepository.Create(ctx, talkSession); err != nil {
			utils.HandleError(ctx, err, "TalkSessionRepository.Create")
			return messages.TalkSessionCreateFailed
		}

		output.TalkSession = *talkSession

		// オーナーのユーザー情報を取得
		ownerUser, err := i.UserRepository.FindByID(ctx, input.OwnerID)
		if err != nil {
			utils.HandleError(ctx, err, "UserRepository.FindByID")
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
