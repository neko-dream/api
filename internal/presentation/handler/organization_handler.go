package handler

import (
	"context"
	"net/http"

	"github.com/neko-dream/server/internal/application/query/organization_query"
	"github.com/neko-dream/server/internal/application/usecase/organization_usecase"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/service"
	domainservice "github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/infrastructure/http/cookie"
	"github.com/neko-dream/server/internal/presentation/oas"
	cookie_utils "github.com/neko-dream/server/pkg/cookie"
	"go.opentelemetry.io/otel"
)

type organizationHandler struct {
	create               organization_usecase.CreateOrganizationCommand
	invite               organization_usecase.InviteOrganizationCommand
	add                  organization_usecase.InviteOrganizationForUserCommand
	list                 organization_query.ListJoinedOrganizationQuery
	listUsers            organization_query.ListOrganizationUsersQuery
	createAlias          *organization_usecase.CreateOrganizationAliasUseCase
	deactivateAlias      *organization_usecase.DeactivateOrganizationAliasUseCase
	listAliases          *organization_usecase.ListOrganizationAliasesUseCase
	orgRepo              organization.OrganizationRepository
	aliasRepo            organization.OrganizationAliasRepository
	orgUserRepo          organization.OrganizationUserRepository
	aliasService         *domainservice.OrganizationAliasService
	authorizationService service.AuthorizationService
	sessionTokenManager  session.TokenManager
	switchOrganization   organization_usecase.SwitchOrganizationUseCase
	cookieManager        cookie.CookieManager
}

func NewOrganizationHandler(
	create organization_usecase.CreateOrganizationCommand,
	invite organization_usecase.InviteOrganizationCommand,
	add organization_usecase.InviteOrganizationForUserCommand,
	list organization_query.ListJoinedOrganizationQuery,
	listUsers organization_query.ListOrganizationUsersQuery,
	createAlias *organization_usecase.CreateOrganizationAliasUseCase,
	deactivateAlias *organization_usecase.DeactivateOrganizationAliasUseCase,
	listAliases *organization_usecase.ListOrganizationAliasesUseCase,
	orgRepo organization.OrganizationRepository,
	aliasRepo organization.OrganizationAliasRepository,
	orgUserRepo organization.OrganizationUserRepository,
	aliasService *domainservice.OrganizationAliasService,
	authorizationService service.AuthorizationService,
	sessionTokenManager session.TokenManager,
	switchOrganization organization_usecase.SwitchOrganizationUseCase,
	cookieManager cookie.CookieManager,
) oas.OrganizationHandler {
	return &organizationHandler{
		create:               create,
		invite:               invite,
		add:                  add,
		list:                 list,
		listUsers:            listUsers,
		createAlias:          createAlias,
		deactivateAlias:      deactivateAlias,
		listAliases:          listAliases,
		orgRepo:              orgRepo,
		aliasRepo:            aliasRepo,
		orgUserRepo:          orgUserRepo,
		aliasService:         aliasService,
		authorizationService: authorizationService,
		sessionTokenManager:  sessionTokenManager,
		switchOrganization:   switchOrganization,
		cookieManager:        cookieManager,
	}
}

// EstablishOrganization implements oas.OrganizationHandler.
func (o *organizationHandler) EstablishOrganization(ctx context.Context, req *oas.EstablishOrganizationReq) (oas.EstablishOrganizationRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "organizationHandler.EstablishOrganization")
	defer span.End()

	if req == nil {
		return nil, messages.BadRequestError
	}
	authCtx, err := o.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}
	if !authCtx.IsInOrganization() {
		return nil, messages.OrganizationContextRequired
	}
	// SuperAdmin権限を持つか確認
	if !authCtx.HasOrganizationRole(organization.OrganizationUserRoleSuperAdmin) {
		return nil, messages.InsufficientPermissionsError
	}

	_, err = o.create.Execute(ctx, organization_usecase.CreateOrganizationInput{
		UserID: authCtx.UserID,
		Name:   req.Name,
		Code:   req.Code,
		Type:   int(req.OrgType),
	})
	if err != nil {
		return nil, err
	}

	res := &oas.EstablishOrganizationOK{}
	return res, nil
}

