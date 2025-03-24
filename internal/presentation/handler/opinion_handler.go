package handler

import (
	"context"
	"time"

	"mime/multipart"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/neko-dream/server/internal/usecase/command/opinion_command"
	opinion_query "github.com/neko-dream/server/internal/usecase/query/opinion"
	http_utils "github.com/neko-dream/server/pkg/http"
	"github.com/neko-dream/server/pkg/sort"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type opinionHandler struct {
	getOpinionByTalkSessionQuery opinion_query.GetOpinionsByTalkSessionQuery
	getOpinionDetailByIDQuery    opinion_query.GetOpinionDetailByIDQuery
	getOpinionRepliesQuery       opinion_query.GetOpinionRepliesQuery
	getSwipeOpinionQuery         opinion_query.GetSwipeOpinionsQuery

	submitOpinionCommand opinion_command.SubmitOpinion
	reportOpinionCommand opinion_command.ReportOpinion

	session.TokenManager
}

func NewOpinionHandler(
	getOpinionByTalkSessionUseCase opinion_query.GetOpinionsByTalkSessionQuery,
	getOpinionDetailUseCase opinion_query.GetOpinionDetailByIDQuery,
	getOpinionRepliesQuery opinion_query.GetOpinionRepliesQuery,
	getSwipeOpinionsQuery opinion_query.GetSwipeOpinionsQuery,

	submitOpinionCommand opinion_command.SubmitOpinion,
	reportOpinionCommand opinion_command.ReportOpinion,

	tokenManager session.TokenManager,
) oas.OpinionHandler {
	return &opinionHandler{
		getOpinionByTalkSessionQuery: getOpinionByTalkSessionUseCase,
		getOpinionDetailByIDQuery:    getOpinionDetailUseCase,
		getOpinionRepliesQuery:       getOpinionRepliesQuery,
		getSwipeOpinionQuery:         getSwipeOpinionsQuery,

		submitOpinionCommand: submitOpinionCommand,
		reportOpinionCommand: reportOpinionCommand,

		TokenManager: tokenManager,
	}
}

// GetOpinionDetail2 GetOpinionDetailは/talksessions/{talkSessionID}/opinions/{opinionID}だが、長いので `/opinion/{opinionID}` なGetOpinionDetail2を作成
func (o *opinionHandler) GetOpinionDetail2(ctx context.Context, params oas.GetOpinionDetail2Params) (oas.GetOpinionDetail2Res, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "opinionHandler.GetOpinionDetail")
	defer span.End()

	claim := session.GetSession(o.SetSession(ctx))
	var userID *shared.UUID[user.User]
	if claim != nil {
		userIDTmp, err := claim.UserID()
		if err != nil {
			return nil, messages.ForbiddenError
		}
		userID = lo.ToPtr(userIDTmp)
	}

	opinionID, err := shared.ParseUUID[opinion.Opinion](params.OpinionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	opinion, err := o.getOpinionDetailByIDQuery.Execute(ctx, opinion_query.GetOpinionDetailByIDInput{
		OpinionID: opinionID,
		UserID:    userID,
	})
	if err != nil {
		return nil, err
	}

	user := &oas.GetOpinionDetail2OKUser{
		DisplayID:   opinion.Opinion.User.DisplayID,
		DisplayName: opinion.Opinion.User.DisplayName,
		IconURL:     utils.ToOptNil[oas.OptNilString](opinion.Opinion.User.IconURL),
	}
	var parentOpinionID oas.OptString
	if opinion.Opinion.Opinion.ParentOpinionID != nil {
		parentOpinionID = utils.ToOpt[oas.OptString](opinion.Opinion.Opinion.ParentOpinionID.String())
	}

	op := &oas.GetOpinionDetail2OKOpinion{
		ID:           opinion.Opinion.Opinion.OpinionID.String(),
		ParentID:     parentOpinionID,
		Title:        utils.ToOpt[oas.OptString](opinion.Opinion.Opinion.Title),
		Content:      opinion.Opinion.Opinion.Content,
		VoteType:     utils.ToOptNil[oas.OptNilString](opinion.Opinion.GetParentVoteType()),
		PictureURL:   utils.ToOptNil[oas.OptNilString](opinion.Opinion.Opinion.PictureURL),
		ReferenceURL: utils.ToOpt[oas.OptString](opinion.Opinion.Opinion.ReferenceURL),
		PostedAt:     opinion.Opinion.Opinion.CreatedAt.Format(time.RFC3339),
	}

	return &oas.GetOpinionDetail2OK{
		User:    *user,
		Opinion: *op,
	}, nil
}

// OpinionComments2 implements oas.OpinionHandler.
func (o *opinionHandler) OpinionComments2(ctx context.Context, params oas.OpinionComments2Params) (oas.OpinionComments2Res, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "opinionHandler.OpinionComments")
	defer span.End()

	claim := session.GetSession(o.SetSession(ctx))
	var userID *shared.UUID[user.User]
	if claim != nil {
		userIDTmp, err := claim.UserID()
		if err != nil {
			return nil, messages.ForbiddenError
		}
		userID = lo.ToPtr(userIDTmp)
	}

	opinionID, err := shared.ParseUUID[opinion.Opinion](params.OpinionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	opinions, err := o.getOpinionRepliesQuery.Execute(ctx, opinion_query.GetOpinionRepliesQueryInput{
		OpinionID: opinionID,
		UserID:    userID,
	})
	if err != nil {
		return nil, err
	}

	var replies []oas.OpinionComments2OKOpinionsItem
	for _, reply := range opinions.Replies {
		user := &oas.OpinionComments2OKOpinionsItemUser{
			DisplayID:   reply.User.DisplayID,
			DisplayName: reply.User.DisplayName,
			IconURL:     utils.ToOptNil[oas.OptNilString](reply.User.IconURL),
		}
		var parentOpinionID oas.OptString
		if reply.Opinion.ParentOpinionID != nil {
			parentOpinionID = utils.ToOpt[oas.OptString](reply.Opinion.ParentOpinionID.String())
		}

		opinion := &oas.OpinionComments2OKOpinionsItemOpinion{
			ID:           reply.Opinion.OpinionID.String(),
			ParentID:     parentOpinionID,
			Title:        utils.ToOpt[oas.OptString](reply.Opinion.Title),
			Content:      reply.Opinion.Content,
			VoteType:     utils.ToOptNil[oas.OptNilString](reply.GetParentVoteType()),
			PictureURL:   utils.ToOptNil[oas.OptNilString](reply.Opinion.PictureURL),
			ReferenceURL: utils.ToOpt[oas.OptString](reply.Opinion.ReferenceURL),
			PostedAt:     reply.Opinion.CreatedAt.Format(time.RFC3339),
		}

		replies = append(replies, oas.OpinionComments2OKOpinionsItem{
			User:       *user,
			Opinion:    *opinion,
			MyVoteType: utils.ToOptNil[oas.OptNilString](reply.GetMyVoteType()),
		})
	}

	return &oas.OpinionComments2OK{
		Opinions: replies,
	}, nil
}

// GetOpinionsForTalkSession implements oas.OpinionHandler.
func (o *opinionHandler) GetOpinionsForTalkSession(ctx context.Context, params oas.GetOpinionsForTalkSessionParams) (oas.GetOpinionsForTalkSessionRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "opinionHandler.GetOpinionsForTalkSession")
	defer span.End()

	claim := session.GetSession(o.SetSession(ctx))
	var userID *shared.UUID[user.User]
	if claim != nil {
		id, err := claim.UserID()
		if err == nil {
			userID = &id
		}
	}

	var sortKey sort.SortKey
	if params.Sort.IsSet() {
		txt, err := params.Sort.Value.MarshalText()
		if err != nil {
			utils.HandleError(ctx, err, "params.Sort.Value.MarshalText")
			return nil, messages.InternalServerError
		}
		sortKey = sort.SortKey(txt)
	}
	var limit, offset *int
	if params.Limit.IsSet() {
		limit = &params.Limit.Value
	}
	if params.Offset.IsSet() {
		offset = &params.Offset.Value
	}
	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	out, err := o.getOpinionByTalkSessionQuery.Execute(ctx, opinion_query.GetOpinionsByTalkSessionInput{
		TalkSessionID: talkSessionID,
		SortKey:       sortKey,
		Limit:         limit,
		Offset:        offset,
		UserID:        userID,
	})
	if err != nil {
		return nil, err
	}
	opinions := make([]oas.GetOpinionsForTalkSessionOKOpinionsItem, 0, len(out.Opinions))
	for _, opinion := range out.Opinions {
		var parentOpinionID oas.OptString
		if opinion.Opinion.ParentOpinionID != nil {
			parentOpinionID = utils.ToOpt[oas.OptString](opinion.Opinion.ParentOpinionID.String())
		}

		opinions = append(opinions, oas.GetOpinionsForTalkSessionOKOpinionsItem{
			Opinion: oas.GetOpinionsForTalkSessionOKOpinionsItemOpinion{
				ID:           opinion.Opinion.OpinionID.String(),
				Title:        utils.ToOpt[oas.OptString](opinion.Opinion.Title),
				Content:      opinion.Opinion.Content,
				VoteType:     utils.ToOptNil[oas.OptNilString](opinion.GetParentVoteType()),
				ParentID:     parentOpinionID,
				ReferenceURL: utils.ToOpt[oas.OptString](opinion.Opinion.ReferenceURL),
				PictureURL:   utils.ToOptNil[oas.OptNilString](opinion.Opinion.PictureURL),
				PostedAt:     opinion.Opinion.CreatedAt.Format(time.RFC3339),
			},
			User: oas.GetOpinionsForTalkSessionOKOpinionsItemUser{
				DisplayID:   opinion.User.DisplayID,
				DisplayName: opinion.User.DisplayName,
				IconURL:     utils.ToOptNil[oas.OptNilString](opinion.User.IconURL),
			},
			ReplyCount: opinion.ReplyCount,
			MyVoteType: utils.ToOptNil[oas.OptNilString](opinion.GetMyVoteType()),
		})
	}

	return &oas.GetOpinionsForTalkSessionOK{
		Opinions: opinions,
		Pagination: oas.GetOpinionsForTalkSessionOKPagination{
			TotalCount: out.TotalCount,
		},
	}, nil
}

