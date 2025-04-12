package handler

import (
	"bytes"
	"context"
	"strings"
	"time"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/neko-dream/server/internal/usecase/command/talksession_command"
	"github.com/neko-dream/server/internal/usecase/query/analysis_query"
	"github.com/neko-dream/server/internal/usecase/query/report_query"
	talksession_query "github.com/neko-dream/server/internal/usecase/query/talksession"
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

	addConclusionCommand    talksession_command.AddConclusionCommand
	startTalkSessionCommand talksession_command.StartTalkSessionCommand
	editTalkSessionCommand  talksession_command.EditCommand

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

	AddConclusionCommand talksession_command.AddConclusionCommand,
	startTalkSessionCommand talksession_command.StartTalkSessionCommand,
	editTalkSessionCommand talksession_command.EditCommand,

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

		addConclusionCommand:    AddConclusionCommand,
		startTalkSessionCommand: startTalkSessionCommand,
		editTalkSessionCommand:  editTalkSessionCommand,

		TokenManager: tokenManager,
	}
}

// PostConclusion implements oas.TalkSessionHandler.
func (t *talkSessionHandler) PostConclusion(ctx context.Context, req oas.OptPostConclusionReq, params oas.PostConclusionParams) (oas.PostConclusionRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.PostConclusion")
	defer span.End()

	claim := session.GetSession(t.SetSession(ctx))
	if !req.IsSet() {
		return nil, messages.RequiredParameterError
	}

	userID, err := claim.UserID()
	if err != nil {
		utils.HandleError(ctx, err, "claim.UserID")
		return nil, messages.ForbiddenError
	}

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	if err := t.addConclusionCommand.Execute(ctx, talksession_command.AddConclusionCommandInput{
		TalkSessionID: talkSessionID,
		UserID:        userID,
		Conclusion:    req.Value.Content,
	}); err != nil {
		return nil, errtrace.Wrap(err)
	}

	res, err := t.getConclusionByIDQuery.Execute(ctx, talksession_query.GetConclusionByIDQueryRequest{
		TalkSessionID: talkSessionID,
	})
	if err != nil {
		return nil, err
	}

	return &oas.PostConclusionOK{
		User: oas.PostConclusionOKUser{
			DisplayID:   res.DisplayID,
			DisplayName: res.DisplayName,
			IconURL:     utils.ToOptNil[oas.OptNilString](res.IconURL),
		},
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

	return &oas.GetConclusionOK{
		User: oas.GetConclusionOKUser{
			DisplayID:   res.DisplayID,
			DisplayName: res.DisplayName,
			IconURL:     utils.ToOptNil[oas.OptNilString](res.IconURL),
		},
		Content: res.Content,
	}, nil
}

// GetOpenedTalkSession implements oas.TalkSessionHandler.
func (t *talkSessionHandler) GetOpenedTalkSession(ctx context.Context, params oas.GetOpenedTalkSessionParams) (oas.GetOpenedTalkSessionRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.GetOpenedTalkSession")
	defer span.End()

	claim := session.GetSession(t.SetSession(ctx))
	if claim == nil {
		return nil, messages.ForbiddenError
	}
	userID, err := claim.UserID()
	if err != nil {
		utils.HandleError(ctx, err, "claim.UserID")
		return nil, messages.ForbiddenError
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
		UserID: userID,
		Limit:  limit,
		Offset: offset,
		Status: talksession_query.Status(status),
	})
	if err != nil {
		return nil, err
	}

	resultTalkSession := make([]oas.GetOpenedTalkSessionOKTalkSessionsItem, 0, len(out.TalkSessions))
	for _, talkSession := range out.TalkSessions {
		owner := oas.GetOpenedTalkSessionOKTalkSessionsItemTalkSessionOwner{
			DisplayID:   talkSession.User.DisplayID,
			DisplayName: talkSession.User.DisplayName,
			IconURL:     utils.ToOptNil[oas.OptNilString](talkSession.User.IconURL),
		}
		var location oas.OptGetOpenedTalkSessionOKTalkSessionsItemTalkSessionLocation
		if talkSession.HasLocation() {
			location = oas.OptGetOpenedTalkSessionOKTalkSessionsItemTalkSessionLocation{
				Value: oas.GetOpenedTalkSessionOKTalkSessionsItemTalkSessionLocation{
					Latitude:  utils.ToOpt[oas.OptFloat64](talkSession.Latitude),
					Longitude: utils.ToOpt[oas.OptFloat64](talkSession.Longitude),
				},
				Set: true,
			}
		} else {
			location = oas.OptGetOpenedTalkSessionOKTalkSessionsItemTalkSessionLocation{}
			location.Set = false
		}
		var restrictions []oas.GetOpenedTalkSessionOKTalkSessionsItemTalkSessionRestrictionsItem
		for _, restriction := range talkSession.TalkSession.Restrictions {
			res := talksession.RestrictionAttributeKey(restriction)
			attr := res.RestrictionAttribute()
			restrictions = append(restrictions, oas.GetOpenedTalkSessionOKTalkSessionsItemTalkSessionRestrictionsItem{
				Key:         string(attr.Key),
				Description: attr.Description,
			})
		}

		resultTalkSession = append(resultTalkSession, oas.GetOpenedTalkSessionOKTalkSessionsItem{
			TalkSession: oas.GetOpenedTalkSessionOKTalkSessionsItemTalkSession{
				ID:               talkSession.TalkSessionID.String(),
				Theme:            talkSession.Theme,
				Owner:            owner,
				Description:      utils.ToOptNil[oas.OptNilString](talkSession.Description),
				ThumbnailURL:     utils.ToOptNil[oas.OptNilString](talkSession.ThumbnailURL),
				Location:         location,
				CreatedAt:        talkSession.CreatedAt.Format(time.RFC3339),
				ScheduledEndTime: talkSession.ScheduledEndTime.Format(time.RFC3339),
				City:             utils.ToOptNil[oas.OptNilString](talkSession.City),
				Prefecture:       utils.ToOptNil[oas.OptNilString](talkSession.Prefecture),
				Restrictions:     restrictions,
			},
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

	return &oas.GetTalkSessionReportOK{
		Report: bytes.NewBufferString(out.Report).String(),
	}, nil
}

// CreateTalkSession トークセッション作成
func (t *talkSessionHandler) CreateTalkSession(ctx context.Context, req oas.OptCreateTalkSessionReq) (oas.CreateTalkSessionRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.CreateTalkSession")
	defer span.End()

	claim := session.GetSession(ctx)
	if !req.IsSet() {
		return nil, messages.RequiredParameterError
	}

	userID, err := claim.UserID()
	if err != nil {
		utils.HandleError(ctx, err, "claim.UserID")
		return nil, messages.ForbiddenError
	}
	var restrictionStrings []string
	if req.Value.Restrictions != nil {
		if sl := strings.Split(strings.Join(req.Value.Restrictions, ","), ","); len(sl) > 0 {
			restrictionStrings = sl
		}
	}

	out, err := t.startTalkSessionCommand.Execute(ctx, talksession_command.StartTalkSessionCommandInput{
		Theme:            req.Value.Theme,
		Description:      utils.ToPtrIfNotNullValue(!req.Value.Description.IsSet(), req.Value.Description.Value),
		ThumbnailURL:     utils.ToPtrIfNotNullValue(!req.Value.ThumbnailURL.IsSet(), req.Value.ThumbnailURL.Value),
		OwnerID:          userID,
		ScheduledEndTime: req.Value.ScheduledEndTime,
		Latitude:         utils.ToPtrIfNotNullValue(!req.Value.Latitude.IsSet(), req.Value.Latitude.Value),
		Longitude:        utils.ToPtrIfNotNullValue(!req.Value.Longitude.IsSet(), req.Value.Longitude.Value),
		City:             utils.ToPtrIfNotNullValue(!req.Value.City.IsSet(), req.Value.City.Value),
		Prefecture:       utils.ToPtrIfNotNullValue(!req.Value.Prefecture.IsSet(), req.Value.Prefecture.Value),
		Restrictions:     restrictionStrings,
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	var location oas.OptCreateTalkSessionOKLocation
	if out.HasLocation() {
		location = oas.OptCreateTalkSessionOKLocation{
			Value: oas.CreateTalkSessionOKLocation{
				Latitude:  utils.ToOpt[oas.OptFloat64](out.Latitude),
				Longitude: utils.ToOpt[oas.OptFloat64](out.Longitude),
			},
			Set: true,
		}
	}
	var restrictions []oas.CreateTalkSessionOKRestrictionsItem
	for _, restriction := range out.TalkSession.Restrictions {
		res := talksession.RestrictionAttributeKey(restriction)
		attr := res.RestrictionAttribute()
		restrictions = append(restrictions, oas.CreateTalkSessionOKRestrictionsItem{
			Key:         string(attr.Key),
			Description: attr.Description,
		})
	}

	res := &oas.CreateTalkSessionOK{
		ID: out.TalkSession.TalkSessionID.String(),
		Owner: oas.CreateTalkSessionOKOwner{
			DisplayID:   out.User.DisplayID,
			DisplayName: out.User.DisplayName,
			IconURL:     utils.ToOptNil[oas.OptNilString](out.User.IconURL),
		},
		Theme:            out.TalkSession.Theme,
		Description:      utils.ToOptNil[oas.OptNilString](out.TalkSession.Description),
		ThumbnailURL:     utils.ToOptNil[oas.OptNilString](out.TalkSession.ThumbnailURL),
		CreatedAt:        clock.Now(ctx).Format(time.RFC3339),
		ScheduledEndTime: out.TalkSession.ScheduledEndTime.Format(time.RFC3339),
		Location:         location,
		Restrictions:     restrictions,
	}
	return res, nil
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

	owner := oas.GetTalkSessionDetailOKOwner{
		DisplayID:   out.User.DisplayID,
		DisplayName: out.User.DisplayName,
		IconURL:     utils.ToOptNil[oas.OptNilString](out.User.IconURL),
	}
	var location oas.OptGetTalkSessionDetailOKLocation
	if out.HasLocation() {
		location = oas.OptGetTalkSessionDetailOKLocation{
			Value: oas.GetTalkSessionDetailOKLocation{
				Latitude:  utils.ToOpt[oas.OptFloat64](out.Latitude),
				Longitude: utils.ToOpt[oas.OptFloat64](out.Longitude),
			},
			Set: true,
		}
	}
	var restrictions []oas.GetTalkSessionDetailOKRestrictionsItem
	for _, restriction := range out.TalkSession.Restrictions {
		res := talksession.RestrictionAttributeKey(restriction)
		attr := res.RestrictionAttribute()
		restrictions = append(restrictions, oas.GetTalkSessionDetailOKRestrictionsItem{
			Key:         string(attr.Key),
			Description: attr.Description,
		})
	}

	return &oas.GetTalkSessionDetailOK{
		ID:               out.TalkSessionID.String(),
		Theme:            out.Theme,
		Description:      utils.ToOptNil[oas.OptNilString](out.Description),
		ThumbnailURL:     utils.ToOptNil[oas.OptNilString](out.ThumbnailURL),
		Owner:            owner,
		CreatedAt:        out.CreatedAt.Format(time.RFC3339),
		ScheduledEndTime: out.ScheduledEndTime.Format(time.RFC3339),
		Location:         location,
		City:             utils.ToOptNil[oas.OptNilString](out.City),
		Prefecture:       utils.ToOptNil[oas.OptNilString](out.Prefecture),
		Restrictions:     restrictions,
	}, nil
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
		owner := oas.GetTalkSessionListOKTalkSessionsItemTalkSessionOwner{
			DisplayID:   talkSession.User.DisplayID,
			DisplayName: talkSession.User.DisplayName,
			IconURL:     utils.ToOptNil[oas.OptNilString](talkSession.User.IconURL),
		}
		var location oas.OptGetTalkSessionListOKTalkSessionsItemTalkSessionLocation
		if talkSession.HasLocation() {
			location = oas.OptGetTalkSessionListOKTalkSessionsItemTalkSessionLocation{
				Value: oas.GetTalkSessionListOKTalkSessionsItemTalkSessionLocation{
					Latitude:  utils.ToOpt[oas.OptFloat64](talkSession.Latitude),
					Longitude: utils.ToOpt[oas.OptFloat64](talkSession.Longitude),
				},
				Set: true,
			}
		}
		var restrictions []oas.GetTalkSessionListOKTalkSessionsItemTalkSessionRestrictionsItem
		for _, restriction := range talkSession.TalkSession.Restrictions {
			res := talksession.RestrictionAttributeKey(restriction)
			attr := res.RestrictionAttribute()
			restrictions = append(restrictions, oas.GetTalkSessionListOKTalkSessionsItemTalkSessionRestrictionsItem{
				Key:         string(attr.Key),
				Description: attr.Description,
			})
		}

		res := oas.GetTalkSessionListOKTalkSessionsItem{
			TalkSession: oas.GetTalkSessionListOKTalkSessionsItemTalkSession{
				ID:               talkSession.TalkSessionID.String(),
				Theme:            talkSession.Theme,
				Description:      utils.ToOptNil[oas.OptNilString](talkSession.Description),
				ThumbnailURL:     utils.ToOptNil[oas.OptNilString](talkSession.ThumbnailURL),
				Owner:            owner,
				CreatedAt:        talkSession.CreatedAt.Format(time.RFC3339),
				ScheduledEndTime: talkSession.ScheduledEndTime.Format(time.RFC3339),
				Location:         location,
				City:             utils.ToOptNil[oas.OptNilString](talkSession.City),
				Prefecture:       utils.ToOptNil[oas.OptNilString](talkSession.Prefecture),
				Restrictions:     restrictions,
			},
			OpinionCount: talkSession.OpinionCount,
		}
		resultTalkSession = append(resultTalkSession, res)
	}

	return &oas.GetTalkSessionListOK{
		TalkSessions: resultTalkSession,
		Pagination: oas.GetTalkSessionListOKPagination{
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

	claim := session.GetSession(t.SetSession(ctx))
	var userID *shared.UUID[user.User]
	if claim != nil {
		id, err := claim.UserID()
		if err == nil {
			userID = &id
		}
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

	var myPosition oas.OptTalkSessionAnalysisOKMyPosition
	if out.MyPosition != nil {
		myPosition = oas.OptTalkSessionAnalysisOKMyPosition{
			Value: oas.TalkSessionAnalysisOKMyPosition{
				PosX:           out.MyPosition.PosX,
				PosY:           out.MyPosition.PosY,
				DisplayId:      out.MyPosition.DisplayID,
				DisplayName:    out.MyPosition.DisplayName,
				IconURL:        utils.ToOptNil[oas.OptNilString](out.MyPosition.IconURL),
				GroupId:        out.MyPosition.GroupID,
				GroupName:      out.MyPosition.GroupName,
				PerimeterIndex: utils.ToOpt[oas.OptInt](out.MyPosition.PerimeterIndex),
			},
			Set: true,
		}
	}

	positions := make([]oas.TalkSessionAnalysisOKPositionsItem, 0, len(out.Positions))
	for _, position := range out.Positions {
		positions = append(positions, oas.TalkSessionAnalysisOKPositionsItem{
			PosX:           position.PosX,
			PosY:           position.PosY,
			DisplayId:      position.DisplayID,
			DisplayName:    position.DisplayName,
			IconURL:        utils.ToOptNil[oas.OptNilString](position.IconURL),
			GroupName:      position.GroupName,
			GroupId:        position.GroupID,
			PerimeterIndex: utils.ToOpt[oas.OptInt](position.PerimeterIndex),
		})
	}

	groupOpinions := make([]oas.TalkSessionAnalysisOKGroupOpinionsItem, 0, len(out.GroupOpinions))
	for _, groupOpinion := range out.GroupOpinions {
		opinions := make([]oas.TalkSessionAnalysisOKGroupOpinionsItemOpinionsItem, 0, len(groupOpinion.Opinions))
		for _, opinion := range groupOpinion.Opinions {
			opinions = append(opinions, oas.TalkSessionAnalysisOKGroupOpinionsItemOpinionsItem{
				Opinion: oas.TalkSessionAnalysisOKGroupOpinionsItemOpinionsItemOpinion{
					ID:           opinion.Opinion.OpinionID.String(),
					Title:        utils.ToOpt[oas.OptString](opinion.Opinion.Title),
					Content:      opinion.Opinion.Content,
					ParentID:     utils.ToOpt[oas.OptString](opinion.Opinion.ParentOpinionID),
					PictureURL:   utils.ToOptNil[oas.OptNilString](opinion.Opinion.PictureURL),
					ReferenceURL: utils.ToOpt[oas.OptString](opinion.Opinion.ReferenceURL),
					PostedAt:     opinion.Opinion.CreatedAt.Format(time.RFC3339),
					IsDeleted:    opinion.Opinion.IsDeleted,
				},
				User: oas.TalkSessionAnalysisOKGroupOpinionsItemOpinionsItemUser{
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
			GroupId:   groupOpinion.GroupID,
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
func (t *talkSessionHandler) EditTalkSession(ctx context.Context, req oas.OptEditTalkSessionReq, params oas.EditTalkSessionParams) (oas.EditTalkSessionRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.EditTalkSession")
	defer span.End()

	claim := session.GetSession(t.SetSession(ctx))
	if claim == nil {
		return nil, messages.ForbiddenError
	}
	userID, err := claim.UserID()
	if err != nil {
		utils.HandleError(ctx, err, "claim.UserID")
		return nil, messages.ForbiddenError
	}

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	if !req.IsSet() {
		return nil, messages.RequiredParameterError
	}

	var restrictionStrings []string
	if req.Value.Restrictions != nil {
		if sl := strings.Split(strings.Join(req.Value.Restrictions, ","), ","); len(sl) > 0 {
			restrictionStrings = sl
		}
	}

	out, err := t.editTalkSessionCommand.Execute(ctx, talksession_command.EditCommandInput{
		TalkSessionID:    talkSessionID,
		UserID:           userID,
		Theme:            req.Value.Theme,
		Description:      utils.ToPtrIfNotNullValue(!req.Value.Description.IsSet(), req.Value.Description.Value),
		ThumbnailURL:     utils.ToPtrIfNotNullValue(!req.Value.ThumbnailURL.IsSet(), req.Value.ThumbnailURL.Value),
		ScheduledEndTime: req.Value.ScheduledEndTime,
		Latitude:         utils.ToPtrIfNotNullValue(!req.Value.Latitude.IsSet(), req.Value.Latitude.Value),
		Longitude:        utils.ToPtrIfNotNullValue(!req.Value.Longitude.IsSet(), req.Value.Longitude.Value),
		City:             utils.ToPtrIfNotNullValue(!req.Value.City.IsSet(), req.Value.City.Value),
		Prefecture:       utils.ToPtrIfNotNullValue(!req.Value.Prefecture.IsSet(), req.Value.Prefecture.Value),
		Restrictions:     restrictionStrings,
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
	owner := oas.EditTalkSessionOKOwner{
		DisplayID:   out.User.DisplayID,
		DisplayName: out.User.DisplayName,
		IconURL:     utils.ToOptNil[oas.OptNilString](talkSessionDetail.User.IconURL),
	}

	var location oas.OptEditTalkSessionOKLocation
	if out.HasLocation() {
		location = oas.OptEditTalkSessionOKLocation{
			Value: oas.EditTalkSessionOKLocation{
				Latitude:  utils.ToOpt[oas.OptFloat64](talkSessionDetail.Latitude),
				Longitude: utils.ToOpt[oas.OptFloat64](talkSessionDetail.Longitude),
			},
			Set: true,
		}
	}
	var restrictions []oas.EditTalkSessionOKRestrictionsItem
	for _, restriction := range out.TalkSession.Restrictions {
		res := talksession.RestrictionAttributeKey(restriction)
		attr := res.RestrictionAttribute()
		restrictions = append(restrictions, oas.EditTalkSessionOKRestrictionsItem{
			Key:         string(attr.Key),
			Description: attr.Description,
		})
	}
	res := &oas.EditTalkSessionOK{
		ID:               out.TalkSession.TalkSessionID.String(),
		Theme:            out.TalkSession.Theme,
		Description:      utils.ToOptNil[oas.OptNilString](out.TalkSession.Description),
		ThumbnailURL:     utils.ToOptNil[oas.OptNilString](out.TalkSession.ThumbnailURL),
		Owner:            owner,
		CreatedAt:        out.CreatedAt.Format(time.RFC3339),
		ScheduledEndTime: out.TalkSession.ScheduledEndTime.Format(time.RFC3339),
		Location:         location,
		City:             utils.ToOptNil[oas.OptNilString](out.City),
		Prefecture:       utils.ToOptNil[oas.OptNilString](out.Prefecture),
		Restrictions:     restrictions,
	}
	return res, nil
}

// GetTalkSessionRestrictionKeys implements oas.TalkSessionHandler.
func (t *talkSessionHandler) GetTalkSessionRestrictionKeys(ctx context.Context) (oas.GetTalkSessionRestrictionKeysRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.GetTalkSessionRestrictionKeys")
	defer span.End()

	out, err := t.getRestrictions.Execute(ctx)
	if err != nil {
		return nil, err
	}

	keys := make([]oas.GetTalkSessionRestrictionKeysOKItem, 0, len(out.Restrictions))
	for _, restriction := range out.Restrictions {
		keys = append(keys, oas.GetTalkSessionRestrictionKeysOKItem{
			Key:         string(restriction.Key),
			Description: restriction.Description,
		})
	}

	res := oas.GetTalkSessionRestrictionKeysOKApplicationJSON(keys)
	return &res, nil
}

// GetTalkSessionRestrictionSatisfied implements oas.TalkSessionHandler.
func (t *talkSessionHandler) GetTalkSessionRestrictionSatisfied(ctx context.Context, params oas.GetTalkSessionRestrictionSatisfiedParams) (oas.GetTalkSessionRestrictionSatisfiedRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.GetTalkSessionRestrictionSatisfied")
	defer span.End()

	claim := session.GetSession(t.SetSession(ctx))
	var userID *shared.UUID[user.User]
	if claim != nil {
		id, err := claim.UserID()
		if err == nil {
			userID = &id
		}
	}
	if userID == nil {
		return nil, messages.ForbiddenError
	}

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	out, err := t.isSatisfied.Execute(ctx, talksession_query.IsTalkSessionSatisfiedInput{
		TalkSessionID: talkSessionID,
		UserID:        *userID,
	})
	if err != nil {
		return nil, err
	}

	attributes := make([]oas.GetTalkSessionRestrictionSatisfiedOKItem, 0, len(out.Attributes))
	for _, attribute := range out.Attributes {
		attributes = append(attributes, oas.GetTalkSessionRestrictionSatisfiedOKItem{
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

	claim := session.GetSession(t.SetSession(ctx))
	var userID *shared.UUID[user.User]
	if claim != nil {
		id, err := claim.UserID()
		if err == nil {
			userID = &id
		}
	}
	if userID == nil {
		return nil, messages.ForbiddenError
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
		UserID:        *userID,
		Status:        status,
	})
	if err != nil {
		return nil, err
	}

	reports := make([]oas.GetReportsForTalkSessionOKReportsItem, 0, len(out.Reports))
	for _, report := range out.Reports {
		var parentOpinionID oas.OptString
		if report.Opinion.ParentOpinionID != nil {
			parentOpinionID = oas.OptString{
				Value: report.Opinion.ParentOpinionID.String(),
				Set:   true,
			}
		}

		reasons := make([]oas.GetReportsForTalkSessionOKReportsItemReasonsItem, 0, len(report.Reasons))
		for _, reason := range report.Reasons {
			reasons = append(reasons, oas.GetReportsForTalkSessionOKReportsItemReasonsItem{
				Reason:  reason.Reason,
				Content: utils.ToOptNil[oas.OptNilString](reason.Content),
			})
		}

		reports = append(reports, oas.GetReportsForTalkSessionOKReportsItem{
			Opinion: oas.GetReportsForTalkSessionOKReportsItemOpinion{
				ID:           report.Opinion.OpinionID.String(),
				Title:        utils.ToOpt[oas.OptString](report.Opinion.Title),
				Content:      report.Opinion.Content,
				ParentID:     parentOpinionID,
				PictureURL:   utils.ToOptNil[oas.OptNilString](report.Opinion.PictureURL),
				ReferenceURL: utils.ToOpt[oas.OptString](report.Opinion.ReferenceURL),
				PostedAt:     report.Opinion.CreatedAt.Format(time.RFC3339),
				IsDeleted:    report.Opinion.IsDeleted,
			},
			User: oas.GetReportsForTalkSessionOKReportsItemUser{
				DisplayID:   report.User.DisplayID,
				DisplayName: report.User.DisplayName,
				IconURL:     utils.ToOptNil[oas.OptNilString](report.User.IconURL),
			},
			Status:      oas.GetReportsForTalkSessionOKReportsItemStatus(report.Status),
			Reasons:     reasons,
			ReportCount: report.ReportCount,
		})
	}

	return &oas.GetReportsForTalkSessionOK{
		Reports: reports,
	}, nil
}

// GetTalkSessionReportCount implements oas.TalkSessionHandler.
func (t *talkSessionHandler) GetTalkSessionReportCount(ctx context.Context, params oas.GetTalkSessionReportCountParams) (oas.GetTalkSessionReportCountRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.GetTalkSessionReportCount")
	defer span.End()

	claim := session.GetSession(t.SetSession(ctx))
	var userID *shared.UUID[user.User]
	if claim != nil {
		id, err := claim.UserID()
		if err == nil {
			userID = &id
		}
	}
	if userID == nil {
		return nil, messages.ForbiddenError
	}

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	out, err := t.getReportCount.Execute(ctx, report_query.GetCountInput{
		TalkSessionID: talkSessionID,
		UserID:        *userID,
		Status:        string(params.Status),
	})
	if err != nil {
		return nil, err
	}

	return &oas.GetTalkSessionReportCountOK{
		Count: out.Count,
	}, nil
}
