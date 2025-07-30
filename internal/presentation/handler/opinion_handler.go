package handler

import (
	"context"
	"mime/multipart"

	opinion_query "github.com/neko-dream/server/internal/application/query/opinion"
	"github.com/neko-dream/server/internal/application/query/report_query"
	"github.com/neko-dream/server/internal/application/usecase/opinion_usecase"
	"github.com/neko-dream/server/internal/application/usecase/report_usecase"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/presentation/oas"
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
	getReportReasons             opinion_query.GetReportReasons
	getOpinionGroupRatio         opinion_query.GetOpinionGroupRatioQuery
	getReportByOpinionID         report_query.GetOpinionReportQuery

	submitOpinionCommand opinion_usecase.SubmitOpinion
	reportOpinionCommand opinion_usecase.ReportOpinion
	solveReportCommand   report_usecase.SolveReportCommand

	authService service.AuthenticationService
	session.TokenManager
}

func NewOpinionHandler(
	getOpinionByTalkSessionUseCase opinion_query.GetOpinionsByTalkSessionQuery,
	getOpinionDetailUseCase opinion_query.GetOpinionDetailByIDQuery,
	getOpinionRepliesQuery opinion_query.GetOpinionRepliesQuery,
	getSwipeOpinionsQuery opinion_query.GetSwipeOpinionsQuery,
	getReportReasons opinion_query.GetReportReasons,
	getOpinionGroupRatio opinion_query.GetOpinionGroupRatioQuery,
	getReportByOpinionID report_query.GetOpinionReportQuery,

	submitOpinionCommand opinion_usecase.SubmitOpinion,
	reportOpinionCommand opinion_usecase.ReportOpinion,
	solveReportCommand report_usecase.SolveReportCommand,

	authService service.AuthenticationService,
	tokenManager session.TokenManager,
) oas.OpinionHandler {
	return &opinionHandler{
		getOpinionByTalkSessionQuery: getOpinionByTalkSessionUseCase,
		getOpinionDetailByIDQuery:    getOpinionDetailUseCase,
		getOpinionRepliesQuery:       getOpinionRepliesQuery,
		getSwipeOpinionQuery:         getSwipeOpinionsQuery,
		getReportReasons:             getReportReasons,
		getOpinionGroupRatio:         getOpinionGroupRatio,
		getReportByOpinionID:         getReportByOpinionID,

		submitOpinionCommand: submitOpinionCommand,
		reportOpinionCommand: reportOpinionCommand,
		solveReportCommand:   solveReportCommand,

		authService:  authService,
		TokenManager: tokenManager,
	}
}

// GetOpinionDetail2 GetOpinionDetailは/talksessions/{talkSessionID}/opinions/{opinionID}だが、長いので `/opinion/{opinionID}` なGetOpinionDetail2を作成
func (o *opinionHandler) GetOpinionDetail2(ctx context.Context, params oas.GetOpinionDetail2Params) (oas.GetOpinionDetail2Res, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "opinionHandler.GetOpinionDetail")
	defer span.End()

	authCtx, err := getAuthenticationContext(o.authService, o.SetSession(ctx))
	if err != nil {
		return nil, err
	}

	opinionID, err := shared.ParseUUID[opinion.Opinion](params.OpinionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	opinion, err := o.getOpinionDetailByIDQuery.Execute(ctx, opinion_query.GetOpinionDetailByIDInput{
		OpinionID: opinionID,
		UserID:    lo.ToPtr(authCtx.UserID),
	})
	if err != nil {
		return nil, err
	}

	// Convert to OpinionWithVote response
	var myVoteType oas.OptNilOpinionWithVoteMyVoteType
	if opinion.Opinion.GetMyVoteType() != nil {
		myVoteType = oas.OptNilOpinionWithVoteMyVoteType{
			Value: oas.OpinionWithVoteMyVoteType(*opinion.Opinion.GetMyVoteType()),
			Set:   true,
			Null:  false,
		}
	}

	result := oas.OpinionWithVote{
		User:       opinion.Opinion.User.ToResponse(),
		Opinion:    opinion.Opinion.Opinion.ToResponse(),
		MyVoteType: myVoteType,
	}
	return &result, nil
}

