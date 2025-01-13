package opinion

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/image"
	"go.opentelemetry.io/otel"
)

type ReferenceImage struct {
	image *image.ImageInfo
	url   *string
}

func NewReferenceImage(
	url *string,
) *ReferenceImage {
	return &ReferenceImage{
		url: url,
	}
}

func (p *ReferenceImage) ImageInfo() *image.ImageInfo {
	return p.image
}
func (p *ReferenceImage) URL() *string {
	return p.url
}

var (
	maxImageSize4MiB = 4194304
	objectPath       = "reference_images/%v/%v/%v/%v.jpg"
)

func (p *ReferenceImage) SetReferenceImage(
	ctx context.Context,
	file *multipart.FileHeader,
) error {
	ctx, span := otel.Tracer("opinion").Start(ctx, "ReferenceImage.SetReferenceImage")
	defer span.End()

	if file == nil {
		return nil
	}
	bytes, ext, err := image.ValidateImage(ctx, file, maxImageSize4MiB)
	if err != nil {
		return err
	}

	now := clock.Now(ctx)
	img := image.NewImage(bytes)
	imageInfo := image.NewImageInfo(fmt.Sprintf(
		objectPath,
		now.Year(),
		now.Month(),
		now.Day(),
		now.UnixNano()), *ext, img)
	p.image = imageInfo

	return nil
}
