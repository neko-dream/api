package auth_usecase

import (
	"context"
	"database/sql"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"go.opentelemetry.io/otel"
)

type (
	// WithdrawUser interface for user withdrawal
	WithdrawUser interface {
		Execute(ctx context.Context, input WithdrawUserInput) error
	}

	// WithdrawUserInput input for user withdrawal
	WithdrawUserInput struct {
		UserID shared.UUID[user.User]
	}

	// withdrawUserInteractor implementation
	withdrawUserInteractor struct {
		*db.DBManager
		userRepository     user.UserRepository
		userAuthRepository user.UserAuthRepository
		sessionService     session.SessionService
	}
)

// NewWithdrawUser creates a new WithdrawUser use case
func NewWithdrawUser(
	dbm *db.DBManager,
	userRepository user.UserRepository,
	userAuthRepository user.UserAuthRepository,
	sessionService session.SessionService,
) WithdrawUser {
	return &withdrawUserInteractor{
		DBManager:          dbm,
		userRepository:     userRepository,
		userAuthRepository: userAuthRepository,
		sessionService:     sessionService,
	}
}

// Execute performs user withdrawal
func (w *withdrawUserInteractor) Execute(ctx context.Context, input WithdrawUserInput) error {
	ctx, span := otel.Tracer("auth_usecase").Start(ctx, "withdrawUserInteractor.Execute")
	defer span.End()

	return w.ExecTx(ctx, func(ctx context.Context) error {
		// Get user and user auth
		user, err := w.userRepository.FindByID(ctx, input.UserID)
		if err != nil {
			return errtrace.Wrap(err)
		}

		userAuth, err := w.userAuthRepository.FindByUserID(ctx, input.UserID)
		if err != nil {
			return errtrace.Wrap(err)
		}

		// Perform withdrawal on domain objects
		user.Withdraw(ctx)
		userAuth.Withdraw(ctx)

		// Update in database using direct queries
		if err := w.GetQueries(ctx).WithdrawUser(ctx, model.WithdrawUserParams{
			UserID:         input.UserID.UUID(),
			WithdrawalDate: sql.NullTime{Time: *userAuth.WithdrawalDate(), Valid: true},
		}); err != nil {
			return errtrace.Wrap(err)
		}

		if err := w.GetQueries(ctx).AnonymizeUser(ctx, input.UserID.UUID()); err != nil {
			return errtrace.Wrap(err)
		}

		// Deactivate all user sessions
		if err := w.sessionService.DeactivateUserSessions(ctx, input.UserID); err != nil {
			return errtrace.Wrap(err)
		}

		return nil
	})
}