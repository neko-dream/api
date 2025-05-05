package handler

import (
	"context"
	"time"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	timelineactions "github.com/neko-dream/server/internal/domain/model/timeline_actions"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/neko-dream/server/internal/usecase/command/timeline_command"
	"github.com/neko-dream/server/internal/usecase/query/timeline_query"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type timelineHandler struct {
	timeline_command.AddTimeLine
	et timeline_command.EditTimeLine
	gt timeline_query.GetTimeLine
}

func NewTimelineHandler(
	addTimeLine timeline_command.AddTimeLine,
	editTimeLine timeline_command.EditTimeLine,
	getTimeLine timeline_query.GetTimeLine,
) oas.TimelineHandler {
	return &timelineHandler{
		AddTimeLine: addTimeLine,
		et:          editTimeLine,
		gt:          getTimeLine,
	}
}

// GetTimeLine implements oas.TimelineHandler.
func (t *timelineHandler) GetTimeLine(ctx context.Context, params oas.GetTimeLineParams) (oas.GetTimeLineRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "timelineHandler.GetTimeLine")
	defer span.End()

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		utils.HandleError(ctx, err, "shared.ParseUUID")
		return nil, messages.InternalServerError
	}

	output, err := t.gt.Execute(ctx, timeline_query.GetTimeLineInput{
		TalkSessionID: talkSessionID,
	})
	if err != nil {
		utils.HandleError(ctx, err, "GetTimeLineUseCase.Execute")
		return nil, err
	}

	actionItems := make([]oas.ActionItem, len(output.ActionItems))
	for i, actionItem := range output.ActionItems {
		actionItems[i] = oas.ActionItem{
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
	ctx, span := otel.Tracer("handler").Start(ctx, "timelineHandler.PostTimeLineItem")
	defer span.End()

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

	output, err := t.AddTimeLine.Execute(ctx, timeline_command.AddTimeLineInput{
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

// EditTimeLine implements oas.TimelineHandler.
func (t *timelineHandler) EditTimeLine(ctx context.Context, req oas.OptEditTimeLineReq, params oas.EditTimeLineParams) (oas.EditTimeLineRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "timelineHandler.EditTimeLine")
	defer span.End()

	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}
	userID, err := claim.UserID()
	if err != nil {
		utils.HandleError(ctx, err, "claim.UserID")
		return nil, messages.InternalServerError
	}
	actionItemID, err := shared.ParseUUID[timelineactions.ActionItem](params.ActionItemID)
	if err != nil {
		utils.HandleError(ctx, err, "shared.ParseUUID")
		return nil, messages.InternalServerError
	}
	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		utils.HandleError(ctx, err, "shared.ParseUUID")
		return nil, messages.InternalServerError
	}
	var content, status *string
	if req.Value.Content.IsSet() {
		content = lo.ToPtr(req.Value.Content.Value)
	}
	if req.Value.Status.IsSet() {
		status = lo.ToPtr(req.Value.Status.Value)
	}

	output, err := t.et.Execute(ctx, timeline_command.EditTimeLineInput{
		OwnerID:       userID,
		TalkSessionID: talkSessionID,
		ActionItemID:  actionItemID,
		Content:       content,
		Status:        status,
	})
	if err != nil {
		utils.HandleError(ctx, err, "EditTimeLineUseCase.Execute")
		return nil, err
	}

	return &oas.ActionItem{
		ActionItemID: output.ActionItem.ActionItemID.String(),
		Content:      output.ActionItem.Content,
		Status:       output.ActionItem.Status.String(),
		CreatedAt:    output.ActionItem.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    output.ActionItem.UpdatedAt.Format(time.RFC3339),
	}, nil
}
