package analysis_usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/neko-dream/server/internal/domain/model/analysis"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/infrastructure/db"
	model "github.com/neko-dream/server/internal/infrastructure/db/sqlc"
	"github.com/neko-dream/server/pkg/utils"
)

type (
	GetReportQuery interface {
		Execute(context.Context, GetReportInput) (*GetReportOutput, error)
	}

	GetReportInput struct {
		TalkSessionID shared.UUID[talksession.TalkSession]
	}

	GetReportOutput struct {
		Report string
	}

	GetReportQueryHandler struct {
		*db.DBManager
		analysis.AnalysisService
	}
)

func NewGetReportQueryHandler(
	tm *db.DBManager,
	as analysis.AnalysisService,
) GetReportQuery {
	return &GetReportQueryHandler{
		DBManager:       tm,
		AnalysisService: as,
	}
}

func (h *GetReportQueryHandler) Execute(ctx context.Context, input GetReportInput) (*GetReportOutput, error) {
	rows, err := h.GetQueries(ctx).GetOpinionsByTalkSessionID(ctx, model.GetOpinionsByTalkSessionIDParams{
		Limit:         10000,
		Offset:        0,
		TalkSessionID: input.TalkSessionID.UUID(),
		SortKey:       sql.NullString{String: "latest", Valid: true},
	})
	if err != nil {
		utils.HandleError(ctx, err, "トークセッションIDに紐づく意見の取得に失敗しました")
		return nil, err
	}
	if len(rows) == 0 {
		return nil, errors.New("opinions not found")
	}

	out, err := h.GetQueries(ctx).GetReportByTalkSessionId(ctx, input.TalkSessionID.UUID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if err := h.AnalysisService.GenerateReport(ctx, input.TalkSessionID); err != nil {
				return nil, err
			}
		}
		utils.HandleError(ctx, err, "レポートの生成に失敗しました")
		return nil, err
	}

	return &GetReportOutput{
		Report: out.Report,
	}, nil
}
