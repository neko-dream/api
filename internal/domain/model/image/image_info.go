package image

import "context"

type (
	ImageRepository interface {
		Create(context.Context, Image) error
	}

	ImageInfo struct {
		FileName string
		URL      string
		Image    Image
	}
)

func NewImageInfo(
	fileName string,
	image Image,
) *ImageInfo {
	return &ImageInfo{
		FileName: fileName,
		Image:    image,
	}
}
