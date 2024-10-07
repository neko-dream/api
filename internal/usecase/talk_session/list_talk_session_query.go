package talk_session_usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/neko-dream/server/internal/domain/model/shared/time"
	"github.com/neko-dream/server/internal/infrastructure/db"
	model "github.com/neko-dream/server/internal/infrastructure/db/sqlc"
	"github.com/neko-dream/server/pkg/utils"
)

type (
	ListTalkSessionQuery interface {
		Execute(context.Context, ListTalkSessionInput) (*ListTalkSessionOutput, error)
	}

	ListTalkSessionInput struct {
		Limit  int
		Offset int
		Theme  *string
		Status string
	}

	ListTalkSessionOutput struct {
		TalkSessions []TalkSessionDTO
	}

	TalkSessionDTO struct {
		ID           string
		Theme        string
		Owner        UserDTO
		OpinionCount int
		FinishedAt   *string
		CreatedAt    string
	}

	UserDTO struct {
		DisplayID   string
		DisplayName string
		IconURL     *string
	}

	listTalkSessionQueryHandler struct {
		*db.DBManager
	}
)

func NewListTalkSessionQueryHandler(
	tm *db.DBManager,
) ListTalkSessionQuery {
	return &listTalkSessionQueryHandler{
		DBManager: tm,
	}
}

func (h *listTalkSessionQueryHandler) Execute(ctx context.Context, input ListTalkSessionInput) (*ListTalkSessionOutput, error) {
	var talkSessionOut ListTalkSessionOutput
	if input.Status == "" {
		input.Status = "open"
	}

	talkSessionRow, err := h.GetQueries(ctx).ListTalkSessions(ctx, model.ListTalkSessionsParams{
		Limit:  int32(input.Limit),
		Offset: int32(input.Offset),
		Theme: utils.IfThenElse(
			input.Theme != nil,
			sql.NullString{String: *input.Theme, Valid: true},
			sql.NullString{},
		),
		Status: sql.NullString{String: input.Status, Valid: true},
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			talkSessionOut.TalkSessions = make([]TalkSessionDTO, 0)
			return &talkSessionOut, nil
		}
		return nil, err
	}
	if len(talkSessionRow) <= 0 {
		talkSessionOut.TalkSessions = make([]TalkSessionDTO, 0)
		return &talkSessionOut, nil
	}

	talkSessionDTOList := make([]TalkSessionDTO, 0, len(talkSessionRow))
	for _, row := range talkSessionRow {
		talkSessionDTOList = append(talkSessionDTOList, TalkSessionDTO{
			ID:    row.TalkSessionID.String(),
			Theme: row.Theme,
			Owner: UserDTO{
				DisplayID:   row.DisplayID.String,
				DisplayName: row.DisplayName.String,
				IconURL: utils.ToPtrIfNotNullValue[string](
					row.IconUrl.Valid,
					row.IconUrl.String,
				),
			},
			OpinionCount: int(row.OpinionCount),
			CreatedAt:    time.NewTime(ctx, row.CreatedAt).Format(ctx),
			FinishedAt: utils.ToPtrIfNotNullFunc[string](!row.FinishedAt.Valid, func() string {
				return time.NewTime(ctx, row.FinishedAt.Time).Format(ctx)
			}),
		})
	}

	talkSessionOut.TalkSessions = talkSessionDTOList
	return &talkSessionOut, nil
}
