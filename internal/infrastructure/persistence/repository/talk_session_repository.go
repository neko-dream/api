package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/sqlc-dev/pqtype"
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

	restrictionRaw := pqtype.NullRawMessage{}
	if len(talkSession.Restrictions()) > 0 {
		restrictions, err := json.Marshal(talkSession.Restrictions())
		if err != nil {
			return err
		}
		restrictionRaw = pqtype.NullRawMessage{
			RawMessage: json.RawMessage(restrictions),
			Valid:      true,
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
		Restrictions:     restrictionRaw,
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
	var description, thumbnailURL, city, preference sql.NullString
	if talkSession.Description() != nil {
		description = sql.NullString{
			String: *talkSession.Description(),
			Valid:  true,
		}
	}
	if talkSession.ThumbnailURL() != nil {
		thumbnailURL = sql.NullString{
			String: *talkSession.ThumbnailURL(),
			Valid:  true,
		}
	}
	if talkSession.City() != nil {
		city = sql.NullString{
			String: *talkSession.City(),
			Valid:  true,
		}
	}
	if talkSession.Prefecture() != nil {
		preference = sql.NullString{
			String: *talkSession.Prefecture(),
			Valid:  true,
		}
	}

	restrictionRaw := pqtype.NullRawMessage{}
	if len(talkSession.Restrictions()) > 0 {
		restrictions, err := json.Marshal(talkSession.Restrictions())
		if err != nil {
			return err
		}
		restrictionRaw = pqtype.NullRawMessage{
			RawMessage: json.RawMessage(restrictions),
			Valid:      true,
		}
	}

	if err := t.GetQueries(ctx).EditTalkSession(ctx, model.EditTalkSessionParams{
		TalkSessionID:    talkSession.TalkSessionID().UUID(),
		Theme:            talkSession.Theme(),
		ScheduledEndTime: talkSession.ScheduledEndTime(),
		Description:      description,
		ThumbnailUrl:     thumbnailURL,
		City:             city,
		Prefecture:       preference,
		Restrictions:     restrictionRaw,
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
	var description, city, prefecture, thumbnailURL *string
	if row.TalkSession.City.Valid {
		city = &row.TalkSession.City.String
	}
	if row.TalkSession.Prefecture.Valid {
		prefecture = &row.TalkSession.Prefecture.String
	}
	if row.TalkSession.Description.Valid {
		description = &row.TalkSession.Description.String
	}
	if row.TalkSession.ThumbnailUrl.Valid {
		thumbnailURL = &row.TalkSession.ThumbnailUrl.String
	}

	ts := talksession.NewTalkSession(
		talkSessionID,
		row.TalkSession.Theme,
		description,
		thumbnailURL,
		shared.UUID[user.User](row.User.UserID),
		row.TalkSession.CreatedAt,
		row.TalkSession.ScheduledEndTime,
		location,
		city,
		prefecture,
	)

	restrictionRaw := row.TalkSession.Restrictions
	if restrictionRaw.Valid {
		// stringの配列になっているはず
		var restrictions []string
		if err := json.Unmarshal([]byte(restrictionRaw.RawMessage), &restrictions); err != nil {
			return nil, errtrace.Wrap(err)
		}

		if err := ts.UpdateRestrictions(ctx, restrictions); err != nil {
			return nil, errtrace.Wrap(err)
		}
	}

	return ts, nil
}
