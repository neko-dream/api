package opinion_query

import (
	"context"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/neko-dream/server/internal/application/query/dto"
	opinion_query "github.com/neko-dream/server/internal/application/query/opinion"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	dto_mapper "github.com/neko-dream/server/internal/infrastructure/persistence/utils"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
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
	ctx, span := otel.Tracer("opinion_query").Start(ctx, "GetOpinionDetailByIDQueryHandler.Execute")
	defer span.End()

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

	var op dto.SwipeOpinion
	if err := copier.CopyWithOption(&op, &opinionRow, copier.Option{
		DeepCopy:    true,
		IgnoreEmpty: true,
	}); err != nil {
		utils.HandleError(ctx, err, "マッピングに失敗")
		return nil, err
	}

	// 通報された意見を処理
	opinionIDs := []uuid.UUID{op.Opinion.OpinionID.UUID()}
	reports, err := g.GetQueries(ctx).FindReportByOpinionIDs(ctx, model.FindReportByOpinionIDsParams{
		OpinionIds: opinionIDs,
		Status:     "deleted",
	})
	if err != nil {
		utils.HandleError(ctx, err, "通報情報の取得に失敗")
		return nil, err
	}

	dto_mapper.ProcessSingleReportedOpinion(&op, reports)

	return &opinion_query.GetOpinionDetailByIDOutput{
		Opinion: op,
	}, nil
}
