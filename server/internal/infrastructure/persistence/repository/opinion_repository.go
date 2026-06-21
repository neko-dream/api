package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/neko-dream/api/internal/domain/model/image"
	"github.com/neko-dream/api/internal/domain/model/opinion"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/api/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type opinionRepository struct {
	*db.DBManager
	image.ImageStorage
}

func NewOpinionRepository(
	dbManager *db.DBManager,
	imageRepo image.ImageStorage,
) opinion.OpinionRepository {
	return &opinionRepository{
		DBManager:    dbManager,
		ImageStorage: imageRepo,
	}
}

// Create Opinion作成
func (o *opinionRepository) Create(ctx context.Context, op opinion.Opinion) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "opinionRepository.Create")
	defer span.End()

	var parentOpinionID uuid.NullUUID
	if op.ParentOpinionID() != nil {
		parentOpinionID = uuid.NullUUID{UUID: op.ParentOpinionID().UUID(), Valid: true}
	}
	var referenceImageURL sql.NullString
	if op.ReferenceImageURL() != nil {
		referenceImageURL = sql.NullString{String: *op.ReferenceImageURL(), Valid: true}
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

func (o *opinionRepository) FindByID(ctx context.Context, opinionID shared.UUID[opinion.Opinion]) (*opinion.Opinion, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "opinionRepository.FindByID")
	defer span.End()

	op, err := o.GetQueries(ctx).GetOpinionByID(ctx, model.GetOpinionByIDParams{
		OpinionID: opinionID.UUID(),
	})
	if err != nil {
		utils.HandleError(ctx, err, "opinionRepository.FindByID")
		return nil, err
	}
	var parentOpinionID *shared.UUID[opinion.Opinion]
	if op.Opinion.ParentOpinionID.Valid {
		parentOpinionID = lo.ToPtr(shared.UUID[opinion.Opinion](op.Opinion.ParentOpinionID.UUID))
	}
	var title *string
	if op.Opinion.Title.Valid {
		title = lo.ToPtr(op.Opinion.Title.String)
	}

	opEntity, err := opinion.NewOpinion(
		opinionID,
		shared.UUID[talksession.TalkSession](op.Opinion.TalkSessionID),
		shared.UUID[user.User](op.Opinion.UserID),
		parentOpinionID,
		title,
		op.Opinion.Content,
		op.Opinion.CreatedAt,
		lo.ToPtr(op.Opinion.ReferenceUrl.String),
	)
	if err != nil {
		utils.HandleError(ctx, err, "opinionRepository.FindByID")
		return nil, err
	}

	return opEntity, nil
}

// FindByParentID implements opinion.OpinionRepository.
func (o *opinionRepository) FindByParentID(ctx context.Context, opinionID shared.UUID[opinion.Opinion]) ([]opinion.Opinion, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "opinionRepository.FindByParentID")
	defer span.End()

	_ = ctx

	panic("unimplemented")
}

// FindByTalkSessionID implements opinion.OpinionRepository.
func (o *opinionRepository) FindByTalkSessionID(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) ([]opinion.Opinion, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "opinionRepository.FindByTalkSessionID")
	defer span.End()

	_ = ctx

	panic("unimplemented")
}

// FindByTalkSessionWithoutVote implements opinion.OpinionRepository.
func (o *opinionRepository) FindByTalkSessionWithoutVote(ctx context.Context, userID shared.UUID[user.User], talkSessionID shared.UUID[talksession.TalkSession], limit int) ([]opinion.Opinion, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "opinionRepository.FindByTalkSessionWithoutVote")
	defer span.End()

	_ = ctx

	panic("unimplemented")
}
