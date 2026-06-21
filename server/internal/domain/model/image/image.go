package image

import (
	"context"
	"mime/multipart"

	"github.com/neko-dream/api/internal/domain/model/image/meta"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"go.opentelemetry.io/otel"
)

//go:generate go tool mockgen -source=$GOFILE -package=mock_${GOPACKAGE}_model -destination=../mock/$GOPACKAGE/$GOFILE
type (
	ImageStorage interface {
		Upload(context.Context, meta.ImageMeta, *multipart.FileHeader) (*string, error)
	}
	ImageRepository interface {
		Create(context.Context, *UserImage) error
		FindByID(context.Context, shared.UUID[UserImage]) (*UserImage, error)
		FindByUserID(context.Context, shared.UUID[user.User]) ([]*UserImage, error)
	}

	UserImage struct {
		UserImageID shared.UUID[UserImage]
		UserID      shared.UUID[user.User]
		Metadata    meta.ImageMeta
		URL         string
	}
)

func NewUserImage(
	ctx context.Context,
	userImageID shared.UUID[UserImage],
	userID shared.UUID[user.User],
	metadata meta.ImageMeta,
	url string,
) *UserImage {
	ctx, span := otel.Tracer("image").Start(ctx, "NewUserImage")
	defer span.End()

	_ = ctx

	return &UserImage{
		UserImageID: userImageID,
		UserID:      userID,
		Metadata:    metadata,
		URL:         url,
	}
}
