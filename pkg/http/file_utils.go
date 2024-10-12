package http_utils

import (
	"bytes"
	"io"
	"mime/multipart"
)

func MakeFileHeader(name string, dateBytes []byte) (*multipart.FileHeader, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", name)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, bytes.NewReader(dateBytes)); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	// ダミーファイルをパースして、multipart.FileHeaderを取得する
	reader := multipart.NewReader(body, writer.Boundary())
	// 2M
	form, err := reader.ReadForm(2 * 1_000_000)
	if err != nil {
		return nil, err
	}
	fh := form.File["file"]

	return fh[0], nil
}
