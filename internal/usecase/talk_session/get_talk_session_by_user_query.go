package talk_session_usecase

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/shared/time"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/db"
	model "github.com/neko-dream/server/internal/infrastructure/db/sqlc"
	"github.com/neko-dream/server/pkg/utils"
)

type (
	GetTalkSessionByUserQuery interface {
		Execute(context.Context, GetTalkSessionByUserInput) (*GetTalkSessionByUserOutput, error)
	}

	GetTalkSessionByUserInput struct {
		UserID shared.UUID[user.User]
		Limit  int
		Offset int
		Status string
		Theme  *string
	}

	GetTalkSessionByUserOutput struct {
		TalkSessions []TalkSessionDTO
	}

	GetTalkSessionByUserQueryHandler struct {
		*db.DBManager
	}
)

func NewGetTalkSessionByUserQueryHandler(tm *db.DBManager) GetTalkSessionByUserQuery {
	return &GetTalkSessionByUserQueryHandler{
		DBManager: tm,
	}
}

func (h *GetTalkSessionByUserQueryHandler) Execute(ctx context.Context, input GetTalkSessionByUserInput) (*GetTalkSessionByUserOutput, error) {

	var limit, offset int
	if input.Limit == 0 {
		limit = 10
	} else {
		limit = input.Limit
	}
	if input.Offset == 0 {
		offset = 0
	} else {
		offset = input.Offset
	}

	var status sql.NullString
	if input.Status == "" {
		status = sql.NullString{String: "", Valid: false}
	} else {
		status = sql.NullString{String: input.Status, Valid: true}
	}

	var theme sql.NullString
	if input.Theme == nil {
		theme = sql.NullString{String: "", Valid: false}
	} else {
		theme = sql.NullString{String: *input.Theme, Valid: true}
	}

	out, err := h.GetQueries(ctx).GetOwnTalkSessionByUserID(ctx, model.GetOwnTalkSessionByUserIDParams{
		UserID: uuid.NullUUID{UUID: input.UserID.UUID(), Valid: true},
		Limit:  int32(limit),
		Offset: int32(offset),
		Theme:  theme,
		Status: status,
	})
	if err != nil {
		return nil, err
	}
	if len(out) <= 0 {
		return &GetTalkSessionByUserOutput{
			TalkSessions: make([]TalkSessionDTO, 0),
		}, nil
	}

	talkSessions := make([]TalkSessionDTO, 0, len(out))
	for _, ts := range out {
		var locationDTO *LocationDTO
		if ts.LocationID.Valid {
			locationDTO = &LocationDTO{
				Latitude:  ts.Latitude,
				Longitude: ts.Longitude,
			}
		}
		talkSessions = append(talkSessions, TalkSessionDTO{
			ID:    ts.TalkSessionID.String(),
			Theme: ts.Theme,
			Owner: UserDTO{
				DisplayID:   ts.DisplayID.String,
				DisplayName: ts.DisplayName.String,
				IconURL: utils.ToPtrIfNotNullValue[string](
					ts.IconUrl.Valid,
					ts.IconUrl.String,
				),
			},
			OpinionCount:     int(ts.OpinionCount),
			CreatedAt:        time.NewTime(ctx, ts.CreatedAt).Format(ctx),
			ScheduledEndTime: time.NewTime(ctx, ts.ScheduledEndTime).Format(ctx),
			Location:         locationDTO,
			City:             utils.ToPtrIfNotNullValue[string](!ts.City.Valid, ts.City.String),
			Prefecture:       utils.ToPtrIfNotNullValue[string](!ts.Prefecture.Valid, ts.Prefecture.String),
		})
	}

	return &GetTalkSessionByUserOutput{
		TalkSessions: talkSessions,
	}, nil
}
