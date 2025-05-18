package opinion_query

import (
	"context"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	dto_mapper "github.com/neko-dream/server/internal/infrastructure/persistence/utils"
	"github.com/neko-dream/server/internal/application/query/dto"
	opinion_query "github.com/neko-dream/server/internal/application/query/opinion"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
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
	ctx, span := otel.Tracer("opinion_query").Start(ctx, "GetOpinionRepliesQueryHandler.Execute")
	defer span.End()

	var userID uuid.NullUUID
	if in.UserID != nil {
		userID = uuid.NullUUID{UUID: in.UserID.UUID(), Valid: true}
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

	// 通報された意見を処理
	if len(replies) > 0 {
		opinionIDs := dto_mapper.ExtractOpinionIDs(replies)
		reports, err := g.GetQueries(ctx).FindReportByOpinionIDs(ctx, model.FindReportByOpinionIDsParams{
			OpinionIds: opinionIDs,
			Status:     "deleted",
		})
		if err != nil {
			utils.HandleError(ctx, err, "通報情報の取得に失敗")
			return nil, err
		}

		replies = dto_mapper.ProcessReportedOpinions(replies, reports)
	}

	return &opinion_query.GetOpinionRepliesQueryOutput{
		Replies: replies,
	}, nil
}
