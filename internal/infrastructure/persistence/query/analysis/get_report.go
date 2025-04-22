package analysis

import (
	"context"
	"database/sql"
	"errors"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/analysis"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/internal/usecase/query/analysis_query"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type GetReportQueryHandler struct {
	*db.DBManager
	analysis.AnalysisService
}

func NewGetReportQueryHandler(
	tm *db.DBManager,
	as analysis.AnalysisService,
) analysis_query.GetReportQuery {
	return &GetReportQueryHandler{
		DBManager:       tm,
		AnalysisService: as,
	}
}

func (h *GetReportQueryHandler) Execute(ctx context.Context, input analysis_query.GetReportInput) (*analysis_query.GetReportOutput, error) {
	ctx, span := otel.Tracer("analysis_usecase").Start(ctx, "GetReportQueryHandler.Execute")
	defer span.End()

	rows, err := h.GetQueries(ctx).GetOpinionsByTalkSessionID(ctx, model.GetOpinionsByTalkSessionIDParams{
		Limit:         10000,
		Offset:        0,
		TalkSessionID: input.TalkSessionID.UUID(),
		SortKey:       sql.NullString{String: "latest", Valid: true},
		IsSeed:        sql.NullBool{Bool: false, Valid: true},
	})
	if err != nil || len(rows) == 0 {
		utils.HandleError(ctx, err, "トークセッションIDに紐づく意見の取得に失敗しました")
		return nil, messages.AnalysisReportOpinionNotFound
	}

	out, err := h.GetQueries(ctx).GetReportByTalkSessionId(ctx, input.TalkSessionID.UUID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if err := h.AnalysisService.GenerateReport(ctx, input.TalkSessionID); err != nil {
				utils.HandleError(ctx, err, "レポートの生成に失敗しました")
				return nil, messages.AnalysisReportNotFound
			}
		}

		utils.HandleError(ctx, err, "レポートの生成に失敗しました")
		return nil, messages.AnalysisReportNotFound
	}

	return &analysis_query.GetReportOutput{
		Report: out.Report,
	}, nil
}
