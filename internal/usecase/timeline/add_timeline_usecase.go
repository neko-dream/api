package timeline_usecase

import (
	"context"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	timelineactions "github.com/neko-dream/server/internal/domain/model/timeline_actions"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type (
	AddTimeLineUseCase interface {
		Execute(context.Context, AddTimeLineInput) (*AddTimeLineOutput, error)
	}

	AddTimeLineInput struct {
		OwnerID        shared.UUID[user.User]
		TalkSessionID  shared.UUID[talksession.TalkSession]
		ParentActionID *shared.UUID[timelineactions.ActionItem]
		Content        string
		Status         string
	}

	AddTimeLineOutput struct {
		ActionItem *timelineactions.ActionItem
	}

	addTimeLineInteractor struct {
		timelineactions.ActionItemRepository
		talksession.TalkSessionRepository
		timelineactions.ActionItemService
		*db.DBManager
	}
)

func NewAddTimeLineUseCase(
	actionItemRepository timelineactions.ActionItemRepository,
	talkSessionRepository talksession.TalkSessionRepository,
	actionItemService timelineactions.ActionItemService,
	dbManager *db.DBManager,
) AddTimeLineUseCase {
	return &addTimeLineInteractor{
		ActionItemRepository:  actionItemRepository,
		TalkSessionRepository: talkSessionRepository,
		ActionItemService:     actionItemService,
		DBManager:             dbManager,
	}
}

func (i *addTimeLineInteractor) Execute(ctx context.Context, input AddTimeLineInput) (*AddTimeLineOutput, error) {
	ctx, span := otel.Tracer("timeline_usecase").Start(ctx, "addTimeLineInteractor.Execute")
	defer span.End()

	talkSession, err := i.TalkSessionRepository.FindByID(ctx, input.TalkSessionID)
	if err != nil {
		return nil, err
	}
	// セッションが存在しなければTimelineは作成できない
	if talkSession == nil {
		return nil, messages.TalkSessionNotFound
	}
	// セッションが終了しているか確認
	if !talkSession.IsFinished(ctx) {
		return nil, messages.TalkSessionNotFinished
	}
	// セッションのオーナーでなければTimelineは作成できない
	if talkSession.OwnerUserID() != input.OwnerID {
		return nil, messages.TalkSessionNotOwner
	}
	now := clock.Now(ctx)

	newActionItem, err := timelineactions.NewActionItem(
		shared.NewUUID[timelineactions.ActionItem](),
		input.TalkSessionID,
		0,
		input.Content,
		timelineactions.ActionStatus(input.Status),
		now,
		now,
	)
	if err != nil {
		utils.HandleError(ctx, err, "timelineactions.NewActionItem")
		return nil, err
	}

	if err := i.ActionItemService.InsertActionItem(ctx, input.ParentActionID, *newActionItem); err != nil {
		utils.HandleError(ctx, err, "ActionItemService.InsertActionItem")
		return nil, err
	}

	return &AddTimeLineOutput{
		ActionItem: newActionItem,
	}, nil
}
