package admin

import (
	"context"
	"time"

	"github.com/neko-dream/api/internal/domain/model/shared"
)

// AdminRole 管理者ロール。数値が小さいほど高権限。
type AdminRole int

const (
	AdminRoleSuperAdmin AdminRole = 1
	AdminRoleOperator   AdminRole = 10
)

// AdminUser 管理画面用の管理者アカウント(許可リスト)。
// ログインは Google 認証で、email により照合する。既存の user.User とは無関係。
type AdminUser struct {
	AdminUserID shared.UUID[AdminUser]
	Email       string
	DisplayName string
	Role        AdminRole
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
