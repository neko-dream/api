package organization_usecase

import (
	"context"
	"errors"
	"time"

	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/internal/presentation/oas"
	"go.opentelemetry.io/otel"
)

var ErrSessionNotFound = errors.New("session not found")

type CreateOrganizationAliasInput struct {
	OrganizationID shared.UUID[organization.Organization]
	AliasName      string
}

type CreateOrganizationAliasOutput struct {
	AliasID   string
	AliasName string
	CreatedAt time.Time
}

func (o *CreateOrganizationAliasOutput) ToResponse() oas.OrganizationAlias {
	return oas.OrganizationAlias{
		AliasID:   o.AliasID,
		AliasName: o.AliasName,
		CreatedAt: o.CreatedAt.Format(time.RFC3339),
	}
}

type CreateOrganizationAliasUseCase struct {
	dbManager       *db.DBManager
	sessionRepo     session.SessionRepository
	orgAliasService *service.OrganizationAliasService
}

func NewCreateOrganizationAliasUseCase(
	dbManager *db.DBManager,
	sessionRepo session.SessionRepository,
	orgAliasService *service.OrganizationAliasService,
) *CreateOrganizationAliasUseCase {
	return &CreateOrganizationAliasUseCase{
		dbManager:       dbManager,
		sessionRepo:     sessionRepo,
		orgAliasService: orgAliasService,
	}
}

func (u *CreateOrganizationAliasUseCase) Execute(
	ctx context.Context,
	sessionID shared.UUID[session.Session],
	input CreateOrganizationAliasInput,
) (*CreateOrganizationAliasOutput, error) {
	ctx, span := otel.Tracer("organization_usecase").Start(ctx, "CreateOrganizationAliasUseCase.Execute")
	defer span.End()

	var output *CreateOrganizationAliasOutput
	err := u.dbManager.ExecTx(ctx, func(ctx context.Context) error {
		sess, err := u.sessionRepo.FindBySessionID(ctx, sessionID)
		if err != nil {
			return err
		}
		if sess == nil {
			return ErrSessionNotFound
		}

		alias, err := u.orgAliasService.CreateAlias(
			ctx,
			input.AliasName,
			input.OrganizationID,
			sess.UserID(),
		)
		if err != nil {
			return err
		}

		output = &CreateOrganizationAliasOutput{
			AliasID:   alias.AliasID().String(),
			AliasName: alias.AliasName(),
			CreatedAt: alias.CreatedAt(),
		}

		return nil
	})

	return output, err
}
