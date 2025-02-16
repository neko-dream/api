package handler

import (
	"context"

	"github.com/neko-dream/server/internal/presentation/oas"
)

type imageHandler struct {
}

func NewImageHandler() oas.ImageHandler {
	return &imageHandler{}
}

// PostImage implements oas.ImageHandler.
func (i *imageHandler) PostImage(ctx context.Context, req oas.OptPostImageReq) (oas.PostImageRes, error) {
	panic("unimplemented")
}
