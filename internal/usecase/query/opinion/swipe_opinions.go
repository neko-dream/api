package opinion_query

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/usecase/query/dto"
)

type (
	GetSwipeOpinionsQuery interface {
		Execute(context.Context, GetSwipeOpinionsQueryInput) (*GetSwipeOpinionsQueryOutput, error)
	}

	GetSwipeOpinionsQueryInput struct {
		UserID        shared.UUID[user.User]
		TalkSessionID shared.UUID[talksession.TalkSession]
		Limit         int
	}

	GetSwipeOpinionsQueryOutput struct {
		Opinions          []dto.SwipeOpinion
		RemainingOpinions int
	}
)
