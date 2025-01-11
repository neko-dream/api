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

type GetOpinionRepliesQueryHandler struct {
	*db.DBManager
}

func NewGetOpinionRepliesQueryHandler(
	dbManager *db.DBManager,
) opinion_query.GetOpinionRepliesQuery {
	return &GetOpinionRepliesQueryHandler{
		DBManager: dbManager,
	}
}

func (g *GetOpinionRepliesQueryHandler) Execute(ctx context.Context, in opinion_query.GetOpinionRepliesQueryInput) (*opinion_query.GetOpinionRepliesQueryOutput, error) {

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

	replyRows, err := g.GetQueries(ctx).GetOpinionReplies(ctx, model.GetOpinionRepliesParams{
		OpinionID: in.OpinionID.UUID(),
		UserID:    userID,
	})
	if err != nil {
		utils.HandleError(ctx, err, "リプライの取得に失敗")
		return nil, err
	}

	var replies []dto.SwipeOpinion
	if err := copier.CopyWithOption(&replies, &replyRows, copier.Option{
		DeepCopy:    true,
		IgnoreEmpty: true,
	}); err != nil {
		utils.HandleError(ctx, err, "マッピングに失敗")
		return nil, err
	}

	return &opinion_query.GetOpinionRepliesQueryOutput{
		RootOpinion: opinion,
		Replies:     replies,
	}, nil
}
