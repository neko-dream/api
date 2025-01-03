package talk_session_usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/pkg/utils"
)

type (
	ListTalkSessionQuery interface {
		Execute(context.Context, ListTalkSessionInput) (*ListTalkSessionOutput, error)
	}

	ListTalkSessionInput struct {
		Limit     int
		Offset    int
		Theme     *string
		Status    string
		SortKey   *string
		Latitude  *float64
		Longitude *float64
	}

	ListTalkSessionOutput struct {
		TalkSessions []TalkSessionDTO
		TotalCount   int
	}

	TalkSessionDTO struct {
		ID               string
		Theme            string
		Description      *string
		Owner            UserDTO
		OpinionCount     int
		CreatedAt        string
		ScheduledEndTime string
		Location         *LocationDTO
		City             *string
		Prefecture       *string
	}

	UserDTO struct {
		DisplayID   string
		DisplayName string
		IconURL     *string
	}
	LocationDTO struct {
		Latitude  float64
		Longitude float64
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
	var sortKey string
	if input.SortKey != nil {
		sortKey = *input.SortKey
	} else {
		sortKey = "latest"
	}
	var latitude, longitude sql.NullFloat64
	if input.Latitude != nil {
		latitude = sql.NullFloat64{Float64: *input.Latitude, Valid: true}
	}
	if input.Longitude != nil {
		longitude = sql.NullFloat64{Float64: *input.Longitude, Valid: true}
	}

	talkSessionRow, err := h.GetQueries(ctx).ListTalkSessions(ctx, model.ListTalkSessionsParams{
		Limit:  int32(input.Limit),
		Offset: int32(input.Offset),
		Theme: utils.IfThenElse(
			input.Theme != nil,
			sql.NullString{String: *input.Theme, Valid: true},
			sql.NullString{},
		),
		Status:    sql.NullString{String: input.Status, Valid: true},
		SortKey:   sql.NullString{String: sortKey, Valid: true},
		Latitude:  latitude,
		Longitude: longitude,
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
		var locationDTO *LocationDTO
		if row.LocationID.Valid {
			locationDTO = &LocationDTO{
				Latitude:  row.Latitude,
				Longitude: row.Longitude,
			}
		}
		talkSessionDTOList = append(talkSessionDTOList, TalkSessionDTO{
			ID:    row.TalkSessionID.String(),
			Theme: row.Theme,
			Description: utils.ToPtrIfNotNullValue[string](
				!row.Description.Valid,
				row.Description.String,
			),
			Owner: UserDTO{
				DisplayID:   row.DisplayID.String,
				DisplayName: row.DisplayName.String,
				IconURL: utils.ToPtrIfNotNullValue[string](
					!row.IconUrl.Valid,
					row.IconUrl.String,
				),
			},
			OpinionCount:     int(row.OpinionCount),
			CreatedAt:        row.CreatedAt.Format(time.RFC3339),
			ScheduledEndTime: row.ScheduledEndTime.Format(time.RFC3339),
			Location:         locationDTO,
			City:             utils.ToPtrIfNotNullValue[string](!row.City.Valid, row.City.String),
			Prefecture:       utils.ToPtrIfNotNullValue[string](!row.Prefecture.Valid, row.Prefecture.String),
		})
	}

	talkSessionCount, err := h.GetQueries(ctx).CountTalkSessions(ctx, model.CountTalkSessionsParams{
		Theme: utils.IfThenElse(
			input.Theme != nil,
			sql.NullString{String: *input.Theme, Valid: true},
			sql.NullString{},
		),
		Status: sql.NullString{String: input.Status, Valid: true},
	})
	if err != nil {
		utils.HandleError(ctx, err, "CountTalkSessions")
		return nil, err
	}

	talkSessionOut.TalkSessions = talkSessionDTOList
	talkSessionOut.TotalCount = int(talkSessionCount.TalkSessionCount)
	return &talkSessionOut, nil
}
