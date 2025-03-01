package service

import (
	"context"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	timelineactions "github.com/neko-dream/server/internal/domain/model/timeline_actions"
	"go.opentelemetry.io/otel"
)

type actionItemService struct {
	timelineactions.ActionItemRepository
	talksession.TalkSessionRepository
}

func NewActionItemService(
	actionItemRepository timelineactions.ActionItemRepository,
	talkSessionRepository talksession.TalkSessionRepository,
) timelineactions.ActionItemService {
	return &actionItemService{
		ActionItemRepository:  actionItemRepository,
		TalkSessionRepository: talkSessionRepository,
	}
}

// CanCreateActionItem
func (a *actionItemService) CanCreateActionItem(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) (*timelineactions.ActionItem, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "actionItemService.CanCreateActionItem")
	defer span.End()

	// TalkSessionが存在するか確認
	talkSession, err := a.TalkSessionRepository.FindByID(ctx, talkSessionID)
	if err != nil {
		return nil, err
	}
	if talkSession == nil {
		return nil, messages.TalkSessionNotFound
	}

	// セッションが終了しているか確認
	if !talkSession.IsFinished(ctx) {
		return nil, messages.TalkSessionNotFinished
	}

	// 最新のアクションアイテムを取得
	actionItems, err := a.ActionItemRepository.FindLatestActionItemByTalkSessionID(ctx, talkSessionID)
	if err != nil {
		return nil, err
	}

	if len(actionItems) == 0 {
		return nil, nil
	}

	// 最新のアクションアイテムを返す
	return &actionItems[0], nil
}

// InsertActionItem implements timelineactions.ActionItemService.
func (a *actionItemService) InsertActionItem(
	ctx context.Context,
	parentItemID *shared.UUID[timelineactions.ActionItem],
	actionItem timelineactions.ActionItem,
) error {
	ctx, span := otel.Tracer("service").Start(ctx, "actionItemService.InsertActionItem")
	defer span.End()

	// 親アクションアイテムが存在する場合、親アクションアイテムを取得
	var parentItem *timelineactions.ActionItem
	if parentItemID != nil {
		item, err := a.ActionItemRepository.FindByID(ctx, *parentItemID)
		if err != nil {
			return err
		}
		if item == nil {
			return messages.ActionItemNotFound
		}
		parentItem = item
	} else {
		// talkSessionIDより親アクションアイテムを取得
		parentItems, err := a.ActionItemRepository.FindLatestActionItemByTalkSessionID(ctx, actionItem.TalkSessionID)
		if err != nil {
			return err
		}
		if len(parentItems) == 0 {
			parentItem = nil
		} else {
			parentItem = &parentItems[0]
		}
	}

	var newSequence int
	if parentItem != nil {
		newSequence = parentItem.Sequence + 1
	} else {
		newSequence = 0
	}

	// アクションアイテムを作成
	newActionItem, err := timelineactions.NewActionItem(
		actionItem.ActionItemID,
		actionItem.TalkSessionID,
		newSequence,
		actionItem.Content,
		actionItem.Status,
		actionItem.CreatedAt,
		actionItem.UpdatedAt,
	)
	if err != nil {
		return err
	}
	actionItem = *newActionItem

	// アクションアイテムを登録
	if err := a.ActionItemRepository.CreateActionItem(ctx, actionItem); err != nil {
		return err
	}

	return nil
}
