package handler

import (
	"context"
	"time"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	timelineactions "github.com/neko-dream/server/internal/domain/model/timeline_actions"
	"github.com/neko-dream/server/internal/presentation/oas"
	timeline_usecase "github.com/neko-dream/server/internal/usecase/timeline"
	"github.com/neko-dream/server/pkg/utils"
)

type timelineHandler struct {
	timeline_usecase.AddTimeLineUseCase
	timeline_usecase.GetTimeLineUseCase
}

func NewTimelineHandler(
	addTimeLineUseCase timeline_usecase.AddTimeLineUseCase,
	getTimeLineUseCase timeline_usecase.GetTimeLineUseCase,
) oas.TimelineHandler {
	return &timelineHandler{
		AddTimeLineUseCase: addTimeLineUseCase,
		GetTimeLineUseCase: getTimeLineUseCase,
	}
}

// GetTimeLine implements oas.TimelineHandler.
func (t *timelineHandler) GetTimeLine(ctx context.Context, params oas.GetTimeLineParams) (oas.GetTimeLineRes, error) {
	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		utils.HandleError(ctx, err, "shared.ParseUUID")
		return nil, messages.InternalServerError
	}

	output, err := t.GetTimeLineUseCase.Execute(ctx, timeline_usecase.GetTimeLineInput{
		TalkSessionID: talkSessionID,
	})
	if err != nil {
		utils.HandleError(ctx, err, "GetTimeLineUseCase.Execute")
		return nil, err
	}

	actionItems := make([]oas.GetTimeLineOKItemsItem, len(output.ActionItems))
	for i, actionItem := range output.ActionItems {
		actionItems[i] = oas.GetTimeLineOKItemsItem{
			ActionItemID: actionItem.ActionItemID.String(),
			Sequence:     actionItem.Sequence,
			Content:      actionItem.Content,
			Status:       actionItem.Status,
			CreatedAt:    actionItem.CreatedAt,
			UpdatedAt:    actionItem.UpdatedAt,
		}
	}

	return &oas.GetTimeLineOK{
		Items: actionItems,
	}, nil
}

// PostTimeLineItem implements oas.TimelineHandler.
func (t *timelineHandler) PostTimeLineItem(ctx context.Context, req oas.OptPostTimeLineItemReq, params oas.PostTimeLineItemParams) (oas.PostTimeLineItemRes, error) {
	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}
	userID, err := claim.UserID()
	if err != nil {
		utils.HandleError(ctx, err, "claim.UserID")
		return nil, messages.InternalServerError
	}
	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		utils.HandleError(ctx, err, "shared.ParseUUID")
		return nil, messages.InternalServerError
	}
	var parentActionID *shared.UUID[timelineactions.ActionItem]
	if req.Value.ParentActionItemID.IsSet() {
		parentActionIDIn, err := shared.ParseUUID[timelineactions.ActionItem](req.Value.ParentActionItemID.Value)
		if err != nil {
			utils.HandleError(ctx, err, "shared.ParseUUID")
			return nil, messages.InternalServerError
		}
		parentActionID = &parentActionIDIn
	}

	output, err := t.AddTimeLineUseCase.Execute(ctx, timeline_usecase.AddTimeLineInput{
		OwnerID:        userID,
		TalkSessionID:  talkSessionID,
		ParentActionID: parentActionID,
		Content:        req.Value.Content,
		Status:         req.Value.Status,
	})
	if err != nil {
		utils.HandleError(ctx, err, "AddTimeLineUseCase.Execute")
		return nil, err
	}

	return &oas.PostTimeLineItemOK{
		ActionItemID: output.ActionItem.ActionItemID.String(),
		Content:      output.ActionItem.Content,
		Status:       output.ActionItem.Status.String(),
		CreatedAt:    output.ActionItem.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    output.ActionItem.UpdatedAt.Format(time.RFC3339),
	}, nil
}
