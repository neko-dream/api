package opinion_query

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	dto_mapper "github.com/neko-dream/server/internal/infrastructure/persistence/utils"
	"github.com/neko-dream/server/internal/usecase/query/dto"
	opinion_query "github.com/neko-dream/server/internal/usecase/query/opinion"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
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
	ctx, span := otel.Tracer("opinion_query").Start(ctx, "GetOpinionsByTalkSessionIDQueryHandler.Execute")
	defer span.End()

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
		IsSeed:        sql.NullBool{Bool: in.IsSeed, Valid: true},
	})
	if err != nil {
		utils.HandleError(ctx, err, "意見の取得に失敗")
		return nil, messages.OpinionContentFailedToFetch
	}

	var opinions []dto.SwipeOpinion
	for _, opinionRow := range opinionRows {
		var op dto.SwipeOpinion
		if err := copier.CopyWithOption(&op, &opinionRow, copier.Option{
			DeepCopy:    true,
			IgnoreEmpty: true,
		}); err != nil {
			utils.HandleError(ctx, err, "マッピングに失敗")
			return nil, messages.OpinionContentFailedToFetch
		}

		if opinionRow.Opinion.ParentOpinionID.Valid {
			op.Opinion.ParentOpinionID = lo.ToPtr(shared.UUID[opinion.Opinion](opinionRow.Opinion.ParentOpinionID.UUID))
		}

		opinions = append(opinions, op)
	}

	// 通報された意見を処理
	if len(opinions) > 0 {
		opinionIDs := dto_mapper.ExtractOpinionIDs(opinions)
		reports, err := g.GetQueries(ctx).FindReportByOpinionIDs(ctx, model.FindReportByOpinionIDsParams{
			OpinionIds: opinionIDs,
			Status:     "deleted",
		})
		if err != nil {
			utils.HandleError(ctx, err, "通報情報の取得に失敗")
			return nil, err
		}

		opinions = dto_mapper.ProcessReportedOpinions(opinions, reports)
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
