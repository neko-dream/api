package handler

import (
	"bytes"
	"context"
	"time"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/neko-dream/server/internal/usecase/command/talksession_command"
	"github.com/neko-dream/server/internal/usecase/query/analysis_query"
	talksession_query "github.com/neko-dream/server/internal/usecase/query/talksession"
	"github.com/neko-dream/server/pkg/sort"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type talkSessionHandler struct {
	browseTalkSessionsQuery       talksession_query.BrowseTalkSessionQuery
	browseOpenedByUserQuery       talksession_query.BrowseOpenedByUserQuery
	getConclusionByIDQuery        talksession_query.GetConclusionByIDQuery
	getTalkSessionDetailByIDQuery talksession_query.GetTalkSessionDetailByIDQuery
	getAnalysisResultQuery        analysis_query.GetAnalysisResult
	getReportQuery                analysis_query.GetReportQuery

	addConclusionCommand    talksession_command.AddConclusionCommand
	startTalkSessionCommand talksession_command.StartTalkSessionCommand

	session.TokenManager
}

func NewTalkSessionHandler(
	browseTalkSessionsQuery talksession_query.BrowseTalkSessionQuery,
	browseOpenedByUserQuery talksession_query.BrowseOpenedByUserQuery,
	getConclusionByIDQuery talksession_query.GetConclusionByIDQuery,
	getTalkSessionDetailByIDQuery talksession_query.GetTalkSessionDetailByIDQuery,
	getAnalysisQuery analysis_query.GetAnalysisResult,
	getReportQuery analysis_query.GetReportQuery,

	AddConclusionCommand talksession_command.AddConclusionCommand,
	startTalkSessionCommand talksession_command.StartTalkSessionCommand,

	tokenManager session.TokenManager,
) oas.TalkSessionHandler {
	return &talkSessionHandler{
		browseTalkSessionsQuery:       browseTalkSessionsQuery,
		browseOpenedByUserQuery:       browseOpenedByUserQuery,
		getConclusionByIDQuery:        getConclusionByIDQuery,
		getTalkSessionDetailByIDQuery: getTalkSessionDetailByIDQuery,
		getAnalysisResultQuery:        getAnalysisQuery,
		getReportQuery:                getReportQuery,

		addConclusionCommand:    AddConclusionCommand,
		startTalkSessionCommand: startTalkSessionCommand,

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

		resultTalkSession = append(resultTalkSession, oas.GetOpenedTalkSessionOKTalkSessionsItem{
			TalkSession: oas.GetOpenedTalkSessionOKTalkSessionsItemTalkSession{
				ID:               talkSession.TalkSessionID.String(),
				Theme:            talkSession.Theme,
				Owner:            owner,
				Location:         location,
				CreatedAt:        talkSession.CreatedAt.Format(time.RFC3339),
				ScheduledEndTime: talkSession.ScheduledEndTime.Format(time.RFC3339),
				City:             utils.ToOptNil[oas.OptNilString](talkSession.City),
				Prefecture:       utils.ToOptNil[oas.OptNilString](talkSession.Prefecture),
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

	out, err := t.getReportQuery.Execute(ctx, analysis_query.GetReportInput{
		TalkSessionID: shared.MustParseUUID[talksession.TalkSession](params.TalkSessionId),
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
	latitude := utils.ToPtrIfNotNullValue(!req.Value.Latitude.IsSet(), req.Value.Latitude.Value)
	longitude := utils.ToPtrIfNotNullValue(!req.Value.Longitude.IsSet(), req.Value.Longitude.Value)
	city := utils.ToPtrIfNotNullValue(!req.Value.City.IsSet(), req.Value.City.Value)
	prefecture := utils.ToPtrIfNotNullValue(!req.Value.Prefecture.IsSet(), req.Value.Prefecture.Value)
	description := utils.ToPtrIfNotNullValue(!req.Value.Description.IsSet(), req.Value.Description.Value)

	out, err := t.startTalkSessionCommand.Execute(ctx, talksession_command.StartTalkSessionCommandInput{
		Theme:            req.Value.Theme,
		Description:      description,
		OwnerID:          userID,
		ScheduledEndTime: req.Value.ScheduledEndTime,
		Latitude:         latitude,
		Longitude:        longitude,
		City:             city,
		Prefecture:       prefecture,
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

	res := &oas.CreateTalkSessionOK{
		Owner: oas.CreateTalkSessionOKOwner{
			DisplayID:   out.User.DisplayID,
			DisplayName: out.User.DisplayName,
			IconURL:     utils.ToOptNil[oas.OptNilString](out.User.IconURL),
		},
		Theme:            out.TalkSession.Theme,
		Description:      utils.ToOptNil[oas.OptNilString](out.TalkSession.Description),
		ID:               out.TalkSession.TalkSessionID.String(),
		CreatedAt:        clock.Now(ctx).Format(time.RFC3339),
		ScheduledEndTime: out.TalkSession.ScheduledEndTime.Format(time.RFC3339),
		Location:         location,
	}
	return res, nil
}

// ViewTalkSessionDetail トークセッション詳細取得
func (t *talkSessionHandler) GetTalkSessionDetail(ctx context.Context, params oas.GetTalkSessionDetailParams) (oas.GetTalkSessionDetailRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.GetTalkSessionDetail")
	defer span.End()

	out, err := t.getTalkSessionDetailByIDQuery.Execute(ctx, talksession_query.GetTalkSessionDetailInput{
		TalkSessionID: shared.MustParseUUID[talksession.TalkSession](params.TalkSessionId),
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

	return &oas.GetTalkSessionDetailOK{
		ID:               out.TalkSessionID.String(),
		Theme:            out.Theme,
		Description:      utils.ToOptNil[oas.OptNilString](out.Description),
		Owner:            owner,
		CreatedAt:        out.CreatedAt.Format(time.RFC3339),
		ScheduledEndTime: out.ScheduledEndTime.Format(time.RFC3339),
		Location:         location,
		City:             utils.ToOptNil[oas.OptNilString](out.City),
		Prefecture:       utils.ToOptNil[oas.OptNilString](out.Prefecture),
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

	var status talksession_query.Status
	if params.Status.IsSet() {
		bytes, err := params.Status.Value.MarshalText()
		if err == nil {
			status = talksession_query.Status(string(bytes))
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
		res := oas.GetTalkSessionListOKTalkSessionsItem{
			TalkSession: oas.GetTalkSessionListOKTalkSessionsItemTalkSession{
				ID:               talkSession.TalkSessionID.String(),
				Theme:            talkSession.Theme,
				Description:      utils.ToOptNil[oas.OptNilString](talkSession.Description),
				Owner:            owner,
				CreatedAt:        talkSession.CreatedAt.Format(time.RFC3339),
				ScheduledEndTime: talkSession.ScheduledEndTime.Format(time.RFC3339),
				Location:         location,
				City:             utils.ToOptNil[oas.OptNilString](talkSession.City),
				Prefecture:       utils.ToOptNil[oas.OptNilString](talkSession.Prefecture),
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

	out, err := t.getAnalysisResultQuery.Execute(ctx, analysis_query.GetAnalysisResultInput{
		UserID:        userID,
		TalkSessionID: shared.MustParseUUID[talksession.TalkSession](params.TalkSessionId),
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
					PictureURL:   utils.ToOpt[oas.OptString](opinion.Opinion.PictureURL),
					ReferenceURL: utils.ToOpt[oas.OptString](opinion.Opinion.ReferenceURL),
					PostedAt:     opinion.Opinion.CreatedAt.Format(time.RFC3339),
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

	_ = ctx

	panic("unimplemented")
}
