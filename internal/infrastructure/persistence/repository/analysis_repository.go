package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/neko-dream/server/internal/domain/model/analysis"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
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

	return &analysis.AnalysisReport{
		Report:    &analysisReport.Report,
		UpdatedAt: analysisReport.UpdatedAt,
		CreatedAt: analysisReport.CreatedAt,
	}, nil
}
