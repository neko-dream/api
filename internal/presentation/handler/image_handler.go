package handler

import (
	"context"

	"github.com/neko-dream/server/internal/application/usecase/image_usecase"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/presentation/oas"
	http_utils "github.com/neko-dream/server/pkg/http"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type imageHandler struct {
	image_usecase.UploadImage
	authService service.AuthenticationService
}

func NewImageHandler(
	uploadImage image_usecase.UploadImage,
	authService service.AuthenticationService,
) oas.ImageHandler {
	return &imageHandler{
		UploadImage: uploadImage,
		authService: authService,
	}
}

// PostImage POST /image
func (i *imageHandler) PostImage(ctx context.Context, req *oas.PostImageReq) (oas.PostImageRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "imageHandler.PostImage")
	defer span.End()

	authCtx, err := requireAuthentication(i.authService, ctx)
	if err != nil {
		return nil, err
	}

	file, err := http_utils.CreateFileHeader(ctx, req.Image.File, req.GetImage().Name)
	if err != nil {
		utils.HandleError(ctx, err, "MakeFileHeader")
		return nil, messages.InternalServerError
	}

	input := image_usecase.UploadImageInput{
		OwnerID: authCtx.UserID,
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
