package opinion_usecase

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/db"
	model "github.com/neko-dream/server/internal/infrastructure/db/sqlc"
	"github.com/neko-dream/server/pkg/utils"
)

type (
	GetUserOpinionListQueryHandler interface {
		Execute(context.Context, GetUserOpinionListQuery) (*GetUserOpinionListOutput, error)
	}

	GetUserOpinionListQuery struct {
		UserID shared.UUID[user.User]
		// latest, mostReply, oldest
		SortKey *string
		Limit   *int
		Offset  *int
	}

	GetUserOpinionListOutput struct {
		Opinions   []SwipeOpinionDTO
		TotalCount int
	}

	getUserOpinionListQueryHandler struct {
		*db.DBManager
	}
)

func NewGetUserOpinionListQueryHandler(
	dbManager *db.DBManager,
) GetUserOpinionListQueryHandler {
	return &getUserOpinionListQueryHandler{
		DBManager: dbManager,
	}
}

func (h *getUserOpinionListQueryHandler) Execute(ctx context.Context, q GetUserOpinionListQuery) (*GetUserOpinionListOutput, error) {
	var limit, offset int
	if q.Limit == nil {
		limit = 10
	} else {
		limit = *q.Limit
	}
	if q.Offset == nil {
		offset = 0
	} else {
		offset = *q.Offset
	}
	sortKey := "latest"
	if q.SortKey != nil {
		sortKey = *q.SortKey
	}

	opinionRows, err := h.GetQueries(ctx).GetOpinionsByUserID(ctx, model.GetOpinionsByUserIDParams{
		UserID:  q.UserID.UUID(),
		SortKey: sql.NullString{String: sortKey, Valid: true},
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
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

	count, err := h.GetQueries(ctx).CountOpinions(ctx, model.CountOpinionsParams{
		UserID: uuid.NullUUID{UUID: q.UserID.UUID(), Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &GetUserOpinionListOutput{
		Opinions:   opinions,
		TotalCount: int(count),
	}, nil
}
