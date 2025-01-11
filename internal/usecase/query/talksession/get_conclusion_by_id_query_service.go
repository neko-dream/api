package talksession

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/usecase/query/dto"
)

type (
	GetConclusionByIDQueryRequest struct {
		TalkSessionID shared.UUID[talksession.TalkSession]
	}

	GetConclusionByIDQueryResponse struct {
		dto.TalkSessionConclusion
		dto.User
	}

	GetConclusionByIDQuery interface {
		Execute(context.Context, GetConclusionByIDQueryRequest) (*GetConclusionByIDQueryResponse, error)
	}
)
