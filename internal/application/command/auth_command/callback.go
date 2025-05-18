package auth_command

import (
	"context"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type (
	// AuthCallback
	AuthCallback interface {
		Execute(ctx context.Context, input CallbackInput) (CallbackOutput, error)
	}

	// CallbackInput コールバック時の入力情報
	CallbackInput struct {
		Provider string
		Code     string
		State    string
	}

	// CallbackOutput 認証トークン
	CallbackOutput struct {
		Token       string
		RedirectURL string
	}

	// authCallbackInteractor
	authCallbackInteractor struct {
		*db.DBManager
		*config.Config
		service.AuthService
		session.SessionRepository
		session.SessionService
		session.TokenManager
		auth.StateRepository
	}
)

// NewAuthCallback
func NewAuthCallback(
	tm *db.DBManager,
	config *config.Config,
	authService service.AuthService,
	sessionRepository session.SessionRepository,
	sessionService session.SessionService,
	tokenManager session.TokenManager,
	stateRepository auth.StateRepository,
) AuthCallback {
	return &authCallbackInteractor{
		DBManager:         tm,
		Config:            config,
		AuthService:       authService,
		SessionRepository: sessionRepository,
		SessionService:    sessionService,
		TokenManager:      tokenManager,
		StateRepository:   stateRepository,
	}
}

// Execute コールバック時の認証・セッション生成処理を行う
// 1. stateの検証（DBから取得・有効性チェック）
// 2. stateの削除（ワンタイム性担保）
// 3. ユーザー認証・セッション生成・トークン発行
func (u *authCallbackInteractor) Execute(ctx context.Context, input CallbackInput) (CallbackOutput, error) {
	ctx, span := otel.Tracer("auth_usecase").Start(ctx, "authCallbackInteractor.Execute")
	defer span.End()

	var tokenRes string
	var redirectURL string
	if err := u.ExecTx(ctx, func(ctx context.Context) error {
		// stateの検証
		state, err := u.StateRepository.Get(ctx, input.State)
		if err != nil {
			utils.HandleError(ctx, err, "stateの取得に失敗しました") // state取得失敗
			return errtrace.Wrap(err)
		}

		if err := state.Validate(input.State); err != nil {
			utils.HandleError(ctx, err, "stateが不正です") // state検証失敗
			return errtrace.Wrap(err)
		}
		redirectURL = state.RedirectURL

		// stateの削除（ワンタイム性担保）
		if err := u.StateRepository.Delete(ctx, input.State); err != nil {
			utils.HandleError(ctx, err, "stateの削除に失敗しました") // state削除失敗
			return errtrace.Wrap(err)
		}

		user, err := u.AuthService.Authenticate(ctx, input.Provider, input.Code)
		if err != nil {
			utils.HandleError(ctx, err, "ユーザー認証に失敗しました") // ユーザー認証失敗
			return errtrace.Wrap(err)
		}
		if user != nil {
			if err := u.SessionService.DeactivateUserSessions(ctx, user.UserID()); err != nil {
				utils.HandleError(ctx, err, "既存セッションの無効化に失敗しました")
			}
		}

		sess := session.NewSession(
			shared.NewUUID[session.Session](),
			user.UserID(),
			user.Provider(),
			session.SESSION_ACTIVE,
			*session.NewExpiresAt(ctx),
			clock.Now(ctx),
		)

		if _, err := u.SessionRepository.Create(ctx, *sess); err != nil {
			utils.HandleError(ctx, err, "セッション生成に失敗しました") // セッション生成失敗
			return errtrace.Wrap(err)
		}

		token, err := u.TokenManager.Generate(ctx, *user, sess.SessionID())
		if err != nil {
			utils.HandleError(ctx, err, "トークン生成に失敗しました") // トークン生成失敗
			return errtrace.Wrap(err)
		}

		tokenRes = token

		return nil
	}); err != nil {
		return CallbackOutput{}, err
	}

	return CallbackOutput{
		Token:       tokenRes,
		RedirectURL: redirectURL,
	}, nil
}
