package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/neko-dream/api/internal/domain/model/opinion"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

type reportRepository struct {
	*db.DBManager
}

func NewReportRepository(dbm *db.DBManager) opinion.ReportRepository {
	return &reportRepository{dbm}
}

// Create 意見に対する通報を作成する
func (r *reportRepository) Create(ctx context.Context, rep opinion.Report) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "reportRepository.Create")
	defer span.End()

	if err := r.DBManager.GetQueries(ctx).CreateReport(
		ctx,
		model.CreateReportParams{
			OpinionReportID: rep.OpinionReportID.UUID(),
			OpinionID:       rep.OpinionID.UUID(),
			TalkSessionID:   rep.TalkSessionID.UUID(),
			ReporterID:      rep.ReporterID.UUID(),
			Reason:          int32(rep.Reason),
			Status:          string(rep.Status),
		},
	); err != nil {
		utils.HandleError(ctx, err, "CreateReport")
		return err
	}

	return nil
}

// UpdateStatus 通報のステータスを更新する
func (r *reportRepository) UpdateStatus(ctx context.Context, reportID shared.UUID[opinion.Report], status opinion.Status) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "reportRepository.UpdateStatus")
	defer span.End()

	if err := r.DBManager.GetQueries(ctx).UpdateReportStatus(
		ctx,
		model.UpdateReportStatusParams{
			Status:          string(status),
			OpinionReportID: reportID.UUID(),
		},
	); err != nil {
		utils.HandleError(ctx, err, "UpdateReportStatus")
		return err
	}

	return nil
}

// FindByOpinionID 意見IDから通報を取得する
func (r *reportRepository) FindByOpinionID(ctx context.Context, opinionID shared.UUID[opinion.Opinion]) ([]opinion.Report, error) {

	ctx, span := otel.Tracer("repository").Start(ctx, "reportRepository.FindByOpinionID")
	defer span.End()

	reports, err := r.DBManager.GetQueries(ctx).FindReportByOpinionID(ctx, uuid.NullUUID{UUID: opinionID.UUID(), Valid: true})
	if err != nil {
		utils.HandleError(ctx, err, "FindReportByOpinionID")
		return nil, err
	}

	result := make([]opinion.Report, 0, len(reports))
	for _, report := range reports {
		var reasonText *string
		if report.OpinionReport.ReasonText.Valid {
			reasonText = &report.OpinionReport.ReasonText.String
		}
		rep := opinion.Report{
			OpinionReportID: shared.UUID[opinion.Report](report.OpinionReport.OpinionReportID),
			OpinionID:       shared.UUID[opinion.Opinion](report.OpinionReport.OpinionID),
			TalkSessionID:   shared.UUID[talksession.TalkSession](report.OpinionReport.TalkSessionID),
			ReporterID:      shared.UUID[user.User](report.OpinionReport.ReporterID),
			Reason:          opinion.Reason(report.OpinionReport.Reason),
			Status:          opinion.Status(report.OpinionReport.Status),
			ReasonText:      reasonText,
			CreatedAt:       report.OpinionReport.CreatedAt,
		}
		result = append(result, rep)
	}

	return result, nil
}

// CountByTalkSessionIDAndStatus
func (r *reportRepository) CountByTalkSessionIDAndStatus(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession], status opinion.Status) (int, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "reportRepository.CountByTalkSessionIDAndStatus")
	defer span.End()
	count, err := r.DBManager.GetQueries(ctx).CountReportsByTalkSession(ctx, model.CountReportsByTalkSessionParams{
		TalkSessionID: uuid.NullUUID{
			UUID:  talkSessionID.UUID(),
			Valid: true,
		},
		Status: sql.NullString{
			String: string(status),
			Valid:  true,
		},
	})
	if err != nil {
		utils.HandleError(ctx, err, "CountReportByTalkSessionIDAndStatus")
		return 0, err
	}

	return int(count), nil
}
