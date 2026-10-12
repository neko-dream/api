package analysis_query

import (
	"context"

	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
)

type (
	GetReportQuery interface {
		Execute(context.Context, GetReportInput) (*GetReportOutput, error)
	}

	GetReportInput struct {
		TalkSessionID shared.UUID[talksession.TalkSession]
	}

	GetReportOutput struct {
		Report *string
	}
)
