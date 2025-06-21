package handler

import (
	"context"
	"time"

	"github.com/neko-dream/server/internal/application/query/timeline_query"
	"github.com/neko-dream/server/internal/application/usecase/timeline_usecase"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	timelineactions "github.com/neko-dream/server/internal/domain/model/timeline_actions"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type timelineHandler struct {
	timeline_usecase.AddTimeLine
	et          timeline_usecase.EditTimeLine
	gt          timeline_query.GetTimeLine
	authService service.AuthenticationService
}

func NewTimelineHandler(
	addTimeLine timeline_usecase.AddTimeLine,
	editTimeLine timeline_usecase.EditTimeLine,
	getTimeLine timeline_query.GetTimeLine,
	authService service.AuthenticationService,
) oas.TimelineHandler {
	return &timelineHandler{
		AddTimeLine: addTimeLine,
		et:          editTimeLine,
		gt:          getTimeLine,
		authService: authService,
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
func (t *timelineHandler) PostTimeLineItem(ctx context.Context, req *oas.PostTimeLineItemReq, params oas.PostTimeLineItemParams) (oas.PostTimeLineItemRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "timelineHandler.PostTimeLineItem")
	defer span.End()

	authCtx, err := requireAuthentication(t.authService, ctx)
	if err != nil {
		return nil, err
	}
	if req == nil {
		utils.HandleError(ctx, nil, "req is nil")
		return nil, messages.RequiredParameterError
	}

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		utils.HandleError(ctx, err, "shared.ParseUUID")
		return nil, messages.InternalServerError
	}
	var parentActionID *shared.UUID[timelineactions.ActionItem]
	if req.ParentActionItemID.IsSet() {
		parentActionIDIn, err := shared.ParseUUID[timelineactions.ActionItem](req.ParentActionItemID.Value)
		if err != nil {
			utils.HandleError(ctx, err, "shared.ParseUUID")
			return nil, messages.InternalServerError
		}
		parentActionID = &parentActionIDIn
	}

	output, err := t.AddTimeLine.Execute(ctx, timeline_usecase.AddTimeLineInput{
		OwnerID:        authCtx.UserID,
		TalkSessionID:  talkSessionID,
		ParentActionID: parentActionID,
		Content:        req.Content,
		Status:         req.Status,
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
func (t *timelineHandler) EditTimeLine(ctx context.Context, req *oas.EditTimeLineReq, params oas.EditTimeLineParams) (oas.EditTimeLineRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "timelineHandler.EditTimeLine")
	defer span.End()

	authCtx, err := requireAuthentication(t.authService, ctx)
	if err != nil {
		return nil, err
	}
	userID := authCtx.UserID
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
	if req.Content.IsSet() {
		content = lo.ToPtr(req.Content.Value)
	}
	if req.Status.IsSet() {
		status = lo.ToPtr(req.Status.Value)
	}

	output, err := t.et.Execute(ctx, timeline_usecase.EditTimeLineInput{
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
