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
	oas.ImageHandler
	oas.PolicyHandler
	oas.OrganizationHandler
	oas.HealthHandler
	oas.AnalysisHandler
}

func NewHandler(
	authHandler oas.AuthHandler,
	voteHandler oas.VoteHandler,
	opinionHandler oas.OpinionHandler,
	talkSessionHandler oas.TalkSessionHandler,
	userHandler oas.UserHandler,
	testHandler oas.TestHandler,
	manageHandler oas.ManageHandler,
	timelineHandler oas.TimelineHandler,
	imageHandler oas.ImageHandler,
	policyHandler oas.PolicyHandler,
	organizationHandler oas.OrganizationHandler,
	healthHandler oas.HealthHandler,
	analysisHandler oas.AnalysisHandler,
) oas.Handler {
	return &handlers{
		AuthHandler:         authHandler,
		VoteHandler:         voteHandler,
		OpinionHandler:      opinionHandler,
		TalkSessionHandler:  talkSessionHandler,
		UserHandler:         userHandler,
		TestHandler:         testHandler,
		ManageHandler:       manageHandler,
		TimelineHandler:     timelineHandler,
		ImageHandler:        imageHandler,
		PolicyHandler:       policyHandler,
		OrganizationHandler: organizationHandler,
		HealthHandler:       healthHandler,
		AnalysisHandler:     analysisHandler,
	}
}
