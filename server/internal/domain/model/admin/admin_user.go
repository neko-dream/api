package admin

import (
	"context"
	"time"

	"github.com/neko-dream/api/internal/domain/model/shared"
)

// AdminUser 管理画面の管理者。ログインは Google 認証で email により照合する。
// 既存の user.User とは無関係。
type AdminUser struct {
	AdminUserID shared.UUID[AdminUser]
	Email       string
	DisplayName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AdminUserRepository interface {
	Create(ctx context.Context, adminUser AdminUser) error
	FindByEmail(ctx context.Context, email string) (*AdminUser, error)
	FindByID(ctx context.Context, id shared.UUID[AdminUser]) (*AdminUser, error)
	List(ctx context.Context) ([]AdminUser, error)
}

// GoogleClaims Google ID トークンから取り出した検証済みクレーム。
type GoogleClaims struct {
	Email string
	Name  string
}

// GoogleTokenVerifier Google ID トークンを検証する。
type GoogleTokenVerifier interface {
	Verify(ctx context.Context, idToken string) (*GoogleClaims, error)
}
