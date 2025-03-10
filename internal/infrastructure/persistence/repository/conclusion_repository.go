package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/neko-dream/server/internal/domain/model/conclusion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"go.opentelemetry.io/otel"
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
	ctx, span := otel.Tracer("repository").Start(ctx, "conclusionRepository.Create")
	defer span.End()

	if err := r.GetQueries(ctx).CreateTalkSessionConclusion(ctx, model.CreateTalkSessionConclusionParams{
		TalkSessionID: conclusion.TalkSessionID().UUID(),
		CreatedBy:     conclusion.CreatedBy().UUID(),
		Content:       conclusion.Conclusion(),
	}); err != nil {
		return err
	}
	return nil
}

func (r *conclusionRepository) FindByTalkSessionID(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) (*conclusion.Conclusion, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "conclusionRepository.FindByTalkSessionID")
	defer span.End()

	res, err := r.GetQueries(ctx).GetTalkSessionConclusionByID(ctx, talkSessionID.UUID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	tsID, err := shared.ParseUUID[talksession.TalkSession](res.TalkSessionID.String())
	if err != nil {
		return nil, err
	}
	uID, err := shared.ParseUUID[user.User](res.UserID.UUID.String())
	if err != nil {
		return nil, err
	}

	conc := conclusion.NewConclusion(
		tsID,
		res.Content,
		uID,
	)
	return conc, nil
}
