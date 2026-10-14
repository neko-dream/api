package repository

import (
	"context"

	password_auth "github.com/neko-dream/api/internal/domain/model/auth/password"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
)

type passwordSetupTokenRepository struct{ db *db.DBManager }

func NewPasswordSetupTokenRepository(db *db.DBManager) password_auth.SetupTokenRepository {
	return &passwordSetupTokenRepository{db: db}
}

func (r *passwordSetupTokenRepository) LockEmailSubject(ctx context.Context, subject string) error {
	return r.db.GetQueries(ctx).LockPasswordSetupEmail(ctx, subject)
}

func (r *passwordSetupTokenRepository) LockUser(ctx context.Context, id shared.UUID[user.User]) error {
	_, err := r.db.GetQueries(ctx).LockPasswordSetupUser(ctx, id.UUID())
	return err
}

func (r *passwordSetupTokenRepository) FindUserIDByTokenHash(ctx context.Context, hash []byte) (shared.UUID[user.User], error) {
	id, err := r.db.GetQueries(ctx).FindPasswordSetupUserByTokenHash(ctx, hash)
	return shared.UUID[user.User](id), err
}

func (r *passwordSetupTokenRepository) FindByUserID(ctx context.Context, id shared.UUID[user.User]) (password_auth.SetupToken, error) {
	row, err := r.db.GetQueries(ctx).GetPasswordSetupToken(ctx, id.UUID())
	if err != nil {
		return password_auth.SetupToken{}, err
	}
	return password_auth.SetupToken{UserID: id, TokenHash: row.TokenHash, EmailSubject: row.EmailSubject, CreatedAt: row.CreatedAt, ExpiresAt: row.ExpiresAt}, nil
}

func (r *passwordSetupTokenRepository) Save(ctx context.Context, token password_auth.SetupToken) error {
	return r.db.GetQueries(ctx).SavePasswordSetupToken(ctx, model.SavePasswordSetupTokenParams{UserID: token.UserID.UUID(), TokenHash: token.TokenHash, EmailSubject: token.EmailSubject, CreatedAt: token.CreatedAt, ExpiresAt: token.ExpiresAt})
}

func (r *passwordSetupTokenRepository) Delete(ctx context.Context, id shared.UUID[user.User]) error {
	return r.db.GetQueries(ctx).DeletePasswordSetupToken(ctx, id.UUID())
}
