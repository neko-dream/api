package handler

import (
	"context"

	"github.com/neko-dream/server/internal/application/usecase/image_usecase"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/presentation/oas"
	http_utils "github.com/neko-dream/server/pkg/http"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type imageHandler struct {
	image_usecase.UploadImage
}

func NewImageHandler(
	uploadImage image_usecase.UploadImage,
) oas.ImageHandler {
	return &imageHandler{
		UploadImage: uploadImage,
	}
}

// PostImage POST /image
func (i *imageHandler) PostImage(ctx context.Context, req oas.OptPostImageReq) (oas.PostImageRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "imageHandler.PostImage")
	defer span.End()

	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}

	userID, err := claim.UserID()
	if err != nil {
		utils.HandleError(ctx, err, "claim.UserID")
		return nil, messages.ForbiddenError
	}

	file, err := http_utils.CreateFileHeader(ctx, req.Value.GetImage().File, req.Value.GetImage().Name)
	if err != nil {
		utils.HandleError(ctx, err, "MakeFileHeader")
		return nil, messages.InternalServerError
	}

	input := image_usecase.UploadImageInput{
		OwnerID: userID,
		Image:   file,
	}

	output, err := i.UploadImage.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return &oas.PostImageOK{
		URL: output.ImageURL,
	}, nil
}
