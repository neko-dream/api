package analysis_query

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/analysis"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
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