// InviteOrganization implements oas.OrganizationHandler.
func (o *organizationHandler) InviteOrganization(ctx context.Context, req *oas.InviteOrganizationReq) (oas.InviteOrganizationRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "organizationHandler.InviteOrganization")
	defer span.End()
	if req == nil {
		return nil, messages.BadRequestError
	}

	authCtx, err := o.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}
	if !authCtx.IsInOrganization() {
		return nil, messages.OrganizationContextRequired
	}
	if !authCtx.HasOrganizationRole(organization.OrganizationUserRoleAdmin) {
		return nil, messages.InsufficientPermissionsError
	}

	_, err = o.invite.Execute(ctx, organization_usecase.InviteOrganizationInput{
		UserID:         authCtx.UserID,
		OrganizationID: *authCtx.OrganizationID,
		Email:          req.Email,
		Role:           int(req.Role),
	})
	if err != nil {
		return nil, err
	}

	res := &oas.InviteOrganizationOK{}
	return res, nil
}

// InviteOrganizationForUser implements oas.OrganizationHandler.
func (o *organizationHandler) InviteOrganizationForUser(ctx context.Context, req *oas.InviteOrganizationForUserReq) (oas.InviteOrganizationForUserRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "organizationHandler.InviteOrganizationForUser")
	defer span.End()

	authCtx, err := o.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}
	if !authCtx.IsInOrganization() {
		return nil, messages.OrganizationContextRequired
	}
	if !authCtx.HasOrganizationRole(organization.OrganizationUserRoleAdmin) {
		return nil, messages.InsufficientPermissionsError
	}

	_, err = o.add.Execute(ctx, organization_usecase.InviteOrganizationForUserInput{
		UserID:         authCtx.UserID,
		OrganizationID: *authCtx.OrganizationID,
		DisplayID:      req.DisplayID,
		Role:           int(req.Role),
	})
	if err != nil {
		return nil, err
	}

	res := &oas.InviteOrganizationForUserOK{
		Success: true,
	}
	return res, nil
}

// GetOrganizations 所属組織一覧
func (o *organizationHandler) GetOrganizations(ctx context.Context) (oas.GetOrganizationsRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "organizationHandler.GetOrganizations")
	defer span.End()

	authCtx, err := o.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	res, err := o.list.Execute(ctx, organization_query.ListJoinedOrganizationInput{
		UserID: authCtx.UserID,
	})
	if err != nil {
		return nil, err
	}
	if len(res.Organizations) == 0 {
		return &oas.GetOrganizationsOK{
			Organizations: []oas.Organization{},
		}, nil
	}

	var orgs []oas.Organization
	for _, org := range res.Organizations {
		orgs = append(orgs, org.ToResponse())
	}

	return &oas.GetOrganizationsOK{
		Organizations: orgs,
	}, nil

}

// ValidateOrganizationCode 組織コード検証
func (o *organizationHandler) ValidateOrganizationCode(ctx context.Context, params oas.ValidateOrganizationCodeParams) (oas.ValidateOrganizationCodeRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "organizationHandler.ValidateOrganizationCode")
	defer span.End()

	org, err := o.orgRepo.FindByCode(ctx, params.Code)
	if err != nil {
		return &oas.ValidateOrganizationCodeOK{
			Valid: false,
		}, nil
	}

	return &oas.ValidateOrganizationCodeOK{
		Valid: true,
		Organization: oas.NewOptOrganization(oas.Organization{
			ID:       org.OrganizationID.String(),
			Name:     org.Name,
			Code:     org.Code,
			Type:     int(org.OrganizationType),
			Role:     0,  // Role is not applicable in this context
			RoleName: "", // RoleName is not applicable in this context
		}),
	}, nil
}

