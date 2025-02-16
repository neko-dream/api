package repository

import (
	"context"
	"database/sql"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"go.opentelemetry.io/otel"
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
	ctx, span := otel.Tracer("repository").Start(ctx, "talkSessionRepository.Create")
	defer span.End()

	var description, city, prefecture sql.NullString
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
	if talkSession.Description() != nil {
		description = sql.NullString{
			String: *talkSession.Description(),
			Valid:  true,
		}
	}

	if err := t.GetQueries(ctx).CreateTalkSession(ctx, model.CreateTalkSessionParams{
		TalkSessionID:    talkSession.TalkSessionID().UUID(),
		Theme:            talkSession.Theme(),
		Description:      description,
		OwnerID:          talkSession.OwnerUserID().UUID(),
		CreatedAt:        talkSession.CreatedAt(),
		ScheduledEndTime: talkSession.ScheduledEndTime(),
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
	ctx, span := otel.Tracer("repository").Start(ctx, "talkSessionRepository.Update")
	defer span.End()

	if talkSession == nil {
		return nil
	}
	var description sql.NullString
	if talkSession.Description() != nil {
		description = sql.NullString{
			String: *talkSession.Description(),
			Valid:  true,
		}
	}

	if err := t.GetQueries(ctx).EditTalkSession(ctx, model.EditTalkSessionParams{
		TalkSessionID:    talkSession.TalkSessionID().UUID(),
		Theme:            talkSession.Theme(),
		ScheduledEndTime: talkSession.ScheduledEndTime(),
		Description:      description,
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
	ctx, span := otel.Tracer("repository").Start(ctx, "talkSessionRepository.FindByID")
	defer span.End()

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
	var description, city, prefecture *string
	if row.TalkSession.City.Valid {
		city = &row.TalkSession.City.String
	}
	if row.TalkSession.Prefecture.Valid {
		prefecture = &row.TalkSession.Prefecture.String
	}
	if row.TalkSession.Description.Valid {
		description = &row.TalkSession.Description.String
	}

	return talksession.NewTalkSession(
		talkSessionID,
		row.TalkSession.Theme,
		description,
		shared.UUID[user.User](row.User.UserID),
		row.TalkSession.CreatedAt,
		row.TalkSession.ScheduledEndTime,
		location,
		city,
		prefecture,
	), nil

}
