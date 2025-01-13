package repository

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/image"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type ImageRepositoryMock struct {
}

// Create implements image.ImageRepository.
func (i *ImageRepositoryMock) Create(ctx context.Context, _ image.ImageInfo) (*string, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "ImageRepositoryMock.Create")
	defer span.End()

	_ = ctx

	return lo.ToPtr("https://image.kotohiro.com/hogehoge"), nil
}

func NewImageRepositoryMock() image.ImageRepository {
	return &ImageRepositoryMock{}
}
