package handler

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/presentation/oas"
	opinion_usecase "github.com/neko-dream/server/internal/usecase/opinion"
	http_utils "github.com/neko-dream/server/pkg/http"
	"github.com/neko-dream/server/pkg/utils"
)

type opinionHandler struct {
	postOpinionUsecase opinion_usecase.PostOpinionUseCase
}

func NewOpinionHandler(
	postOpinionUsecase opinion_usecase.PostOpinionUseCase,
) oas.OpinionHandler {
	return &opinionHandler{
		postOpinionUsecase: postOpinionUsecase,
	}
}

// GetTopOpinions implements oas.OpinionHandler.
func (o *opinionHandler) GetTopOpinions(ctx context.Context, params oas.GetTopOpinionsParams) (oas.GetTopOpinionsRes, error) {
	panic("unimplemented")
}

// ListOpinions implements oas.OpinionHandler.
func (o *opinionHandler) ListOpinions(ctx context.Context, params oas.ListOpinionsParams) (oas.ListOpinionsRes, error) {
	panic("unimplemented")
}

// OpinionComments implements oas.OpinionHandler.
func (o *opinionHandler) OpinionComments(ctx context.Context, params oas.OpinionCommentsParams) (oas.OpinionCommentsRes, error) {
	panic("unimplemented")
}

// PostOpinionPost implements oas.OpinionHandler.
func (o *opinionHandler) PostOpinionPost(ctx context.Context, req oas.OptPostOpinionPostReq, params oas.PostOpinionPostParams) (oas.PostOpinionPostRes, error) {
	claim := session.GetSession(ctx)
	if !req.IsSet() {
		return nil, messages.RequiredParameterError
	}

	talkSessionID := shared.MustParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err := req.Value.Validate(); err != nil {
		return nil, err
	}
	value := req.Value

	userID, err := claim.UserID()
	if err != nil {
		return nil, messages.ForbiddenError
	}
	var file *multipart.FileHeader
	if value.Picture.IsSet() {
		content, err := io.ReadAll(value.Picture.Value.File)
		if err != nil {
			utils.HandleError(ctx, err, "io.ReadAll")
			return nil, messages.InternalServerError
		}
		file, err = http_utils.MakeFileHeader(value.Picture.Value.Name, content)
		if err != nil {
			utils.HandleError(ctx, err, "MakeFileHeader")
			return nil, messages.InternalServerError
		}
	}

	_, err = o.postOpinionUsecase.Execute(ctx, opinion_usecase.PostOpinionInput{
		TalkSessionID:   talkSessionID,
		OwnerID:         userID,
		ParentOpinionID: utils.ToPtrIfNotNullValue(req.Value.ParentOpinionID.Null, shared.MustParseUUID[opinion.Opinion](value.ParentOpinionID.Value)),
		Title:           utils.ToPtrIfNotNullValue(req.Value.Title.Null, value.Title.Value),
		Content:         req.Value.OpinionContent,
		ReferenceURL:    utils.ToPtrIfNotNullValue(req.Value.ReferenceURL.Null, value.ReferenceURL.Value),
		Picture:         file,
	})
	if err != nil {
		return nil, err
	}

	res := &oas.PostOpinionPostOK{}
	return res, nil
}
