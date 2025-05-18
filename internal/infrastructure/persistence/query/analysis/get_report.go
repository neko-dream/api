package analysis

import (
	"context"
	"database/sql"
	"errors"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/analysis"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/internal/application/query/analysis_query"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type GetReportQueryHandler struct {
	*db.DBManager
	analysis.AnalysisService
	talksession.TalkSessionRepository
}

func NewGetReportQueryHandler(
	tm *db.DBManager,
	as analysis.AnalysisService,
	talkSessionRep talksession.TalkSessionRepository,
) analysis_query.GetReportQuery {
	return &GetReportQueryHandler{
		DBManager:             tm,
		AnalysisService:       as,
		TalkSessionRepository: talkSessionRep,
	}
}

func (h *GetReportQueryHandler) Execute(ctx context.Context, input analysis_query.GetReportInput) (*analysis_query.GetReportOutput, error) {
	ctx, span := otel.Tracer("analysis_usecase").Start(ctx, "GetReportQueryHandler.Execute")
	defer span.End()

	// セッションのレポートを取得するかどうか
	talkSession, err := h.TalkSessionRepository.FindByID(ctx, input.TalkSessionID)
	if err != nil {
		utils.HandleError(ctx, err, "トークセッションの取得に失敗しました")
		return nil, messages.TalkSessionNotFound
	}
	if talkSession.HideReport() {
		return &analysis_query.GetReportOutput{
			Report: nil,
		}, nil
	}

	rows, err := h.GetQueries(ctx).GetOpinionsByTalkSessionID(ctx, model.GetOpinionsByTalkSessionIDParams{
		Limit:         10000,
		Offset:        0,
		TalkSessionID: input.TalkSessionID.UUID(),
		SortKey:       sql.NullString{String: "latest", Valid: true},
		IsSeed:        sql.NullBool{Bool: false, Valid: true},
	})
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			utils.HandleError(ctx, err, "トークセッションの取得に失敗しました")
			return nil, messages.TalkSessionNotFound
		}
	}
	if len(rows) == 0 {
		return &analysis_query.GetReportOutput{
			Report: nil,
		}, nil
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
		Report: &out.Report,
	}, nil
}
