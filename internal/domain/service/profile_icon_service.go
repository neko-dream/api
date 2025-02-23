package service

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

type ProfileIconService interface {
	// UploadProfileIcon プロフィールアイコンをアップロード
	UploadProfileIcon(ctx context.Context, userID shared.UUID[user.User], input *multipart.FileHeader) (*string, error)
}

type profileIconService struct {
	imageStorage image.ImageStorage
	imageRep     image.ImageRepository
}

func NewProfileIconService(
	imageStorage image.ImageStorage,
	imageRep image.ImageRepository,
) ProfileIconService {
	return &profileIconService{
		imageStorage: imageStorage,
		imageRep:     imageRep,
	}
}

func (i *profileIconService) UploadProfileIcon(ctx context.Context, userID shared.UUID[user.User], input *multipart.FileHeader) (*string, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "profileIconService.UploadProfileIcon")
	defer span.End()

	file, err := input.Open()
	if err != nil {
		utils.HandleError(ctx, err, "input.Icon.Open")
		return nil, messages.UserUpdateError
	}
	defer file.Close()

	imageMeta, err := meta.NewImageForProfile(ctx, userID, file)
	if err != nil {
		utils.HandleError(ctx, err, "meta.NewImageForProfile")
		return nil, messages.UserUpdateError
	}

	if err := imageMeta.Validate(ctx, meta.ProfileImageValidationRule); err != nil {
		utils.HandleError(ctx, err, "ImageMeta.Validate")
		msg := messages.UserUpdateError
		msg.Message = err.Error()
		return nil, msg
	}

	// 画像をアップロード
	url, err := i.imageStorage.Upload(ctx, *imageMeta, input)
	if err != nil {
		utils.HandleError(ctx, err, "ImageRepository.Upload")
		return nil, messages.UserUpdateError
	}

	if err := i.imageRep.Create(ctx, image.NewUserImage(
		ctx,
		shared.NewUUID[image.UserImage](),
		userID,
		*imageMeta,
		*url,
	)); err != nil {
		utils.HandleError(ctx, err, "ImageRepository.Create")
		return nil, messages.UserUpdateError
	}

	return url, nil
}
