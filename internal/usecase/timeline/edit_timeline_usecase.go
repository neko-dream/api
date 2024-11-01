package timeline_usecase

import (
	"context"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	timelineactions "github.com/neko-dream/server/internal/domain/model/timeline_actions"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/db"
)

type (
	EditTimeLineUseCase interface {
		Execute(context.Context, EditTimeLineInput) (*EditTimeLineOutput, error)
	}

	EditTimeLineInput struct {
		OwnerID       shared.UUID[user.User]
		TalkSessionID shared.UUID[talksession.TalkSession]
		ActionItemID  shared.UUID[timelineactions.ActionItem]
		Content       *string
		Status        *string
	}

	EditTimeLineOutput struct {
		ActionItem *timelineactions.ActionItem
	}

	EditTimeLineInteractor struct {
		timelineactions.ActionItemRepository
		talksession.TalkSessionRepository
		timelineactions.ActionItemService
		*db.DBManager
	}
)

func NewEditTimeLineUseCase(
	actionItemRepository timelineactions.ActionItemRepository,
	talkSessionRepository talksession.TalkSessionRepository,
	actionItemService timelineactions.ActionItemService,
	dbManager *db.DBManager,
) EditTimeLineUseCase {
	return &EditTimeLineInteractor{
		ActionItemRepository:  actionItemRepository,
		TalkSessionRepository: talkSessionRepository,
		ActionItemService:     actionItemService,
		DBManager:             dbManager,
	}
}

func (i *EditTimeLineInteractor) Execute(ctx context.Context, input EditTimeLineInput) (*EditTimeLineOutput, error) {
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

	actionItem, err := i.ActionItemRepository.FindByID(ctx, input.ActionItemID)
	if err != nil {
		return nil, err
	}
	if actionItem == nil {
		return nil, messages.ActionItemNotFound
	}
	if input.Content != nil {
		actionItem.Content = *input.Content
	}
	if input.Status != nil {
		actionItem.UpdateStatus(timelineactions.ActionStatus(*input.Status))
	}

	if err := i.ActionItemRepository.UpdateActionItem(ctx, *actionItem); err != nil {
		return nil, err
	}

	actionItem, err = i.ActionItemRepository.FindByID(ctx, input.ActionItemID)
	if err != nil {
		return nil, err
	}

	return &EditTimeLineOutput{
		ActionItem: actionItem,
	}, nil
}
