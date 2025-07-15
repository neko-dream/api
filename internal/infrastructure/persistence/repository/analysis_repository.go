package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/neko-dream/server/internal/domain/model/analysis"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type analysisRepository struct {
	*db.DBManager
}

func NewAnalysisRepository(dbManager *db.DBManager) analysis.AnalysisRepository {
	return &analysisRepository{
		DBManager: dbManager,
	}
}

func (r *analysisRepository) FindByTalkSessionID(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) (*analysis.AnalysisReport, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "analysisRepository.FindByTalkSessionID")
	defer span.End()

	analysisReport, err := r.GetQueries(ctx).GetReportByTalkSessionId(ctx, talkSessionID.UUID())
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			utils.HandleError(ctx, err, "failed to retrieve analysis report")
		}
		return nil, err
	}

	// feedbackRows, err := r.DBManager.GetQueries(ctx).GetFeedbackByReportHistoryID(ctx, analysisReport.AnalysisReportHistoryID)
	// if err != nil {
	// 	if !errors.Is(err, sql.ErrNoRows) {
	// 		utils.HandleError(ctx, err, "failed to retrieve feedback by report history ID")
	// 	}
	// 	return nil, err
	// }

	// feedbacks := make([]analysis.Feedback, len(feedbackRows))
	// for i, row := range feedbackRows {
	// 	feedbacks[i] = analysis.Feedback{
	// 		FeedbackID: shared.UUID[analysis.Feedback](row.ReportFeedbackID),
	// 		UserID:     shared.UUID[user.User](row.UserID),
	// 		Type:       analysis.FeedbackType(row.FeedbackType),
	// 		CreatedAt:  row.CreatedAt,
	// 	}
	// }

	reportModel := &analysis.AnalysisReport{
		AnalysisReportID: shared.UUID[analysis.AnalysisReport](analysisReport.TalkSessionID),
		Report:           &analysisReport.Report,
		CreatedAt:        analysisReport.CreatedAt,
	}
	// if len(feedbacks) > 0 {
	// 	reportModel.Feedbacks = feedbacks
	// }

	return reportModel, nil
}

// FindByID implements analysis.AnalysisRepository.
func (r *analysisRepository) FindByID(ctx context.Context, analysisReportID shared.UUID[analysis.AnalysisReport]) (*analysis.AnalysisReport, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "analysisRepository.FindByID")
	defer span.End()

	analysisReport, err := r.GetQueries(ctx).FindReportByID(ctx, analysisReportID.UUID())
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			utils.HandleError(ctx, err, "failed to retrieve analysis report by ID")
		}
		return nil, err
	}

	// feedbackRows, err := r.DBManager.GetQueries(ctx).GetFeedbackByReportHistoryID(ctx, analysisReport.AnalysisReportHistoryID)
	// if err != nil {
	// 	if !errors.Is(err, sql.ErrNoRows) {
	// 		utils.HandleError(ctx, err, "failed to retrieve feedback by report history ID")
	// 	}
	// 	return nil, err
	// }

	// feedbacks := make([]analysis.Feedback, len(feedbackRows))
	// for i, row := range feedbackRows {
	// 	feedbacks[i] = analysis.Feedback{
	// 		FeedbackID: shared.UUID[analysis.Feedback](row.ReportFeedbackID),
	// 		UserID:     shared.UUID[user.User](row.UserID),
	// 		Type:       analysis.FeedbackType(row.FeedbackType),
	// 		CreatedAt:  row.CreatedAt,
	// 	}
	// }

	reportModel := &analysis.AnalysisReport{
		AnalysisReportID: shared.UUID[analysis.AnalysisReport](analysisReport.TalkSessionID),
		Report:           &analysisReport.Report,
		CreatedAt:        analysisReport.CreatedAt,
	}
	// if len(feedbacks) > 0 {
	// 	reportModel.Feedbacks = feedbacks
	// }

	return reportModel, nil
}

// SaveReport
func (r *analysisRepository) SaveReport(ctx context.Context, report *analysis.AnalysisReport) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "analysisRepository.SaveReport")
	defer span.End()

	if len(report.Feedbacks) > 0 {
		for _, feedback := range report.Feedbacks {
			if err := r.GetQueries(ctx).SaveReportFeedback(ctx, model.SaveReportFeedbackParams{
				ReportFeedbackID:           feedback.FeedbackID.UUID(),
				TalkSessionReportHistoryID: report.AnalysisReportID.UUID(),
				UserID:                     feedback.UserID.UUID(),
				FeedbackType:               int32(feedback.Type),
			}); err != nil {
				utils.HandleError(ctx, err, "failed to save report feedback")
				return err
			}
		}
	}

	return nil
}