// GetOpinionDetail implements oas.OpinionHandler.
func (o *opinionHandler) GetOpinionDetail(ctx context.Context, params oas.GetOpinionDetailParams) (oas.GetOpinionDetailRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "opinionHandler.GetOpinionDetail")
	defer span.End()

	claim := session.GetSession(o.SetSession(ctx))
	var userID *shared.UUID[user.User]
	if claim != nil {
		userIDTmp, err := claim.UserID()
		if err != nil {
			return nil, messages.ForbiddenError
		}
		userID = lo.ToPtr(userIDTmp)
	}

	opinionID, err := shared.ParseUUID[opinion.Opinion](params.OpinionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	opinion, err := o.getOpinionDetailByIDQuery.Execute(ctx, opinion_query.GetOpinionDetailByIDInput{
		OpinionID: opinionID,
		UserID:    userID,
	})
	if err != nil {
		return nil, err
	}

	user := &oas.GetOpinionDetailOKUser{
		DisplayID:   opinion.Opinion.User.DisplayID,
		DisplayName: opinion.Opinion.User.DisplayName,
		IconURL:     utils.ToOptNil[oas.OptNilString](opinion.Opinion.User.IconURL),
	}

	var parentOpinionID oas.OptString
	if opinion.Opinion.Opinion.ParentOpinionID != nil {
		parentOpinionID = utils.ToOpt[oas.OptString](opinion.Opinion.Opinion.ParentOpinionID.String())
	}

	op := &oas.GetOpinionDetailOKOpinion{
		ID:           opinion.Opinion.Opinion.OpinionID.String(),
		ParentID:     parentOpinionID,
		Title:        utils.ToOpt[oas.OptString](opinion.Opinion.Opinion.Title),
		Content:      opinion.Opinion.Opinion.Content,
		VoteType:     utils.ToOptNil[oas.OptNilString](opinion.Opinion.GetParentVoteType()),
		PictureURL:   utils.ToOptNil[oas.OptNilString](opinion.Opinion.Opinion.PictureURL),
		ReferenceURL: utils.ToOpt[oas.OptString](opinion.Opinion.Opinion.ReferenceURL),
		PostedAt:     opinion.Opinion.Opinion.CreatedAt.Format(time.RFC3339),
	}

	return &oas.GetOpinionDetailOK{
		User:    *user,
		Opinion: *op,
	}, nil
}

