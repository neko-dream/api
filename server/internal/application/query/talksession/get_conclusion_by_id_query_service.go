package talksession

import (
	"context"

	"github.com/neko-dream/api/internal/application/query/dto"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
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
