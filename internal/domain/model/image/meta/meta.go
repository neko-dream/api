package meta

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/h2non/filetype/types"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type ImageMeta struct {
	Key       string
	Size      int
	Extension types.MIME
	Width     int
	Height    int
	Archived  bool
}

var (
	ProfileImageKeyPattern = "users/%s.jpg"
	// ReferenceImageKeyPattern Year/Month/Day/opinionID.jpg
	ReferenceImageKeyPattern = "reference_images/%v/%v/%v/%v.jpg"
	// 種類-talkSessionID-時間.jpg
	AnalysisImageKeyPattern = "generated/%v-%v-%v.png"
)

// ProfileIcon用の画像メタデータを生成
func NewImageForProfile(
	ctx context.Context,
	userID shared.UUID[user.User],
	file io.Reader,
) (*ImageMeta, error) {
	ctx, span := otel.Tracer("meta").Start(ctx, "NewImageForProfile")
	defer span.End()

	key := fmt.Sprintf(ProfileImageKeyPattern, userID.String())
	// bytes[]に変換
	imageBytes, err := io.ReadAll(file)
	if err != nil {
		utils.HandleError(ctx, err, "io.ReadAll")
		return nil, err
	}

	// 複製したデータを再利用
	x, y, err := GetBounds(ctx, bytes.NewReader(imageBytes))
	if err != nil {
		utils.HandleError(ctx, err, "GetBounds")
		return nil, err
	}

	size, err := GetImageSize(ctx, bytes.NewReader(imageBytes))
	if err != nil {
		utils.HandleError(ctx, err, "GetImageSize")
		return nil, err
	}

	// 複製したデータを再利用
	ext, err := GetExtension(ctx, bytes.NewReader(imageBytes))
	if err != nil {
		utils.HandleError(ctx, err, "GetExtension")
		return nil, err
	}

	return &ImageMeta{
		Key:       key,
		Size:      size,
		Extension: ext,
		Width:     x,
		Height:    y,
	}, nil
}

// ReferenceImage用の画像メタデータを生成
func NewImageForReference(
	ctx context.Context,
	opinionID shared.UUID[opinion.Opinion],
	file io.Reader,
) (*ImageMeta, error) {
	ctx, span := otel.Tracer("meta").Start(ctx, "NewImageForReference")
	defer span.End()

	now := clock.Now(ctx)
	key := fmt.Sprintf(
		ReferenceImageKeyPattern,
		now.Year(),
		now.Month(),
		now.Day(),
		opinionID.String(),
	)
	// bytes[]に変換
	imageBytes, err := io.ReadAll(file)
	if err != nil {
		utils.HandleError(ctx, err, "io.ReadAll")
		return nil, err
	}

	size, err := GetImageSize(ctx, bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}
	ext, err := GetExtension(ctx, bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}

	x, y, err := GetBounds(ctx, bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}

	return &ImageMeta{
		Key:       key,
		Size:      size,
		Extension: ext,
		Width:     x,
		Height:    y,
	}, nil
}

func NewImageForAnalysis(
	ctx context.Context,
	file io.Reader,
) (*ImageMeta, error) {
	ctx, span := otel.Tracer("meta").Start(ctx, "NewImageForAnalysis")
	defer span.End()

	now := clock.Now(ctx)
	key := fmt.Sprintf(
		AnalysisImageKeyPattern,
		"analysis",
		now.Format("20060102150405"),
	)
	// bytes[]に変換
	imageBytes, err := io.ReadAll(file)
	if err != nil {
		utils.HandleError(ctx, err, "io.ReadAll")
		return nil, err
	}

	size, err := GetImageSize(ctx, bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}
	ext, err := GetExtension(ctx, bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}

	x, y, err := GetBounds(ctx, bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}

	return &ImageMeta{
		Key:       key,
		Size:      size,
		Extension: ext,
		Width:     x,
		Height:    y,
	}, nil
}

func (m *ImageMeta) Validate(ctx context.Context, rule ImageValidationRule) error {
	ctx, span := otel.Tracer("meta").Start(ctx, "ImageMeta.Validate")
	defer span.End()

	var err error

	if !rule.ValidExtension(ctx, *m) {
		err = errors.Join(err, errors.New("サポートされていないフォーマットです。"))
	}

	if !rule.ValidFileSize(ctx, *m) {
		err = errors.Join(err, errors.New("ファイルサイズが大きすぎます。"))
	}

	if !rule.ValidBounds(ctx, *m) {
		err = errors.Join(err, errors.New("画像のサイズが大きすぎます。"))
	}

	if !rule.ValidAspectRatio(ctx, *m) {
		err = errors.Join(err, errors.New("アスペクト比が不正です。"))
	}

	return err
}
