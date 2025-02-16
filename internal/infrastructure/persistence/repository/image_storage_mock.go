package repository

import (
	"context"
	"io"

	"github.com/neko-dream/server/internal/domain/model/image"
	"github.com/neko-dream/server/internal/domain/model/image/meta"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type ImageRepositoryMock struct {
}

// Upload implements image.ImageRepository.
func (i *ImageRepositoryMock) Upload(ctx context.Context, _ meta.ImageMeta, _ io.Reader) (*string, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "ImageRepositoryMock.Upload")
	defer span.End()

	_ = ctx

	return lo.ToPtr("https://image.kotohiro.com/hogehoge"), nil
}

func NewImageRepositoryMock() image.ImageStorage {
	return &ImageRepositoryMock{}
}
