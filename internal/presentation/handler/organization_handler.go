package handler

import (
	"context"
	"time"

	"github.com/neko-dream/server/internal/application/query/organization_query"
	"github.com/neko-dream/server/internal/application/usecase/organization_usecase"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	domainservice "github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/presentation/oas"
	"go.opentelemetry.io/otel"
)

type organizationHandler struct {
	create              organization_usecase.CreateOrganizationCommand
	invite              organization_usecase.InviteOrganizationCommand
	add                 organization_usecase.InviteOrganizationForUserCommand
	list                organization_query.ListJoinedOrganizationQuery
	orgRepo             organization.OrganizationRepository
	aliasRepo           organization.OrganizationAliasRepository
	orgUserRepo         organization.OrganizationUserRepository
	aliasService        *domainservice.OrganizationAliasService
	sessionTokenManager session.TokenManager
}

func NewOrganizationHandler(
	create organization_usecase.CreateOrganizationCommand,
	invite organization_usecase.InviteOrganizationCommand,
	add organization_usecase.InviteOrganizationForUserCommand,
	list organization_query.ListJoinedOrganizationQuery,
	orgRepo organization.OrganizationRepository,
	aliasRepo organization.OrganizationAliasRepository,
	orgUserRepo organization.OrganizationUserRepository,
	aliasService *domainservice.OrganizationAliasService,
	sessionTokenManager session.TokenManager,
) oas.OrganizationHandler {
	return &organizationHandler{
		create:              create,
		invite:              invite,
		add:                 add,
		list:                list,
		orgRepo:             orgRepo,
		aliasRepo:           aliasRepo,
		orgUserRepo:         orgUserRepo,
		aliasService:        aliasService,
		sessionTokenManager: sessionTokenManager,
	}
}

// CreateOrganizations implements oas.OrganizationHandler.
func (o *organizationHandler) CreateOrganizations(ctx context.Context, req *oas.CreateOrganizationsReq) (oas.CreateOrganizationsRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "organizationHandler.CreateOrganizations")
	defer span.End()

	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}
	if req == nil {
		return nil, messages.BadRequestError
	}
	userID, err := claim.UserID()
	if err != nil {
		return nil, messages.ForbiddenError
	}

	_, err = o.create.Execute(ctx, organization_usecase.CreateOrganizationInput{
		UserID: userID,
		Name:   req.Name,
		Code:   req.Code,
		Type:   int(req.OrgType),
	})
	if err != nil {
		return nil, err
	}

	res := &oas.CreateOrganizationsOK{}
	return res, nil
}

