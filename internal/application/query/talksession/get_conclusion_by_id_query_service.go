package talksession

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/application/query/dto"
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
