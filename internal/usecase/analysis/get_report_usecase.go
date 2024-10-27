package analysis_usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/neko-dream/server/internal/domain/model/analysis"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/infrastructure/db"
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

	out, err := h.GetQueries(ctx).GetReportByTalkSessionId(ctx, input.TalkSessionID.UUID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if err := h.AnalysisService.GenerateReport(ctx, input.TalkSessionID); err != nil {
				return nil, err
			}
		}
		return nil, err
	}

	return &GetReportOutput{
		Report: out.Report,
	}, nil
}
