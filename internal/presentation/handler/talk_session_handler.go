package handler

import (
	"bytes"
	"context"
	"strings"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/application/query/analysis_query"
	"github.com/neko-dream/server/internal/application/query/report_query"
	talksession_query "github.com/neko-dream/server/internal/application/query/talksession"
	"github.com/neko-dream/server/internal/application/usecase/talksession_usecase"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/neko-dream/server/pkg/sort"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type talkSessionHandler struct {
	browseTalkSessionsQuery       talksession_query.BrowseTalkSessionQuery
	browseOpenedByUserQuery       talksession_query.BrowseOpenedByUserQuery
	getConclusionByIDQuery        talksession_query.GetConclusionByIDQuery
	getTalkSessionDetailByIDQuery talksession_query.GetTalkSessionDetailByIDQuery
	getRestrictions               talksession_query.GetRestrictionsQuery
	getAnalysisResultQuery        analysis_query.GetAnalysisResult
	getReportQuery                analysis_query.GetReportQuery
	isSatisfied                   talksession_query.IsTalkSessionSatisfiedQuery
	getReports                    report_query.GetByTalkSessionQuery
	getReportCount                report_query.GetCountQuery
	hasConsent                    talksession_query.HasConsentQuery

	addConclusionCommand    talksession_usecase.AddConclusionCommand
	startTalkSessionCommand talksession_usecase.StartTalkSessionUseCase
	editTalkSessionCommand  talksession_usecase.EditTalkSessionUseCase
	takeConsentCommand      talksession_usecase.TakeConsentUseCase

	authorizationService service.AuthorizationService
	session.TokenManager
}

func NewTalkSessionHandler(
	browseTalkSessionsQuery talksession_query.BrowseTalkSessionQuery,
	browseOpenedByUserQuery talksession_query.BrowseOpenedByUserQuery,
	getConclusionByIDQuery talksession_query.GetConclusionByIDQuery,
	getTalkSessionDetailByIDQuery talksession_query.GetTalkSessionDetailByIDQuery,
	getRestrictionsQuery talksession_query.GetRestrictionsQuery,
	getAnalysisQuery analysis_query.GetAnalysisResult,
	getReportQuery analysis_query.GetReportQuery,
	isSatisfied talksession_query.IsTalkSessionSatisfiedQuery,
	getReports report_query.GetByTalkSessionQuery,
	getReportCount report_query.GetCountQuery,
	hasConsent talksession_query.HasConsentQuery,

	AddConclusionCommand talksession_usecase.AddConclusionCommand,
	startTalkSessionCommand talksession_usecase.StartTalkSessionUseCase,
	editTalkSessionCommand talksession_usecase.EditTalkSessionUseCase,
	takeConsentCommand talksession_usecase.TakeConsentUseCase,

	authorizationService service.AuthorizationService,
	tokenManager session.TokenManager,
) oas.TalkSessionHandler {
	return &talkSessionHandler{
		browseTalkSessionsQuery:       browseTalkSessionsQuery,
		browseOpenedByUserQuery:       browseOpenedByUserQuery,
		getConclusionByIDQuery:        getConclusionByIDQuery,
		getTalkSessionDetailByIDQuery: getTalkSessionDetailByIDQuery,
		getRestrictions:               getRestrictionsQuery,
		getAnalysisResultQuery:        getAnalysisQuery,
		getReportQuery:                getReportQuery,
		isSatisfied:                   isSatisfied,
		getReports:                    getReports,
		getReportCount:                getReportCount,
		hasConsent:                    hasConsent,

		addConclusionCommand:    AddConclusionCommand,
		startTalkSessionCommand: startTalkSessionCommand,
		editTalkSessionCommand:  editTalkSessionCommand,
		takeConsentCommand:      takeConsentCommand,

		authorizationService: authorizationService,
		TokenManager:         tokenManager,
	}
}

// GetUserTalkSessions implements oas.TalkSessionHandler.
func (t *talkSessionHandler) GetUserTalkSessions(ctx context.Context, params oas.GetUserTalkSessionsParams) (oas.GetUserTalkSessionsRes, error) {
	panic("unimplemented")
}

