package image

import (
	"context"

	"io"

	"github.com/neko-dream/server/internal/domain/model/image/meta"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"go.opentelemetry.io/otel"
)

type (
	ImageStorage interface {
		Upload(context.Context, meta.ImageMeta, io.Reader) (*string, error)
	}
	ImageRepository interface {
		Save(context.Context, *UserImage) error
		FindByID(context.Context, shared.UUID[UserImage]) (*UserImage, error)
		FindByUserID(context.Context, shared.UUID[user.User]) ([]*UserImage, error)
	}

	UserImage struct {
		UserImageID shared.UUID[UserImage]
		UserID      shared.UUID[user.User]
		Metadata    meta.ImageMeta
	}
)

func NewUserImage(
	ctx context.Context,
	userImageID shared.UUID[UserImage],
	userID shared.UUID[user.User],
	metadata meta.ImageMeta,
) *UserImage {
	ctx, span := otel.Tracer("image").Start(ctx, "NewUserImage")
	defer span.End()

	_ = ctx

	return &UserImage{
		UserImageID: userImageID,
		UserID:      userID,
		Metadata:    metadata,
	}
}
