package repository

import (
	"context"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/image"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.opentelemetry.io/otel"
)

type imageRepository struct {
	s3Client *s3.Client
	conf     *config.Config
}

func NewImageRepository(
	s3Client *s3.Client,
	conf *config.Config,
) image.ImageRepository {
	return &imageRepository{
		s3Client: s3Client,
		conf:     conf,
	}
}

func (i *imageRepository) Create(ctx context.Context, imageInfo image.ImageInfo) (*string, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "imageRepository.Create")
	defer span.End()

	if imageInfo.FilePath() == "" {
		return nil, errtrace.Wrap(messages.ImageFilePathEmptyError)
	}

	uploader := manager.NewUploader(i.s3Client, func(u *manager.Uploader) {
		u.BufferProvider = manager.NewBufferedReadSeekerWriteToPool(25 * 1024 * 1024)
	})

	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(i.conf.AWS_S3_BUCKET),
		Key:         aws.String(imageInfo.FilePath()),
		Body:        imageInfo.ImageReader(),
		ContentType: aws.String(imageInfo.ContentType()),
	})
	if err != nil {
		utils.HandleError(ctx, err, "画像のアップロードに失敗")
		return nil, errtrace.Wrap(err)
	}

	url := i.conf.IMAGE_DOMAIN + "/" + imageInfo.FilePath()
	return lo.ToPtr(url), nil
}
