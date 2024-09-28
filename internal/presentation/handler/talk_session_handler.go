package handler

import (
	"context"
	"os/user"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/presentation/oas"
	talk_session_usecase "github.com/neko-dream/server/internal/usecase/talk_session"
)

type talkSessionHandler struct {
	createTalkSessionUsecase talk_session_usecase.CreateTalkSessionUseCase
}

// CreateTalkSession implements oas.TalkSessionHandler.
func (t *talkSessionHandler) CreateTalkSession(ctx context.Context, req oas.OptCreateTalkSessionReq) (*oas.CreateTalkSessionOK, error) {
	out, err := t.createTalkSessionUsecase.Execute(ctx, talk_session_usecase.CreateTalkSessionInput{
		Theme:   req.Value.Theme.Value,
		OwnerID: shared.NewUUID[user.User]().String(),
	})

	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	res := &oas.CreateTalkSessionOK{}
	res.TalkSessionID = out.TalkSession.TalkSessionID().String()
	res.TalkSessionTheme = out.TalkSession.Theme()
	res.TalkSessionStatus = out.TalkSession.OwnerUserID().String()

	return res, nil
}

// GetTalkSessionDetail implements oas.TalkSessionHandler.
func (t *talkSessionHandler) GetTalkSessionDetail(ctx context.Context, params oas.GetTalkSessionDetailParams) (*oas.GetTalkSessionDetailOK, error) {
	panic("unimplemented")
}

// GetTalkSessions implements oas.TalkSessionHandler.
func (t *talkSessionHandler) GetTalkSessions(ctx context.Context) error {
	panic("unimplemented")
}

func NewTalkSessionHandler(
	createTalkSessionUsecase talk_session_usecase.CreateTalkSessionUseCase,
) oas.TalkSessionHandler {
	return &talkSessionHandler{
		createTalkSessionUsecase: createTalkSessionUsecase,
	}
}
