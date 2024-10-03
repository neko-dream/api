package user

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/neko-dream/server/internal/domain/model/image"
)

type ProfileIcon struct {
	image image.ImageInfo
	url   string
}

func NewProfileIcon(
	url string,
) *ProfileIcon {
	return &ProfileIcon{
		url: url,
	}
}

var (
	maxImageSize4MiB = 4194304
	objectPath       = "users/%s/profile_icon/%v.jpg"
)

func (p *ProfileIcon) SetProfileIconImage(
	ctx context.Context,
	file *multipart.FileHeader,
	user User,
) error {
	bytes, err := image.ValidateImage(ctx, file, maxImageSize4MiB)
	if err != nil {
		return err
	}

	if err := image.CheckSize(bytes, 300, 300); err != nil {
		return err
	}

	img := image.NewImage(bytes)
	imageInfo := image.NewImageInfo(fmt.Sprintf(objectPath, user.UserID().String(), time.Now().Unix()), img)
	p.image = *imageInfo
	return nil
}
