package di

import "github.com/neko-dream/server/internal/domain/service"

func domainDeps() []ProvideArg {
	return []ProvideArg{
		{service.NewAuthService, nil},
		{service.NewSessionService, nil},
		{service.NewUserService, nil},
		{service.NewOpinionService, nil},
		{service.NewActionItemService, nil},
		{service.NewStateGenerator, nil},
	}
}
