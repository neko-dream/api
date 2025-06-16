package organization_usecase

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"go.opentelemetry.io/otel"
)

// DeactivateOrganizationAliasInput 入力
type DeactivateOrganizationAliasInput struct {
	AliasID string
}

// DeactivateOrganizationAliasUseCase エイリアス論理削除ユースケース
type DeactivateOrganizationAliasUseCase struct {
	dbManager       *db.DBManager
	sessionRepo     session.SessionRepository
	orgAliasService *service.OrganizationAliasService
}

// NewDeactivateOrganizationAliasUseCase コンストラクタ
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

// Execute エイリアス論理削除を実行
func (u *DeactivateOrganizationAliasUseCase) Execute(
	ctx context.Context,
	sessionID string,
	input DeactivateOrganizationAliasInput,
) error {
	ctx, span := otel.Tracer("organization_usecase").Start(ctx, "DeactivateOrganizationAliasUseCase.Execute")
	defer span.End()

	return u.dbManager.ExecTx(ctx, func(ctx context.Context) error {
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

		// エイリアスID解析
		aliasID, err := shared.ParseUUID[organization.OrganizationAlias](input.AliasID)
		if err != nil {
			return err
		}

		// エイリアス論理削除
		return u.orgAliasService.DeactivateAlias(ctx, aliasID, sess.UserID())
	})
}
