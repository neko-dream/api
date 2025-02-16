package repository

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/image"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"go.opentelemetry.io/otel"
)

type imageRepository struct {
	*db.DBManager
}

// FindByID implements image.ImageRepository.
func (i *imageRepository) FindByID(ctx context.Context, _ shared.UUID[image.UserImage]) (*image.UserImage, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "imageRepository.FindByID")
	defer span.End()

	_ = ctx

	panic("unimplemented")
}

// FindByUserID implements image.ImageRepository.
func (i *imageRepository) FindByUserID(ctx context.Context, _ shared.UUID[user.User]) ([]*image.UserImage, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "imageRepository.FindByUserID")
	defer span.End()

	_ = ctx

	panic("unimplemented")
}

// Save implements image.ImageRepository.
func (i *imageRepository) Save(ctx context.Context, _ *image.UserImage) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "imageRepository.Save")
	defer span.End()

	_ = ctx

	panic("unimplemented")
}

func NewImageRepository(DBManager *db.DBManager) image.ImageRepository {
	return &imageRepository{
		DBManager: DBManager,
	}
}
