package opinion_query

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/neko-dream/api/internal/application/query/dto"
	opinion_query "github.com/neko-dream/api/internal/application/query/opinion"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/opinion"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
	dto_mapper "github.com/neko-dream/api/internal/infrastructure/persistence/utils"
	"github.com/neko-dream/api/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
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
	ctx, span := otel.Tracer("opinion_query").Start(ctx, "GetMyOpinionsQueryHandler.Execute")
	defer span.End()

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
	for _, opinionRow := range opinionRows {
		var op dto.SwipeOpinion
		if err := copier.CopyWithOption(&op, &opinionRow, copier.Option{
			DeepCopy: true,
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
