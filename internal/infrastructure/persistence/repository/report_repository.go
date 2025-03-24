package repository

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/pkg/utils"
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
