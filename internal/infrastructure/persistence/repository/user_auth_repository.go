package repository

import (
	"context"
	"database/sql"
	"time"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type userAuthRepository struct {
	*db.DBManager
}

// NewUserAuthRepository creates a new UserAuthRepository
func NewUserAuthRepository(dbManager *db.DBManager) user.UserAuthRepository {
	return &userAuthRepository{
		DBManager: dbManager,
	}
}

// FindByUserID implements user.UserAuthRepository
func (r *userAuthRepository) FindByUserID(ctx context.Context, userID shared.UUID[user.User]) (*user.UserAuth, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "userAuthRepository.FindByUserID")
	defer span.End()

	result, err := r.GetQueries(ctx).GetUserAuthByUserID(ctx, userID.UUID())
	if err != nil {
		utils.HandleError(ctx, err, "GetUserAuthByUserID")
		return nil, errtrace.Wrap(err)
	}

	return r.toDomainModel(result.UserAuth), nil
}

// FindBySubject implements user.UserAuthRepository
func (r *userAuthRepository) FindBySubject(ctx context.Context, subject string) (*user.UserAuth, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "userAuthRepository.FindBySubject")
	defer span.End()

	result, err := r.GetQueries(ctx).GetUserBySubject(ctx, subject)
	if err != nil {
		utils.HandleError(ctx, err, "GetUserBySubject")
		return nil, errtrace.Wrap(err)
	}

	return r.toDomainModel(result.UserAuth), nil
}

// Update implements user.UserAuthRepository
func (r *userAuthRepository) Update(ctx context.Context, userAuth *user.UserAuth) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "userAuthRepository.Update")
	defer span.End()

	var withdrawalDate sql.NullTime
	if userAuth.WithdrawalDate() != nil {
		withdrawalDate = sql.NullTime{
			Time:  *userAuth.WithdrawalDate(),
			Valid: true,
		}
	}

	err := r.GetQueries(ctx).WithdrawUser(ctx, model.WithdrawUserParams{
		UserID:         userAuth.UserID().UUID(),
		WithdrawalDate: withdrawalDate,
	})
	if err != nil {
		utils.HandleError(ctx, err, "WithdrawUser")
		return errtrace.Wrap(err)
	}

	return nil
}

// CheckReregistrationAllowed implements user.UserAuthRepository
func (r *userAuthRepository) CheckReregistrationAllowed(ctx context.Context, subject string) (bool, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "userAuthRepository.CheckReregistrationAllowed")
	defer span.End()

	allowed, err := r.GetQueries(ctx).CheckReregistrationAllowed(ctx, subject)
	if err != nil {
		utils.HandleError(ctx, err, "CheckReregistrationAllowed")
		return false, errtrace.Wrap(err)
	}

	return allowed, nil
}

// toDomainModel converts database model to domain model
func (r *userAuthRepository) toDomainModel(dbModel model.UserAuth) *user.UserAuth {
	var withdrawalDate *time.Time
	if dbModel.WithdrawalDate.Valid {
		withdrawalDate = &dbModel.WithdrawalDate.Time
	}

	return user.NewUserAuthWithWithdrawal(
		shared.UUID[user.UserAuth](dbModel.UserAuthID),
		shared.UUID[user.User](dbModel.UserID),
		dbModel.Provider,
		dbModel.Subject,
		dbModel.IsVerified,
		dbModel.CreatedAt,
		withdrawalDate,
	)
}