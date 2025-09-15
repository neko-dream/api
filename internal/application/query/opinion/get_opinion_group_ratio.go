package opinion_query

import (
	"context"

	"github.com/neko-dream/api/internal/application/query/dto"
	"github.com/neko-dream/api/internal/domain/model/opinion"
	"github.com/neko-dream/api/internal/domain/model/shared"
)

type GetOpinionGroupRatioQuery interface {
	Execute(ctx context.Context, input GetOpinionGroupRatioInput) ([]dto.OpinionGroupRatio, error)
}

type GetOpinionGroupRatioInput struct {
	OpinionID shared.UUID[opinion.Opinion]
}
