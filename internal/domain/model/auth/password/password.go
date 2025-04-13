package password_auth

import (
	"context"

	"time"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type PasswordAuthRepository interface {
	// CreatePasswordAuth 新しいパスワード認証情報を作成
	CreatePasswordAuth(ctx context.Context, userID shared.UUID[user.User], passwordHash, salt string) error

	// GetPasswordAuthByUserID ユーザーIDからパスワード認証情報を取得
	GetPasswordAuthByUserID(ctx context.Context, userID shared.UUID[user.User]) (PasswordAuth, error)

	// UpdatePasswordAuth パスワード認証情報を更新
	UpdatePasswordAuth(ctx context.Context, userID shared.UUID[user.User], passwordHash, salt string) error

	// DeletePasswordAuth パスワード認証情報を削除
	DeletePasswordAuth(ctx context.Context, userID shared.UUID[user.User]) error
}

type PasswordAuthManager interface {
	// RegisterPassword 新しいパスワードを登録
	RegisterPassword(ctx context.Context, userID shared.UUID[user.User], plainPassword string) error

	// VerifyPassword パスワードを検証
	VerifyPassword(ctx context.Context, userID shared.UUID[user.User], plainPassword string) (bool, error)

	// UpdatePassword パスワードを更新
	UpdatePassword(ctx context.Context, userID shared.UUID[user.User], plainPassword string) error
}

type PasswordAuth struct {
	PasswordAuthID shared.UUID[PasswordAuth]
	UserID         shared.UUID[user.User]
	PasswordHash   string
	Salt           string
	LastChanged    time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
