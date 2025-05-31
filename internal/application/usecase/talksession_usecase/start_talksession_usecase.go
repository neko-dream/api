package talksession_usecase

import (
	"context"
	"time"
	"unicode/utf8"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/application/query/dto"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type (
	StartTalkSessionUseCase interface {
		Execute(context.Context, StartTalkSessionUseCaseInput) (StartTalkSessionUseCaseOutput, error)
	}

	StartTalkSessionUseCaseInput struct {
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

	StartTalkSessionUseCaseOutput struct {
		dto.TalkSessionWithDetail
	}

	startTalkSessionHandler struct {
		talksession.TalkSessionRepository
		user.UserRepository
		organization.OrganizationUserRepository
		*db.DBManager
	}
)

func (in *StartTalkSessionUseCaseInput) Validate() error {
	// Themeは20文字
	if utf8.RuneCountInString(in.Theme) > 100 {
		return messages.TalkSessionThemeTooLong
	}
	// Descriptionは400文字
	if in.Description != nil && utf8.RuneCountInString(*in.Description) > 40000 {
		return messages.TalkSessionDescriptionTooLong
	}

	return nil
}

func NewStartTalkSessionUseCase(
	talkSessionRepository talksession.TalkSessionRepository,
	userRepository user.UserRepository,
	organizationUserRepository organization.OrganizationUserRepository,
	DBManager *db.DBManager,
) StartTalkSessionUseCase {
	return &startTalkSessionHandler{
		TalkSessionRepository:      talkSessionRepository,
		UserRepository:             userRepository,
		OrganizationUserRepository: organizationUserRepository,
		DBManager:                  DBManager,
	}
}

func (i *startTalkSessionHandler) Execute(ctx context.Context, input StartTalkSessionUseCaseInput) (StartTalkSessionUseCaseOutput, error) {
	ctx, span := otel.Tracer("talksession_command").Start(ctx, "startTalkSessionHandler.Execute")
	defer span.End()
	// Organizationに所属していないとセッションを開始できない
	orgs, err := i.OrganizationUserRepository.FindByUserID(ctx, input.OwnerID)
	if err != nil {
		utils.HandleError(ctx, err, "OrganizationUserRepository.FindByUserID")
		return StartTalkSessionUseCaseOutput{}, messages.ForbiddenError
	}
	if len(orgs) == 0 {
		return StartTalkSessionUseCaseOutput{}, messages.ForbiddenError
	}

	var output StartTalkSessionUseCaseOutput

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

		if len(input.Restrictions) > 0 {
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
