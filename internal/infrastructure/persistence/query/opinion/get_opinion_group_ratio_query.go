package opinion_query

import (
	"context"

	"github.com/google/uuid"
	"github.com/neko-dream/server/internal/application/query/dto"
	opinion_query "github.com/neko-dream/server/internal/application/query/opinion"
	"github.com/neko-dream/server/internal/domain/model/analysis"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"go.opentelemetry.io/otel"
)

type getOpinionGroupRatioInteractor struct {
	*db.DBManager
}

func NewGetOpinionGroupRatioInteractor(dbm *db.DBManager) opinion_query.GetOpinionGroupRatioQuery {
	return &getOpinionGroupRatioInteractor{dbm}
}

// Execute
func (g *getOpinionGroupRatioInteractor) Execute(ctx context.Context, input opinion_query.GetOpinionGroupRatioInput) ([]dto.OpinionGroupRatio, error) {
	ctx, span := otel.Tracer("opinion_query").Start(ctx, "getOpinionGroupRatioInteractor.Execute")
	defer span.End()

	opinionID := uuid.NullUUID{
		UUID:  input.OpinionID.UUID(),
		Valid: true,
	}

	res, err := g.DBManager.GetQueries(ctx).GetGroupRatioByOpinionID(ctx, opinionID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.OpinionGroupRatio, len(res))
	for i, r := range res {
		result[i] = dto.OpinionGroupRatio{
			GroupName:     analysis.NewGroupIDFromInt(int(r.RepresentativeOpinion.GroupID)).String(),
			GroupID:       int(r.RepresentativeOpinion.GroupID),
			AgreeCount:    int(r.RepresentativeOpinion.AgreeCount),
			DisagreeCount: int(r.RepresentativeOpinion.DisagreeCount),
			PassCount:     int(r.RepresentativeOpinion.PassCount),
		}
	}

	return result, nil
}
