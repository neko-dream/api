package talk_session_usecase

import (
	"context"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/conclusion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/utils"
)

type (
	CreateTalkSessionConclusionUseCase interface {
		Execute(context.Context, CreateTalkSessionConclusionInput) (*CreateTalkSessionConclusionOutput, error)
	}

	CreateTalkSessionConclusionInput struct {
		TalkSessionID shared.UUID[talksession.TalkSession]
		UserID        shared.UUID[user.User]
		Conclusion    string
	}

	CreateTalkSessionConclusionOutput struct {
		User    UserDTO
		Content string
	}
	createTalkSessionConclusionInteractor struct {
		*db.DBManager
		conclusion.ConclusionRepository
	}
)

func NewCreateTalkSessionConclusionUseCase(
	DBManager *db.DBManager,
	concRepo conclusion.ConclusionRepository,
) CreateTalkSessionConclusionUseCase {
	return &createTalkSessionConclusionInteractor{
		DBManager:            DBManager,
		ConclusionRepository: concRepo,
	}
}

func (i *createTalkSessionConclusionInteractor) Execute(ctx context.Context, input CreateTalkSessionConclusionInput) (*CreateTalkSessionConclusionOutput, error) {
	res, err := i.DBManager.GetQueries(ctx).GetTalkSessionByID(ctx, input.TalkSessionID.UUID())
	if err != nil {
		return nil, err
	}

	// オーナーでなければ結論を作成できない
	if res.UserID.UUID != input.UserID.UUID() {
		return nil, messages.TalkSessionNotOwner
	}

	// まだ終了していないトークセッションに対しては結論を作成できない
	if res.ScheduledEndTime.After(clock.Now(ctx)) {
		return nil, messages.TalkSessionNotFinished
	}

	conc, err := i.ConclusionRepository.FindByTalkSessionID(ctx, input.TalkSessionID)
	if err != nil {
		return nil, err
	}
	if conc != nil {
		return nil, messages.TalkSessionConclusionAlreadySet
	}

	conclusion := conclusion.NewConclusion(
		input.TalkSessionID,
		input.Conclusion,
		input.UserID,
	)
	if err := i.ConclusionRepository.Create(ctx, *conclusion); err != nil {
		return nil, err
	}

	return &CreateTalkSessionConclusionOutput{
		User: UserDTO{
			DisplayID:   res.DisplayID.String,
			DisplayName: res.DisplayName.String,
			IconURL: utils.ToPtrIfNotNullValue(
				!res.IconUrl.Valid,
				res.IconUrl.String,
			),
		},
		Content: input.Conclusion,
	}, nil
}
