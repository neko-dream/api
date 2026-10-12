package repository

import (
	"context"
	"time"

	password_auth "github.com/neko-dream/api/internal/domain/model/auth/password"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
)

type passwordLoginChallengeRepository struct{ manager *db.DBManager }

func NewPasswordLoginChallengeRepository(manager *db.DBManager) password_auth.LoginChallengeRepository {
	return &passwordLoginChallengeRepository{manager: manager}
}

func (r *passwordLoginChallengeRepository) Save(ctx context.Context, c password_auth.LoginChallenge, resendBefore time.Time) error {
	_, err := r.manager.GetQueries(ctx).SavePasswordLoginChallenge(ctx, model.SavePasswordLoginChallengeParams{
		UserID: c.UserID.UUID(), TokenHash: c.TokenHash, BrowserHash: c.BrowserHash, CredentialHash: c.CredentialHash, PinHash: c.PINHash,
		CreatedAt: c.CreatedAt, ExpiresAt: c.ExpiresAt, ResendBefore: resendBefore,
	})
	return err
}

func (r *passwordLoginChallengeRepository) FindForPIN(ctx context.Context, browser []byte) (password_auth.LoginChallenge, error) {
	row, err := r.manager.GetQueries(ctx).FindPasswordLoginChallengeForPIN(ctx, browser)
	return password_auth.LoginChallenge{UserID: shared.UUID[user.User](row.UserID), TokenHash: row.TokenHash, BrowserHash: row.BrowserHash, CredentialHash: row.CredentialHash, PINHash: row.PinHash, PINAttempts: int(row.PinAttempts), Consumed: row.ConsumedAt.Valid, CreatedAt: row.CreatedAt, ExpiresAt: row.ExpiresAt}, err
}

func (r *passwordLoginChallengeRepository) RecordPINFailure(ctx context.Context, browser []byte) (int, error) {
	count, err := r.manager.GetQueries(ctx).RecordPasswordLoginPINFailure(ctx, browser)
	return int(count), err
}

func (r *passwordLoginChallengeRepository) Consume(ctx context.Context, token, browser []byte, now time.Time) (password_auth.LoginChallenge, error) {
	row, err := r.manager.GetQueries(ctx).ConsumePasswordLoginChallenge(ctx, model.ConsumePasswordLoginChallengeParams{TokenHash: token, BrowserHash: browser, Now: now})
	return password_auth.LoginChallenge{UserID: shared.UUID[user.User](row.UserID), CredentialHash: row.CredentialHash}, err
}
