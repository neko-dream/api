package di

import "github.com/neko-dream/server/internal/presentation/handler"

func presentationDeps() []ProvideArg {
	return []ProvideArg{
		{handler.NewSecurityHandler, nil},
		{handler.NewAuthHandler, nil},
		{handler.NewUserHandler, nil},
		{handler.NewVoteHandler, nil},
		{handler.NewOpinionHandler, nil},
		{handler.NewTalkSessionHandler, nil},
		{handler.NewHandler, nil},
		{handler.NewTestHandler, nil},
		{handler.NewManageHandler, nil},
		{handler.NewTimelineHandler, nil},
		{handler.NewImageHandler, nil},
		{handler.NewPolicyHandler, nil},
		{handler.NewOrganizationHandler, nil},
		{handler.NewHealthHandler, nil},
	}
}
