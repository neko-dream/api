package handler

import (
	"context"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/presentation/oas"
	talk_session_usecase "github.com/neko-dream/server/internal/usecase/talk_session"
	"github.com/neko-dream/server/pkg/utils"
)

type talkSessionHandler struct {
	createTalkSessionUsecase talk_session_usecase.CreateTalkSessionUseCase
}

// CreateTalkSession implements oas.TalkSessionHandler.
func (t *talkSessionHandler) CreateTalkSession(ctx context.Context, req oas.OptCreateTalkSessionReq) (*oas.CreateTalkSessionOK, error) {
	claim := session.GetSession(ctx)
	if !req.IsSet() {
		return nil, messages.RequiredParameterError
	}
	userID, err := claim.UserID()
	if err != nil {
		utils.HandleError(ctx, err, "claim.UserID")
		return nil, messages.ForbiddenError
	}

	out, err := t.createTalkSessionUsecase.Execute(ctx, talk_session_usecase.CreateTalkSessionInput{
		Theme:   req.Value.Theme.Value,
		OwnerID: userID,
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	return &oas.CreateTalkSessionOK{
		Owner: oas.CreateTalkSessionOKOwner{
			DisplayID:   *out.OwnerUser.DisplayID(),
			DisplayName: *out.OwnerUser.DisplayName(),
			IconURL: utils.IfThenElse(
				out.OwnerUser.ProfileIconURL() != nil,
				oas.OptString{Value: *out.OwnerUser.ProfileIconURL()},
				oas.OptString{},
			),
		},
		Theme: out.TalkSession.Theme(),
		ID:    out.TalkSession.TalkSessionID().String(),
	}, nil
}

// GetTalkSessionDetail implements oas.TalkSessionHandler.
func (t *talkSessionHandler) GetTalkSessionDetail(ctx context.Context, params oas.GetTalkSessionDetailParams) (*oas.GetTalkSessionDetailOK, error) {
	panic("unimplemented")
}

// GetTalkSessions implements oas.TalkSessionHandler.
func (t *talkSessionHandler) GetTalkSessions(context.Context) (*oas.GetTalkSessionsOK, error) {
	panic("unimplemented")
}

func NewTalkSessionHandler(
	createTalkSessionUsecase talk_session_usecase.CreateTalkSessionUseCase,
) oas.TalkSessionHandler {
	return &talkSessionHandler{
		createTalkSessionUsecase: createTalkSessionUsecase,
	}
}
