package repository

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"go.opentelemetry.io/otel"
)

type authStateRepository struct {
	*db.DBManager
}

// NewAuthStateRepository
func NewAuthStateRepository(db *db.DBManager) auth.StateRepository {
	return &authStateRepository{
		DBManager: db,
	}
}

// Create 新しいstateをDBに保存。
func (r *authStateRepository) Create(ctx context.Context, state *auth.State) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "authStateRepository.Create")
	defer span.End()

	return r.ExecTx(ctx, func(ctx context.Context) error {
		_, err := r.GetQueries(ctx).CreateAuthState(ctx, model.CreateAuthStateParams{
			State:       state.State,
			Provider:    state.Provider,
			ExpiresAt:   state.ExpiresAt,
			RedirectUrl: state.RedirectURL,
		})
		return err // 失敗時はそのまま返す
	})
}

// Get 指定したstateをDBから取得
func (r *authStateRepository) Get(ctx context.Context, state string) (*auth.State, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "authStateRepository.Get")
	defer span.End()

	var result *auth.State
	err := r.ExecTx(ctx, func(ctx context.Context) error {
		s, err := r.GetQueries(ctx).GetAuthState(ctx, state)
		if err != nil {
			return err // 取得失敗時はエラーを返す
		}

		result = &auth.State{
			ID:          int(s.ID),
			State:       s.State,
			Provider:    s.Provider,
			RedirectURL: s.RedirectUrl,
			CreatedAt:   s.CreatedAt,
			ExpiresAt:   s.ExpiresAt,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Delete 指定したstateをDBから削除
func (r *authStateRepository) Delete(ctx context.Context, state string) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "authStateRepository.Delete")
	defer span.End()

	return r.ExecTx(ctx, func(ctx context.Context) error {
		return r.GetQueries(ctx).DeleteAuthState(ctx, state)
	})
}

// DeleteExpired 期限切れのstateをDBから一括削除
func (r *authStateRepository) DeleteExpired(ctx context.Context) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "authStateRepository.DeleteExpired")
	defer span.End()

	return r.ExecTx(ctx, func(ctx context.Context) error {
		return r.GetQueries(ctx).DeleteExpiredAuthStates(ctx)
	})
}
