package repository

import (
	"context"
	"log"
	"mime/multipart"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/image"
	"github.com/neko-dream/server/internal/domain/model/image/meta"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.opentelemetry.io/otel"
)

type imageStorage struct {
	s3Client *s3.Client
	conf     *config.Config
}

func NewImageStorage(
	s3Client *s3.Client,
	conf *config.Config,
) image.ImageStorage {
	return &imageStorage{
		s3Client: s3Client,
		conf:     conf,
	}
}

func (i *imageStorage) Upload(ctx context.Context, meta meta.ImageMeta, file *multipart.FileHeader) (*string, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "imageRepository.Upload")
	defer span.End()

	reader, err := file.Open()
	if err != nil {
		utils.HandleError(ctx, err, "file.Open")
		return nil, errtrace.Wrap(err)
	}
	defer reader.Close()

	uploader := manager.NewUploader(i.s3Client, func(u *manager.Uploader) {
		u.BufferProvider = manager.NewBufferedReadSeekerWriteToPool(25 * 1024 * 1024)
	})
	out, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(i.conf.AWS_S3_BUCKET),
		Key:         aws.String(meta.Key),
		Body:        reader,
		ContentType: aws.String(meta.Extension.Value),
	})
	if err != nil {
		utils.HandleError(ctx, err, "画像のアップロードに失敗")
		return nil, errtrace.Wrap(err)
	}
	log.Println("アップロード成功", out)

	url := i.conf.IMAGE_DOMAIN + "/" + meta.Key
	return lo.ToPtr(url), nil
}
