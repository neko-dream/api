package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/db"
	model "github.com/neko-dream/server/internal/infrastructure/db/sqlc"
	"github.com/neko-dream/server/pkg/utils"
)

type opinionRepository struct {
	*db.DBManager
}

func NewOpinionRepository(dbManager *db.DBManager) opinion.OpinionRepository {
	return &opinionRepository{
		DBManager: dbManager,
	}
}

// Create Opinion作成
func (o *opinionRepository) Create(ctx context.Context, op opinion.Opinion) error {
	if err := o.GetQueries(ctx).CreateOpinion(ctx, model.CreateOpinionParams{
		TalkSessionID: op.TalkSessionID().UUID(),
		UserID:        op.UserID().UUID(),
		ParentOpinionID: utils.IfThenElse[uuid.NullUUID](
			op.ParentOpinionID() == nil,
			uuid.NullUUID{},
			uuid.NullUUID{UUID: op.ParentOpinionID().UUID(), Valid: true},
		),
		Title: utils.IfThenElse(
			op.Title() == nil,
			sql.NullString{},
			sql.NullString{String: *op.Title(), Valid: true},
		),
		Content:   op.Content(),
		CreatedAt: op.CreatedAt(),
	}); err != nil {
		return err
	}
	return nil
}

// FindByParentID implements opinion.OpinionRepository.
func (o *opinionRepository) FindByParentID(ctx context.Context, opinionID shared.UUID[opinion.Opinion]) ([]opinion.Opinion, error) {
	panic("unimplemented")
}

// FindByTalkSessionID implements opinion.OpinionRepository.
func (o *opinionRepository) FindByTalkSessionID(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) ([]opinion.Opinion, error) {
	panic("unimplemented")
}

// FindByTalkSessionWithoutVote implements opinion.OpinionRepository.
func (o *opinionRepository) FindByTalkSessionWithoutVote(ctx context.Context, userID shared.UUID[user.User], talkSessionID shared.UUID[talksession.TalkSession], limit int) ([]opinion.Opinion, error) {
	panic("unimplemented")
}