// SwipeOpinions スワイプ用の意見取得
// 自分が投稿した意見は取得しない
func (o *opinionHandler) SwipeOpinions(ctx context.Context, params oas.SwipeOpinionsParams) (oas.SwipeOpinionsRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "opinionHandler.SwipeOpinions")
	defer span.End()

	claim := session.GetSession(ctx)
	userID, err := claim.UserID()
	if err != nil {
		return nil, messages.ForbiddenError
	}
	var limit int
	if params.Limit.IsSet() {
		limit = params.Limit.Value
	} else {
		limit = 10
	}

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	opinions, err := o.getSwipeOpinionQuery.Execute(ctx, opinion_query.GetSwipeOpinionsQueryInput{
		UserID:        userID,
		TalkSessionID: talkSessionID,
		Limit:         limit,
	})
	if err != nil {
		return nil, err
	}

	var res oas.SwipeOpinionsOKApplicationJSON
	for _, opinion := range opinions.Opinions {
		user := &oas.SwipeOpinionsOKItemUser{
			DisplayID:   opinion.User.DisplayID,
			DisplayName: opinion.User.DisplayName,
			IconURL:     utils.ToOptNil[oas.OptNilString](opinion.User.IconURL),
		}
		var parentOpinionID oas.OptString
		if opinion.Opinion.ParentOpinionID != nil {
			parentOpinionID = utils.ToOpt[oas.OptString](opinion.Opinion.ParentOpinionID.String())
		}
		ops := &oas.SwipeOpinionsOKItemOpinion{
			ID:           opinion.Opinion.OpinionID.String(),
			ParentID:     parentOpinionID,
			Title:        utils.ToOpt[oas.OptString](opinion.Opinion.Title),
			Content:      opinion.Opinion.Content,
			VoteType:     utils.ToOptNil[oas.OptNilString](opinion.GetParentVoteType()),
			PictureURL:   utils.ToOptNil[oas.OptNilString](opinion.Opinion.PictureURL),
			ReferenceURL: utils.ToOpt[oas.OptString](opinion.Opinion.ReferenceURL),
			PostedAt:     opinion.Opinion.CreatedAt.Format(time.RFC3339),
		}
		res = append(res, oas.SwipeOpinionsOKItem{
			User:       *user,
			Opinion:    *ops,
			ReplyCount: opinion.ReplyCount,
		})
	}

	return &res, nil
}

