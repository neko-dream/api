package handler

import (
	"context"
	"io"
	"time"

	"mime/multipart"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/presentation/oas"
	opinion_usecase "github.com/neko-dream/server/internal/usecase/opinion"
	opinion_query "github.com/neko-dream/server/internal/usecase/query/opinion"
	http_utils "github.com/neko-dream/server/pkg/http"
	"github.com/neko-dream/server/pkg/sort"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
)

type opinionHandler struct {
	postOpinionUsecase           opinion_usecase.PostOpinionUseCase
	getSwipeOpinionsUseCase      opinion_usecase.GetSwipeOpinionsQueryHandler
	getOpinionByTalkSessionQuery opinion_query.GetOpinionsByTalkSessionQuery
	getOpinionDetailByIDQuery    opinion_query.GetOpinionDetailByIDQuery
	getOpinionRepliesQuery       opinion_query.GetOpinionRepliesQuery

	session.TokenManager
}

func NewOpinionHandler(
	postOpinionUsecase opinion_usecase.PostOpinionUseCase,
	getSwipeOpinionsUseCase opinion_usecase.GetSwipeOpinionsQueryHandler,
	getOpinionByTalkSessionUseCase opinion_query.GetOpinionsByTalkSessionQuery,
	getOpinionDetailUseCase opinion_query.GetOpinionDetailByIDQuery,
	getOpinionRepliesQuery opinion_query.GetOpinionRepliesQuery,

	tokenManager session.TokenManager,
) oas.OpinionHandler {
	return &opinionHandler{
		postOpinionUsecase:           postOpinionUsecase,
		getSwipeOpinionsUseCase:      getSwipeOpinionsUseCase,
		getOpinionByTalkSessionQuery: getOpinionByTalkSessionUseCase,
		getOpinionDetailByIDQuery:    getOpinionDetailUseCase,
		getOpinionRepliesQuery:       getOpinionRepliesQuery,

		TokenManager: tokenManager,
	}
}

