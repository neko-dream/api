package auth_usecase

import (
	"context"
	"database/sql"
	"errors"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type (
	LoginForDev interface {
		Execute(context.Context, LoginForDevInput) (LoginForDevOutput, error)
	}

	LoginForDevInput struct {
		Subject string
		// 組織コード（組織ログインの場合）
		OrganizationCode *string
	}

	LoginForDevOutput struct {
		Token string
	}

	loginForDevInteractor struct {
		*db.DBManager
		*config.Config
		service.AuthService
		session.SessionRepository
		session.SessionService
		session.TokenManager
		organizationRepo     organization.OrganizationRepository
		organizationUserRepo organization.OrganizationUserRepository
	}
)

func NewLoginForDev(
	tm *db.DBManager,
	config *config.Config,
	authService service.AuthService,
	sessionRepository session.SessionRepository,
	sessionService session.SessionService,
	tokenManager session.TokenManager,
	organizationRepo organization.OrganizationRepository,
	organizationUserRepo organization.OrganizationUserRepository,
) LoginForDev {
	return &loginForDevInteractor{
		DBManager:            tm,
		Config:               config,
		AuthService:          authService,
		SessionRepository:    sessionRepository,
		SessionService:       sessionService,
		TokenManager:         tokenManager,
		organizationRepo:     organizationRepo,
		organizationUserRepo: organizationUserRepo,
	}
}

func (a *loginForDevInteractor) Execute(ctx context.Context, input LoginForDevInput) (LoginForDevOutput, error) {
	ctx, span := otel.Tracer("auth_usecase").Start(ctx, "loginForDevInteractor.Execute")
	defer span.End()

	if a.Config.Env == config.PROD {
		utils.HandleError(ctx, errtrace.New("このエンドポイントは開発でのみ有効です。"), "failed to login for dev")
		return LoginForDevOutput{}, errtrace.New("このエンドポイントは開発でのみ有効です。")
	}

	var (
		tok string
	)

	if err := a.ExecTx(ctx, func(ctx context.Context) error {
		newUser, err := a.AuthService.Authenticate(ctx, "dev", input.Subject)
		if err != nil {
			return err
		}
		if newUser != nil {
			if err := a.SessionService.DeactivateUserSessions(ctx, newUser.UserID()); err != nil {
				utils.HandleError(ctx, err, "failed to deactivate user sessions")
			}

			// 組織コードが指定されている場合、ユーザーを組織に追加
			if input.OrganizationCode != nil && *input.OrganizationCode != "" {
				org, err := a.organizationRepo.FindByCode(ctx, *input.OrganizationCode)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						utils.HandleError(ctx, err, "FindOrganizationByCode")
						return errtrace.Wrap(err)
					}
					// 組織が見つからない場合は無視して通常のログインを続行
				} else {
					// ユーザーが既に組織のメンバーかチェック
					_, err := a.organizationUserRepo.FindByOrganizationIDAndUserID(ctx, org.OrganizationID, newUser.UserID())
					if err != nil {
						if errors.Is(err, sql.ErrNoRows) {
							// ユーザーが組織のメンバーでない場合
							// 開発環境でも組織コードを知っているだけでは自動的にメンバーにはしない
							// 招待制であるべきなので、ここでは何もしない
						} else {
							utils.HandleError(ctx, err, "failed to check organization membership")
							// エラーが発生してもログインは続行
						}
					}
					// 既にメンバーの場合は通常通りログイン
				}
			}
		}

		// セッション作成（組織IDの処理）
		var sess *session.Session
		var orgID *shared.UUID[any]
		if input.OrganizationCode != nil && *input.OrganizationCode != "" {
			// 組織コードから組織を検索（エラーは無視）
			if org, err := a.organizationRepo.FindByCode(ctx, *input.OrganizationCode); err == nil {
				// ユーザーが組織のメンバーかチェック
				if _, err := a.organizationUserRepo.FindByOrganizationIDAndUserID(ctx, org.OrganizationID, newUser.UserID()); err == nil {
					// メンバーの場合のみ組織IDを設定
					orgIDAny := shared.UUID[any](org.OrganizationID)
					orgID = &orgIDAny
				}
			}
		}

		if orgID != nil {
			sess = session.NewSessionWithOrganization(
				shared.NewUUID[session.Session](),
				newUser.UserID(),
				newUser.Provider(),
				session.SESSION_ACTIVE,
				*session.NewExpiresAt(ctx),
				clock.Now(ctx),
				orgID,
			)
		} else {
			sess = session.NewSession(
				shared.NewUUID[session.Session](),
				newUser.UserID(),
				newUser.Provider(),
				session.SESSION_ACTIVE,
				*session.NewExpiresAt(ctx),
				clock.Now(ctx),
			)
		}

		if _, err := a.SessionRepository.Create(ctx, *sess); err != nil {
			utils.HandleError(ctx, err, "failed to create session")
			return errtrace.Wrap(err)
		}

		token, err := a.TokenManager.Generate(ctx, *newUser, sess.SessionID())
		if err != nil {
			utils.HandleError(ctx, err, "failed to generate token")
			return errtrace.Wrap(err)
		}

		tok = token

		return nil
	}); err != nil {
		return LoginForDevOutput{}, errtrace.Wrap(err)
	}

	return LoginForDevOutput{
		Token: tok,
	}, nil
}