// OpinionComments 意見に対するリプライ意見取得
func (o *opinionHandler) OpinionComments(ctx context.Context, params oas.OpinionCommentsParams) (oas.OpinionCommentsRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "opinionHandler.OpinionComments")
	defer span.End()

	claim := session.GetSession(o.SetSession(ctx))
	var userID *shared.UUID[user.User]
	if claim != nil {
		userIDTmp, err := claim.UserID()
		if err != nil {
			return nil, messages.ForbiddenError
		}
		userID = lo.ToPtr(userIDTmp)
	}

	opinionID, err := shared.ParseUUID[opinion.Opinion](params.OpinionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	opinions, err := o.getOpinionRepliesQuery.Execute(ctx, opinion_query.GetOpinionRepliesQueryInput{
		OpinionID: opinionID,
		UserID:    userID,
	})
	if err != nil {
		return nil, err
	}

	var replies []oas.OpinionCommentsOKOpinionsItem
	for _, reply := range opinions.Replies {
		user := &oas.OpinionCommentsOKOpinionsItemUser{
			DisplayID:   reply.User.DisplayID,
			DisplayName: reply.User.DisplayName,
			IconURL:     utils.ToOptNil[oas.OptNilString](reply.User.IconURL),
		}
		var parentOpinionID oas.OptString
		if reply.Opinion.ParentOpinionID != nil {
			parentOpinionID = utils.ToOpt[oas.OptString](reply.Opinion.ParentOpinionID.String())
		}
		opinion := &oas.OpinionCommentsOKOpinionsItemOpinion{
			ID:           reply.Opinion.OpinionID.String(),
			ParentID:     parentOpinionID,
			Title:        utils.ToOpt[oas.OptString](reply.Opinion.Title),
			Content:      reply.Opinion.Content,
			VoteType:     utils.ToOptNil[oas.OptNilString](reply.GetParentVoteType()),
			PictureURL:   utils.ToOptNil[oas.OptNilString](reply.Opinion.PictureURL),
			ReferenceURL: utils.ToOpt[oas.OptString](reply.Opinion.ReferenceURL),
			PostedAt:     reply.Opinion.CreatedAt.Format(time.RFC3339),
		}

		replies = append(replies, oas.OpinionCommentsOKOpinionsItem{
			User:       *user,
			Opinion:    *opinion,
			MyVoteType: utils.ToOptNil[oas.OptNilString](reply.GetMyVoteType()),
		})
	}

	return &oas.OpinionCommentsOK{
		Opinions: replies,
	}, nil

}

