package opinion_query

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/internal/usecase/query/dto"
	opinion_query "github.com/neko-dream/server/internal/usecase/query/opinion"
	"github.com/neko-dream/server/pkg/utils"
)

type GetOpinionsByTalkSessionIDQueryHandler struct {
	*db.DBManager
}

func NewGetOpinionsByTalkSessionIDQueryHandler(
	dbManager *db.DBManager,
) opinion_query.GetOpinionsByTalkSessionQuery {
	return &GetOpinionsByTalkSessionIDQueryHandler{
		DBManager: dbManager,
	}
}

func (g *GetOpinionsByTalkSessionIDQueryHandler) Execute(ctx context.Context, in opinion_query.GetOpinionsByTalkSessionInput) (*opinion_query.GetOpinionsByTalkSessionOutput, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	var userID uuid.NullUUID
	if in.UserID != nil {
		userID = uuid.NullUUID{UUID: in.UserID.UUID(), Valid: true}
	}

	opinionRows, err := g.GetQueries(ctx).GetOpinionsByTalkSessionID(ctx, model.GetOpinionsByTalkSessionIDParams{
		TalkSessionID: in.TalkSessionID.UUID(),
		Limit:         int32(*in.Limit),
		Offset:        int32(*in.Offset),
		SortKey:       sql.NullString{String: in.SortKey.String(), Valid: true},
		UserID:        userID,
	})
	if err != nil {
		utils.HandleError(ctx, err, "意見の取得に失敗")
		return nil, messages.OpinionContentFailedToFetch
	}

	var opinions []dto.SwipeOpinion
	if err := copier.CopyWithOption(&opinions, &opinionRows, copier.Option{
		DeepCopy:    true,
		IgnoreEmpty: true,
	}); err != nil {
		utils.HandleError(ctx, err, "マッピングに失敗")
		return nil, messages.OpinionContentFailedToFetch
	}

	count, err := g.GetQueries(ctx).CountOpinions(ctx, model.CountOpinionsParams{
		TalkSessionID: uuid.NullUUID{UUID: in.TalkSessionID.UUID(), Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &opinion_query.GetOpinionsByTalkSessionOutput{
		Opinions:   opinions,
		TotalCount: int(count),
	}, nil
}
