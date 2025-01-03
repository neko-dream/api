package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/neko-dream/server/internal/domain/model/image"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/pkg/utils"
)

type opinionRepository struct {
	*db.DBManager
	image.ImageRepository
}

func NewOpinionRepository(
	dbManager *db.DBManager,
	imageRepo image.ImageRepository,
) opinion.OpinionRepository {
	return &opinionRepository{
		DBManager:       dbManager,
		ImageRepository: imageRepo,
	}
}

// Create Opinion作成
func (o *opinionRepository) Create(ctx context.Context, op opinion.Opinion) error {
	var parentOpinionID uuid.NullUUID
	if op.ParentOpinionID() != nil {
		parentOpinionID = uuid.NullUUID{UUID: op.ParentOpinionID().UUID(), Valid: true}
	}
	var referenceImageURL sql.NullString
	if op.IsReferenceImageUpdateRequired() {
		res, err := o.ImageRepository.Create(ctx, *op.ReferenceImage().ImageInfo())
		if err != nil {
			utils.HandleError(ctx, err, "opinionRepository.Create")
			return err
		}
		referenceImageURL = sql.NullString{String: *res, Valid: true}
	}

	var title sql.NullString
	if op.Title() != nil {
		title = sql.NullString{String: *op.Title(), Valid: true}
	}
	var referenceURL sql.NullString
	if op.ReferenceURL() != nil {
		referenceURL = sql.NullString{String: *op.ReferenceURL(), Valid: true}
	}

	if err := o.DBManager.GetQueries(ctx).CreateOpinion(ctx, model.CreateOpinionParams{
		OpinionID:       op.OpinionID().UUID(),
		TalkSessionID:   op.TalkSessionID().UUID(),
		UserID:          op.UserID().UUID(),
		ParentOpinionID: parentOpinionID,
		Title:           title,
		Content:         op.Content(),
		ReferenceUrl:    referenceURL,
		PictureUrl:      referenceImageURL,
		CreatedAt:       op.CreatedAt(),
	}); err != nil {
		utils.HandleError(ctx, err, "opinionRepository.Create")
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
