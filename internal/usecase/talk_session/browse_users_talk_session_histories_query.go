package talk_session_usecase

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/pkg/utils"
)

type (
	BrowseUsersTalkSessionHistoriesQuery interface {
		Execute(context.Context, BrowseUsersTalkSessionHistoriesInput) (*BrowseUsersTalkSessionHistoriesOutput, error)
	}

	BrowseUsersTalkSessionHistoriesInput struct {
		UserID shared.UUID[user.User]
		Limit  int
		Offset int
		Status string
		Theme  *string
	}

	BrowseUsersTalkSessionHistoriesOutput struct {
		TalkSessions []TalkSessionDTO
	}

	BrowseUsersTalkSessionHistoriesQueryHandler struct {
		*db.DBManager
	}
)

func NewBrowseUsersTalkSessionHistoriesQueryHandler(tm *db.DBManager) BrowseUsersTalkSessionHistoriesQuery {
	return &BrowseUsersTalkSessionHistoriesQueryHandler{
		DBManager: tm,
	}
}

func (h *BrowseUsersTalkSessionHistoriesQueryHandler) Execute(ctx context.Context, input BrowseUsersTalkSessionHistoriesInput) (*BrowseUsersTalkSessionHistoriesOutput, error) {

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
		return nil, messages.TalkSessionNotFound
	}
	if len(out) <= 0 {
		return &BrowseUsersTalkSessionHistoriesOutput{
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
			Description: utils.ToPtrIfNotNullValue[string](
				!ts.Description.Valid,
				ts.Description.String,
			),
			Owner: UserDTO{
				DisplayID:   ts.DisplayID.String,
				DisplayName: ts.DisplayName.String,
				IconURL: utils.ToPtrIfNotNullValue[string](
					!ts.IconUrl.Valid,
					ts.IconUrl.String,
				),
			},
			OpinionCount:     int(ts.OpinionCount),
			CreatedAt:        ts.CreatedAt.Format(time.RFC3339),
			ScheduledEndTime: ts.ScheduledEndTime.Format(time.RFC3339),
			Location:         locationDTO,
			City:             utils.ToPtrIfNotNullValue[string](!ts.City.Valid, ts.City.String),
			Prefecture:       utils.ToPtrIfNotNullValue[string](!ts.Prefecture.Valid, ts.Prefecture.String),
		})
	}

	return &BrowseUsersTalkSessionHistoriesOutput{
		TalkSessions: talkSessions,
	}, nil
}
