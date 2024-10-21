package opinion_usecase

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/infrastructure/db"
	model "github.com/neko-dream/server/internal/infrastructure/db/sqlc"
	"github.com/neko-dream/server/pkg/utils"
)

type (
	GetOpinionsByTalkSessionUseCase interface {
		Execute(context.Context, GetOpinionsByTalkSessionInput) (*GetOpinionsByTalkSessionOutput, error)
	}

	GetOpinionsByTalkSessionInput struct {
		TalkSessionID shared.UUID[talksession.TalkSession]
		SortKey       *string
		Limit         *int
		Offset        *int
	}

	GetOpinionsByTalkSessionOutput struct {
		Opinions   []SwipeOpinionDTO
		TotalCount int
	}

	getOpinionsByTalkSessionInteractor struct {
		*db.DBManager
	}
)

func NewGetOpinionsByTalkSessionUseCase(
	dbManager *db.DBManager,
) GetOpinionsByTalkSessionUseCase {
	return &getOpinionsByTalkSessionInteractor{
		DBManager: dbManager,
	}
}

func (i *getOpinionsByTalkSessionInteractor) Execute(ctx context.Context, input GetOpinionsByTalkSessionInput) (*GetOpinionsByTalkSessionOutput, error) {
	var limit, offset int
	if input.Limit == nil {
		limit = 10
	} else {
		limit = *input.Limit
	}
	if input.Offset == nil {
		offset = 0
	} else {
		offset = *input.Offset
	}
	sortKey := "latest"
	if input.SortKey != nil {
		sortKey = *input.SortKey
	}

	// 親意見を取得
	opinionRows, err := i.GetQueries(ctx).GetOpinionsByTalkSessionID(ctx, model.GetOpinionsByTalkSessionIDParams{
		TalkSessionID: input.TalkSessionID.UUID(),
		Limit:         int32(limit),
		Offset:        int32(offset),
		SortKey:       sql.NullString{String: sortKey, Valid: true},
	})
	if err != nil {
		utils.HandleError(ctx, err, "GetOpinionsByTalkSessionInteractor.Execute")
		return nil, err
	}

	opinions := make([]SwipeOpinionDTO, 0, len(opinionRows))
	for _, row := range opinionRows {
		opinions = append(opinions, SwipeOpinionDTO{
			Opinion: OpinionDTO{
				OpinionID:       row.OpinionID.String(),
				TalkSessionID:   row.TalkSessionID.String(),
				UserID:          row.UserID.String(),
				ParentOpinionID: utils.ToPtrIfNotNullValue(!row.ParentOpinionID.Valid, row.ParentOpinionID.UUID.String()),
				Title:           utils.ToPtrIfNotNullValue(!row.Title.Valid, row.Title.String),
				Content:         row.Content,
				CreatedAt:       row.CreatedAt,
			},
			User: UserDTO{
				ID:   row.DisplayID.String,
				Name: row.DisplayName.String,
				Icon: utils.ToPtrIfNotNullValue(!row.IconUrl.Valid, row.IconUrl.String),
			},
			ReplyCount: int(row.ReplyCount),
		})
	}

	count, err := i.GetQueries(ctx).CountOpinions(ctx, model.CountOpinionsParams{
		TalkSessionID: uuid.NullUUID{UUID: input.TalkSessionID.UUID(), Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &GetOpinionsByTalkSessionOutput{
		Opinions:   opinions,
		TotalCount: int(count),
	}, nil
}