// PostOpinionPost implements oas.OpinionHandler.
func (o *opinionHandler) PostOpinionPost(ctx context.Context, req oas.OptPostOpinionPostReq, params oas.PostOpinionPostParams) (oas.PostOpinionPostRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "opinionHandler.PostOpinionPost")
	defer span.End()

	claim := session.GetSession(ctx)
	userID, err := claim.UserID()
	if err != nil {
		return nil, messages.ForbiddenError
	}
	if !req.IsSet() {
		return nil, messages.RequiredParameterError
	}

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	if err := req.Value.Validate(); err != nil {
		return nil, err
	}
	value := req.Value

	var file *multipart.FileHeader
	if value.Picture.IsSet() {
		file, err = http_utils.CreateFileHeader(ctx, value.Picture.Value.File, value.Picture.Value.Name)
		if err != nil {
			utils.HandleError(ctx, err, "MakeFileHeader")
			return nil, messages.InternalServerError
		}
	}
	var parentOpinionID *shared.UUID[opinion.Opinion]
	if value.ParentOpinionID.IsSet() {
		id, err := shared.ParseUUID[opinion.Opinion](value.ParentOpinionID.Value)
		if err != nil {
			return nil, messages.BadRequestError
		}
		parentOpinionID = &id
	}

	if err = o.submitOpinionCommand.Execute(ctx, opinion_command.SubmitOpinionInput{
		TalkSessionID:   &talkSessionID,
		OwnerID:         userID,
		UserID:          userID,
		ParentOpinionID: parentOpinionID,
		Title:           utils.ToPtrIfNotNullValue(!req.Value.Title.IsSet(), value.Title.Value),
		Content:         req.Value.OpinionContent,
		ReferenceURL:    utils.ToPtrIfNotNullValue(!req.Value.ReferenceURL.IsSet(), value.ReferenceURL.Value),
		Picture:         file,
	}); err != nil {
		return nil, err
	}

	res := &oas.PostOpinionPostOK{}
	return res, nil
}

// PostOpinionPost2 TalkSessionIDをBodyで受け取るタイプのやつ
func (o *opinionHandler) PostOpinionPost2(ctx context.Context, req oas.OptPostOpinionPost2Req) (oas.PostOpinionPost2Res, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "opinionHandler.PostOpinionPost2")
	defer span.End()

	claim := session.GetSession(ctx)
	userID, err := claim.UserID()
	if err != nil {
		return nil, messages.ForbiddenError
	}
	if !req.IsSet() {
		return nil, messages.RequiredParameterError
	}

	if err := req.Value.Validate(); err != nil {
		return nil, err
	}
	value := req.Value

	var talkSessionID *shared.UUID[talksession.TalkSession]
	if value.TalkSessionID.IsSet() {
		id, err := shared.ParseUUID[talksession.TalkSession](value.TalkSessionID.Value)
		if err != nil {
			return nil, messages.BadRequestError
		}
		talkSessionID = &id
	}

	var file *multipart.FileHeader
	if value.Picture.IsSet() {
		file, err = http_utils.CreateFileHeader(ctx, value.Picture.Value.File, value.Picture.Value.Name)
		if err != nil {
			utils.HandleError(ctx, err, "MakeFileHeader")
			return nil, messages.InternalServerError
		}
	}
	var parentOpinionID *shared.UUID[opinion.Opinion]
	if value.ParentOpinionID.IsSet() {
		id, err := shared.ParseUUID[opinion.Opinion](value.ParentOpinionID.Value)
		if err != nil {
			return nil, messages.BadRequestError
		}
		parentOpinionID = &id
	}

	if err = o.submitOpinionCommand.Execute(ctx, opinion_command.SubmitOpinionInput{
		TalkSessionID:   talkSessionID,
		OwnerID:         userID,
		UserID:          userID,
		ParentOpinionID: parentOpinionID,
		Title:           utils.ToPtrIfNotNullValue(!req.Value.Title.IsSet(), value.Title.Value),
		Content:         req.Value.OpinionContent,
		ReferenceURL:    utils.ToPtrIfNotNullValue(!req.Value.ReferenceURL.IsSet(), value.ReferenceURL.Value),
		Picture:         file,
	}); err != nil {
		return nil, err
	}

	res := &oas.PostOpinionPost2OK{}
	return res, nil
}

// ReportOpinion 意見を通報する
func (o *opinionHandler) ReportOpinion(ctx context.Context, req oas.OptReportOpinionReq, params oas.ReportOpinionParams) (oas.ReportOpinionRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "opinionHandler.ReportOpinion")
	defer span.End()

	claim := session.GetSession(ctx)
	userID, err := claim.UserID()
	if err != nil {
		return nil, messages.ForbiddenError
	}
	if !req.IsSet() {
		return nil, messages.RequiredParameterError
	}

	opinionID, err := shared.ParseUUID[opinion.Opinion](params.OpinionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	if err := o.reportOpinionCommand.Execute(ctx, opinion_command.ReportOpinionInput{
		ReporterID: userID,
		OpinionID:  opinionID,
		Reason:     int32(req.Value.Reason.Value),
	}); err != nil {
		utils.HandleError(ctx, err, "reportOpinionCommand.Execute")
		return nil, err
	}

	res := &oas.ReportOpinionOK{}
	return res, nil
}
