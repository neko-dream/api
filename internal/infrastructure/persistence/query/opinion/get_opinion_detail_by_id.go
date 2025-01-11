package opinion_query

import (
	"context"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/internal/usecase/query/dto"
	opinion_query "github.com/neko-dream/server/internal/usecase/query/opinion"
	"github.com/neko-dream/server/pkg/utils"
)

type GetOpinionDetailByIDQueryHandler struct {
	*db.DBManager
}

func NewGetOpinionDetailByIDQueryHandler(
	dbManager *db.DBManager,
) opinion_query.GetOpinionDetailByIDQuery {
	return &GetOpinionDetailByIDQueryHandler{
		DBManager: dbManager,
	}
}

func (g *GetOpinionDetailByIDQueryHandler) Execute(ctx context.Context, in opinion_query.GetOpinionDetailByIDInput) (*opinion_query.GetOpinionDetailByIDOutput, error) {
	var userID uuid.NullUUID
	if in.UserID != nil {
		userID = uuid.NullUUID{UUID: in.UserID.UUID(), Valid: true}
	}

	opinionRow, err := g.GetQueries(ctx).GetOpinionByID(ctx, model.GetOpinionByIDParams{
		OpinionID: in.OpinionID.UUID(),
		UserID:    userID,
	})
	if err != nil {
		utils.HandleError(ctx, err, "意見の取得に失敗")
		return nil, err
	}

	var opinion dto.SwipeOpinion
	if err := copier.CopyWithOption(&opinion, &opinionRow, copier.Option{
		DeepCopy:    true,
		IgnoreEmpty: true,
	}); err != nil {
		utils.HandleError(ctx, err, "マッピングに失敗")
		return nil, err
	}

	return &opinion_query.GetOpinionDetailByIDOutput{
		Opinion: opinion,
	}, nil
}
