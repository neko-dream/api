package opinion_usecase

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/domain/model/vote"
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
		UserID        *shared.UUID[user.User]
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
	userID := uuid.NullUUID{}
	if input.UserID != nil {
		userID = uuid.NullUUID{UUID: input.UserID.UUID(), Valid: true}
	}

	// 親意見を取得
	opinionRows, err := i.GetQueries(ctx).GetOpinionsByTalkSessionID(ctx, model.GetOpinionsByTalkSessionIDParams{
		TalkSessionID: input.TalkSessionID.UUID(),
		Limit:         int32(limit),
		Offset:        int32(offset),
		SortKey:       sql.NullString{String: sortKey, Valid: true},
		UserID:        userID,
	})
	if err != nil {
		utils.HandleError(ctx, err, "GetOpinionsByTalkSessionInteractor.Execute")
		return nil, messages.OpinionContentBadLengthForUpdate
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
				VoteType:        vote.VoteTypeFromInt(int(row.VoteType)).String(),
				PictureURL:      utils.ToPtrIfNotNullValue(!row.PictureUrl.Valid, row.PictureUrl.String),
				ReferenceURL:    utils.ToPtrIfNotNullValue(!row.ReferenceUrl.Valid, row.ReferenceUrl.String),
			},
			User: UserDTO{
				ID:   row.DisplayID.String,
				Name: row.DisplayName.String,
				Icon: utils.ToPtrIfNotNullValue(!row.IconUrl.Valid, row.IconUrl.String),
			},
			ReplyCount: int(row.ReplyCount),
			MyVoteType: vote.VoteTypeFromInt(int(row.CurrentVoteType)).String(),
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
