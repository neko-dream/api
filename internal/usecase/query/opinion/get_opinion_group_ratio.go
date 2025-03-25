package opinion_query

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/usecase/query/dto"
)

type GetOpinionGroupRatioQuery interface {
	Execute(ctx context.Context, input GetOpinionGroupRatioInput) ([]dto.OpinionGroupRatio, error)
}

type GetOpinionGroupRatioInput struct {
	OpinionID shared.UUID[opinion.Opinion]
}
