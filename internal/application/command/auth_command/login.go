package auth_command

import (
	"context"
	"time"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type (
	// AuthLogin OAuth認証開始時のリダイレクトURL生成
	AuthLogin interface {
		Execute(context.Context, AuthLoginInput) (AuthLoginOutput, error)
	}

	// AuthLoginInput 認証プロバイダー名を受け取る
	AuthLoginInput struct {
		Provider    string
		RedirectURL string
	}

	// AuthLoginOutput
	AuthLoginOutput struct {
		RedirectURL string
		State       string
	}

	// authLoginInteractor 認証開始ユースケースの実装
	authLoginInteractor struct {
		*db.DBManager
		*config.Config
		service.AuthService
		authProviderFactory auth.AuthProviderFactory
		stateRepository     auth.StateRepository
	}
)

// NewAuthLogin
func NewAuthLogin(
	tm *db.DBManager,
	config *config.Config,
	authService service.AuthService,
	authProviderFactory auth.AuthProviderFactory,
	stateRepository auth.StateRepository,
) AuthLogin {
	return &authLoginInteractor{
		DBManager:           tm,
		Config:              config,
		AuthService:         authService,
		authProviderFactory: authProviderFactory,
		stateRepository:     stateRepository,
	}
}

// Execute 認証プロバイダーの認可URLとstateを生成し返す
// 1. プロバイダーのインスタンス生成
// 2. state生成・DB保存（CSRF対策）
// 3. 認可URL生成
func (a *authLoginInteractor) Execute(ctx context.Context, input AuthLoginInput) (AuthLoginOutput, error) {
	ctx, span := otel.Tracer("auth_usecase").Start(ctx, "authLoginInteractor.Execute")
	defer span.End()

	var (
		s  string
		au string
	)

	if err := a.ExecTx(ctx, func(ctx context.Context) error {
		provider, err := a.authProviderFactory.NewAuthProvider(ctx, input.Provider)
		if err != nil {
			utils.HandleError(ctx, err, "NewAuthProvider") // プロバイダー生成失敗
			return errtrace.Wrap(err)
		}

		state, err := a.GenerateState(ctx)
		if err != nil {
			utils.HandleError(ctx, err, "GenerateState") // state生成失敗
			return errtrace.Wrap(err)
		}

		// stateをDBに保存（CSRF対策）
		err = a.stateRepository.Create(ctx, &auth.State{
			State:       state,
			Provider:    input.Provider,
			RedirectURL: input.RedirectURL,
			ExpiresAt:   time.Now().Add(auth.StateExpirationDuration), // 15分後に期限切れ
		})
		if err != nil {
			utils.HandleError(ctx, err, "CreateAuthState") // DB保存失敗
			return errtrace.Wrap(err)
		}

		s = state
		au = provider.GetAuthorizationURL(ctx, state)
		return nil
	}); err != nil {
		return AuthLoginOutput{}, errtrace.Wrap(err)
	}

	return AuthLoginOutput{
		RedirectURL: au,
		State:       s,
	}, nil
}
