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

// AdminUser 管理画面用の独立した管理者アカウント。
// 既存の user.User / session.Session（OAuth・組織）とは無関係。
type AdminUser struct {
	AdminUserID  shared.UUID[AdminUser]
	Email        string
	DisplayName  string
	PasswordHash string
	Salt         string
	Role         AdminRole
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type AdminUserRepository interface {
	Create(ctx context.Context, adminUser AdminUser) error
	FindByEmail(ctx context.Context, email string) (*AdminUser, error)
	FindByID(ctx context.Context, id shared.UUID[AdminUser]) (*AdminUser, error)
}
