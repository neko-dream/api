package handler

import (
	"bytes"
	"context"
	"log"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/shared/time"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/presentation/oas"
	analysis_usecase "github.com/neko-dream/server/internal/usecase/analysis"
	talk_session_usecase "github.com/neko-dream/server/internal/usecase/talk_session"
	"github.com/neko-dream/server/pkg/utils"
)

type talkSessionHandler struct {
	createTalkSessionUsecase  talk_session_usecase.CreateTalkSessionUseCase
	listTalkSessionQuery      talk_session_usecase.ListTalkSessionQuery
	getTalkSessionDetailQuery talk_session_usecase.GetTalkSessionDetailUseCase
	getAnalysisResultUseCase  analysis_usecase.GetAnalysisResultUseCase
	getReportUseCase          analysis_usecase.GetReportQuery
	session.TokenManager
}

func NewTalkSessionHandler(
	createTalkSessionUsecase talk_session_usecase.CreateTalkSessionUseCase,
	listTalkSessionQuery talk_session_usecase.ListTalkSessionQuery,
	getTalkSessionDetailQuery talk_session_usecase.GetTalkSessionDetailUseCase,
	getAnalysisResultUseCase analysis_usecase.GetAnalysisResultUseCase,
	getReportUseCase analysis_usecase.GetReportQuery,
	tokenManager session.TokenManager,
) oas.TalkSessionHandler {
	return &talkSessionHandler{
		createTalkSessionUsecase:  createTalkSessionUsecase,
		listTalkSessionQuery:      listTalkSessionQuery,
		getTalkSessionDetailQuery: getTalkSessionDetailQuery,
		getAnalysisResultUseCase:  getAnalysisResultUseCase,
		getReportUseCase:          getReportUseCase,
		TokenManager:              tokenManager,
	}
}

// GetTalkSessionReport implements oas.TalkSessionHandler.
func (t *talkSessionHandler) GetTalkSessionReport(ctx context.Context, params oas.GetTalkSessionReportParams) (oas.GetTalkSessionReportRes, error) {
	out, err := t.getReportUseCase.Execute(ctx, analysis_usecase.GetReportInput{
		TalkSessionID: shared.MustParseUUID[talksession.TalkSession](params.TalkSessionId),
	})
	if err != nil {
		return nil, err
	}

	res := &oas.GetTalkSessionReportOKHeaders{}
	buf := make([]byte, len(out.Report))
	copy(buf, out.Report)
	dt := &oas.GetTalkSessionReportOK{}
	dt.Data = bytes.NewBuffer(buf)
	res.SetResponse(*dt)
	return res, nil

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
	var latitude, longitude *float64
	if req.Value.Latitude.IsSet() {
		latitude = &req.Value.Latitude.Value
	}
	if req.Value.Longitude.IsSet() {
		longitude = &req.Value.Longitude.Value
	}

	out, err := t.createTalkSessionUsecase.Execute(ctx, talk_session_usecase.CreateTalkSessionInput{
		Theme:            req.Value.Theme,
		OwnerID:          userID,
		ScheduledEndTime: time.NewTime(ctx, req.Value.ScheduledEndTime),
		Latitude:         latitude,
		Longitude:        longitude,
		City:             utils.ToPtrIfNotNullValue(req.Value.City.Null, req.Value.City.Value),
		Prefecture:       utils.ToPtrIfNotNullValue(req.Value.Prefecture.Null, req.Value.Prefecture.Value),
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	var location oas.OptCreateTalkSessionOKLocation
	if out.Location != nil {
		location = oas.OptCreateTalkSessionOKLocation{}
		location.Value.CreateTalkSessionOKLocation0.Latitude = out.Location.Latitude()
		location.Value.CreateTalkSessionOKLocation0.Longitude = out.Location.Longitude()
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
		ID:               out.TalkSession.TalkSessionID().String(),
		CreatedAt:        time.Now(ctx).Format(ctx),
		ScheduledEndTime: out.TalkSession.ScheduledEndTime().Format(ctx),
		Location:         location,
	}
	return res, nil
}

// GetTalkSessionDetail トークセッション詳細取得
func (t *talkSessionHandler) GetTalkSessionDetail(ctx context.Context, params oas.GetTalkSessionDetailParams) (oas.GetTalkSessionDetailRes, error) {
	out, err := t.getTalkSessionDetailQuery.Execute(ctx, talk_session_usecase.GetTalkSessionDetailInput{
		TalkSessionID: shared.MustParseUUID[talksession.TalkSession](params.TalkSessionId),
	})
	if err != nil {
		return nil, err
	}

	owner := oas.GetTalkSessionDetailOKOwner{
		DisplayID:   out.Owner.DisplayID,
		DisplayName: out.Owner.DisplayName,
		IconURL:     utils.ToOptNil[oas.OptNilString](out.Owner.IconURL),
	}
	var location oas.OptGetTalkSessionDetailOKLocation
	if out.Location != nil {
		location.Value.GetTalkSessionDetailOKLocation0.Latitude = out.Location.Latitude
		location.Value.GetTalkSessionDetailOKLocation0.Longitude = out.Location.Longitude
	}

	return &oas.GetTalkSessionDetailOK{
		ID:               out.ID,
		Theme:            out.Theme,
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
	theme := utils.ToPtrIfNotNullValue(params.Theme.Null, params.Theme.Value)
	out, err := t.listTalkSessionQuery.Execute(ctx, talk_session_usecase.ListTalkSessionInput{
		Limit:  limit,
		Offset: offset,
		Theme:  theme,
		Status: status,
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
		resultTalkSession = append(resultTalkSession, oas.GetTalkSessionListOKTalkSessionsItem{
			TalkSession: oas.GetTalkSessionListOKTalkSessionsItemTalkSession{
				ID:               talkSession.ID,
				Theme:            talkSession.Theme,
				Owner:            owner,
				CreatedAt:        talkSession.CreatedAt,
				ScheduledEndTime: talkSession.ScheduledEndTime,
			},
			OpinionCount: talkSession.OpinionCount,
		})
	}

	return &oas.GetTalkSessionListOK{
		TalkSessions: resultTalkSession,
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
				GroupId:        out.MyPosition.GroupID,
				PerimeterIndex: utils.ToOpt[oas.OptInt](out.MyPosition.PerimeterIndex),
			},
			Set: true,
		}
	}
	log.Println("out.Positions", myPosition)

	positions := make([]oas.TalkSessionAnalysisOKPositionsItem, 0, len(out.Positions))
	for _, position := range out.Positions {
		positions = append(positions, oas.TalkSessionAnalysisOKPositionsItem{
			PosX:           position.PosX,
			PosY:           position.PosY,
			DisplayId:      position.DisplayID,
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
			GroupId:  groupOpinion.GroupID,
			Opinions: opinions,
		})
	}

	return &oas.TalkSessionAnalysisOK{
		MyPosition:    myPosition,
		Positions:     positions,
		GroupOpinions: groupOpinions,
	}, nil
}
