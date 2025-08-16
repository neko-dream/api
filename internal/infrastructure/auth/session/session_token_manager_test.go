package session_test

import (
	"context"
	"testing"
	"time"

	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/organization"
	domainSession "github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/auth/session"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/crypto"
	"github.com/neko-dream/server/internal/infrastructure/di"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/internal/infrastructure/persistence/repository"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionTokenManager(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(ctx context.Context, tm domainSession.TokenManager, sessionRepo domainSession.SessionRepository, userRepo user.UserRepository) (string, error)
		wantErr bool
	}{
		{
			name: "正常にセッションIDからClaimを取得できる",
			setup: func(ctx context.Context, tm domainSession.TokenManager, sessionRepo domainSession.SessionRepository, userRepo user.UserRepository) (string, error) {
				// ユーザー作成
				testUser := user.NewUser(
					shared.NewUUID[user.User](),
					lo.ToPtr("testid"),
					lo.ToPtr("TestUser"),
					"test-subject",
					shared.ProviderGoogle,
					lo.ToPtr("test-icon.png"),
				)
				testUser.SetEmailVerified(true)
				err := userRepo.Create(ctx, testUser)
				if err != nil {
					return "", err
				}

				// CreateはdisplayIDとdisplayNameを保存しないので、Updateで設定する
				err = userRepo.Update(ctx, testUser)
				if err != nil {
					return "", err
				}

				// セッション作成
				sess := domainSession.NewSession(
					shared.NewUUID[domainSession.Session](),
					testUser.UserID(),
					shared.ProviderGoogle,
					domainSession.SESSION_ACTIVE,
					clock.Now(ctx).Add(24*time.Hour),
					clock.Now(ctx),
				)
				createdSess, err := sessionRepo.Create(ctx, *sess)
				if err != nil {
					return "", err
				}

				// トークン生成（SessionIDを返すだけ）
				token, err := tm.Generate(ctx, testUser, createdSess.SessionID())
				if err != nil {
					return "", err
				}

				return token, nil
			},
			wantErr: false,
		},
		{
			name: "存在しないセッションIDでエラーになる",
			setup: func(ctx context.Context, tm domainSession.TokenManager, sessionRepo domainSession.SessionRepository, userRepo user.UserRepository) (string, error) {
				// 存在しないセッションID
				return shared.NewUUID[domainSession.Session]().String(), nil
			},
			wantErr: true,
		},
		{
			name: "無効なセッションステータスでエラーになる",
			setup: func(ctx context.Context, tm domainSession.TokenManager, sessionRepo domainSession.SessionRepository, userRepo user.UserRepository) (string, error) {
				// ユーザー作成
				testUser := user.NewUser(
					shared.NewUUID[user.User](),
					lo.ToPtr("testid2"),
					lo.ToPtr("TestUser2"),
					"test-subject2",
					shared.ProviderGoogle,
					lo.ToPtr("test-icon2.png"),
				)
				err := userRepo.Create(ctx, testUser)
				if err != nil {
					return "", err
				}

				// CreateはdisplayIDとdisplayNameを保存しないので、Updateで設定する
				err = userRepo.Update(ctx, testUser)
				if err != nil {
					return "", err
				}

				// 無効なセッション作成
				sess := domainSession.NewSession(
					shared.NewUUID[domainSession.Session](),
					testUser.UserID(),
					shared.ProviderGoogle,
					domainSession.SESSION_INACTIVE,
					clock.Now(ctx).Add(24*time.Hour),
					clock.Now(ctx),
				)
				createdSess, err := sessionRepo.Create(ctx, *sess)
				if err != nil {
					return "", err
				}

				return createdSess.SessionID().String(), nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cont := di.BuildContainer()
			dbm := di.Invoke[*db.DBManager](cont)
			encryptor, _ := crypto.NewEncryptor(lo.ToPtr(config.Config{
				ENCRYPTION_VERSION: crypto.Version1,
				ENCRYPTION_SECRET:  "12345678901234567890123456789012",
			}))

			userRepo := repository.NewUserRepository(
				dbm,
				repository.NewImageRepositoryMock(),
				encryptor,
			)
			sessionRepo := di.Invoke[domainSession.SessionRepository](cont)
			orgUserRepo := di.Invoke[organization.OrganizationUserRepository](cont)
			orgRepo := di.Invoke[organization.OrganizationRepository](cont)

			err := dbm.ExecTx(context.Background(), func(ctx context.Context) error {

				cfg := &config.Config{TokenSecret: "test-secret"}
				tm := session.NewSessionTokenManager(cfg, dbm, sessionRepo, userRepo, orgUserRepo, orgRepo)

				token, err := tt.setup(ctx, tm, sessionRepo, userRepo)
				require.NoError(t, err)

				claim, err := tm.Parse(ctx, token)

				if tt.wantErr {
					assert.Error(t, err)
					return nil
				}

				assert.NoError(t, err)
				assert.NotNil(t, claim)

				// Claimの内容を検証
				userID, err := claim.UserID()
				assert.NoError(t, err)
				assert.NotEmpty(t, userID)

				sessionID, err := claim.SessionID()
				assert.NoError(t, err)
				assert.NotEmpty(t, sessionID)

				assert.True(t, claim.IsRegistered)
				assert.False(t, claim.IsExpired(ctx))
				return nil
			})
			require.NoError(t, err)
		})
	}
}
