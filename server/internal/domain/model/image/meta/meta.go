package meta

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/h2non/filetype/types"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/opinion"
	"github.com/neko-dream/api/internal/domain/model/organization"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/pkg/utils"
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
	ProfileImageKeyPattern = "u/%s.%v"
	// ReferenceImageKeyPattern Year/Month/Day/opinionID.jpg
	ReferenceImageKeyPattern = "ref/%v/%v/%v/%v.%v"
	// 種類-talkSessionID-時間.jpg
	AnalysisImageKeyPattern = "gen/%v-%v-%v.%v"
	// 画像ID.jpg
	CommonImageKeyPattern = "i/%v.%v"
	// 組織アイコン.jpg
	OrganizationImageKeyPattern = "o/%s.%v"
)

func NewOrganizationImage(
	ctx context.Context,
	organizationID shared.UUID[organization.Organization],
	file io.Reader,
) (*ImageMeta, error) {
	ctx, span := otel.Tracer("meta").Start(ctx, "NewOrganizationImage")
	defer span.End()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		utils.HandleError(ctx, err, "io.ReadAll")
		return nil, err
	}

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

	ext, err := GetExtension(ctx, bytes.NewReader(imageBytes))
	if err != nil {
		utils.HandleError(ctx, err, "GetExtension")
		return nil, err
	}

	key := fmt.Sprintf(OrganizationImageKeyPattern, organizationID.String(), ext.Subtype)

	return &ImageMeta{
		Key:       key,
		Size:      size,
		Extension: ext,
		Width:     x,
		Height:    y,
	}, nil
}

// ProfileIcon用の画像メタデータを生成
func NewImageForProfile(
	ctx context.Context,
	userID shared.UUID[user.User],
	file io.Reader,
) (*ImageMeta, error) {
	ctx, span := otel.Tracer("meta").Start(ctx, "NewImageForProfile")
	defer span.End()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		utils.HandleError(ctx, err, "io.ReadAll")
		return nil, err
	}

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

	ext, err := GetExtension(ctx, bytes.NewReader(imageBytes))
	if err != nil {
		utils.HandleError(ctx, err, "GetExtension")
		return nil, err
	}

	key := fmt.Sprintf(ProfileImageKeyPattern, userID.String(), ext.Subtype)

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

	key := fmt.Sprintf(
		ReferenceImageKeyPattern,
		now.Year(),
		int(now.Month()),
		now.Day(),
		opinionID.String(),
		ext.Subtype,
	)

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

	now := clock.Now(ctx)
	key := fmt.Sprintf(
		AnalysisImageKeyPattern,
		"analysis",
		now.Format("20060102150405"),
		ext.Subtype,
	)

	return &ImageMeta{
		Key:       key,
		Size:      size,
		Extension: ext,
		Width:     x,
		Height:    y,
	}, nil
}

func NewImageForCommon(
	ctx context.Context,
	imageID uuid.UUID,
	file io.Reader,
) (*ImageMeta, error) {
	ctx, span := otel.Tracer("meta").Start(ctx, "NewImageForCommon")
	defer span.End()

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

	key := fmt.Sprintf(
		CommonImageKeyPattern,
		imageID.String(),
		ext.Subtype,
	)

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
