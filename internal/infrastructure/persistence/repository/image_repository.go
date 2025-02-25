package repository

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/image"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/pkg/utils"
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
func (i *imageRepository) Create(ctx context.Context, img *image.UserImage) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "imageRepository.Save")
	defer span.End()

	_ = ctx

	if err := i.DBManager.GetQueries(ctx).CreateUserImage(
		ctx,
		model.CreateUserImageParams{
			UserImagesID: img.UserImageID.UUID(),
			UserID:       img.UserID.UUID(),
			Key:          img.Metadata.Key,
			Width:        int32(img.Metadata.Width),
			Height:       int32(img.Metadata.Height),
			Extension:    img.Metadata.Extension.Value,
			Archived:     img.Metadata.Archived,
			Url:          img.URL,
		},
	); err != nil {
		utils.HandleError(ctx, err, "CreateUserImage")
		return err
	}

	return nil
}

func NewImageRepository(DBManager *db.DBManager) image.ImageRepository {
	return &imageRepository{
		DBManager: DBManager,
	}
}
