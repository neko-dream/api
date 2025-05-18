package opinion_query

import (
	"context"

	"github.com/neko-dream/server/internal/application/query/dto"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type (
	GetOpinionRepliesQuery interface {
		Execute(context.Context, GetOpinionRepliesQueryInput) (*GetOpinionRepliesQueryOutput, error)
	}

	GetOpinionRepliesQueryInput struct {
		OpinionID shared.UUID[opinion.Opinion]
		UserID    *shared.UUID[user.User]
	}

	GetOpinionRepliesQueryOutput struct {
		Replies []dto.SwipeOpinion
	}
)
