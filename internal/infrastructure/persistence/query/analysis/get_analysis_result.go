package analysis

import (
	"context"
	"errors"

	"github.com/jinzhu/copier"
	"github.com/neko-dream/server/internal/application/query/analysis_query"
	"github.com/neko-dream/server/internal/application/query/dto"
	"github.com/neko-dream/server/internal/domain/model/analysis"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	dto_mapper "github.com/neko-dream/server/internal/infrastructure/persistence/utils"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type GetAnalysisResultHandler struct {
	*db.DBManager
}

func NewGetAnalysisResultHandler(
	dbManager *db.DBManager,
) analysis_query.GetAnalysisResult {
	return &GetAnalysisResultHandler{
		DBManager: dbManager,
	}
}

// Execute implements GetAnalysisResultUseCase.
func (g *GetAnalysisResultHandler) Execute(ctx context.Context, input analysis_query.GetAnalysisResultInput) (*analysis_query.GetAnalysisResultOutput, error) {
	ctx, span := otel.Tracer("analysis_usecase").Start(ctx, "getAnalysisResultInteractor.Execute")
	defer span.End()

	groupInfoRows, err := g.GetQueries(ctx).GetGroupInfoByTalkSessionId(ctx, input.TalkSessionID.UUID())
	if err != nil {
		return nil, err
	}

	var myPosition *dto.UserPosition
	positions := make([]dto.UserPosition, 0, len(groupInfoRows))
	for _, row := range groupInfoRows {
		var position dto.UserPosition
		err = errors.Join(err, copier.CopyWithOption(&position, row, copier.Option{
			DeepCopy:    true,
			IgnoreEmpty: true,
		}))
		positions = append(positions, position)
		if input.UserID != nil && row.UserID == input.UserID.UUID() {
			myPosition = &position
		}
	}
	if err != nil {
		utils.HandleError(ctx, err, "copier.CopyWithOptionでエラー")
		return nil, err
	}

	groupIDs, err := g.GetQueries(ctx).GetGroupListByTalkSessionId(ctx, input.TalkSessionID.UUID())
	if err != nil {
		return nil, err
	}

	groupOpinionsMap := make(map[int32][]dto.OpinionWithRepresentative)
	for _, groupID := range groupIDs {
		groupOpinionsMap[groupID] = make([]dto.OpinionWithRepresentative, 0)
	}

	representativeRows, err := g.GetQueries(ctx).GetRepresentativeOpinionsByTalkSessionId(ctx, input.TalkSessionID.UUID())
	if err != nil {
		return nil, err
	}
	representatives := make([]dto.OpinionWithRepresentative, 0, len(representativeRows))
	for _, row := range representativeRows {
		res := dto.OpinionWithRepresentative{}

		err = copier.CopyWithOption(&res, row, copier.Option{
			DeepCopy:    true,
			IgnoreEmpty: true,
		})
		if err != nil {
			utils.HandleError(ctx, err, "copier.CopyWithOptionでエラー")
			return nil, err
		}
		representatives = append(representatives, res)
	}

	opinionIDs := dto_mapper.ExtractOpinionIDsWithRepresentative(representatives)
	reports, err := g.GetQueries(ctx).FindReportByOpinionIDs(ctx, model.FindReportByOpinionIDsParams{
		OpinionIds: opinionIDs,
		Status:     "deleted",
	})
	if err != nil {
		utils.HandleError(ctx, err, "通報情報の取得に失敗")
		return nil, err
	}
	representatives = dto_mapper.ProcessReportedOpinionsWithRepresentative(representatives, reports)
	for _, row := range representatives {
		groupOpinionsMap[int32(row.GroupID)] = append(groupOpinionsMap[int32(row.GroupID)], row)
	}

	groupOpinions := make([]dto.OpinionGroup, 0, len(groupOpinionsMap))
	for groupID, opinions := range groupOpinionsMap {
		groupOpinions = append(groupOpinions, dto.OpinionGroup{
			GroupName: analysis.NewGroupIDFromInt(int(groupID)).String(),
			GroupID:   int(groupID),
			Opinions:  opinions,
		})
	}

	return &analysis_query.GetAnalysisResultOutput{
		MyPosition:    myPosition,
		Positions:     positions,
		GroupOpinions: groupOpinions,
	}, nil
}
