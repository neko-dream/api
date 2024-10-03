package image

import (
	"context"
	"io"
)

type (
	ImageRepository interface {
		Create(context.Context, ImageInfo) (*string, error)
	}

	ImageInfo struct {
		filePath    string
		contentType string
		image       Image
	}
)

func NewImageInfo(
	filePath string,
	contentType string,
	image Image,
) *ImageInfo {
	return &ImageInfo{
		filePath:    filePath,
		contentType: contentType,
		image:       image,
	}
}

func (i *ImageInfo) FilePath() string {
	return i.filePath
}

func (i *ImageInfo) ImageReader() io.Reader {
	return i.image.GetImageReader()
}

func (i *ImageInfo) ContentType() string {
	return i.contentType
}
