package di

import (
	"github.com/neko-dream/server/internal/domain/service"
	organization_svc "github.com/neko-dream/server/internal/domain/service/organization"
)

func domainDeps() []ProvideArg {
	return []ProvideArg{
		{service.NewAuthService, nil},
		{service.NewSessionService, nil},
		{service.NewUserService, nil},
		{service.NewOpinionService, nil},
		{service.NewActionItemService, nil},
		{service.NewStateGenerator, nil},
		{service.NewProfileIconService, nil},
		{service.NewTalkSessionAccessControl, nil},
		{service.NewConsentService, nil},
		{service.NewPasswordAuthManager, nil},
		{organization_svc.NewOrganizationService, nil},
		{organization_svc.NewOrganizationMemberManager, nil},
	}
}
