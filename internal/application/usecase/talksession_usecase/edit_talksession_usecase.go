package talksession_usecase

import (
	"context"
	"time"
	"unicode/utf8"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/application/query/dto"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type (
	EditTalkSessionUseCase interface {
		Execute(context.Context, EditTalkSessionInput) (*EditTalkSessionOutput, error)
	}

	// EditTalkSessionInput セッション編集に必要な入力データ
	EditTalkSessionInput struct {
		TalkSessionID    shared.UUID[talksession.TalkSession] // セッションのID
		UserID           shared.UUID[user.User]               // 編集するユーザーのID
		Theme            string                               // セッションのテーマ
		Description      *string                              // セッションの説明文
		ThumbnailURL     *string                              // サムネイル画像のURL
		ScheduledEndTime time.Time                            // 予定終了時刻
		Latitude         *float64                             // 緯度
		Longitude        *float64                             // 経度
		City             *string                              // 市区町村
		Prefecture       *string                              // 都道府県
	}

	EditTalkSessionOutput struct {
		dto.TalkSessionWithDetail
	}

	editTalkSessionHandler struct {
		talksession.TalkSessionRepository
		user.UserRepository
		*db.DBManager
	}
)

// Validate 入力データのバリデーションを行う
// - テーマは20文字以内
// - 説明文は400文字以内
func (in *EditTalkSessionInput) Validate() error {
	// テーマは20文字まで
	if utf8.RuneCountInString(in.Theme) > 100 {
		return messages.TalkSessionThemeTooLong
	}
	// 説明文は400文字まで
	if in.Description != nil && utf8.RuneCountInString(*in.Description) > 40000 {
		return messages.TalkSessionDescriptionTooLong
	}

	return nil
}

func NewEditTalkSessionUseCase(
	talkSessionRepository talksession.TalkSessionRepository,
	userRepository user.UserRepository,
	DBManager *db.DBManager,
) EditTalkSessionUseCase {
	return &editTalkSessionHandler{
		TalkSessionRepository: talkSessionRepository,
		UserRepository:        userRepository,
		DBManager:             DBManager,
	}
}

func (i *editTalkSessionHandler) Execute(ctx context.Context, input EditTalkSessionInput) (*EditTalkSessionOutput, error) {
	ctx, span := otel.Tracer("talksession_command").Start(ctx, "startTalkSessionCommandHandler.Execute")
	defer span.End()

	// talkSessionを探してくる
	talkSession, err := i.TalkSessionRepository.FindByID(ctx, input.TalkSessionID)
	if err != nil {
		utils.HandleError(ctx, err, "TalkSessionRepository.FindByID")
		return nil, messages.TalkSessionNotFound
	}

	// 編集するユーザーがオーナーかどうかを確認
	if talkSession.OwnerUserID() != input.UserID {
		return nil, messages.ForbiddenError
	}

	var output EditTalkSessionOutput
	if err := input.Validate(); err != nil {
		return nil, errtrace.Wrap(err)
	}
	if input.ScheduledEndTime.Before(clock.Now(ctx)) {
		return nil, messages.InvalidScheduledEndTime
	}

	if err := i.ExecTx(ctx, func(ctx context.Context) error {

		talkSession.ChangeTheme(input.Theme)
		talkSession.ChangeDescription(input.Description)
		talkSession.ChangeThumbnailURL(input.ThumbnailURL)
		talkSession.ChangeScheduledEndTime(input.ScheduledEndTime)
		talkSession.ChangeCity(input.City)
		talkSession.ChangePrefecture(input.Prefecture)
		if input.Latitude != nil && input.Longitude != nil {
			talkSession.ChangeLocation(talksession.NewLocation(
				input.TalkSessionID,
				*input.Latitude,
				*input.Longitude,
			))
		}

		if err := i.TalkSessionRepository.Update(ctx, talkSession); err != nil {
			utils.HandleError(ctx, err, "TalkSessionRepository.Update")
			return messages.TalkSessionUpdateFailed
		}

		output.TalkSession = dto.TalkSession{
			TalkSessionID:    input.TalkSessionID,
			Theme:            talkSession.Theme(),
			ThumbnailURL:     talkSession.ThumbnailURL(),
			ScheduledEndTime: talkSession.ScheduledEndTime(),
			OwnerID:          talkSession.OwnerUserID(),
			CreatedAt:        talkSession.CreatedAt(),
			Description:      talkSession.Description(),
			City:             talkSession.City(),
			Prefecture:       talkSession.Prefecture(),
		}
		output.Latitude = input.Latitude
		output.Longitude = input.Longitude
		output.Restrictions = talkSession.RestrictionList()

		// オーナーのユーザー情報を取得
		ownerUser, err := i.UserRepository.FindByID(ctx, input.UserID)
		if err != nil {
			utils.HandleError(ctx, err, "UserRepository.FindByID")
			return messages.ForbiddenError
		}

		// ユーザー情報をDTOに設定
		if ownerUser.DisplayID() != nil && ownerUser.DisplayName() != nil {
			output.User = dto.User{
				DisplayID:   *ownerUser.DisplayID(),
				DisplayName: *ownerUser.DisplayName(),
				IconURL:     ownerUser.IconURL(),
			}
		} else {
			// DisplayIDまたはDisplayNameがnilの場合はデフォルト値を使用
			output.User = dto.User{
				DisplayID:   "",
				DisplayName: "",
				IconURL:     ownerUser.IconURL(),
			}
		}

		return nil
	}); err != nil {
		return nil, errtrace.Wrap(err)
	}

	return &output, nil
}
