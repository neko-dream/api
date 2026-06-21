package admin

import (
	"context"
	"time"

	"github.com/neko-dream/api/internal/domain/model/shared"
)

// AdminClaim 管理者トークン(JWT)のクレーム。user.Session とは無関係の独立認証コンテキスト。
type AdminClaim struct {
	Sub         string // admin_user_id
	Email       string
	DisplayName string
	Role        int
	Iat         int64
	Exp         int64
}

func (c *AdminClaim) AdminUserID() (shared.UUID[AdminUser], error) {
	return shared.ParseUUID[AdminUser](c.Sub)
}

func (c *AdminClaim) IsExpired(now time.Time) bool {
	return now.Unix() > c.Exp
}

// AdminTokenManager 管理者トークンの生成・検証を行う。
type AdminTokenManager interface {
	Generate(ctx context.Context, adminUser AdminUser) (string, error)
	Parse(ctx context.Context, token string) (*AdminClaim, error)
}

type adminContextKey string

var adminKey adminContextKey = "adminSession"

func SetAdminSession(ctx context.Context, claim *AdminClaim) context.Context {
	return context.WithValue(ctx, adminKey, claim)
}

func GetAdminSession(ctx context.Context) *AdminClaim {
	if ctx == nil {
		return nil
	}
	value := ctx.Value(adminKey)
	if value == nil {
		return nil
	}
	claim, ok := value.(*AdminClaim)
	if !ok {
		return nil
	}
	return claim
}
