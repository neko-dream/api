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
}

func NewHandler(
	authHandler oas.AuthHandler,
	voteHandler oas.VoteHandler,
	opinionHandler oas.OpinionHandler,
	talkSessionHandler oas.TalkSessionHandler,
	userHandler oas.UserHandler,
) oas.Handler {
	return &handlers{
		AuthHandler:        authHandler,
		VoteHandler:        voteHandler,
		OpinionHandler:     opinionHandler,
		TalkSessionHandler: talkSessionHandler,
		UserHandler:        userHandler,
		TestHandler:        oas.UnimplementedHandler{},
	}
}