// OpinionComments2 implements oas.OpinionHandler.
func (o *opinionHandler) OpinionComments2(ctx context.Context, params oas.OpinionComments2Params) (oas.OpinionComments2Res, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "opinionHandler.OpinionComments")
	defer span.End()

	authCtx, err := getAuthenticationContext(o.authService, o.SetSession(ctx))
	var userID *shared.UUID[user.User]
	if err == nil {
		userID = &authCtx.UserID
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

	var replies []oas.OpinionWithVote
	for _, reply := range opinions.Replies {
		// Convert SwipeOpinion to OpinionWithVote
		var myVoteType oas.OptNilOpinionWithVoteMyVoteType
		if reply.GetMyVoteType() != nil {
			myVoteType = oas.OptNilOpinionWithVoteMyVoteType{
				Value: oas.OpinionWithVoteMyVoteType(*reply.GetMyVoteType()),
				Set:   true,
				Null:  false,
			}
		}

		replies = append(replies, oas.OpinionWithVote{
			User:       reply.User.ToResponse(),
			Opinion:    reply.Opinion.ToResponse(),
			MyVoteType: myVoteType,
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

	authCtx, err := getAuthenticationContext(o.authService, o.SetSession(ctx))
	var userID *shared.UUID[user.User]
	if err == nil {
		userID = &authCtx.UserID
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
	var seed bool
	if params.Seed.IsSet() {
		seed = params.Seed.Value
	} else {
		seed = false
	}

	out, err := o.getOpinionByTalkSessionQuery.Execute(ctx, opinion_query.GetOpinionsByTalkSessionInput{
		TalkSessionID: talkSessionID,
		SortKey:       sortKey,
		Limit:         limit,
		Offset:        offset,
		UserID:        userID,
		IsSeed:        seed,
	})
	if err != nil {
		return nil, err
	}
	opinions := make([]oas.OpinionWithReplyAndVote, 0, len(out.Opinions))
	for _, opinion := range out.Opinions {
		// Convert SwipeOpinion to OpinionWithReplyAndVote
		var myVoteType oas.OptNilOpinionWithReplyAndVoteMyVoteType
		if opinion.GetMyVoteType() != nil {
			myVoteType = oas.OptNilOpinionWithReplyAndVoteMyVoteType{
				Value: oas.OpinionWithReplyAndVoteMyVoteType(*opinion.GetMyVoteType()),
				Set:   true,
				Null:  false,
			}
		}

		opinions = append(opinions, oas.OpinionWithReplyAndVote{
			User:       opinion.User.ToResponse(),
			Opinion:    opinion.Opinion.ToResponse(),
			ReplyCount: opinion.ReplyCount,
			MyVoteType: myVoteType,
		})
	}

	return &oas.GetOpinionsForTalkSessionOK{
		Opinions: opinions,
		Pagination: oas.GetOpinionsForTalkSessionOKPagination{
			TotalCount: out.TotalCount,
		},
	}, nil
}

// SwipeOpinions スワイプ用の意見取得
// 自分が投稿した意見は取得しない
func (o *opinionHandler) SwipeOpinions(ctx context.Context, params oas.SwipeOpinionsParams) (oas.SwipeOpinionsRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "opinionHandler.SwipeOpinions")
	defer span.End()

	authCtx, err := requireAuthentication(o.authService, ctx)
	if err != nil {
		return nil, err
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
		UserID:        authCtx.UserID,
		TalkSessionID: talkSessionID,
		Limit:         limit,
	})
	if err != nil {
		return nil, err
	}

	var ress []oas.OpinionWithReplyCount
	for _, opinion := range opinions.Opinions {
		// Convert SwipeOpinion to OpinionWithReplyCount
		ress = append(ress, oas.OpinionWithReplyCount{
			User:       opinion.User.ToResponse(),
			Opinion:    opinion.Opinion.ToResponse(),
			ReplyCount: opinion.ReplyCount,
		})
	}

	return &oas.SwipeOpinionsOK{
		Opinions:       ress,
		RemainingCount: opinions.RemainingOpinions,
	}, nil
}

// PostOpinionPost2 TalkSessionIDをBodyで受け取るタイプのやつ
func (o *opinionHandler) PostOpinionPost2(ctx context.Context, req *oas.PostOpinionPost2Req) (oas.PostOpinionPost2Res, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "opinionHandler.PostOpinionPost2")
	defer span.End()

	authCtx, err := requireAuthentication(o.authService, ctx)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, messages.RequiredParameterError
	}

	value := req
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
		file, err = http_utils.CreateFileHeader(ctx, value.Picture.Value.File, value.GetPicture().Value.Name)
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
	var isSeed bool
	if value.IsSeed.IsSet() {
		isSeed = value.IsSeed.Value
	} else {
		isSeed = false
	}

	if err = o.submitOpinionCommand.Execute(ctx, opinion_usecase.SubmitOpinionInput{
		TalkSessionID:   talkSessionID,
		UserID:          authCtx.UserID,
		ParentOpinionID: parentOpinionID,
		Title:           utils.ToPtrIfNotNullValue(!req.Title.IsSet(), value.Title.Value),
		Content:         req.OpinionContent,
		ReferenceURL:    utils.ToPtrIfNotNullValue(!req.ReferenceURL.IsSet(), value.ReferenceURL.Value),
		Picture:         file,
		IsSeed:          isSeed,
	}); err != nil {
		return nil, err
	}

	res := &oas.PostOpinionPost2OK{}
	return res, nil
}

