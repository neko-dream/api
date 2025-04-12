package repository

import (
	"context"
	"database/sql"

	password_auth "github.com/neko-dream/server/internal/domain/model/auth/password"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type passwordAuthRepository struct {
	*db.DBManager
}

func NewPasswordAuthRepository(dbManager *db.DBManager) password_auth.PasswordAuthRepository {
	return &passwordAuthRepository{
		DBManager: dbManager,
	}
}

// CreatePasswordAuth
func (p *passwordAuthRepository) CreatePasswordAuth(ctx context.Context, userID shared.UUID[user.User], passwordHash string, salt string) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "passwordAuthRepository.CreatePasswordAuth")
	defer span.End()

	if err := p.DBManager.GetQueries(ctx).CreatePasswordAuth(ctx, model.CreatePasswordAuthParams{
		PasswordAuthID: shared.NewUUID[password_auth.PasswordAuth]().UUID(),
		UserID:         userID.UUID(),
		PasswordHash:   passwordHash,
		Salt:           sql.NullString{String: salt, Valid: true},
		LastChanged:    clock.Now(ctx),
		CreatedAt:      clock.Now(ctx),
		UpdatedAt:      clock.Now(ctx),
	}); err != nil {
		utils.HandleError(ctx, err, "CreatePasswordAuth")
		return err
	}

	return nil
}

// DeletePasswordAuth
func (p *passwordAuthRepository) DeletePasswordAuth(ctx context.Context, userID shared.UUID[user.User]) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "passwordAuthRepository.DeletePasswordAuth")
	defer span.End()

	if err := p.DBManager.GetQueries(ctx).DeletePasswordAuth(ctx, userID.UUID()); err != nil {
		utils.HandleError(ctx, err, "DeletePasswordAuth")
		return err
	}

	return nil
}

// GetPasswordAuthByUserID
func (p *passwordAuthRepository) GetPasswordAuthByUserID(ctx context.Context, userID shared.UUID[user.User]) (password_auth.PasswordAuth, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "passwordAuthRepository.GetPasswordAuthByUserID")
	defer span.End()

	passwordAuth, err := p.DBManager.GetQueries(ctx).GetPasswordAuthByUserId(ctx, userID.UUID())
	if err != nil {
		utils.HandleError(ctx, err, "GetPasswordAuthByUserID")
		return password_auth.PasswordAuth{}, err
	}
	return password_auth.PasswordAuth{
		PasswordAuthID: shared.UUID[password_auth.PasswordAuth](passwordAuth.PasswordAuth.PasswordAuthID),
		UserID:         shared.UUID[user.User](passwordAuth.PasswordAuth.UserID),
		PasswordHash:   passwordAuth.PasswordAuth.PasswordHash,
		Salt:           passwordAuth.PasswordAuth.Salt.String,
		LastChanged:    passwordAuth.PasswordAuth.LastChanged,
		CreatedAt:      passwordAuth.PasswordAuth.CreatedAt,
		UpdatedAt:      passwordAuth.PasswordAuth.UpdatedAt,
	}, nil
}

// UpdatePasswordAuth
func (p *passwordAuthRepository) UpdatePasswordAuth(ctx context.Context, userID shared.UUID[user.User], passwordHash string, salt string) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "passwordAuthRepository.UpdatePasswordAuth")
	defer span.End()

	if err := p.DBManager.GetQueries(ctx).UpdatePasswordAuth(ctx, model.UpdatePasswordAuthParams{
		UserID:       userID.UUID(),
		PasswordHash: passwordHash,
		Salt:         sql.NullString{String: salt, Valid: true},
		LastChanged:  clock.Now(ctx),
	}); err != nil {
		utils.HandleError(ctx, err, "UpdatePasswordAuth")
		return err
	}
	return nil
}
