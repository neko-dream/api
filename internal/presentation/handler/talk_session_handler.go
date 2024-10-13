package handler

import (
	"context"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared/time"
	"github.com/neko-dream/server/internal/presentation/oas"
	talk_session_usecase "github.com/neko-dream/server/internal/usecase/talk_session"
	"github.com/neko-dream/server/pkg/utils"
)

type talkSessionHandler struct {
	createTalkSessionUsecase talk_session_usecase.CreateTalkSessionUseCase
	listTalkSessionQuery     talk_session_usecase.ListTalkSessionQuery
}

func NewTalkSessionHandler(
	createTalkSessionUsecase talk_session_usecase.CreateTalkSessionUseCase,
	listTalkSessionQuery talk_session_usecase.ListTalkSessionQuery,
) oas.TalkSessionHandler {
	return &talkSessionHandler{
		createTalkSessionUsecase: createTalkSessionUsecase,
		listTalkSessionQuery:     listTalkSessionQuery,
	}
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
	panic("unimplemented")
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
