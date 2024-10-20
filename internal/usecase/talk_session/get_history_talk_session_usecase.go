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
	GetTalkSessionHistoriesQuery interface {
		Execute(context.Context, GetTalkSessionHistoriesInput) (*GetTalkSessionHistoriesOutput, error)
	}

	GetTalkSessionHistoriesInput struct {
		UserID shared.UUID[user.User]
		Status string
		Theme  *string
		Limit  *int
		Offset *int
	}

	GetTalkSessionHistoriesOutput struct {
		TalkSessions []TalkSessionDTO
		TotalCount   int
		Limit        int
		Offset       int
	}

	getTalkSessionHistoriesQueryHandler struct {
		*db.DBManager
	}
)

func NewGetTalkSessionHistoriesQuery(
	dbm *db.DBManager,
) GetTalkSessionHistoriesQuery {
	return &getTalkSessionHistoriesQueryHandler{
		DBManager: dbm,
	}
}

func (h *getTalkSessionHistoriesQueryHandler) Execute(ctx context.Context, q GetTalkSessionHistoriesInput) (*GetTalkSessionHistoriesOutput, error) {
	var limit, offset int
	if q.Limit == nil {
		limit = 10
	}
	if q.Offset == nil {
		offset = 0
	}
	if q.Status == "" {
		q.Status = "open"
	}
	var theme sql.NullString
	if q.Theme != nil {
		theme = sql.NullString{String: *q.Theme, Valid: true}
	}

	talkSessions, err := h.GetQueries(ctx).GetTalkSessionByUserID(ctx, model.GetTalkSessionByUserIDParams{
		Limit:  int32(limit),
		Offset: int32(offset),
		UserID: uuid.NullUUID{UUID: q.UserID.UUID(), Valid: true},
		Status: sql.NullString{String: q.Status, Valid: true},
		Theme:  theme,
	})
	if err != nil {
		return nil, err
	}
	totalCount, err := h.GetQueries(ctx).CountTalkSessions(ctx, model.CountTalkSessionsParams{
		Theme:  theme,
		Status: sql.NullString{String: q.Status, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	talkSessionDTOs := make([]TalkSessionDTO, 0, len(talkSessions))
	for _, talkSession := range talkSessions {
		talkSessionDTOs = append(talkSessionDTOs, TalkSessionDTO{
			ID:    talkSession.TalkSessionID.String(),
			Theme: talkSession.Theme,
			Owner: UserDTO{
				DisplayID:   talkSession.DisplayID.String,
				DisplayName: talkSession.DisplayName.String,
				IconURL:     utils.ToPtrIfNotNullValue[string](!talkSession.IconUrl.Valid, talkSession.IconUrl.String),
			},
			OpinionCount:     int(talkSession.OpinionCount),
			CreatedAt:        time.NewTime(ctx, talkSession.CreatedAt).Format(ctx),
			ScheduledEndTime: time.NewTime(ctx, talkSession.ScheduledEndTime).Format(ctx),
			Location: utils.ToPtrIfNotNullValue(
				!talkSession.LocationID.Valid,
				LocationDTO{
					Latitude:  talkSession.Latitude,
					Longitude: talkSession.Longitude,
				}),
			City:       utils.ToPtrIfNotNullValue[string](talkSession.City.Valid, talkSession.City.String),
			Prefecture: utils.ToPtrIfNotNullValue[string](talkSession.Prefecture.Valid, talkSession.Prefecture.String),
		})
	}

	return &GetTalkSessionHistoriesOutput{
		TalkSessions: talkSessionDTOs,
		TotalCount:   int(totalCount),
		Limit:        limit,
		Offset:       offset,
	}, nil
}
