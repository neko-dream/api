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
	analysis_usecase "github.com/neko-dream/server/internal/usecase/analysis"
	talk_session_usecase "github.com/neko-dream/server/internal/usecase/talk_session"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
)

type talkSessionHandler struct {
	StartTalkSessionCommand         talk_session_usecase.StartTalkSessionCommand
	listTalkSessionQuery            talk_session_usecase.ListTalkSessionQuery
	ViewTalkSessionDetailQuery      talk_session_usecase.ViewTalkSessionDetailQuery
	getAnalysisResultUseCase        analysis_usecase.GetAnalysisResultUseCase
	getReportUseCase                analysis_usecase.GetReportQuery
	getOwnTalkSession               talk_session_usecase.BrowseUsersTalkSessionHistoriesQuery
	AddTalkSessionConclusionCommand talk_session_usecase.AddTalkSessionConclusionCommand
	getTalkSessionConclusionQuery   talk_session_usecase.GetTalkSessionConclusionQuery
	session.TokenManager
}

func NewTalkSessionHandler(
	StartTalkSessionCommand talk_session_usecase.StartTalkSessionCommand,
	listTalkSessionQuery talk_session_usecase.ListTalkSessionQuery,
	ViewTalkSessionDetailQuery talk_session_usecase.ViewTalkSessionDetailQuery,
	getAnalysisResultUseCase analysis_usecase.GetAnalysisResultUseCase,
	getReportUseCase analysis_usecase.GetReportQuery,
	getOwnTalkSession talk_session_usecase.BrowseUsersTalkSessionHistoriesQuery,
	AddTalkSessionConclusionCommand talk_session_usecase.AddTalkSessionConclusionCommand,
	getTalkSessionConclusionQuery talk_session_usecase.GetTalkSessionConclusionQuery,
	tokenManager session.TokenManager,
) oas.TalkSessionHandler {
	return &talkSessionHandler{
		StartTalkSessionCommand:         StartTalkSessionCommand,
		listTalkSessionQuery:            listTalkSessionQuery,
		ViewTalkSessionDetailQuery:      ViewTalkSessionDetailQuery,
		getAnalysisResultUseCase:        getAnalysisResultUseCase,
		getReportUseCase:                getReportUseCase,
		getOwnTalkSession:               getOwnTalkSession,
		AddTalkSessionConclusionCommand: AddTalkSessionConclusionCommand,
		getTalkSessionConclusionQuery:   getTalkSessionConclusionQuery,
		TokenManager:                    tokenManager,
	}
}

