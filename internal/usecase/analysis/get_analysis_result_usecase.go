package analysis_usecase

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/db"
	"github.com/neko-dream/server/pkg/utils"
)

type (
	GetAnalysisResultUseCase interface {
		Execute(context.Context, GetAnalysisResultInput) (*GetAnalysisResultOutput, error)
	}

	GetAnalysisResultInput struct {
		UserID        *shared.UUID[user.User]
		TalkSessionID shared.UUID[talksession.TalkSession]
	}

	GetAnalysisResultOutput struct {
		// ユーザーがログインしていれば、自分の位置情報を返す
		MyPosition *PositionDTO
		// トークセッションの全てのポジション情報を返す
		Positions []PositionDTO
		// トークセッションの全てのグループの意見情報を返す
		GroupOpinions []GroupOpinionDTO
	}
	PositionDTO struct {
		PosX           float64
		PosY           float64
		DisplayID      string
		GroupID        int
		PerimeterIndex *int
	}
	GroupOpinionDTO struct {
		GroupID  int
		Opinions []OpinionRootDTO
	}
	OpinionRootDTO struct {
		User    UserDTO
		Opinion OpinionDTO
	}
	UserDTO struct {
		ID   string
		Name string
		Icon *string
	}
	OpinionDTO struct {
		ID           string
		Title        *string
		Content      string
		ParentID     *string
		PictureURL   *string
		ReferenceURL *string
	}

	getAnalysisResultInteractor struct {
		*db.DBManager
	}
)

func NewGetAnalysisResultUseCase(
	dbManager *db.DBManager,
) GetAnalysisResultUseCase {
	return &getAnalysisResultInteractor{
		DBManager: dbManager,
	}
}

// Execute implements GetAnalysisResultUseCase.
func (g *getAnalysisResultInteractor) Execute(ctx context.Context, input GetAnalysisResultInput) (*GetAnalysisResultOutput, error) {
	groupInfoRows, err := g.GetQueries(ctx).GetGroupInfoByTalkSessionId(ctx, input.TalkSessionID.UUID())
	if err != nil {
		return nil, err
	}

	var myPosition *PositionDTO
	positions := make([]PositionDTO, 0, len(groupInfoRows))
	for _, row := range groupInfoRows {
		if input.UserID != nil && row.UserID == input.UserID.UUID() {
			myPosition = &PositionDTO{
				PosX:           row.PosX,
				PosY:           row.PosY,
				DisplayID:      row.DisplayID.String,
				GroupID:        int(row.GroupID),
				PerimeterIndex: utils.ToPtrIfNotNullValue(!row.PerimeterIndex.Valid, int(row.PerimeterIndex.Int32)),
			}
		}

		positions = append(positions, PositionDTO{
			PosX:      row.PosX,
			PosY:      row.PosY,
			DisplayID: row.DisplayID.String,
			GroupID:   int(row.GroupID),
		})
	}
	groupIDs, err := g.GetQueries(ctx).GetGroupListByTalkSessionId(ctx, input.TalkSessionID.UUID())
	if err != nil {
		return nil, err
	}

	groupOpinionsMap := make(map[int32][]OpinionRootDTO)
	for _, groupID := range groupIDs {
		groupOpinionsMap[groupID] = make([]OpinionRootDTO, 0)
	}

	representativeRows, err := g.GetQueries(ctx).GetRepresentativeOpinionsByTalkSessionId(ctx, input.TalkSessionID.UUID())
	if err != nil {
		return nil, err
	}

	for _, row := range representativeRows {
		groupOpinionsMap[row.GroupID] = append(groupOpinionsMap[row.GroupID], OpinionRootDTO{
			Opinion: OpinionDTO{
				ID:           row.OpinionID.UUID.String(),
				Title:        utils.ToPtrIfNotNullValue(!row.Title.Valid, row.Title.String),
				Content:      row.Content.String,
				ParentID:     utils.ToPtrIfNotNullValue(!row.ParentOpinionID.Valid, row.ParentOpinionID.UUID.String()),
				PictureURL:   utils.ToPtrIfNotNullValue(!row.PictureUrl.Valid, row.PictureUrl.String),
				ReferenceURL: utils.ToPtrIfNotNullValue(!row.ReferenceUrl.Valid, row.ReferenceUrl.String),
			},
			User: UserDTO{
				ID:   row.DisplayID.String,
				Name: row.DisplayName.String,
				Icon: utils.ToPtrIfNotNullValue(!row.IconUrl.Valid, row.IconUrl.String),
			},
		})
	}

	groupOpinions := make([]GroupOpinionDTO, 0, len(groupOpinionsMap))
	for groupID, opinions := range groupOpinionsMap {
		groupOpinions = append(groupOpinions, GroupOpinionDTO{
			GroupID:  int(groupID),
			Opinions: opinions,
		})
	}

	return &GetAnalysisResultOutput{
		MyPosition:    myPosition,
		Positions:     positions,
		GroupOpinions: groupOpinions,
	}, nil
}
