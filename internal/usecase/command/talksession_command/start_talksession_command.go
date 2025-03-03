package talksession_command

import (
	"context"
	"time"
	"unicode/utf8"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/internal/usecase/query/dto"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type (
	StartTalkSessionCommand interface {
		Execute(context.Context, StartTalkSessionCommandInput) (StartTalkSessionCommandOutput, error)
	}

	StartTalkSessionCommandInput struct {
		OwnerID          shared.UUID[user.User]
		Theme            string
		Description      *string
		ThumbnailURL     *string
		ScheduledEndTime time.Time
		Latitude         *float64
		Longitude        *float64
		City             *string
		Prefecture       *string
		Restrictions     []string
	}

	StartTalkSessionCommandOutput struct {
		dto.TalkSessionWithDetail
	}

	startTalkSessionCommandHandler struct {
		talksession.TalkSessionRepository
		user.UserRepository
		*db.DBManager
	}
)

func (in *StartTalkSessionCommandInput) Validate() error {
	// Themeは20文字
	if utf8.RuneCountInString(in.Theme) > 20 {
		return messages.TalkSessionThemeTooLong
	}
	// Descriptionは400文字
	if in.Description != nil && utf8.RuneCountInString(*in.Description) > 400 {
		return messages.TalkSessionDescriptionTooLong
	}

	return nil
}

func NewStartTalkSessionCommand(
	talkSessionRepository talksession.TalkSessionRepository,
	userRepository user.UserRepository,
	DBManager *db.DBManager,
) StartTalkSessionCommand {
	return &startTalkSessionCommandHandler{
		TalkSessionRepository: talkSessionRepository,
		UserRepository:        userRepository,
		DBManager:             DBManager,
	}
}

func (i *startTalkSessionCommandHandler) Execute(ctx context.Context, input StartTalkSessionCommandInput) (StartTalkSessionCommandOutput, error) {
	ctx, span := otel.Tracer("talksession_command").Start(ctx, "startTalkSessionCommandHandler.Execute")
	defer span.End()

	var output StartTalkSessionCommandOutput

	if err := input.Validate(); err != nil {
		return output, errtrace.Wrap(err)
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
		if input.ScheduledEndTime.Before(clock.Now(ctx)) {
			return messages.InvalidScheduledEndTime
		}
		talkSession := talksession.NewTalkSession(
			talkSessionID,
			input.Theme,
			input.Description,
			input.ThumbnailURL,
			input.OwnerID,
			clock.Now(ctx),
			input.ScheduledEndTime,
			location,
			input.City,
			input.Prefecture,
		)
		if input.Restrictions != nil {
			if err := talkSession.UpdateRestrictions(ctx, input.Restrictions); err != nil {
				return errtrace.Wrap(err)
			}
		}

		if err := i.TalkSessionRepository.Create(ctx, talkSession); err != nil {
			utils.HandleError(ctx, err, "TalkSessionRepository.Create")
			return messages.TalkSessionCreateFailed
		}

		output.TalkSession = dto.TalkSession{
			TalkSessionID:    talkSessionID,
			Theme:            input.Theme,
			ThumbnailURL:     talkSession.ThumbnailURL(),
			ScheduledEndTime: input.ScheduledEndTime,
			OwnerID:          talkSession.OwnerUserID(),
			CreatedAt:        talkSession.CreatedAt(),
			Description:      input.Description,
			City:             input.City,
			Prefecture:       input.Prefecture,
		}
		output.Latitude = input.Latitude
		output.Longitude = input.Longitude
		if input.Restrictions != nil {
			output.Restrictions = input.Restrictions
		}

		// オーナーのユーザー情報を取得
		ownerUser, err := i.UserRepository.FindByID(ctx, input.OwnerID)
		if err != nil {
			utils.HandleError(ctx, err, "UserRepository.FindByID")
			return messages.ForbiddenError
		}
		output.User = dto.User{
			DisplayID:   *ownerUser.DisplayID(),
			DisplayName: *ownerUser.DisplayName(),
			IconURL:     ownerUser.IconURL(),
		}

		return nil
	}); err != nil {
		return output, errtrace.Wrap(err)
	}

	return output, nil
}
