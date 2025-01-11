package talk_session_usecase

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/pkg/utils"
)

type (
	SearchTalkSessionsQuery interface {
		Execute(context.Context, SearchTalkSessionsInput) (*SearchTalkSessionsOutput, error)
	}

	SearchTalkSessionsInput struct {
		UserID shared.UUID[user.User]
		Status string
		Theme  *string
		Limit  *int
		Offset *int
	}

	SearchTalkSessionsOutput struct {
		TalkSessions []TalkSessionDTO
		TotalCount   int
		Limit        *int
		Offset       *int
	}

	SearchTalkSessionsQueryHandler struct {
		*db.DBManager
	}
)

func NewSearchTalkSessionsQuery(
	dbm *db.DBManager,
) SearchTalkSessionsQuery {
	return &SearchTalkSessionsQueryHandler{
		DBManager: dbm,
	}
}

func (h *SearchTalkSessionsQueryHandler) Execute(ctx context.Context, q SearchTalkSessionsInput) (*SearchTalkSessionsOutput, error) {
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
	var status sql.NullString
	if q.Status != "" {
		status = sql.NullString{String: q.Status, Valid: true}
	} else {
		status = sql.NullString{String: "open", Valid: true}
	}

	var theme sql.NullString
	if q.Theme != nil {
		theme = sql.NullString{String: *q.Theme, Valid: true}
	}

	talkSessions, err := h.GetQueries(ctx).GetRespondTalkSessionByUserID(ctx, model.GetRespondTalkSessionByUserIDParams{
		Limit:  int32(limit),
		Offset: int32(offset),
		UserID: uuid.NullUUID{UUID: q.UserID.UUID(), Valid: true},
		Status: status,
		Theme:  theme,
	})
	if err != nil {
		utils.HandleError(ctx, err, "SearchTalkSessionsQueryHandler.Execute")
		return nil, err
	}
	totalCount, err := h.GetQueries(ctx).CountTalkSessions(ctx, model.CountTalkSessionsParams{
		Theme:  theme,
		Status: status,
		UserID: uuid.NullUUID{UUID: q.UserID.UUID(), Valid: true},
	})
	if err != nil {
		utils.HandleError(ctx, err, "SearchTalkSessionsQueryHandler.Execute")
		return nil, err
	}

	talkSessionDTOs := make([]TalkSessionDTO, 0, len(talkSessions))
	for _, talkSession := range talkSessions {
		talkSessionDTOs = append(talkSessionDTOs, TalkSessionDTO{
			ID:          talkSession.TalkSessionID.String(),
			Theme:       talkSession.Theme,
			Description: utils.ToPtrIfNotNullValue[string](!talkSession.Description.Valid, talkSession.Description.String),
			Owner: UserDTO{
				DisplayID:   talkSession.DisplayID.String,
				DisplayName: talkSession.DisplayName.String,
				IconURL:     utils.ToPtrIfNotNullValue[string](!talkSession.IconUrl.Valid, talkSession.IconUrl.String),
			},
			OpinionCount:     int(talkSession.OpinionCount),
			CreatedAt:        talkSession.CreatedAt.Format(time.RFC3339),
			ScheduledEndTime: talkSession.ScheduledEndTime.Format(time.RFC3339),
			Location: utils.ToPtrIfNotNullValue(
				!talkSession.LocationID.Valid,
				LocationDTO{
					Latitude:  talkSession.Latitude,
					Longitude: talkSession.Longitude,
				}),
			City:       utils.ToPtrIfNotNullValue[string](!talkSession.City.Valid, talkSession.City.String),
			Prefecture: utils.ToPtrIfNotNullValue[string](!talkSession.Prefecture.Valid, talkSession.Prefecture.String),
		})
	}
	return &SearchTalkSessionsOutput{
		TalkSessions: talkSessionDTOs,
		TotalCount:   int(totalCount.TalkSessionCount),
		Limit:        limit,
		Offset:       offset,
	}, nil
}
