package timelineactions

import (
	"context"
	"time"
	"unicode/utf8"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
)

type (
	ActionItemRepository interface {
		CreateActionItem(context.Context, ActionItem) error
		UpdateActionItem(context.Context, ActionItem) error
		// FindLatestActionItemByTalkSessionID トークセッションIDに紐づく最新のアクションアイテムを取得
		FindLatestActionItemByTalkSessionID(context.Context, shared.UUID[talksession.TalkSession]) ([]ActionItem, error)
		// FindActionItemByActionItemID アクションアイテムIDに紐づくアクションアイテムを取得
		FindByID(context.Context, shared.UUID[ActionItem]) (*ActionItem, error)
	}

	ActionItemService interface {
		// ActionItemを作成できるか、できるなら親ActionItemが存在するか、存在すれば返す。存在しなければnilを返す。
		CanCreateActionItem(context.Context, shared.UUID[talksession.TalkSession]) (*ActionItem, error)
		InsertActionItem(context.Context, *shared.UUID[ActionItem], ActionItem) error
	}

	ActionItem struct {
		ActionItemID  shared.UUID[ActionItem]
		TalkSessionID shared.UUID[talksession.TalkSession]
		Sequence      int
		Content       string
		Status        ActionStatus
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}
)

func NewActionItem(
	actionItemID shared.UUID[ActionItem],
	talkSessionID shared.UUID[talksession.TalkSession],
	sequence int,
	content string,
	status ActionStatus,
	createdAt time.Time,
	updatedAt time.Time,
) (*ActionItem, error) {
	if sequence < 0 {
		return nil, messages.ActionItemInvalidSequence
	}

	if !status.Valid() {
		return nil, messages.ActionItemInvalidStatus
	}

	if utf8.RuneCountInString(content) < 1 || utf8.RuneCountInString(content) > 40 {
		return nil, messages.ActionItemInvalidContent
	}

	return &ActionItem{
		ActionItemID:  actionItemID,
		TalkSessionID: talkSessionID,
		Sequence:      sequence,
		Content:       content,
		Status:        status,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}, nil
}

func (a *ActionItem) UpdateContent(content string) error {
	if utf8.RuneCountInString(content) < 1 || utf8.RuneCountInString(content) > 40 {
		return messages.ActionItemInvalidContent
	}

	a.Content = content
	return nil
}

func (a *ActionItem) UpdateStatus(status ActionStatus) error {
	if !status.Valid() {
		return messages.ActionItemInvalidStatus
	}

	a.Status = status
	return nil
}
