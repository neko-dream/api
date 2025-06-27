package auth_usecase

import (
	"context"
	"time"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/domain/service"
	organizationService "github.com/neko-dream/server/internal/domain/service/organization"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type (
	// AuthLogin OAuth認証開始時のリダイレクトURL生成
	AuthLogin interface {
		// Execute 認証プロバイダーの認可URLとstateを生成し返す
		//
		// ここでRegistrationURLがある場合はログイン、ない場合は登録
		Execute(context.Context, AuthLoginInput) (AuthLoginOutput, error)
	}

	// AuthLoginInput 認証プロバイダー名を受け取る
	AuthLoginInput struct {
		Provider    string
		RedirectURL string
		// 登録URLがある場合はログイン
		RegistrationURL *string
		// 組織コード（組織ログインの場合）
		OrganizationCode *string
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
		organizationService organizationService.OrganizationService
	}
)

// NewAuthLogin
func NewAuthLogin(
	tm *db.DBManager,
	config *config.Config,
	authService service.AuthService,
	authProviderFactory auth.AuthProviderFactory,
	stateRepository auth.StateRepository,
	organizationService organizationService.OrganizationService,
) AuthLogin {
	return &authLoginInteractor{
		DBManager:           tm,
		Config:              config,
		AuthService:         authService,
		authProviderFactory: authProviderFactory,
		stateRepository:     stateRepository,
		organizationService: organizationService,
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

		stateString, err := a.GenerateState(ctx)
		if err != nil {
			utils.HandleError(ctx, err, "GenerateState") // state生成失敗
			return errtrace.Wrap(err)
		}

		// 組織コードが指定されている場合、組織IDを取得
		organizationID, err := a.organizationService.ResolveOrganizationIDFromCode(ctx, input.OrganizationCode)
		if err != nil {
			utils.HandleError(ctx, err, "ResolveOrganizationIDFromCode")
			return errtrace.Wrap(err)
		}
		if input.OrganizationCode != nil && organizationID == nil {
			// 組織コードが指定されているが、組織が見つからない場合はエラー
			return messages.OrganizationNotFound
		}

		state := auth.NewState(stateString, input.Provider, input.RedirectURL, time.Now().Add(auth.StateExpirationDuration), input.RegistrationURL, organizationID)
		// stateをDBに保存（cookieじゃないのは一部ブラウザでうまく動作しないため）
		err = a.stateRepository.Create(ctx, state)
		if err != nil {
			utils.HandleError(ctx, err, "CreateAuthState") // DB保存失敗
			return errtrace.Wrap(err)
		}

		s = stateString
		au = provider.GetAuthorizationURL(ctx, stateString)
		return nil
	}); err != nil {
		return AuthLoginOutput{}, errtrace.Wrap(err)
	}

	return AuthLoginOutput{
		RedirectURL: au,
		State:       s,
	}, nil
}
