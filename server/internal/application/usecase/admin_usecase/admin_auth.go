package admin_usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/pkg/hash"
)

// refresh トークンの有効期限。access(短命JWT)に対し長命・revocable。
const refreshTokenTTL = time.Hour * 24 * 30

// AdminAuthOutput ログイン/リフレッシュ共通の出力。
type AdminAuthOutput struct {
	AccessToken  string
	RefreshToken string
	AdminUser    admin.AdminUser
}

// hashRefreshToken は保存・検索用に refresh トークンを HMAC でハッシュ化する。
// newRefreshToken はリフレッシュトークンを作る。Cookie に入れる値で、ogen は Set-Cookie の値を
// URL エンコードするので、URL で使える文字だけ（base64url、パディングなし）にする。
func newRefreshToken() (string, error) {
	b := make([]byte, 48)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashRefreshToken(raw string, cfg *config.Config) (string, error) {
	return hash.HMACHash(raw, cfg.HASH_PEPPER)
}

// issueAccessAndRefresh は access JWT を生成し、新しい refresh セッションを作成して raw refresh token を返す。
func issueAccessAndRefresh(
	ctx context.Context,
	tm admin.AdminTokenManager,
	sessRepo admin.AdminSessionRepository,
	cfg *config.Config,
	adminUser admin.AdminUser,
) (*AdminAuthOutput, error) {
	access, err := tm.Generate(ctx, adminUser)
	if err != nil {
		return nil, err
	}
	refresh, err := newRefreshToken()
	if err != nil {
		return nil, err
	}
	tokenHash, err := hashRefreshToken(refresh, cfg)
	if err != nil {
		return nil, err
	}
	session := admin.AdminSession{
		AdminSessionID: shared.NewUUID[admin.AdminSession](),
		AdminUserID:    adminUser.AdminUserID,
		TokenHash:      tokenHash,
		ExpiresAt:      clock.Now(ctx).Add(refreshTokenTTL),
	}
	if err := sessRepo.Create(ctx, session); err != nil {
		return nil, err
	}
	return &AdminAuthOutput{AccessToken: access, RefreshToken: refresh, AdminUser: adminUser}, nil
}
