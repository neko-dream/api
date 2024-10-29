package repository

import (
	"context"
	"database/sql"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/shared/time"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/db"
	model "github.com/neko-dream/server/internal/infrastructure/db/sqlc"
)

type talkSessionRepository struct {
	*db.DBManager
}

func NewTalkSessionRepository(
	DBManager *db.DBManager,
) talksession.TalkSessionRepository {
	return &talkSessionRepository{
		DBManager: DBManager,
	}
}

func (t *talkSessionRepository) Create(ctx context.Context, talkSession *talksession.TalkSession) error {
	var city, prefecture sql.NullString
	if talkSession.City() != nil {
		city = sql.NullString{
			String: *talkSession.City(),
			Valid:  true,
		}
	}
	if talkSession.Prefecture() != nil {
		prefecture = sql.NullString{
			String: *talkSession.Prefecture(),
			Valid:  true,
		}
	}

	if err := t.GetQueries(ctx).CreateTalkSession(ctx, model.CreateTalkSessionParams{
		TalkSessionID:    talkSession.TalkSessionID().UUID(),
		Theme:            talkSession.Theme(),
		OwnerID:          talkSession.OwnerUserID().UUID(),
		CreatedAt:        talkSession.CreatedAt().Time,
		ScheduledEndTime: talkSession.ScheduledEndTime().Time,
		Prefecture:       prefecture,
		City:             city,
	}); err != nil {
		return errtrace.Wrap(err)
	}
	// 位置情報がある場合は登録
	if talkSession.Location() != nil {
		if err := t.GetQueries(ctx).CreateTalkSessionLocation(ctx, model.CreateTalkSessionLocationParams{
			TalkSessionID:       talkSession.TalkSessionID().UUID(),
			StGeographyfromtext: talkSession.Location().ToGeographyText(),
		}); err != nil {
			return errtrace.Wrap(err)
		}
	}

	return nil
}

// Update implements talksession.TalkSessionRepository.
func (t *talkSessionRepository) Update(ctx context.Context, talkSession *talksession.TalkSession) error {
	if talkSession == nil {
		return nil
	}

	if err := t.GetQueries(ctx).EditTalkSession(ctx, model.EditTalkSessionParams{
		TalkSessionID:    talkSession.TalkSessionID().UUID(),
		Theme:            talkSession.Theme(),
		ScheduledEndTime: talkSession.ScheduledEndTime().Time,
	}); err != nil {
		return errtrace.Wrap(err)
	}

	if talkSession.Location() != nil {
		if err := t.GetQueries(ctx).UpdateTalkSessionLocation(ctx, model.UpdateTalkSessionLocationParams{
			TalkSessionID:       talkSession.TalkSessionID().UUID(),
			StGeographyfromtext: talkSession.Location().ToGeographyText(),
		}); err != nil {
			return errtrace.Wrap(err)
		}
	}

	return nil
}

// FindByID implements talksession.TalkSessionRepository.
func (t *talkSessionRepository) FindByID(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) (*talksession.TalkSession, error) {
	row, err := t.GetQueries(ctx).GetTalkSessionByID(ctx, talkSessionID.UUID())
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	var location *talksession.Location
	if row.LocationID.Valid {
		location = talksession.NewLocation(
			talkSessionID,
			row.Latitude,
			row.Longitude,
		)
	}
	var city, prefecture *string
	if row.City.Valid {
		city = &row.City.String
	}
	if row.Prefecture.Valid {
		prefecture = &row.Prefecture.String
	}

	return talksession.NewTalkSession(
		talkSessionID,
		row.Theme,
		shared.UUID[user.User](row.UserID.UUID),
		time.NewTime(ctx, row.CreatedAt),
		time.NewTime(ctx, row.ScheduledEndTime),
		location,
		city,
		prefecture,
	), nil

}
