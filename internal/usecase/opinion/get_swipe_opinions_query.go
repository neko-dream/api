package opinion_usecase

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/domain/model/vote"
	"github.com/neko-dream/server/internal/infrastructure/db"
	model "github.com/neko-dream/server/internal/infrastructure/db/sqlc"
	"github.com/neko-dream/server/pkg/utils"
)

type (
	GetSwipeOpinionsQueryHandler interface {
		Execute(context.Context, GetSwipeOpinionsQuery) (*GetSwipeOpinionsOutput, error)
	}

	GetSwipeOpinionsQuery struct {
		UserID        shared.UUID[user.User]
		TalkSessionID shared.UUID[talksession.TalkSession]
		Limit         int
	}

	GetSwipeOpinionsOutput struct {
		Opinions []SwipeOpinionDTO
	}

	SwipeOpinionDTO struct {
		Opinion    OpinionDTO
		User       UserDTO
		ReplyCount int
	}

	getSwipeOpinionsQueryHandler struct {
		*db.DBManager
	}
)

func NewGetSwipeOpinionsQueryHandler(
	dbManager *db.DBManager,
) GetSwipeOpinionsQueryHandler {
	return &getSwipeOpinionsQueryHandler{
		DBManager: dbManager,
	}
}

func (h *getSwipeOpinionsQueryHandler) Execute(ctx context.Context, q GetSwipeOpinionsQuery) (*GetSwipeOpinionsOutput, error) {

	// limitが3以上の場合、2件はtop, 1件はrandomで取得する
	// 1:2の割合を算出する
	topLimit := q.Limit / 3
	randomLimit := q.Limit - topLimit
	var swipeRows []model.GetRandomOpinionsRow
	if q.Limit >= 3 {
		swipeRow, err := h.GetQueries(ctx).GetRandomOpinions(ctx, model.GetRandomOpinionsParams{
			UserID:        q.UserID.UUID(),
			TalkSessionID: q.TalkSessionID.UUID(),
			Limit:         int32(topLimit),
		})
		if err != nil {
			return nil, err
		}
		swipeRows = swipeRow

		randomSwipeRow, err := h.GetQueries(ctx).GetRandomOpinions(ctx, model.GetRandomOpinionsParams{
			UserID:        q.UserID.UUID(),
			TalkSessionID: q.TalkSessionID.UUID(),
			Limit:         int32(randomLimit),
		})
		if err != nil {
			return nil, err
		}
		swipeRows = append(swipeRows, randomSwipeRow...)
	} else {
		swipeRow, err := h.GetQueries(ctx).GetRandomOpinions(ctx, model.GetRandomOpinionsParams{
			UserID:        q.UserID.UUID(),
			TalkSessionID: q.TalkSessionID.UUID(),
			Limit:         int32(q.Limit),
		})
		if err != nil {
			return nil, err
		}
		swipeRows = swipeRow
	}

	opinions := make([]SwipeOpinionDTO, 0, len(swipeRows))
	for _, row := range swipeRows {
		opinionDTO := OpinionDTO{
			OpinionID:       row.OpinionID.String(),
			TalkSessionID:   row.TalkSessionID.String(),
			UserID:          row.UserID.String(),
			ParentOpinionID: utils.ToPtrIfNotNullValue[string](!row.ParentOpinionID.Valid, row.ParentOpinionID.UUID.String()),
			Title:           utils.ToPtrIfNotNullValue[string](!row.Title.Valid, row.Title.String),
			Content:         row.Content,
			CreatedAt:       row.CreatedAt,
			VoteType:        vote.VoteTypeFromInt(int(row.VoteType)).String(),
			ReferenceURL:    utils.ToPtrIfNotNullValue[string](!row.ReferenceUrl.Valid, row.ReferenceUrl.String),
			PictureURL:      utils.ToPtrIfNotNullValue[string](!row.PictureUrl.Valid, row.PictureUrl.String),
		}
		userDTO := UserDTO{
			ID:   row.DisplayID.String,
			Name: row.DisplayName.String,
			Icon: utils.ToPtrIfNotNullValue[string](!row.IconUrl.Valid, row.IconUrl.String),
		}
		opinions = append(opinions, SwipeOpinionDTO{
			Opinion:    opinionDTO,
			User:       userDTO,
			ReplyCount: int(row.ReplyCount),
		})
	}

	return &GetSwipeOpinionsOutput{
		Opinions: opinions,
	}, nil
}
