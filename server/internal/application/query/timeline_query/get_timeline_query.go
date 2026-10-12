package timeline_query

import (
	"context"
	"time"

	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	timelineactions "github.com/neko-dream/api/internal/domain/model/timeline_actions"
	"go.opentelemetry.io/otel"
)

type (
	GetTimeLine interface {
		Execute(context.Context, GetTimeLineInput) (*GetTimeLineOutput, error)
	}

	GetTimeLineInput struct {
		TalkSessionID shared.UUID[talksession.TalkSession]
	}

	GetTimeLineOutput struct {
		ActionItems []ActionItemDTO
	}

	ActionItemDTO struct {
		ActionItemID shared.UUID[timelineactions.ActionItem]
		Sequence     int
		Content      string
		Status       string
		CreatedAt    string
		UpdatedAt    string
	}

	getTimeLineInteractor struct {
		timelineactions.ActionItemRepository
		talksession.TalkSessionRepository
	}
)

func NewGetTimeLine(
	actionItemRepository timelineactions.ActionItemRepository,
	talkSessionRepository talksession.TalkSessionRepository,
) GetTimeLine {
	return &getTimeLineInteractor{
		ActionItemRepository:  actionItemRepository,
		TalkSessionRepository: talkSessionRepository,
	}
}

// Execute implements GetTimeLine.
func (g *getTimeLineInteractor) Execute(ctx context.Context, input GetTimeLineInput) (*GetTimeLineOutput, error) {
	ctx, span := otel.Tracer("timeline_").Start(ctx, "getTimeLineInteractor.Execute")
	defer span.End()

	talkSession, err := g.TalkSessionRepository.FindByID(ctx, input.TalkSessionID)
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

	actionItems, err := g.ActionItemRepository.FindLatestActionItemByTalkSessionID(ctx, input.TalkSessionID)
	if err != nil {
		return nil, err
	}

	actionItemDTOList := make([]ActionItemDTO, 0, len(actionItems))
	for _, actionItem := range actionItems {
		actionItemDTOList = append(actionItemDTOList, ActionItemDTO{
			ActionItemID: actionItem.ActionItemID,
			Sequence:     actionItem.Sequence,
			Content:      actionItem.Content,
			Status:       string(actionItem.Status),
			CreatedAt:    actionItem.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    actionItem.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &GetTimeLineOutput{
		ActionItems: actionItemDTOList,
	}, nil
}
