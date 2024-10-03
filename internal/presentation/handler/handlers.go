package handler

import (
	"github.com/neko-dream/server/internal/presentation/oas"
)

type handlers struct {
	oas.AuthHandler
	oas.IntentionHandler
	oas.OpinionHandler
	oas.TalkSessionHandler
	oas.UserHandler
}

func NewHandler(
	authHandler oas.AuthHandler,
	intentionHandler oas.IntentionHandler,
	opinionHandler oas.OpinionHandler,
	talkSessionHandler oas.TalkSessionHandler,
	userHandler oas.UserHandler,
) oas.Handler {
	return &handlers{
		AuthHandler:        authHandler,
		IntentionHandler:   intentionHandler,
		OpinionHandler:     opinionHandler,
		TalkSessionHandler: talkSessionHandler,
		UserHandler:        userHandler,
	}
}
