package user

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/image"
)

type ProfileIcon struct {
	image *image.ImageInfo
	url   *string
}

func NewProfileIcon(
	url *string,
) *ProfileIcon {
	return &ProfileIcon{
		url: url,
	}
}

func (p *ProfileIcon) ImageInfo() *image.ImageInfo {
	return p.image
}
func (p *ProfileIcon) URL() *string {
	return p.url
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
	if file == nil {
		return nil
	}
	bytes, ext, err := image.ValidateImage(ctx, file, maxImageSize4MiB)
	if err != nil {
		return err
	}

	if err := image.CheckSize(bytes, 300, 300); err != nil {
		return err
	}

	img := image.NewImage(bytes)
	imageInfo := image.NewImageInfo(fmt.Sprintf(objectPath, *user.DisplayID(), clock.Now(ctx).Unix()), *ext, img)
	p.image = imageInfo

	return nil
}