// PostConclusion implements oas.TalkSessionHandler.
func (t *talkSessionHandler) PostConclusion(ctx context.Context, req *oas.PostConclusionReq, params oas.PostConclusionParams) (oas.PostConclusionRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.PostConclusion")
	defer span.End()

	authCtx, err := t.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	if err := t.addConclusionCommand.Execute(ctx, talksession_usecase.AddConclusionCommandInput{
		TalkSessionID: talkSessionID,
		UserID:        authCtx.UserID,
		Conclusion:    req.Content,
	}); err != nil {
		return nil, errtrace.Wrap(err)
	}

	res, err := t.getConclusionByIDQuery.Execute(ctx, talksession_query.GetConclusionByIDQueryRequest{
		TalkSessionID: talkSessionID,
	})
	if err != nil {
		return nil, err
	}

	return &oas.Conclusion{
		User:    oas.ConclusionUser(res.ToResponse()),
		Content: res.Content,
	}, nil
}

// GetConclusion implements oas.TalkSessionHandler.
func (t *talkSessionHandler) GetConclusion(ctx context.Context, params oas.GetConclusionParams) (oas.GetConclusionRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.GetConclusion")
	defer span.End()

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	res, err := t.getConclusionByIDQuery.Execute(ctx, talksession_query.GetConclusionByIDQueryRequest{
		TalkSessionID: talkSessionID,
	})
	if err != nil {
		return nil, err
	}
	// まだ結論が出ていない場合はエラーを返す
	if res == nil {
		return nil, messages.TalkSessionConclusionNotSet
	}

	return &oas.Conclusion{
		User:    oas.ConclusionUser(res.ToResponse()),
		Content: res.Content,
	}, nil
}

// GetOpenedTalkSession implements oas.TalkSessionHandler.
func (t *talkSessionHandler) GetOpenedTalkSession(ctx context.Context, params oas.GetOpenedTalkSessionParams) (oas.GetOpenedTalkSessionRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.GetOpenedTalkSession")
	defer span.End()

	authCtx, err := t.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	var limit, offset *int
	if params.Limit.IsSet() {
		limit = &params.Limit.Value
	}
	if params.Offset.IsSet() {
		offset = &params.Offset.Value
	}

	status := ""
	if params.Status.IsSet() {
		bytes, err := params.Status.Value.MarshalText()
		if err == nil {
			status = string(bytes)
		}
	}
	out, err := t.browseOpenedByUserQuery.Execute(ctx, talksession_query.BrowseOpenedByUserInput{
		UserID: authCtx.UserID,
		Limit:  limit,
		Offset: offset,
		Status: talksession_query.Status(status),
	})
	if err != nil {
		return nil, err
	}

	resultTalkSession := make([]oas.GetOpenedTalkSessionOKTalkSessionsItem, 0, len(out.TalkSessions))
	for _, talkSession := range out.TalkSessions {
		resultTalkSession = append(resultTalkSession, oas.GetOpenedTalkSessionOKTalkSessionsItem{
			TalkSession:  talkSession.ToResponse(),
			OpinionCount: talkSession.OpinionCount,
		})
	}

	return &oas.GetOpenedTalkSessionOK{
		TalkSessions: resultTalkSession,
	}, nil
}

// GetTalkSessionReport implements oas.TalkSessionHandler.
func (t *talkSessionHandler) GetTalkSessionReport(ctx context.Context, params oas.GetTalkSessionReportParams) (oas.GetTalkSessionReportRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.GetTalkSessionReport")
	defer span.End()

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	out, err := t.getReportQuery.Execute(ctx, analysis_query.GetReportInput{
		TalkSessionID: talkSessionID,
	})
	if err != nil {
		return nil, err
	}

	var report oas.OptNilString
	if out.Report != nil {
		report = oas.OptNilString{
			Value: bytes.NewBufferString(*out.Report).String(),
			Set:   true,
		}
	}

	return &oas.GetTalkSessionReportOK{
		Report: report,
	}, nil
}

