package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/neko-dream/server/internal/domain/model/conclusion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
)

type conclusionRepository struct {
	*db.DBManager
}

func NewConclusionRepository(
	dbManager *db.DBManager,
) conclusion.ConclusionRepository {
	return &conclusionRepository{
		DBManager: dbManager,
	}
}

func (r *conclusionRepository) Create(ctx context.Context, conclusion conclusion.Conclusion) error {
	r.GetQueries(ctx).CreateTalkSessionConclusion(ctx, model.CreateTalkSessionConclusionParams{
		TalkSessionID: conclusion.TalkSessionID().UUID(),
		CreatedBy:     conclusion.CreatedBy().UUID(),
		Content:       conclusion.Conclusion(),
	})
	return nil
}

func (r *conclusionRepository) FindByTalkSessionID(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) (*conclusion.Conclusion, error) {
	res, err := r.GetQueries(ctx).GetTalkSessionConclusionByID(ctx, talkSessionID.UUID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	conc := conclusion.NewConclusion(
		shared.MustParseUUID[talksession.TalkSession](res.TalkSessionID.String()),
		res.Content,
		shared.MustParseUUID[user.User](res.UserID.UUID.String()),
	)
	return conc, nil
}
