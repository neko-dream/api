package talksession_usecase

import (
	"context"
	"time"
	"unicode/utf8"

	"braces.dev/errtrace"
	"github.com/neko-dream/api/internal/application/query/dto"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/organization"
	"github.com/neko-dream/api/internal/domain/model/session"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"github.com/neko-dream/api/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type (
	StartTalkSessionUseCase interface {
		Execute(context.Context, StartTalkSessionUseCaseInput) (StartTalkSessionUseCaseOutput, error)
	}

	StartTalkSessionUseCaseInput struct {
		OwnerID             shared.UUID[user.User]
		Theme               string
		Description         *string
		ThumbnailURL        *string
		ScheduledEndTime    time.Time
		Latitude            *float64
		Longitude           *float64
		City                *string
		Prefecture          *string
		Restrictions        []string
		SessionClaim        *session.Claim // セッション情報を追加
		OrganizationAliasID *shared.UUID[organization.OrganizationAlias]
		HideTop             *bool   // トップに表示するかどうか
		Passphrase          *string // 合言葉（設定するとプライベートセッションになる）
	}

	StartTalkSessionUseCaseOutput struct {
		dto.TalkSessionWithDetail
	}

	startTalkSessionHandler struct {
		talksession.TalkSessionRepository
		user.UserRepository
		organization.OrganizationUserRepository
		creationPermissions user.SessionCreationPermissionRepository
		organization.OrganizationAliasRepository
		talksession.PassphraseManager
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
	organizationAliasRepository organization.OrganizationAliasRepository,
	passphraseManager talksession.PassphraseManager,
	DBManager *db.DBManager,
	creationPermissions user.SessionCreationPermissionRepository,
) StartTalkSessionUseCase {
	return &startTalkSessionHandler{
		TalkSessionRepository:       talkSessionRepository,
		UserRepository:              userRepository,
		OrganizationUserRepository:  organizationUserRepository,
		OrganizationAliasRepository: organizationAliasRepository,
		PassphraseManager:           passphraseManager,
		DBManager:                   DBManager,
		creationPermissions:         creationPermissions,
	}
}

func (i *startTalkSessionHandler) Execute(ctx context.Context, input StartTalkSessionUseCaseInput) (StartTalkSessionUseCaseOutput, error) {
	ctx, span := otel.Tracer("talksession_command").Start(ctx, "startTalkSessionHandler.Execute")
	defer span.End()

	// ゲストユーザーはセッションを作成できない
	owner, err := i.UserRepository.FindByID(ctx, input.OwnerID)
	if err != nil {
		utils.HandleError(ctx, err, "UserRepository.FindByID")
		return StartTalkSessionUseCaseOutput{}, messages.UserNotFoundError
	}
	if owner == nil {
		return StartTalkSessionUseCaseOutput{}, messages.UserNotFoundError
	}
	if owner.IsGuest() {
		return StartTalkSessionUseCaseOutput{}, messages.GuestCannotCreateSessionError
	}
	granted, err := i.creationPermissions.Get(ctx, input.OwnerID)
	if err != nil {
		return StartTalkSessionUseCaseOutput{}, err
	}
	if !owner.CanCreateTalkSession(granted) {
		return StartTalkSessionUseCaseOutput{}, messages.ForbiddenError
	}

	// セッションから組織IDとエイリアスIDを取得
	var organizationID *shared.UUID[organization.Organization]
	var organizationAliasID *shared.UUID[organization.OrganizationAlias]

	// 明示的にエイリアスIDが指定されている場合はそれを使用
	if input.OrganizationAliasID != nil {
		organizationAliasID = input.OrganizationAliasID
		// エイリアスから組織IDを取得
		alias, err := i.OrganizationAliasRepository.FindByID(ctx, *input.OrganizationAliasID)
		if err != nil {
			return StartTalkSessionUseCaseOutput{}, err
		}
		if alias == nil {
			return StartTalkSessionUseCaseOutput{}, messages.BadRequestError
		}
		orgID := alias.OrganizationID()
		organizationID = &orgID
	} else if input.SessionClaim != nil && input.SessionClaim.OrganizationID != nil {
		// セッションに組織IDがある場合
		orgID, err := shared.ParseUUID[organization.Organization](*input.SessionClaim.OrganizationID)
		if err != nil {
			return StartTalkSessionUseCaseOutput{}, messages.BadRequestError
		}
		organizationID = &orgID
	}
	// Creation permission does not authorize using another organization's identity.
	if organizationID != nil {
		membership, err := i.OrganizationUserRepository.FindByOrganizationIDAndUserID(ctx, *organizationID, input.OwnerID)
		if err != nil {
			return StartTalkSessionUseCaseOutput{}, err
		}
		if membership == nil {
			return StartTalkSessionUseCaseOutput{}, messages.ForbiddenError
		}
	}

	var output StartTalkSessionUseCaseOutput

	if err := input.Validate(); err != nil {
		return output, errtrace.Wrap(err)
	}

	if err := i.DBManager.ExecTx(ctx, func(ctx context.Context) error {
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

		// 合言葉が設定されている場合はハッシュ化
		var passphraseHash *string
		if input.Passphrase != nil && *input.Passphrase != "" {
			hash, err := i.PassphraseManager.Hash(*input.Passphrase)
			if err != nil {
				return errtrace.Wrap(err)
			}
			passphraseHash = &hash
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
			lo.FromPtrOr(input.HideTop, false),
			organizationID,
			organizationAliasID,
			passphraseHash,
		)

		if len(input.Restrictions) > 0 {
			if err := talkSession.UpdateRestrictions(ctx, input.Restrictions); err != nil {
				return errtrace.Wrap(err)
			}
		}

		if err := talkSession.StartSession(); err != nil {
			return errtrace.Wrap(err)
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
			HideTop:          talkSession.HideTop(),
			IsLocked:         talkSession.IsLocked(),
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
		// 組織エイリアス情報を取得
		if organizationAliasID != nil {
			alias, err := i.OrganizationAliasRepository.FindByID(ctx, *organizationAliasID)
			if err != nil {
				utils.HandleError(ctx, err, "OrganizationAliasRepository.FindByID")
				return messages.TalkSessionNotFound
			}
			if alias != nil {
				output.OrganizationAlias = &dto.OrganizationAlias{
					AliasID:   alias.AliasID().String(),
					AliasName: alias.AliasName(),
					CreatedAt: lo.ToPtr(alias.CreatedAt()),
				}
			} else {
				output.OrganizationAlias = nil
			}
		} else {
			output.OrganizationAlias = nil
		}

		return nil
	}); err != nil {
		return output, errtrace.Wrap(err)
	}

	return output, nil
}