// InviteOrganization implements oas.OrganizationHandler.
func (o *organizationHandler) InviteOrganization(ctx context.Context, req *oas.InviteOrganizationReq, params oas.InviteOrganizationParams) (oas.InviteOrganizationRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "organizationHandler.InviteOrganization")
	defer span.End()

	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}
	if req == nil {
		return nil, messages.BadRequestError
	}

	userID, err := claim.UserID()
	if err != nil {
		return nil, messages.ForbiddenError
	}
	organizationID, err := shared.ParseUUID[organization.Organization](params.OrganizationID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	_, err = o.invite.Execute(ctx, organization_usecase.InviteOrganizationInput{
		UserID:         userID,
		OrganizationID: organizationID,
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
func (o *organizationHandler) InviteOrganizationForUser(ctx context.Context, req *oas.InviteOrganizationForUserReq, params oas.InviteOrganizationForUserParams) (oas.InviteOrganizationForUserRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "organizationHandler.InviteOrganizationForUser")
	defer span.End()

	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}
	userID, err := claim.UserID()
	if err != nil {
		return nil, messages.ForbiddenError
	}
	organizationID, err := shared.ParseUUID[organization.Organization](params.OrganizationID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	_, err = o.add.Execute(ctx, organization_usecase.InviteOrganizationForUserInput{
		UserID:         userID,
		OrganizationID: organizationID,
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

	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}
	userID, err := claim.UserID()
	if err != nil {
		return nil, messages.ForbiddenError
	}

	res, err := o.list.Execute(ctx, organization_query.ListJoinedOrganizationInput{
		UserID: userID,
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
		args := oas.Organization{
			ID:       org.Organization.ID,
			Name:     org.Organization.Name,
			Code:     org.Organization.Code,
			Type:     org.Organization.OrganizationType,
			Role:     org.OrganizationUser.Role,
			RoleName: org.OrganizationUser.RoleName,
		}
		orgs = append(orgs, args)
	}

	return &oas.GetOrganizationsOK{
		Organizations: orgs,
	}, nil

}

// ValidateOrganizationCode 組織コード検証
func (o *organizationHandler) ValidateOrganizationCode(ctx context.Context, params oas.ValidateOrganizationCodeParams) (oas.ValidateOrganizationCodeRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "organizationHandler.ValidateOrganizationCode")
	defer span.End()

	// Find organization by code
	org, err := o.orgRepo.FindByCode(ctx, params.Code)
	if err != nil {
		// Organization not found
		return &oas.ValidateOrganizationCodeOK{
			Valid: false,
		}, nil
	}

	// Organization found
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
func (o *organizationHandler) GetOrganizationAliases(ctx context.Context, params oas.GetOrganizationAliasesParams) (oas.GetOrganizationAliasesRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "organizationHandler.GetOrganizationAliases")
	defer span.End()

	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}
	userID, err := claim.UserID()
	if err != nil {
		return nil, messages.ForbiddenError
	}

	organizationID, err := shared.ParseUUID[organization.Organization](params.OrganizationID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	// 権限チェック - ユーザーがAdmin以上の権限を持っているか確認
	orgUser, err := o.orgUserRepo.FindByOrganizationIDAndUserID(ctx, organizationID, userID)
	if err != nil {
		return nil, messages.ForbiddenError
	}
	if orgUser.Role < organization.OrganizationUserRoleAdmin {
		return nil, messages.ForbiddenError
	}

	// エイリアス一覧取得
	aliases, err := o.aliasRepo.FindActiveByOrganizationID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	var aliasResponses []oas.GetOrganizationAliasesOKAliasesItem
	for _, alias := range aliases {
		aliasResponses = append(aliasResponses, oas.GetOrganizationAliasesOKAliasesItem{
			AliasID:   alias.AliasID().String(),
			AliasName: alias.AliasName(),
			CreatedAt: alias.CreatedAt().Format(time.RFC3339),
		})
	}

	return &oas.GetOrganizationAliasesOK{
		Aliases: aliasResponses,
	}, nil
}

// CreateOrganizationAlias 組織エイリアス作成
func (o *organizationHandler) CreateOrganizationAlias(ctx context.Context, req *oas.CreateOrganizationAliasReq, params oas.CreateOrganizationAliasParams) (oas.CreateOrganizationAliasRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "organizationHandler.CreateOrganizationAlias")
	defer span.End()

	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}
	userID, err := claim.UserID()
	if err != nil {
		return nil, messages.ForbiddenError
	}

	organizationID, err := shared.ParseUUID[organization.Organization](params.OrganizationID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	// エイリアス作成
	alias, err := o.aliasService.CreateAlias(ctx, req.AliasName, organizationID, userID)
	if err != nil {
		return nil, err
	}

	return &oas.CreateOrganizationAliasOK{
		AliasID:   alias.AliasID().String(),
		AliasName: alias.AliasName(),
		CreatedAt: alias.CreatedAt().Format(time.RFC3339),
	}, nil
}

// DeleteOrganizationAlias 組織エイリアス削除
func (o *organizationHandler) DeleteOrganizationAlias(ctx context.Context, params oas.DeleteOrganizationAliasParams) (oas.DeleteOrganizationAliasRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "organizationHandler.DeleteOrganizationAlias")
	defer span.End()

	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}
	userID, err := claim.UserID()
	if err != nil {
		return nil, messages.ForbiddenError
	}

	organizationID, err := shared.ParseUUID[organization.Organization](params.OrganizationID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	aliasID, err := shared.ParseUUID[organization.OrganizationAlias](params.AliasID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	// 権限チェック - ユーザーがAdmin以上の権限を持っているか確認
	orgUser, err := o.orgUserRepo.FindByOrganizationIDAndUserID(ctx, organizationID, userID)
	if err != nil {
		return nil, messages.ForbiddenError
	}
	if orgUser.Role < organization.OrganizationUserRoleAdmin {
		return nil, messages.ForbiddenError
	}

	// エイリアス削除
	err = o.aliasService.DeactivateAlias(ctx, aliasID, userID)
	if err != nil {
		return nil, err
	}

	return &oas.DeleteOrganizationAliasOK{}, nil
}
