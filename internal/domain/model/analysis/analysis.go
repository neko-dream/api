package analysis

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
)

type (
	AnalysisService interface {
		StartAnalysis(context.Context, shared.UUID[talksession.TalkSession]) error
		GenerateReport(context.Context, shared.UUID[talksession.TalkSession]) error
	}
)
