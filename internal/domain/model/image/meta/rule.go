package meta

import (
	"bytes"
	"context"
	"image"
	"io"
	"strings"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/samber/lo"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"go.opentelemetry.io/otel"
)

type ImageValidationRule struct {
	Domain         imageDomain
	maxFileSize    int
	allowedFormats []string
	minAspectRatio *float64
	maxAspectRatio *float64
	width          *int
	height         *int
}

type imageDomain string

var (
	ProfileImageValidationRule = ImageValidationRule{
		maxFileSize:    4194304,
		allowedFormats: []string{"jpeg", "png"},
		width:          lo.ToPtr(300),
		height:         lo.ToPtr(300),
	}
	ReferenceImageValidationRule = ImageValidationRule{
		maxFileSize:    4194304,
		allowedFormats: []string{"jpeg", "png"},
	}
	TalkSessionImageValidationRule = ImageValidationRule{
		maxFileSize:    4194304,
		allowedFormats: []string{"jpeg", "png"},
	}
	NoValidationRule = ImageValidationRule{
		maxFileSize:    4194304,
		allowedFormats: []string{"jpeg", "png", "gif"},
	}
)

func (r *ImageValidationRule) ValidExtension(ctx context.Context, meta ImageMeta) bool {
	ctx, span := otel.Tracer("meta").Start(ctx, "ImageValidationRule.ValidExtension")
	defer span.End()

	_ = ctx

	if len(r.allowedFormats) == 0 {
		return true
	}

	for _, format := range r.allowedFormats {
		if strings.Contains(meta.Extension.Value, format) {
			return true
		}
	}

	return false
}

func (r *ImageValidationRule) ValidBounds(ctx context.Context, meta ImageMeta) bool {
	ctx, span := otel.Tracer("meta").Start(ctx, "ImageValidationRule.ValidBounds")
	defer span.End()

	_ = ctx

	if r.width == nil && r.height == nil {
		return true
	}

	if r.width != nil && meta.Width > *r.width {
		return false
	}
	if r.height != nil && meta.Height > *r.height {
		return false
	}

	return true
}

func (r *ImageValidationRule) ValidAspectRatio(ctx context.Context, meta ImageMeta) bool {
	ctx, span := otel.Tracer("meta").Start(ctx, "ImageValidationRule.ValidAspectRatio")
	defer span.End()

	_ = ctx

	if r.minAspectRatio == nil && r.maxAspectRatio == nil {
		return true
	}

	ratio := float64(meta.Width) / float64(meta.Height)
	if r.minAspectRatio != nil && r.maxAspectRatio != nil {
		if ratio < *r.minAspectRatio || ratio > *r.maxAspectRatio {
			return false
		}
	} else if r.minAspectRatio != nil {
		if ratio < *r.minAspectRatio {
			return false
		}
	} else if r.maxAspectRatio != nil {
		if ratio > *r.maxAspectRatio {
			return false
		}
	}

	return true
}

func (r *ImageValidationRule) ValidFileSize(ctx context.Context, meta ImageMeta) bool {
	ctx, span := otel.Tracer("meta").Start(ctx, "ImageValidationRule.ValidFileSize")
	defer span.End()

	_ = ctx

	if r.maxFileSize == 0 {
		return true
	}

	return meta.Size <= r.maxFileSize
}

func GetImageSize(ctx context.Context, file io.Reader) (int, error) {
	ctx, span := otel.Tracer("meta").Start(ctx, "GetImageSize")
	defer span.End()

	_ = ctx

	imageByte := new(bytes.Buffer)
	if _, err := io.Copy(imageByte, file); err != nil {
		return 0, err
	}

	return imageByte.Len(), nil
}

func GetExtension(ctx context.Context, file io.Reader) (types.MIME, error) {
	ctx, span := otel.Tracer("meta").Start(ctx, "GetExtension")
	defer span.End()

	_ = ctx

	header := make([]byte, 261)
	if _, err := file.Read(header); err != nil {
		return types.Unknown.MIME, err
	}
	t, _ := filetype.Match(header)

	return t.MIME, nil
}

func GetBounds(ctx context.Context, file io.Reader) (int, int, error) {
	ctx, span := otel.Tracer("meta").Start(ctx, "GetBounds")
	defer span.End()

	_ = ctx
	img, _, err := image.Decode(file)
	if err != nil {
		return 0, 0, messages.ImageDecodeFailedError
	}

	return img.Bounds().Dx(), img.Bounds().Dy(), nil
}