// GetOpinionsForTalkSession implements oas.OpinionHandler.
func (o *opinionHandler) GetOpinionsForTalkSession(ctx context.Context, params oas.GetOpinionsForTalkSessionParams) (oas.GetOpinionsForTalkSessionRes, error) {
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

	out, err := o.getOpinionByTalkSessionQuery.Execute(ctx, opinion_query.GetOpinionsByTalkSessionInput{
		TalkSessionID: shared.MustParseUUID[talksession.TalkSession](params.TalkSessionID),
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
		opinions = append(opinions, oas.GetOpinionsForTalkSessionOKOpinionsItem{
			Opinion: oas.GetOpinionsForTalkSessionOKOpinionsItemOpinion{
				ID:      opinion.Opinion.OpinionID.String(),
				Title:   utils.ToOpt[oas.OptString](opinion.Opinion.Title),
				Content: opinion.Opinion.Content,
				VoteType: oas.OptGetOpinionsForTalkSessionOKOpinionsItemOpinionVoteType{
					Set:   true,
					Value: oas.GetOpinionsForTalkSessionOKOpinionsItemOpinionVoteType(opinion.GetParentVoteType()),
				},
				ParentID:     utils.ToOpt[oas.OptString](opinion.Opinion.ParentOpinionID),
				ReferenceURL: utils.ToOpt[oas.OptString](opinion.Opinion.ReferenceURL),
				PictureURL:   utils.ToOpt[oas.OptString](opinion.Opinion.PictureURL),
				PostedAt:     opinion.Opinion.CreatedAt.Format(time.RFC3339),
			},
			User: oas.GetOpinionsForTalkSessionOKOpinionsItemUser{
				DisplayID:   opinion.User.DisplayID,
				DisplayName: opinion.User.DisplayName,
				IconURL:     utils.ToOptNil[oas.OptNilString](opinion.User.IconURL),
			},
			ReplyCount: opinion.ReplyCount,
			MyVoteType: oas.GetOpinionsForTalkSessionOKOpinionsItemMyVoteType(opinion.GetMyVoteType()),
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
	claim := session.GetSession(o.SetSession(ctx))
	var userID *shared.UUID[user.User]
	if claim != nil {
		userIDTmp, err := claim.UserID()
		if err != nil {
			return nil, messages.ForbiddenError
		}
		userID = lo.ToPtr(userIDTmp)
	}

	opinionID := shared.MustParseUUID[opinion.Opinion](params.OpinionID)
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
	op := &oas.GetOpinionDetailOKOpinion{
		ID:       opinion.Opinion.Opinion.OpinionID.String(),
		ParentID: utils.ToOpt[oas.OptString](opinion.Opinion.Opinion.ParentOpinionID),
		Title:    utils.ToOpt[oas.OptString](opinion.Opinion.Opinion.Title),
		Content:  opinion.Opinion.Opinion.Content,
		VoteType: oas.OptGetOpinionDetailOKOpinionVoteType{
			Value: oas.GetOpinionDetailOKOpinionVoteType(opinion.Opinion.GetParentVoteType()),
			Set:   true,
		},
		PictureURL:   utils.ToOpt[oas.OptString](opinion.Opinion.Opinion.PictureURL),
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

	opinions, err := o.getSwipeOpinionsUseCase.Execute(ctx, opinion_usecase.GetSwipeOpinionsQuery{
		UserID:        userID,
		TalkSessionID: shared.MustParseUUID[talksession.TalkSession](params.TalkSessionID),
		Limit:         limit,
	})
	if err != nil {
		return nil, err
	}

	var res oas.SwipeOpinionsOKApplicationJSON
	for _, opinion := range opinions.Opinions {
		user := &oas.SwipeOpinionsOKItemUser{
			DisplayID:   opinion.User.ID,
			DisplayName: opinion.User.Name,
			IconURL:     utils.ToOptNil[oas.OptNilString](opinion.User.Icon),
		}

		ops := &oas.SwipeOpinionsOKItemOpinion{
			ID:       opinion.Opinion.OpinionID,
			ParentID: utils.ToOpt[oas.OptString](opinion.Opinion.ParentOpinionID),
			Title:    utils.ToOpt[oas.OptString](opinion.Opinion.Title),
			Content:  opinion.Opinion.Content,
			VoteType: oas.OptSwipeOpinionsOKItemOpinionVoteType{
				Value: oas.SwipeOpinionsOKItemOpinionVoteType(opinion.Opinion.VoteType),
				Set:   true,
			},
			PictureURL:   utils.ToOpt[oas.OptString](opinion.Opinion.PictureURL),
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
	claim := session.GetSession(o.SetSession(ctx))
	var userID *shared.UUID[user.User]
	if claim != nil {
		userIDTmp, err := claim.UserID()
		if err != nil {
			return nil, messages.ForbiddenError
		}
		userID = lo.ToPtr(userIDTmp)
	}

	opinions, err := o.getOpinionRepliesQuery.Execute(ctx, opinion_query.GetOpinionRepliesQueryInput{
		OpinionID: shared.MustParseUUID[opinion.Opinion](params.OpinionID),
		UserID:    userID,
	})
	if err != nil {
		return nil, err
	}
	rootUser := &oas.OpinionCommentsOKRootOpinionUser{
		DisplayID:   opinions.RootOpinion.User.DisplayID,
		DisplayName: opinions.RootOpinion.User.DisplayName,
		IconURL:     utils.ToOptNil[oas.OptNilString](opinions.RootOpinion.User.IconURL),
	}
	rootOpinion := &oas.OpinionCommentsOKRootOpinionOpinion{
		ID:      opinions.RootOpinion.Opinion.OpinionID.String(),
		Title:   utils.ToOpt[oas.OptString](opinions.RootOpinion.Opinion.Title),
		Content: opinions.RootOpinion.Opinion.Content,
		VoteType: oas.OptOpinionCommentsOKRootOpinionOpinionVoteType{
			Value: oas.OpinionCommentsOKRootOpinionOpinionVoteType(opinions.RootOpinion.GetParentVoteType()),
			Set:   true,
		},
		ParentID:     utils.ToOpt[oas.OptString](opinions.RootOpinion.Opinion.ParentOpinionID),
		PictureURL:   utils.ToOpt[oas.OptString](opinions.RootOpinion.Opinion.PictureURL),
		ReferenceURL: utils.ToOpt[oas.OptString](opinions.RootOpinion.Opinion.ReferenceURL),
		PostedAt:     opinions.RootOpinion.Opinion.CreatedAt.Format(time.RFC3339),
	}
	root := oas.OpinionCommentsOKRootOpinion{
		User:    *rootUser,
		Opinion: *rootOpinion,
		MyVoteType: oas.OptOpinionCommentsOKRootOpinionMyVoteType{
			Value: oas.OpinionCommentsOKRootOpinionMyVoteType(opinions.RootOpinion.GetMyVoteType()),
			Set:   true,
		},
	}

	var replies []oas.OpinionCommentsOKOpinionsItem
	for _, reply := range opinions.Replies {
		user := &oas.OpinionCommentsOKOpinionsItemUser{
			DisplayID:   reply.User.DisplayID,
			DisplayName: reply.User.DisplayName,
			IconURL:     utils.ToOptNil[oas.OptNilString](reply.User.IconURL),
		}

		opinion := &oas.OpinionCommentsOKOpinionsItemOpinion{
			ID:       reply.Opinion.OpinionID.String(),
			ParentID: utils.ToOpt[oas.OptString](reply.Opinion.ParentOpinionID),
			Title:    utils.ToOpt[oas.OptString](reply.Opinion.Title),
			Content:  reply.Opinion.Content,
			VoteType: oas.OptOpinionCommentsOKOpinionsItemOpinionVoteType{
				Value: oas.OpinionCommentsOKOpinionsItemOpinionVoteType(reply.GetParentVoteType()),
				Set:   true,
			},
			PictureURL:   utils.ToOpt[oas.OptString](reply.Opinion.PictureURL),
			ReferenceURL: utils.ToOpt[oas.OptString](reply.Opinion.ReferenceURL),
			PostedAt:     reply.Opinion.CreatedAt.Format(time.RFC3339),
		}
		replies = append(replies, oas.OpinionCommentsOKOpinionsItem{
			User:    *user,
			Opinion: *opinion,
			MyVoteType: oas.OptOpinionCommentsOKOpinionsItemMyVoteType{
				Value: oas.OpinionCommentsOKOpinionsItemMyVoteType(reply.GetMyVoteType()),
				Set:   true,
			},
		})
	}

	return &oas.OpinionCommentsOK{
		RootOpinion: root,
		Opinions:    replies,
	}, nil

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
