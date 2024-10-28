package handler

import (
	"github.com/neko-dream/server/internal/presentation/oas"
)

type handlers struct {
	oas.AuthHandler
	oas.VoteHandler
	oas.OpinionHandler
	oas.TalkSessionHandler
	oas.UserHandler
	oas.TestHandler
	oas.ManageHandler
	oas.TimelineHandler
}

func NewHandler(
	authHandler oas.AuthHandler,
	voteHandler oas.VoteHandler,
	opinionHandler oas.OpinionHandler,
	talkSessionHandler oas.TalkSessionHandler,
	userHandler oas.UserHandler,
	testHandler oas.TestHandler,
	manageHandler oas.ManageHandler,
) oas.Handler {
	return &handlers{
		AuthHandler:        authHandler,
		VoteHandler:        voteHandler,
		OpinionHandler:     opinionHandler,
		TalkSessionHandler: talkSessionHandler,
		UserHandler:        userHandler,
		TestHandler:        testHandler,
		ManageHandler:      manageHandler,
		TimelineHandler:    oas.UnimplementedHandler{},
	}
}
