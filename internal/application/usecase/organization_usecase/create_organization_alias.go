package organization_usecase

import (
	"context"
	"errors"

	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"go.opentelemetry.io/otel"
)

var ErrSessionNotFound = errors.New("session not found")

// CreateOrganizationAliasInput 入力
type CreateOrganizationAliasInput struct {
	OrganizationID string
	AliasName      string
}

// CreateOrganizationAliasOutput 出力
type CreateOrganizationAliasOutput struct {
	AliasID   string
	AliasName string
}

// CreateOrganizationAliasUseCase エイリアス作成ユースケース
type CreateOrganizationAliasUseCase struct {
	dbManager       *db.DBManager
	sessionRepo     session.SessionRepository
	orgAliasService *service.OrganizationAliasService
}

// NewCreateOrganizationAliasUseCase コンストラクタ
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

// Execute エイリアス作成を実行
func (u *CreateOrganizationAliasUseCase) Execute(
	ctx context.Context,
	sessionID string,
	input CreateOrganizationAliasInput,
) (*CreateOrganizationAliasOutput, error) {
	ctx, span := otel.Tracer("organization_usecase").Start(ctx, "CreateOrganizationAliasUseCase.Execute")
	defer span.End()

	var output *CreateOrganizationAliasOutput
	err := u.dbManager.ExecTx(ctx, func(ctx context.Context) error {
		// セッション取得
		sessID, err := shared.ParseUUID[session.Session](sessionID)
		if err != nil {
			return err
		}
		sess, err := u.sessionRepo.FindBySessionID(ctx, sessID)
		if err != nil {
			return err
		}
		if sess == nil {
			return ErrSessionNotFound
		}

		// 組織ID解析
		orgID, err := shared.ParseUUID[organization.Organization](input.OrganizationID)
		if err != nil {
			return err
		}

		// エイリアス作成
		alias, err := u.orgAliasService.CreateAlias(
			ctx,
			input.AliasName,
			orgID,
			sess.UserID(),
		)
		if err != nil {
			return err
		}

		output = &CreateOrganizationAliasOutput{
			AliasID:   alias.AliasID().String(),
			AliasName: alias.AliasName(),
		}

		return nil
	})

	return output, err
}
