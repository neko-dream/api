package repository

import (
	"context"
	"database/sql"
	"log"

	"braces.dev/errtrace"
	"github.com/google/uuid"
	"github.com/neko-dream/server/internal/domain/model/event"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type talkSessionRepository struct {
	*db.DBManager
	eventStore event.EventStore
}

func NewTalkSessionRepository(
	DBManager *db.DBManager,
	eventStore event.EventStore,
) talksession.TalkSessionRepository {
	return &talkSessionRepository{
		DBManager:  DBManager,
		eventStore: eventStore,
	}
}

func (t *talkSessionRepository) Create(ctx context.Context, talkSession *talksession.TalkSession) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "talkSessionRepository.Create")
	defer span.End()

	var restrictions []string
	if len(talkSession.Restrictions()) > 0 {
		for _, restriction := range talkSession.Restrictions() {
			restrictions = append(restrictions, string(restriction.Key))
		}
	}

	if err := t.GetQueries(ctx).CreateTalkSession(ctx, model.CreateTalkSessionParams{
		TalkSessionID:       talkSession.TalkSessionID().UUID(),
		Theme:               talkSession.Theme(),
		Description:         utils.ToSQLNull[sql.NullString](talkSession.Description()),
		OwnerID:             talkSession.OwnerUserID().UUID(),
		CreatedAt:           talkSession.CreatedAt(),
		ScheduledEndTime:    talkSession.ScheduledEndTime(),
		ThumbnailUrl:        utils.ToSQLNull[sql.NullString](talkSession.ThumbnailURL()),
		Prefecture:          utils.ToSQLNull[sql.NullString](talkSession.Prefecture()),
		City:                utils.ToSQLNull[sql.NullString](talkSession.City()),
		Restrictions:        talksession.Restrictions(restrictions),
		HideReport:          utils.ToSQLNull[sql.NullBool](talkSession.HideReport()),
		OrganizationAliasID: utils.ToSQLNull[uuid.NullUUID](talkSession.OrganizationAliasID()),
		OrganizationID:      utils.ToSQLNull[uuid.NullUUID](talkSession.OrganizationID()),
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

	// イベントがある場合は保存
	events := talkSession.GetRecordedEvents()
	if len(events) > 0 {
		if err := t.eventStore.StoreBatch(ctx, events); err != nil {
			return errtrace.Wrap(err)
		}
		talkSession.ClearRecordedEvents()
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

	var restrictions []string
	if len(talkSession.Restrictions()) > 0 {
		for _, restriction := range talkSession.Restrictions() {
			restrictions = append(restrictions, string(restriction.Key))
		}
	}

	if err := t.DBManager.GetQueries(ctx).EditTalkSession(ctx, model.EditTalkSessionParams{
		TalkSessionID:    talkSession.TalkSessionID().UUID(),
		Theme:            talkSession.Theme(),
		ScheduledEndTime: talkSession.ScheduledEndTime(),
		Description:      utils.ToSQLNull[sql.NullString](talkSession.Description()),
		ThumbnailUrl:     utils.ToSQLNull[sql.NullString](talkSession.ThumbnailURL()),
		City:             utils.ToSQLNull[sql.NullString](talkSession.City()),
		Prefecture:       utils.ToSQLNull[sql.NullString](talkSession.Prefecture()),
		Restrictions:     talksession.Restrictions(restrictions),
		HideReport:       utils.ToSQLNull[sql.NullBool](talkSession.HideReport()),
	}); err != nil {
		return errtrace.Wrap(err)
	}

	if talkSession.Location() != nil {
		if err := t.DBManager.GetQueries(ctx).UpdateTalkSessionLocation(ctx, model.UpdateTalkSessionLocationParams{
			TalkSessionID:       talkSession.TalkSessionID().UUID(),
			StGeographyfromtext: talkSession.Location().ToGeographyText(),
		}); err != nil {
			return errtrace.Wrap(err)
		}
	}

	// イベントがある場合は保存
	events := talkSession.GetRecordedEvents()
	if len(events) > 0 {
		if err := t.eventStore.StoreBatch(ctx, events); err != nil {
			return errtrace.Wrap(err)
		}
		talkSession.ClearRecordedEvents()
	}

	return nil
}

// FindByID implements talksession.TalkSessionRepository.
func (t *talkSessionRepository) FindByID(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) (*talksession.TalkSession, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "talkSessionRepository.FindByID")
	defer span.End()

	row, err := t.DBManager.GetQueries(ctx).GetTalkSessionByID(ctx, talkSessionID.UUID())
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
		nil,
		nil,
	)
	ts.SetReportVisibility(row.TalkSession.HideReport.Bool)

	if len(row.TalkSession.Restrictions) > 0 {
		if err := ts.UpdateRestrictions(ctx, row.TalkSession.Restrictions); err != nil {
			return nil, errtrace.Wrap(err)
		}
	}

	return ts, nil
}

func (t *talkSessionRepository) GetUnprocessedEndedSessions(ctx context.Context, limit int) ([]*talksession.TalkSession, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "talkSessionRepository.GetUnprocessedEndedSessions")
	defer span.End()

	rows, err := t.GetQueries(ctx).GetUnprocessedEndedSessions(ctx, int32(limit))
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	var sessions []*talksession.TalkSession
	for _, row := range rows {
		var description, thumbnailURL, city, prefecture *string
		if row.Description.Valid {
			description = &row.Description.String
		}
		if row.ThumbnailUrl.Valid {
			thumbnailURL = &row.ThumbnailUrl.String
		}
		if row.City.Valid {
			city = &row.City.String
		}
		if row.Prefecture.Valid {
			prefecture = &row.Prefecture.String
		}

		var organizationID *shared.UUID[organization.Organization]
		if row.OrganizationID.Valid {
			orgID := shared.UUID[organization.Organization](row.OrganizationID.UUID)
			organizationID = &orgID
		}

		var organizationAliasID *shared.UUID[organization.OrganizationAlias]
		if row.OrganizationAliasID.Valid {
			aliasID := shared.UUID[organization.OrganizationAlias](row.OrganizationAliasID.UUID)
			organizationAliasID = &aliasID
		}

		session := talksession.NewTalkSession(
			shared.UUID[talksession.TalkSession](row.TalkSessionID),
			row.Theme,
			description,
			thumbnailURL,
			shared.UUID[user.User](row.OwnerID),
			row.CreatedAt,
			row.ScheduledEndTime,
			nil, // Location は別途取得が必要な場合
			city,
			prefecture,
			organizationID,
			organizationAliasID,
		)
		// 制限事項を設定
		if len(row.Restrictions) > 0 {
			if err := session.UpdateRestrictions(ctx, row.Restrictions); err != nil {
				log.Printf("failed to update restrictions: %v", err)
			}
		}
		// レポート表示設定
		session.SetReportVisibility(row.HideReport.Bool)
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (t *talkSessionRepository) GetParticipantIDs(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) ([]shared.UUID[user.User], error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "talkSessionRepository.GetParticipantIDs")
	defer span.End()

	userIDs, err := t.GetQueries(ctx).GetTalkSessionParticipants(ctx, talkSessionID.UUID())
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	var participantIDs []shared.UUID[user.User]
	for _, userID := range userIDs {
		participantIDs = append(participantIDs, shared.UUID[user.User](userID))
	}

	return participantIDs, nil
}
