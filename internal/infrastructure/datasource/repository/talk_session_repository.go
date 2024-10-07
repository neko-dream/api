package repository

import (
	"context"
	"database/sql"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
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
	if err := t.GetQueries(ctx).CreateTalkSession(ctx, model.CreateTalkSessionParams{
		TalkSessionID:    talkSession.TalkSessionID().UUID(),
		Theme:            talkSession.Theme(),
		OwnerID:          talkSession.OwnerUserID().UUID(),
		CreatedAt:        talkSession.CreatedAt().Time,
		ScheduledEndTime: talkSession.ScheduledEndTime().Time,
	}); err != nil {
		return errtrace.Wrap(err)
	}

	return nil
}

// FindByID implements talksession.TalkSessionRepository.
func (t *talkSessionRepository) FindByID(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) (*talksession.TalkSession, error) {
	panic("unimplemented")
}

// Update implements talksession.TalkSessionRepository.
func (t *talkSessionRepository) Update(ctx context.Context, talkSession *talksession.TalkSession) error {
	if talkSession == nil {
		return nil
	}

	var finishedAt sql.NullTime
	if talkSession.FinishedAt() != nil {
		finishedAt.Time = talkSession.FinishedAt().Time
		finishedAt.Valid = true
	}

	if err := t.GetQueries(ctx).EditTalkSession(ctx, model.EditTalkSessionParams{
		TalkSessionID:    talkSession.TalkSessionID().UUID(),
		Theme:            talkSession.Theme(),
		FinishedAt:       finishedAt,
		ScheduledEndTime: talkSession.ScheduledEndTime().Time,
	}); err != nil {
		return errtrace.Wrap(err)
	}

	return nil
}
