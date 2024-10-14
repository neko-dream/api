package analysis

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	user "github.com/neko-dream/server/internal/domain/model/user"
)

type (
	AnalysisService interface {
		StartAnalysis(context.Context, shared.UUID[talksession.TalkSession], shared.UUID[user.User]) error
	}
)