// PostConclusion implements oas.TalkSessionHandler.
func (t *talkSessionHandler) PostConclusion(ctx context.Context, req oas.OptPostConclusionReq, params oas.PostConclusionParams) (oas.PostConclusionRes, error) {
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

	out, err := t.AddTalkSessionConclusionCommand.Execute(ctx, talk_session_usecase.CreateTalkSessionConclusionInput{
		TalkSessionID: talkSessionID,
		UserID:        userID,
		Conclusion:    req.Value.Content,
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	return &oas.PostConclusionOK{
		User: oas.PostConclusionOKUser{
			DisplayID:   out.User.DisplayID,
			DisplayName: out.User.DisplayName,
			IconURL:     utils.ToOptNil[oas.OptNilString](out.User.IconURL),
		},
		Content: out.Content,
	}, nil
}

// GetConclusion implements oas.TalkSessionHandler.
func (t *talkSessionHandler) GetConclusion(ctx context.Context, params oas.GetConclusionParams) (oas.GetConclusionRes, error) {
	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	out, err := t.getTalkSessionConclusionQuery.Execute(ctx, talk_session_usecase.GetTalkSessionConclusionInput{
		TalkSessionID: talkSessionID,
	})
	if err != nil {
		return nil, err
	}
	// まだ結論が出ていない場合はエラーを返す
	if out == nil {
		return nil, messages.TalkSessionConclusionNotSet
	}

	return &oas.GetConclusionOK{
		User: oas.GetConclusionOKUser{
			DisplayID:   out.User.DisplayID,
			DisplayName: out.User.DisplayName,
			IconURL:     utils.ToOptNil[oas.OptNilString](out.User.IconURL),
		},
		Content: out.Conclusion,
	}, nil
}

// GetOpenedTalkSession implements oas.TalkSessionHandler.
func (t *talkSessionHandler) GetOpenedTalkSession(ctx context.Context, params oas.GetOpenedTalkSessionParams) (oas.GetOpenedTalkSessionRes, error) {
	claim := session.GetSession(t.SetSession(ctx))
	if claim == nil {
		return nil, messages.ForbiddenError
	}
	userID, err := claim.UserID()
	if err != nil {
		utils.HandleError(ctx, err, "claim.UserID")
		return nil, messages.ForbiddenError
	}

	limit := utils.IfThenElse(params.Limit.IsSet(),
		params.Limit.Value,
		10,
	)
	offset := utils.IfThenElse(params.Offset.IsSet(),
		params.Offset.Value,
		0,
	)

	status := ""
	if params.Status.IsSet() {
		bytes, err := params.Status.Value.MarshalText()
		if err == nil {
			status = string(bytes)
		}
	}
	out, err := t.getOwnTalkSession.Execute(ctx, talk_session_usecase.BrowseUsersTalkSessionHistoriesInput{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
		Status: status,
	})
	if err != nil {
		return nil, err
	}

	resultTalkSession := make([]oas.GetOpenedTalkSessionOKTalkSessionsItem, 0, len(out.TalkSessions))
	for _, talkSession := range out.TalkSessions {
		owner := oas.GetOpenedTalkSessionOKTalkSessionsItemTalkSessionOwner{
			DisplayID:   talkSession.Owner.DisplayID,
			DisplayName: talkSession.Owner.DisplayName,
			IconURL:     utils.ToOptNil[oas.OptNilString](talkSession.Owner.IconURL),
		}
		var location oas.OptGetOpenedTalkSessionOKTalkSessionsItemTalkSessionLocation
		if talkSession.Location != nil {
			location = oas.OptGetOpenedTalkSessionOKTalkSessionsItemTalkSessionLocation{
				Value: oas.GetOpenedTalkSessionOKTalkSessionsItemTalkSessionLocation{
					Latitude:  utils.ToOpt[oas.OptFloat64](talkSession.Location.Latitude),
					Longitude: utils.ToOpt[oas.OptFloat64](talkSession.Location.Longitude),
				},
				Set: true,
			}
		} else {
			location = oas.OptGetOpenedTalkSessionOKTalkSessionsItemTalkSessionLocation{}
			location.Set = false
		}

		resultTalkSession = append(resultTalkSession, oas.GetOpenedTalkSessionOKTalkSessionsItem{
			TalkSession: oas.GetOpenedTalkSessionOKTalkSessionsItemTalkSession{
				ID:               talkSession.ID,
				Theme:            talkSession.Theme,
				Owner:            owner,
				Location:         location,
				CreatedAt:        talkSession.CreatedAt,
				ScheduledEndTime: talkSession.ScheduledEndTime,
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
	out, err := t.getReportUseCase.Execute(ctx, analysis_usecase.GetReportInput{
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

	out, err := t.StartTalkSessionCommand.Execute(ctx, talk_session_usecase.CreateTalkSessionInput{
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
	if out.Location != nil {
		location = oas.OptCreateTalkSessionOKLocation{
			Value: oas.CreateTalkSessionOKLocation{
				Latitude:  utils.ToOpt[oas.OptFloat64](out.Location.Latitude),
				Longitude: utils.ToOpt[oas.OptFloat64](out.Location.Longitude),
			},
			Set: true,
		}
	}

	res := &oas.CreateTalkSessionOK{
		Owner: oas.CreateTalkSessionOKOwner{
			DisplayID:   *out.OwnerUser.DisplayID(),
			DisplayName: *out.OwnerUser.DisplayName(),
			IconURL: utils.IfThenElse(
				out.OwnerUser.ProfileIconURL() != nil,
				oas.OptNilString{Value: *out.OwnerUser.ProfileIconURL()},
				oas.OptNilString{},
			),
		},
		Theme:            out.TalkSession.Theme(),
		Description:      utils.ToOptNil[oas.OptNilString](out.TalkSession.Description()),
		ID:               out.TalkSession.TalkSessionID().String(),
		CreatedAt:        clock.Now(ctx).Format(time.RFC3339),
		ScheduledEndTime: out.TalkSession.ScheduledEndTime().Format(time.RFC3339),
		Location:         location,
	}
	return res, nil
}

// ViewTalkSessionDetail トークセッション詳細取得
func (t *talkSessionHandler) ViewTalkSessionDetail(ctx context.Context, params oas.ViewTalkSessionDetailParams) (oas.ViewTalkSessionDetailRes, error) {
	out, err := t.ViewTalkSessionDetailQuery.Execute(ctx, talk_session_usecase.ViewTalkSessionDetailInput{
		TalkSessionID: shared.MustParseUUID[talksession.TalkSession](params.TalkSessionId),
	})
	if err != nil {
		return nil, err
	}

	owner := oas.ViewTalkSessionDetailOKOwner{
		DisplayID:   out.Owner.DisplayID,
		DisplayName: out.Owner.DisplayName,
		IconURL:     utils.ToOptNil[oas.OptNilString](out.Owner.IconURL),
	}
	var location oas.OptViewTalkSessionDetailOKLocation
	if out.Location != nil {
		location = oas.OptViewTalkSessionDetailOKLocation{
			Value: oas.ViewTalkSessionDetailOKLocation{
				Latitude:  utils.ToOpt[oas.OptFloat64](out.Location.Latitude),
				Longitude: utils.ToOpt[oas.OptFloat64](out.Location.Longitude),
			},
			Set: true,
		}

	}

	return &oas.ViewTalkSessionDetailOK{
		ID:               out.ID,
		Theme:            out.Theme,
		Description:      utils.ToOptNil[oas.OptNilString](out.Description),
		Owner:            owner,
		CreatedAt:        out.CreatedAt,
		ScheduledEndTime: out.ScheduledEndTime,
		Location:         location,
		City:             utils.ToOptNil[oas.OptNilString](out.City),
		Prefecture:       utils.ToOptNil[oas.OptNilString](out.Prefecture),
	}, nil
}

// GetTalkSessionList セッション一覧取得
func (t *talkSessionHandler) GetTalkSessionList(ctx context.Context, params oas.GetTalkSessionListParams) (oas.GetTalkSessionListRes, error) {
	limit := utils.IfThenElse[int](params.Limit.IsSet(),
		params.Limit.Value,
		10,
	)
	offset := utils.IfThenElse[int](params.Offset.IsSet(),
		params.Offset.Value,
		0,
	)
	status := ""
	if params.Status.IsSet() {
		bytes, err := params.Status.Value.MarshalText()
		if err == nil {
			status = string(bytes)
		}
	}
	var sortKey *string
	if params.SortKey.IsSet() {
		bytes, err := params.SortKey.Value.MarshalText()
		if err == nil {
			sortKey = lo.ToPtr(string(bytes))
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
	out, err := t.listTalkSessionQuery.Execute(ctx, talk_session_usecase.ListTalkSessionInput{
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
			DisplayID:   talkSession.Owner.DisplayID,
			DisplayName: talkSession.Owner.DisplayName,
			IconURL:     utils.ToOptNil[oas.OptNilString](talkSession.Owner.IconURL),
		}
		var location oas.OptGetTalkSessionListOKTalkSessionsItemTalkSessionLocation
		if talkSession.Location != nil {
			location = oas.OptGetTalkSessionListOKTalkSessionsItemTalkSessionLocation{
				Value: oas.GetTalkSessionListOKTalkSessionsItemTalkSessionLocation{
					Latitude:  utils.ToOpt[oas.OptFloat64](talkSession.Location.Latitude),
					Longitude: utils.ToOpt[oas.OptFloat64](talkSession.Location.Longitude),
				},
				Set: true,
			}
		}
		res := oas.GetTalkSessionListOKTalkSessionsItem{
			TalkSession: oas.GetTalkSessionListOKTalkSessionsItemTalkSession{
				ID:               talkSession.ID,
				Theme:            talkSession.Theme,
				Description:      utils.ToOptNil[oas.OptNilString](talkSession.Description),
				Owner:            owner,
				CreatedAt:        talkSession.CreatedAt,
				ScheduledEndTime: talkSession.ScheduledEndTime,
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
			Limit:      limit,
			Offset:     offset,
		},
	}, nil
}

// TalkSessionAnalysis 分析結果取得
func (t *talkSessionHandler) TalkSessionAnalysis(ctx context.Context, params oas.TalkSessionAnalysisParams) (oas.TalkSessionAnalysisRes, error) {
	claim := session.GetSession(t.SetSession(ctx))
	var userID *shared.UUID[user.User]
	if claim != nil {
		id, err := claim.UserID()
		if err == nil {
			userID = &id
		}
	}

	out, err := t.getAnalysisResultUseCase.Execute(ctx, analysis_usecase.GetAnalysisResultInput{
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
					ID:           opinion.Opinion.ID,
					Title:        utils.ToOpt[oas.OptString](opinion.Opinion.Title),
					Content:      opinion.Opinion.Content,
					ParentID:     utils.ToOpt[oas.OptString](opinion.Opinion.ParentID),
					PictureURL:   utils.ToOpt[oas.OptString](opinion.Opinion.PictureURL),
					ReferenceURL: utils.ToOpt[oas.OptString](opinion.Opinion.ReferenceURL),
					PostedAt:     opinion.Opinion.PostedAt,
				},
				User: oas.TalkSessionAnalysisOKGroupOpinionsItemOpinionsItemUser{
					DisplayID:   opinion.User.ID,
					DisplayName: opinion.User.Name,
					IconURL:     utils.ToOptNil[oas.OptNilString](opinion.User.Icon),
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
	panic("unimplemented")
}
