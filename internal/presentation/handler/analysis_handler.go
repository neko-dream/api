package handler

import (
	"context"

	"github.com/neko-dream/server/internal/application/usecase/analysis_usecase"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/analysis"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/presentation/oas"
	"go.opentelemetry.io/otel"
)

type analysisHandler struct {
	applyFeedbackUseCase analysis_usecase.ApplyFeedbackUseCase
	authorizationService service.AuthorizationService
}

func NewAnalysisHandler(
	applyFeedbackUseCase analysis_usecase.ApplyFeedbackUseCase,
	authorizationService service.AuthorizationService,
) oas.AnalysisHandler {
	return &analysisHandler{
		applyFeedbackUseCase: applyFeedbackUseCase,
		authorizationService: authorizationService,
	}
}

// ApplyFeedbackToReport セッションのレポートにフィードバックを適用する.
func (a *analysisHandler) ApplyFeedbackToReport(ctx context.Context, req *oas.ApplyFeedbackToReportReq) (oas.ApplyFeedbackToReportRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "analysisHandler.ApplyFeedbackToReport")
	defer span.End()

	authCtx, err := a.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}
	userID := authCtx.UserID

	reportID, err := shared.ParseUUID[analysis.AnalysisReport](req.ReportID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	feedbackType := analysis.NewFeedbackTypeFromString(req.FeedbackType)
	if feedbackType == analysis.FeedbackTypeUnknown {
		return nil, messages.BadRequestError
	}

	if err := a.applyFeedbackUseCase.Execute(ctx, analysis_usecase.ApplyFeedbackInput{
		UserID:           userID,
		AnalysisReportID: reportID,
		FeedbackType:     feedbackType,
	}); err != nil {
		return nil, err
	}

	res := oas.ApplyFeedbackToReportOK{}
	return &res, nil
}
