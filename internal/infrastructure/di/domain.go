package di

import (
	"github.com/neko-dream/server/internal/domain/model/talksession/talksession_consent"
	"github.com/neko-dream/server/internal/domain/service"
	organization_svc "github.com/neko-dream/server/internal/domain/service/organization"
)

// このファイルはドメイン層（サービス等）のコンストラクタを管理します。
// 新しいドメインサービス等を追加した場合は必ずここに追記してください。

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
		{talksession_consent.NewTalkSessionConsentService, nil},
	}
}
