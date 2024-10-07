package auth_usecase

import (
	"context"
	"net/http"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/db"
)

type (
	RevokeUseCase interface {
		Execute(context.Context, RevokeInput) (RevokeOutput, error)
	}

	RevokeInput struct {
		SessID shared.UUID[session.Session]
	}

	RevokeOutput struct {
		Cookies []*http.Cookie
	}

	revokeInteractor struct {
		*db.DBManager
		*config.Config
		session.SessionRepository
	}
)

func NewRevokeUseCase(
	tm *db.DBManager,
	config *config.Config,
	sessRepository session.SessionRepository,
) RevokeUseCase {
	return &revokeInteractor{
		DBManager:         tm,
		Config:            config,
		SessionRepository: sessRepository,
	}
}

func (a *revokeInteractor) Execute(ctx context.Context, input RevokeInput) (RevokeOutput, error) {
	sess, err := a.SessionRepository.FindBySessionID(ctx, input.SessID)
	if err != nil {
		return RevokeOutput{}, errtrace.Wrap(err)
	}

	// アクティブなセッションの場合のみ無効化
	if sess.IsActive() {
		sess.Deactivate()
		_, err = a.SessionRepository.Update(ctx, *sess)
		if err != nil {
			return RevokeOutput{}, errtrace.Wrap(err)
		}
	}

	// Revoke session cookie
	sessionCookie := http.Cookie{
		Name:     "SessionId",
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Domain:   a.DOMAIN,
		MaxAge:   -1,
	}
	return RevokeOutput{
		Cookies: []*http.Cookie{&sessionCookie},
	}, nil
}