// GetOrganizationAliases 組織エイリアス一覧取得
func (o *organizationHandler) GetOrganizationAliases(ctx context.Context) (oas.GetOrganizationAliasesRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "organizationHandler.GetOrganizationAliases")
	defer span.End()

	authCtx, err := o.authorizationService.RequireOrganizationRole(ctx, organization.OrganizationUserRoleMember)
	if err != nil {
		return nil, err
	}

	// エイリアス一覧取得
	output, err := o.listAliases.Execute(ctx, organization_usecase.ListOrganizationAliasesInput{
		OrganizationID: *authCtx.OrganizationID,
	})
	if err != nil {
		return nil, err
	}

	var aliasResponses []oas.OrganizationAlias
	for _, alias := range output.Aliases {
		aliasResponses = append(aliasResponses, alias.ToResponse())
	}

	return &oas.GetOrganizationAliasesOK{
		Aliases: aliasResponses,
	}, nil
}

// CreateOrganizationAlias 組織エイリアス作成
func (o *organizationHandler) CreateOrganizationAlias(ctx context.Context, req *oas.CreateOrganizationAliasReq) (oas.CreateOrganizationAliasRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "organizationHandler.CreateOrganizationAlias")
	defer span.End()

	authCtx, err := o.authorizationService.RequireAdmin(ctx)
	if err != nil {
		return nil, err
	}
	if req == nil || req.AliasName == "" {
		return nil, messages.BadRequestError
	}

	// エイリアス作成
	output, err := o.createAlias.Execute(ctx, authCtx.SessionID, organization_usecase.CreateOrganizationAliasInput{
		OrganizationID: *authCtx.OrganizationID,
		AliasName:      req.AliasName,
	})
	if err != nil {
		return nil, err
	}

	res := output.ToResponse()
	return &res, nil
}

// DeleteOrganizationAlias 組織エイリアス削除
func (o *organizationHandler) DeleteOrganizationAlias(ctx context.Context, params oas.DeleteOrganizationAliasParams) (oas.DeleteOrganizationAliasRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "organizationHandler.DeleteOrganizationAlias")
	defer span.End()

	authCtx, err := o.authorizationService.RequireAdmin(ctx)
	if err != nil {
		return nil, err
	}

	aliasID, err := shared.ParseUUID[organization.OrganizationAlias](params.AliasID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	// エイリアス削除
	err = o.deactivateAlias.Execute(ctx, authCtx.SessionID, organization_usecase.DeactivateOrganizationAliasInput{
		AliasID: aliasID,
	})
	if err != nil {
		return nil, err
	}

	return &oas.DeleteOrganizationAliasOK{}, nil
}

// GetOrganizationUsers 現在の組織のユーザー一覧取得
func (o *organizationHandler) GetOrganizationUsers(ctx context.Context) (oas.GetOrganizationUsersRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "organizationHandler.GetOrganizationUsers")
	defer span.End()

	authCtx, err := o.authorizationService.RequireOrganizationRole(ctx, organization.OrganizationUserRoleMember)
	if err != nil {
		return nil, err
	}

	// クエリを実行
	result, err := o.listUsers.Execute(ctx, organization_query.ListOrganizationUsersInput{
		OrganizationID: *authCtx.OrganizationID,
	})
	if err != nil {
		return nil, err
	}

	organizationUsers := make([]oas.OrganizationUser, 0, len(result.Organizations))
	for _, org := range result.Organizations {
		organizationUsers = append(organizationUsers, org.ToUserResponse())
	}

	return &oas.GetOrganizationUsersOK{
		Users: organizationUsers,
	}, nil
}

// SwitchOrganization 組織を切り替える
func (o *organizationHandler) SwitchOrganization(ctx context.Context, params oas.SwitchOrganizationParams) (oas.SwitchOrganizationRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "organizationHandler.SwitchOrganization")
	defer span.End()

	output, err := o.switchOrganization.Execute(ctx, organization_usecase.SwitchOrganizationUseCaseInput{
		Code: params.Code,
	})
	if err != nil {
		return nil, err
	}

	res := oas.SwitchOrganizationOK{}
	res.SetSetCookie(cookie_utils.EncodeCookies([]*http.Cookie{o.cookieManager.CreateSessionCookie(output.SessionTokenStr)}))
	return &res, nil
}
