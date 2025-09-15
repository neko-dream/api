package organization_usecase

import (
	"context"

	"github.com/neko-dream/api/internal/domain/model/organization"
	"github.com/neko-dream/api/internal/domain/model/session"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/service"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"go.opentelemetry.io/otel"
)

type DeactivateOrganizationAliasInput struct {
	AliasID shared.UUID[organization.OrganizationAlias]
}

type DeactivateOrganizationAliasUseCase struct {
	dbManager       *db.DBManager
	sessionRepo     session.SessionRepository
	orgAliasService *service.OrganizationAliasService
}

func NewDeactivateOrganizationAliasUseCase(
	dbManager *db.DBManager,
	sessionRepo session.SessionRepository,
	orgAliasService *service.OrganizationAliasService,
) *DeactivateOrganizationAliasUseCase {
	return &DeactivateOrganizationAliasUseCase{
		dbManager:       dbManager,
		sessionRepo:     sessionRepo,
		orgAliasService: orgAliasService,
	}
}

func (u *DeactivateOrganizationAliasUseCase) Execute(
	ctx context.Context,
	sessionID shared.UUID[session.Session],
	input DeactivateOrganizationAliasInput,
) error {
	ctx, span := otel.Tracer("organization_usecase").Start(ctx, "DeactivateOrganizationAliasUseCase.Execute")
	defer span.End()

	return u.dbManager.ExecTx(ctx, func(ctx context.Context) error {
		sess, err := u.sessionRepo.FindBySessionID(ctx, sessionID)
		if err != nil {
			return err
		}
		if sess == nil {
			return ErrSessionNotFound
		}

		return u.orgAliasService.DeactivateAlias(ctx, input.AliasID, sess.UserID())
	})
}
