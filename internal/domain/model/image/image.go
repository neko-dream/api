package image

import (
	"bytes"
	"context"
	"image"

	"io"
	"mime/multipart"
	"strings"

	"github.com/neko-dream/server/internal/domain/messages"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
)

type (
	Image []byte
)

func NewImage(image []byte) Image {
	return Image(image)
}

func (i *Image) GetImage() []byte {
	return *i
}
func (i *Image) GetImageReader() io.Reader {
	return bytes.NewReader(*i)
}

func CheckSize(imageByte []byte, x, y int) error {
	// 画像のサイズを取得
	img, _, err := image.Decode(bytes.NewReader(imageByte))
	if err != nil {
		return messages.ImageDecodeFailedError
	}

	// 画像のサイズが指定したサイズより大きい場合はエラー
	if img.Bounds().Dx() > x || img.Bounds().Dy() > y {
		return messages.ImageSizeTooLargeError
	}

	return nil
}

var supportedExtension = "png jpg"

// ValidateImage 拡張子やファイルサイズなどの画像のバリデーションを行う
func ValidateImage(ctx context.Context, file *multipart.FileHeader, maxSize int) ([]byte, *string, error) {
	img, err := file.Open()
	if err != nil {
		return nil, nil, messages.ImageOpenFailedError
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(img); err != nil {
		return nil, nil, messages.ImageOpenFailedError
	}

	// 最大サイズに収まっているか
	if buf.Len() >= maxSize {
		return nil, nil, messages.ImageSizeTooLargeError
	}

	// 対応画像形式か
	reader, t, err := detectExt(buf)
	// ファイル形式を検知するのに必要なバイト数だけ先に読む
	if err != nil || !strings.Contains(supportedExtension, t.Extension) {
		return nil, nil, messages.ImageUnsupportedExtError
	}
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, nil, messages.ImageOpenFailedError
	}

	return body, &t.MIME.Value, nil
}

func detectExt(file *bytes.Buffer) (io.Reader, types.Type, error) {
	header := bytes.NewBuffer(nil)
	reader := io.TeeReader(file, header)

	head := make([]byte, 261)
	if _, err := reader.Read(head); err != nil {
		return nil, types.Unknown, err
	}
	t, e := filetype.Match(head)

	multiReader := io.MultiReader(header, file)
	return multiReader, t, e

}
