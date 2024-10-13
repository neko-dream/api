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
	"github.com/samber/lo"
)

type opinionHandler struct {
	postOpinionUsecase       opinion_usecase.PostOpinionUseCase
	getOpinionRepliesUsecase opinion_usecase.GetOpinionRepliesUseCase
}

func NewOpinionHandler(
	postOpinionUsecase opinion_usecase.PostOpinionUseCase,
	getOpinionRepliesUsecase opinion_usecase.GetOpinionRepliesUseCase,
) oas.OpinionHandler {
	return &opinionHandler{
		postOpinionUsecase:       postOpinionUsecase,
		getOpinionRepliesUsecase: getOpinionRepliesUsecase,
	}
}

// GetTopOpinions 代表意見取得
func (o *opinionHandler) GetTopOpinions(ctx context.Context, params oas.GetTopOpinionsParams) (oas.GetTopOpinionsRes, error) {
	panic("unimplemented")
}

// SwipeOpinions スワイプ用の意見取得
// 自分が投稿した意見は取得しない
func (o *opinionHandler) SwipeOpinions(ctx context.Context, params oas.SwipeOpinionsParams) (oas.SwipeOpinionsRes, error) {
	panic("unimplemented")
}

// OpinionComments 意見に対するリプライ意見取得
func (o *opinionHandler) OpinionComments(ctx context.Context, params oas.OpinionCommentsParams) (oas.OpinionCommentsRes, error) {
	_, err := o.getOpinionRepliesUsecase.Execute(ctx, opinion_usecase.GetOpinionRepliesInput{
		OpinionID: shared.MustParseUUID[opinion.Opinion](params.OpinionID),
	})
	if err != nil {
		return nil, err
	}

	out := &oas.OpinionCommentsOK{}
	return out, nil

}

// PostOpinionPost implements oas.OpinionHandler.
func (o *opinionHandler) PostOpinionPost(ctx context.Context, req oas.OptPostOpinionPostReq, params oas.PostOpinionPostParams) (oas.PostOpinionPostRes, error) {
	claim := session.GetSession(ctx)
	userID, err := claim.UserID()
	if err != nil {
		return nil, messages.ForbiddenError
	}

	if !req.IsSet() {
		return nil, messages.RequiredParameterError
	}

	talkSessionID := shared.MustParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err := req.Value.Validate(); err != nil {
		return nil, err
	}
	value := req.Value

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
	var parentOpinionID *shared.UUID[opinion.Opinion]
	if value.ParentOpinionID.IsSet() {
		parentOpinionID = lo.ToPtr(shared.MustParseUUID[opinion.Opinion](value.ParentOpinionID.Value))
	}

	_, err = o.postOpinionUsecase.Execute(ctx, opinion_usecase.PostOpinionInput{
		TalkSessionID:   talkSessionID,
		OwnerID:         userID,
		ParentOpinionID: parentOpinionID,
		Title:           utils.ToPtrIfNotNullValue(!req.Value.Title.IsSet(), value.Title.Value),
		Content:         req.Value.OpinionContent,
		ReferenceURL:    utils.ToPtrIfNotNullValue(!req.Value.ReferenceURL.IsSet(), value.ReferenceURL.Value),
		Picture:         file,
	})
	if err != nil {
		return nil, err
	}

	res := &oas.PostOpinionPostOK{}
	return res, nil
}