// InitiateTalkSession トークセッション作成
func (t *talkSessionHandler) InitiateTalkSession(ctx context.Context, req *oas.InitiateTalkSessionReq) (oas.InitiateTalkSessionRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.InitiateTalkSession")
	defer span.End()

	if req == nil {
		return nil, messages.RequiredParameterError
	}

	authCtx, err := t.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}
	userID := authCtx.UserID
	var restrictionStrings []string
	if len(req.Restrictions) > 0 && req.Restrictions[0] != "" {
		if sl := strings.Split(strings.Join(req.Restrictions, ","), ","); len(sl) > 0 {
			restrictionStrings = sl
		}
	}

	var organizationAliasID *shared.UUID[organization.OrganizationAlias]
	if req.GetAliasId().IsSet() && req.GetAliasId().Value != "" {
		aliasID, err := shared.ParseUUID[organization.OrganizationAlias](req.GetAliasId().Value)
		if err == nil {
			organizationAliasID = &aliasID
		}
	}

	out, err := t.startTalkSessionCommand.Execute(ctx, talksession_usecase.StartTalkSessionUseCaseInput{
		Theme:               req.Theme,
		Description:         utils.ToPtrIfNotNullValue(!req.Description.IsSet(), req.Description.Value),
		ThumbnailURL:        utils.ToPtrIfNotNullValue(!req.ThumbnailURL.IsSet(), req.ThumbnailURL.Value),
		OwnerID:             userID,
		ScheduledEndTime:    req.ScheduledEndTime,
		Latitude:            utils.ToPtrIfNotNullValue(!req.Latitude.IsSet(), req.Latitude.Value),
		Longitude:           utils.ToPtrIfNotNullValue(!req.Longitude.IsSet(), req.Longitude.Value),
		City:                utils.ToPtrIfNotNullValue(!req.City.IsSet(), req.City.Value),
		Prefecture:          utils.ToPtrIfNotNullValue(!req.Prefecture.IsSet(), req.Prefecture.Value),
		Restrictions:        restrictionStrings,
		SessionClaim:        session.GetSession(ctx),
		OrganizationAliasID: organizationAliasID,
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	res := out.ToResponse()
	return &res, nil
}

// ViewTalkSessionDetail トークセッション詳細取得
func (t *talkSessionHandler) GetTalkSessionDetail(ctx context.Context, params oas.GetTalkSessionDetailParams) (oas.GetTalkSessionDetailRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.GetTalkSessionDetail")
	defer span.End()

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	out, err := t.getTalkSessionDetailByIDQuery.Execute(ctx, talksession_query.GetTalkSessionDetailInput{
		TalkSessionID: talkSessionID,
	})
	if err != nil {
		return nil, err
	}

	res := out.ToResponse()
	return &res, nil
}

// GetTalkSessionList セッション一覧取得
func (t *talkSessionHandler) GetTalkSessionList(ctx context.Context, params oas.GetTalkSessionListParams) (oas.GetTalkSessionListRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.GetTalkSessionList")
	defer span.End()

	var limit, offset *int
	if params.Limit.IsSet() {
		limit = &params.Limit.Value
	}
	if params.Offset.IsSet() {
		offset = &params.Offset.Value
	}

	var status *talksession_query.Status
	if params.Status.IsSet() {
		bytes, err := params.Status.Value.MarshalText()
		if err == nil {
			status = lo.ToPtr(talksession_query.Status(string(bytes)))
		}
	}
	var sortKey sort.SortKey
	if params.SortKey.IsSet() {
		if bytes, err := params.SortKey.Value.MarshalText(); err == nil {
			sortKey = sort.SortKey(string(bytes))
		}
	}
	var latitude, longitude *float64
	if params.Latitude.IsSet() {
		latitude = utils.ToPtrIfNotNullValue(params.Latitude.Null, params.Latitude.Value)
	}
	if params.Longitude.IsSet() {
		longitude = utils.ToPtrIfNotNullValue(params.Longitude.Null, params.Longitude.Value)
	}

	theme := utils.ToPtrIfNotNullValue(params.Theme.Null, params.Theme.Value)
	out, err := t.browseTalkSessionsQuery.Execute(ctx, talksession_query.BrowseTalkSessionQueryInput{
		Limit:     limit,
		Offset:    offset,
		Theme:     theme,
		Status:    status,
		SortKey:   sortKey,
		Latitude:  latitude,
		Longitude: longitude,
	})
	if err != nil {
		return nil, err
	}

	resultTalkSession := make([]oas.GetTalkSessionListOKTalkSessionsItem, 0, len(out.TalkSessions))
	for _, talkSession := range out.TalkSessions {
		res := oas.GetTalkSessionListOKTalkSessionsItem{
			TalkSession:  talkSession.ToResponse(),
			OpinionCount: talkSession.OpinionCount,
		}
		resultTalkSession = append(resultTalkSession, res)
	}

	return &oas.GetTalkSessionListOK{
		TalkSessions: resultTalkSession,
		Pagination: oas.OffsetPagination{
			TotalCount: out.TotalCount,
			Limit:      out.Limit,
			Offset:     out.Offset,
		},
	}, nil
}

// TalkSessionAnalysis 分析結果取得
func (t *talkSessionHandler) TalkSessionAnalysis(ctx context.Context, params oas.TalkSessionAnalysisParams) (oas.TalkSessionAnalysisRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.TalkSessionAnalysis")
	defer span.End()

	userID, err := t.authorizationService.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	out, err := t.getAnalysisResultQuery.Execute(ctx, analysis_query.GetAnalysisResultInput{
		UserID:        userID,
		TalkSessionID: talkSessionID,
	})
	if err != nil {
		return nil, err
	}

	var myPosition oas.OptUserGroupPosition
	if out.MyPosition != nil {
		myPosition = oas.OptUserGroupPosition{
			Value: out.MyPosition.ToResponse(),
			Set:   true,
		}
	}

	positions := make([]oas.UserGroupPosition, 0, len(out.Positions))
	for _, position := range out.Positions {
		positions = append(positions, position.ToResponse())
	}

	groupOpinions := make([]oas.TalkSessionAnalysisOKGroupOpinionsItem, 0, len(out.GroupOpinions))
	for _, groupOpinion := range out.GroupOpinions {
		opinions := make([]oas.TalkSessionAnalysisOKGroupOpinionsItemOpinionsItem, 0, len(groupOpinion.Opinions))
		for _, opinion := range groupOpinion.Opinions {
			// Convert OpinionWithRepresentative to TalkSessionAnalysisOKGroupOpinionsItemOpinionsItem
			opinions = append(opinions, oas.TalkSessionAnalysisOKGroupOpinionsItemOpinionsItem{
				Opinion: opinion.Opinion.ToResponse(),
				User: oas.User{
					DisplayID:   opinion.User.DisplayID,
					DisplayName: opinion.User.DisplayName,
					IconURL:     utils.ToOptNil[oas.OptNilString](opinion.User.IconURL),
				},
				AgreeCount:    opinion.AgreeCount,
				DisagreeCount: opinion.DisagreeCount,
				PassCount:     opinion.PassCount,
			})
		}
		groupOpinions = append(groupOpinions, oas.TalkSessionAnalysisOKGroupOpinionsItem{
			GroupName: groupOpinion.GroupName,
			GroupID:   groupOpinion.GroupID,
			Opinions:  opinions,
		})
	}

	return &oas.TalkSessionAnalysisOK{
		MyPosition:    myPosition,
		Positions:     positions,
		GroupOpinions: groupOpinions,
	}, nil
}

// EditTalkSession implements oas.TalkSessionHandler.
func (t *talkSessionHandler) EditTalkSession(ctx context.Context, req *oas.EditTalkSessionReq, params oas.EditTalkSessionParams) (oas.EditTalkSessionRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.EditTalkSession")
	defer span.End()

	authCtx, err := t.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	if req == nil {
		return nil, messages.RequiredParameterError
	}

	_, err = t.editTalkSessionCommand.Execute(ctx, talksession_usecase.EditTalkSessionInput{
		TalkSessionID:    talkSessionID,
		UserID:           authCtx.UserID,
		Theme:            req.Theme,
		Description:      utils.ToPtrIfNotNullValue(!req.Description.IsSet(), req.Description.Value),
		ThumbnailURL:     utils.ToPtrIfNotNullValue(!req.ThumbnailURL.IsSet(), req.ThumbnailURL.Value),
		ScheduledEndTime: req.ScheduledEndTime,
		Latitude:         utils.ToPtrIfNotNullValue(!req.Latitude.IsSet(), req.Latitude.Value),
		Longitude:        utils.ToPtrIfNotNullValue(!req.Longitude.IsSet(), req.Longitude.Value),
		City:             utils.ToPtrIfNotNullValue(!req.City.IsSet(), req.City.Value),
		Prefecture:       utils.ToPtrIfNotNullValue(!req.Prefecture.IsSet(), req.Prefecture.Value),
	})
	if err != nil {
		return nil, err
	}
	talkSessionDetail, err := t.getTalkSessionDetailByIDQuery.Execute(ctx, talksession_query.GetTalkSessionDetailInput{
		TalkSessionID: talkSessionID,
	})
	if err != nil {
		return nil, err
	}
	res := talkSessionDetail.ToResponse()
	return &res, nil
}

// GetTalkSessionRestrictionKeys implements oas.TalkSessionHandler.
func (t *talkSessionHandler) GetTalkSessionRestrictionKeys(ctx context.Context) (oas.GetTalkSessionRestrictionKeysRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.GetTalkSessionRestrictionKeys")
	defer span.End()

	out, err := t.getRestrictions.Execute(ctx)
	if err != nil {
		return nil, err
	}

	keys := make([]oas.Restriction, 0, len(out.Restrictions))
	for _, restriction := range out.Restrictions {
		keys = append(keys, oas.Restriction{
			Key:         string(restriction.Key),
			Description: restriction.Description,
			DependsOn:   lo.Map(restriction.DependsOn, func(item talksession.RestrictionAttributeKey, _ int) string { return string(item) }),
		})
	}

	res := oas.GetTalkSessionRestrictionKeysOKApplicationJSON(keys)
	return &res, nil
}

// GetTalkSessionRestrictionSatisfied implements oas.TalkSessionHandler.
func (t *talkSessionHandler) GetTalkSessionRestrictionSatisfied(ctx context.Context, params oas.GetTalkSessionRestrictionSatisfiedParams) (oas.GetTalkSessionRestrictionSatisfiedRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.GetTalkSessionRestrictionSatisfied")
	defer span.End()

	authCtx, err := t.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	out, err := t.isSatisfied.Execute(ctx, talksession_query.IsTalkSessionSatisfiedInput{
		TalkSessionID: talkSessionID,
		UserID:        authCtx.UserID,
	})
	if err != nil {
		return nil, err
	}

	attributes := make([]oas.Restriction, 0, len(out.Attributes))
	for _, attribute := range out.Attributes {
		attributes = append(attributes, oas.Restriction{
			Key:         string(attribute.Key),
			Description: attribute.Description,
		})
	}

	res := oas.GetTalkSessionRestrictionSatisfiedOKApplicationJSON(attributes)
	return &res, nil
}

// GetReportsForTalkSession implements oas.TalkSessionHandler.
func (t *talkSessionHandler) GetReportsForTalkSession(ctx context.Context, params oas.GetReportsForTalkSessionParams) (oas.GetReportsForTalkSessionRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.GetReportsForTalkSession")
	defer span.End()

	authCtx, err := t.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}
	var status string
	if params.Status.IsSet() {
		bytes, err := params.Status.Value.MarshalText()
		if err == nil {
			status = string(bytes)
		}
	} else {
		status = "unsolved"
	}

	out, err := t.getReports.Execute(ctx, report_query.GetByTalkSessionInput{
		TalkSessionID: talkSessionID,
		UserID:        authCtx.UserID,
		Status:        status,
	})
	if err != nil {
		return nil, err
	}

	reports := make([]oas.ReportDetail, 0, len(out.Reports))
	for _, report := range out.Reports {
		reports = append(reports, report.ToResponse())
	}

	return &oas.GetReportsForTalkSessionOK{
		Reports: reports,
	}, nil
}

