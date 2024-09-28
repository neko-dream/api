package repository

import (
	"context"
	"time"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/infrastructure/db"
	model "github.com/neko-dream/server/internal/infrastructure/db/sqlc"
)

type talkSessionRepository struct {
	*db.DBManager
}

func (t *talkSessionRepository) Create(ctx context.Context, talkSession *talksession.TalkSession) error {

	if err := t.GetQueries(ctx).CreateTalkSession(ctx, model.CreateTalkSessionParams{
		TalkSessionID: talkSession.TalkSessionID().UUID(),
		Theme:         talkSession.Theme(),
		CreatedAt:     time.Now(),
		OwnerID:       talkSession.OwnerUserID().UUID(),
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
	panic("unimplemented")
}

func NewTalkSessionRepository(
	DBManager *db.DBManager,
) talksession.TalkSessionRepository {
	return &talkSessionRepository{
		DBManager: DBManager,
	}
}
