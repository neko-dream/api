package repository

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/image"
	"github.com/samber/lo"
)

type ImageRepositoryMock struct {
}

// Create implements image.ImageRepository.
func (i *ImageRepositoryMock) Create(context.Context, image.ImageInfo) (*string, error) {
	return lo.ToPtr("https://image.kotohiro.com/hogehoge"), nil
}

func NewImageRepositoryMock() image.ImageRepository {
	return &ImageRepositoryMock{}
}
