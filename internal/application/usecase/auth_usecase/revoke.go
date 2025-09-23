package auth_usecase

import (
	"context"

	"braces.dev/errtrace"
	"github.com/neko-dream/api/internal/domain/model/session"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"go.opentelemetry.io/otel"
)

type (
	Revoke interface {
		Execute(context.Context, RevokeInput) (RevokeOutput, error)
	}

	RevokeInput struct {
		SessID shared.UUID[session.Session]
	}

	RevokeOutput struct {
	}

	revokeInteractor struct {
		*db.DBManager
		*config.Config
		session.SessionRepository
	}
)

func NewRevoke(
	tm *db.DBManager,
	config *config.Config,
	sessRepository session.SessionRepository,
) Revoke {
	return &revokeInteractor{
		DBManager:         tm,
		Config:            config,
		SessionRepository: sessRepository,
	}
}

func (a *revokeInteractor) Execute(ctx context.Context, input RevokeInput) (RevokeOutput, error) {
	ctx, span := otel.Tracer("auth_usecase").Start(ctx, "revokeInteractor.Execute")
	defer span.End()

	sess, err := a.SessionRepository.FindBySessionID(ctx, input.SessID)
	if err != nil {
		return RevokeOutput{}, errtrace.Wrap(err)
	}

	// アクティブなセッションの場合のみ無効化
	if sess.IsActive(ctx) {
		sess.Deactivate(ctx)
		_, err = a.SessionRepository.Update(ctx, *sess)
		if err != nil {
			return RevokeOutput{}, errtrace.Wrap(err)
		}
	}

	return RevokeOutput{}, nil
}
