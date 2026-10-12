package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	timelineactions "github.com/neko-dream/api/internal/domain/model/timeline_actions"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

type actionItemRepository struct {
	*db.DBManager
}

func NewActionItemRepository(
	dbManager *db.DBManager,
) timelineactions.ActionItemRepository {
	return &actionItemRepository{
		DBManager: dbManager,
	}
}

// CreateActionItem implements timelineactions.ActionItemRepository.
func (a *actionItemRepository) CreateActionItem(ctx context.Context, actionItem timelineactions.ActionItem) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "actionItemRepository.CreateActionItem")
	defer span.End()

	err := a.GetQueries(ctx).CreateActionItem(ctx, model.CreateActionItemParams{
		ActionItemID:  actionItem.ActionItemID.UUID(),
		TalkSessionID: actionItem.TalkSessionID.UUID(),
		Sequence:      int32(actionItem.Sequence),
		Content:       actionItem.Content,
		Status:        string(actionItem.Status),
		CreatedAt:     actionItem.CreatedAt,
		UpdatedAt:     actionItem.UpdatedAt,
	})
	if err != nil {
		utils.HandleError(ctx, err, "ActionItemの作成に失敗しました")
		return err
	}

	return nil
}

// UpdateActionItem implements timelineactions.ActionItemRepository.
func (a *actionItemRepository) UpdateActionItem(ctx context.Context, actionItem timelineactions.ActionItem) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "actionItemRepository.UpdateActionItem")
	defer span.End()

	_ = ctx

	panic("unimplemented")
}

// FindActionItemByActionItemID implements timelineactions.ActionItemRepository.
func (a *actionItemRepository) FindByID(ctx context.Context, actionItemID shared.UUID[timelineactions.ActionItem]) (*timelineactions.ActionItem, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "actionItemRepository.FindByID")
	defer span.End()

	row, err := a.GetQueries(ctx).GetActionItemByID(ctx, actionItemID.UUID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		utils.HandleError(ctx, err, "アクションアイテムIDに紐づくアクションアイテムの取得に失敗しました")
		return nil, err
	}

	return &timelineactions.ActionItem{
		ActionItemID:  shared.UUID[timelineactions.ActionItem](row.ActionItemID),
		TalkSessionID: shared.UUID[talksession.TalkSession](row.TalkSessionID),
		Sequence:      int(row.Sequence),
		Content:       row.Content,
		Status:        timelineactions.ActionStatus(row.Status),
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}, nil
}

// FindLatestActionItemByTalkSessionID implements timelineactions.ActionItemRepository.
func (a *actionItemRepository) FindLatestActionItemByTalkSessionID(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) ([]timelineactions.ActionItem, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "actionItemRepository.FindLatestActionItemByTalkSessionID")
	defer span.End()

	row, err := a.GetQueries(ctx).GetActionItemsByTalkSessionID(ctx, talkSessionID.UUID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.HandleError(ctx, err, "トークセッションIDに紐づくアクションアイテムの取得に失敗しました")
			return nil, nil
		}
		return nil, err
	}

	if len(row) == 0 {
		return nil, nil
	}

	actionItems := make([]timelineactions.ActionItem, 0, len(row))
	for _, r := range row {
		actionItems = append(actionItems, timelineactions.ActionItem{
			ActionItemID:  shared.UUID[timelineactions.ActionItem](r.ActionItemID),
			TalkSessionID: shared.UUID[talksession.TalkSession](r.TalkSessionID),
			Sequence:      int(r.Sequence),
			Content:       r.Content,
			Status:        timelineactions.ActionStatus(r.Status),
			CreatedAt:     r.CreatedAt,
			UpdatedAt:     r.UpdatedAt,
		})
	}

	return actionItems, nil
}
