package opinion_query

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/internal/usecase/query/dto"
	opinion_query "github.com/neko-dream/server/internal/usecase/query/opinion"
)

type GetMyOpinionsQueryHandler struct {
	*db.DBManager
}

func NewGetMyOpinionsQueryHandler(
	dbManager *db.DBManager,
) opinion_query.GetMyOpinionsQuery {
	return &GetMyOpinionsQueryHandler{
		DBManager: dbManager,
	}
}

func (g *GetMyOpinionsQueryHandler) Execute(ctx context.Context, in opinion_query.GetMyOpinionsQueryInput) (*opinion_query.GetMyOpinionsQueryOutput, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	opinionRows, err := g.GetQueries(ctx).GetOpinionsByUserID(ctx, model.GetOpinionsByUserIDParams{
		UserID:  in.UserID.UUID(),
		SortKey: sql.NullString{String: in.SortKey.String(), Valid: true},
		Limit:   int32(*in.Limit),
		Offset:  int32(*in.Offset),
	})
	if err != nil {
		return nil, err
	}

	var opinions []dto.SwipeOpinion
	if err := copier.CopyWithOption(&opinions, &opinionRows, copier.Option{
		DeepCopy:    true,
		IgnoreEmpty: true,
	}); err != nil {
		return nil, err
	}

	count, err := g.GetQueries(ctx).CountOpinions(ctx, model.CountOpinionsParams{
		UserID: uuid.NullUUID{UUID: in.UserID.UUID(), Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &opinion_query.GetMyOpinionsQueryOutput{
		Opinions:   opinions,
		TotalCount: int(count),
	}, nil
}