// ReportOpinion 意見を通報する
func (o *opinionHandler) ReportOpinion(ctx context.Context, req *oas.ReportOpinionReq, params oas.ReportOpinionParams) (oas.ReportOpinionRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "opinionHandler.ReportOpinion")
	defer span.End()

	authCtx, err := requireAuthentication(o.authService, ctx)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, messages.RequiredParameterError
	}

	opinionID, err := shared.ParseUUID[opinion.Opinion](params.OpinionID)
	if err != nil {
		return nil, messages.BadRequestError
	}
	var reasonText *string
	if req.Content.IsSet() && req.Content.Value != "" {
		reasonText = lo.ToPtr(req.Content.Value)
	}

	if err := o.reportOpinionCommand.Execute(ctx, opinion_usecase.ReportOpinionInput{
		ReporterID: authCtx.UserID,
		OpinionID:  opinionID,
		Reason:     int32(req.Reason.Value),
		ReasonText: reasonText,
	}); err != nil {
		utils.HandleError(ctx, err, "reportOpinionCommand.Execute")
		return nil, err
	}

	res := &oas.ReportOpinionOK{}
	return res, nil
}

// GetOpinionReportReasons 通報理由一覧取得
func (o *opinionHandler) GetOpinionReportReasons(ctx context.Context) (oas.GetOpinionReportReasonsRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "opinionHandler.GetOpinionReportReasons")
	defer span.End()

	reasons, err := o.getReportReasons.Execute(ctx)
	if err != nil {
		return nil, err
	}

	var res oas.GetOpinionReportReasonsOKApplicationJSON
	for _, reason := range reasons {
		res = append(res, oas.ReportReason{
			ReasonID: reason.ReasonID,
			Reason:   reason.Reason,
		})
	}

	return &res, nil
}

// GetOpinionAnalysis 意見の集計結果取得
func (o *opinionHandler) GetOpinionAnalysis(ctx context.Context, params oas.GetOpinionAnalysisParams) (oas.GetOpinionAnalysisRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "opinionHandler.GetOpinionAnalysis")
	defer span.End()

	opinionID, err := shared.ParseUUID[opinion.Opinion](params.OpinionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	out, err := o.getOpinionGroupRatio.Execute(ctx, opinion_query.GetOpinionGroupRatioInput{
		OpinionID: opinionID,
	})
	if err != nil {
		return nil, err
	}

	var res oas.GetOpinionAnalysisOKApplicationJSON
	for _, r := range out {
		res = append(res, r.ToResponse())
	}

	return &res, nil
}

// GetOpinionReports implements oas.OpinionHandler.
func (o *opinionHandler) GetOpinionReports(ctx context.Context, params oas.GetOpinionReportsParams) (oas.GetOpinionReportsRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "opinionHandler.GetOpinionReports")
	defer span.End()

	authCtx, err := requireAuthentication(o.authService, ctx)
	if err != nil {
		return nil, err
	}

	opinionID, err := shared.ParseUUID[opinion.Opinion](params.OpinionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	reports, err := o.getReportByOpinionID.Execute(ctx, report_query.GetOpinionReportInput{
		OpinionID: opinionID,
		UserID:    authCtx.UserID,
	})
	if err != nil {
		return nil, err
	}
	reportDetail := reports.Report.ToResponse()
	return &reportDetail, nil
}

// SolveOpinionReport 通報を解決する
func (o *opinionHandler) SolveOpinionReport(ctx context.Context, req *oas.SolveOpinionReportReq, params oas.SolveOpinionReportParams) (oas.SolveOpinionReportRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "opinionHandler.SolveOpinionReport")
	defer span.End()

	authCtx, err := requireAuthentication(o.authService, ctx)
	if err != nil {
		return nil, err
	}

	opinionID, err := shared.ParseUUID[opinion.Opinion](params.OpinionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	if req == nil {
		return nil, messages.RequiredParameterError
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}
	status, err := opinion.NewStatus(string(req.GetAction()))
	if err != nil {
		return nil, messages.BadRequestError
	}

	if err := o.solveReportCommand.Execute(ctx, report_usecase.SolveReportInput{
		OpinionID: opinionID,
		UserID:    authCtx.UserID,
		Status:    status,
	}); err != nil {
		return nil, err
	}

	res := &oas.SolveOpinionReportOK{}
	return res, nil
}
