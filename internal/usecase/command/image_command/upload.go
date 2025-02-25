package image_command

import (
	"context"
	"mime/multipart"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/image"
	"github.com/neko-dream/server/internal/domain/model/image/meta"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type (
	UploadImage interface {
		Execute(context.Context, UploadImageInput) (*UploadImageOutput, error)
	}

	UploadImageInput struct {
		OwnerID shared.UUID[user.User]
		Image   *multipart.FileHeader
	}

	UploadImageOutput struct {
		ImageURL string
	}

	uploadImageHandler struct {
		imageRepository image.ImageRepository
		imageStorage    image.ImageStorage
	}
)

func NewUploadImageHandler(
	imageRepository image.ImageRepository,
	imageStorage image.ImageStorage,
) UploadImage {
	return &uploadImageHandler{
		imageRepository: imageRepository,
		imageStorage:    imageStorage,
	}
}

func (h *uploadImageHandler) Execute(ctx context.Context, input UploadImageInput) (*UploadImageOutput, error) {
	ctx, span := otel.Tracer("image_command").Start(ctx, "uploadImageHandler.Execute")
	defer span.End()

	file, err := input.Image.Open()
	if err != nil {
		utils.HandleError(ctx, err, "input.Icon.Open")
		return nil, messages.ImageOpenFailedError
	}
	defer file.Close()

	imageMeta, err := meta.NewImageForProfile(ctx, input.OwnerID, file)
	if err != nil {
		utils.HandleError(ctx, err, "meta.NewImageForProfile")
		return nil, messages.ImageOpenFailedError
	}

	if err := imageMeta.Validate(ctx, meta.ProfileImageValidationRule); err != nil {
		utils.HandleError(ctx, err, "ImageMeta.Validate")
		msg := messages.UserUpdateError
		msg.Message = err.Error()
		return nil, msg
	}

	// 画像をアップロード
	url, err := h.imageStorage.Upload(ctx, *imageMeta, input.Image)
	if err != nil {
		utils.HandleError(ctx, err, "ImageRepository.Upload")
		return nil, messages.ImageUploadFailedError
	}

	if err := h.imageRepository.Create(ctx, image.NewUserImage(
		ctx,
		shared.NewUUID[image.UserImage](),
		input.OwnerID,
		*imageMeta,
		*url,
	)); err != nil {
		utils.HandleError(ctx, err, "ImageRepository.Create")
		return nil, messages.ImageUploadFailedError
	}

	return &UploadImageOutput{
		ImageURL: *url,
	}, nil
}