// GetTalkSessionReportCount implements oas.TalkSessionHandler.
func (t *talkSessionHandler) GetTalkSessionReportCount(ctx context.Context, params oas.GetTalkSessionReportCountParams) (oas.GetTalkSessionReportCountRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.GetTalkSessionReportCount")
	defer span.End()

	authCtx, err := t.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	out, err := t.getReportCount.Execute(ctx, report_query.GetCountInput{
		TalkSessionID: talkSessionID,
		UserID:        authCtx.UserID,
		Status:        string(params.Status),
	})
	if err != nil {
		return nil, err
	}

	return &oas.GetTalkSessionReportCountOK{
		Count: out.Count,
	}, nil
}

// ConsentTalkSession implements oas.TalkSessionHandler.
func (t *talkSessionHandler) ConsentTalkSession(ctx context.Context, params oas.ConsentTalkSessionParams) (oas.ConsentTalkSessionRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.ConsentTalkSession")
	defer span.End()

	authCtx, err := t.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	if err := t.takeConsentCommand.Execute(ctx, talksession_usecase.TakeConsentUseCaseInput{
		TalkSessionID: talkSessionID,
		UserID:        authCtx.UserID,
	}); err != nil {
		return nil, errtrace.Wrap(err)
	}

	res := &oas.ConsentTalkSessionOK{}
	return res, nil
}

// HasConsent implements oas.TalkSessionHandler.
func (t *talkSessionHandler) HasConsent(ctx context.Context, params oas.HasConsentParams) (oas.HasConsentRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.HasConsent")
	defer span.End()

	authCtx, err := t.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	hasConsent, err := t.hasConsent.Execute(ctx, talksession_query.HasConsentQueryInput{
		TalkSessionID: talkSessionID,
		UserID:        authCtx.UserID,
	})
	if err != nil {
		return nil, err
	}

	return &oas.HasConsentOK{
		HasConsent: hasConsent,
	}, nil
}
