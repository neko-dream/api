package talk_session_usecase

import (
	"context"
	"time"

	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/utils"
)

type (
	GetTalkSessionDetailUseCase interface {
		Execute(context.Context, GetTalkSessionDetailInput) (*GetTalkSessionDetailOutput, error)
	}

	GetTalkSessionDetailInput struct {
		TalkSessionID shared.UUID[talksession.TalkSession]
	}

	GetTalkSessionDetailOutput struct {
		TalkSessionDTO
	}

	getTalkSessionDetailInteractor struct {
		*db.DBManager
	}
)

func NewGetTalkSessionDetailUseCase(
	dbManager *db.DBManager,
) GetTalkSessionDetailUseCase {
	return &getTalkSessionDetailInteractor{
		DBManager: dbManager,
	}
}

func (i *getTalkSessionDetailInteractor) Execute(ctx context.Context, input GetTalkSessionDetailInput) (*GetTalkSessionDetailOutput, error) {

	talkSessionRow, err := i.DBManager.GetQueries(ctx).GetTalkSessionByID(ctx, input.TalkSessionID.UUID())
	if err != nil {
		return nil, err
	}

	owner := UserDTO{
		DisplayID:   talkSessionRow.DisplayID.String,
		DisplayName: talkSessionRow.DisplayName.String,
		IconURL:     utils.ToPtrIfNotNullValue(!talkSessionRow.IconUrl.Valid, talkSessionRow.IconUrl.String),
	}
	var location *LocationDTO
	if talkSessionRow.LocationID.Valid {
		location = &LocationDTO{
			Latitude:  talkSessionRow.Latitude,
			Longitude: talkSessionRow.Longitude,
		}
	}

	talkSessionRes := TalkSessionDTO{
		ID:               talkSessionRow.TalkSessionID.String(),
		Theme:            talkSessionRow.Theme,
		Owner:            owner,
		Description:      utils.ToPtrIfNotNullValue(!talkSessionRow.Description.Valid, talkSessionRow.Description.String),
		ScheduledEndTime: talkSessionRow.ScheduledEndTime.Format(time.RFC3339),
		CreatedAt:        talkSessionRow.CreatedAt.Format(time.RFC3339),
		Location:         location,
		City:             utils.ToPtrIfNotNullValue(!talkSessionRow.City.Valid, talkSessionRow.City.String),
		Prefecture:       utils.ToPtrIfNotNullValue(!talkSessionRow.Prefecture.Valid, talkSessionRow.Prefecture.String),
	}

	out := &GetTalkSessionDetailOutput{
		TalkSessionDTO: talkSessionRes,
	}
	return out, nil
}
