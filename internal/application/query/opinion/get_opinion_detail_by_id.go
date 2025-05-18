package opinion_query

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/application/query/dto"
)

type (
	GetOpinionDetailByIDQuery interface {
		Execute(context.Context, GetOpinionDetailByIDInput) (*GetOpinionDetailByIDOutput, error)
	}

	GetOpinionDetailByIDInput struct {
		OpinionID shared.UUID[opinion.Opinion]
		UserID    *shared.UUID[user.User]
	}

	GetOpinionDetailByIDOutput struct {
		Opinion dto.SwipeOpinion
	}
)
