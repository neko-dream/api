package http_utils

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"

	"github.com/neko-dream/server/pkg/utils"
)

func CreateFileHeader(ctx context.Context, reader io.Reader, filename string) (*multipart.FileHeader, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		utils.HandleError(ctx, err, "CreateFormFile")
		return nil, err
	}

	if _, err := io.Copy(fw, reader); err != nil {
		utils.HandleError(ctx, err, "Copy")
		return nil, err
	}

	w.Close()

	r := multipart.NewReader(&b, w.Boundary())
	form, err := r.ReadForm(2 << 20) // 2MBのメモリ制限
	if err != nil {
		utils.HandleError(ctx, err, "ReadForm")
		return nil, err
	}

	return form.File["file"][0], nil
}
