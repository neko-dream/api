package crypto

import (
	"context"
	"testing"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/config"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptUserDemographics(t *testing.T) {
	ctx := context.Background()
	encryptor, err := NewEncryptor(&config.Config{
		ENCRYPTION_VERSION: Version1,
		ENCRYPTION_SECRET:  "12345678901234567890123456789012",
	})
	require.NoError(t, err)

	userID := shared.NewUUID[user.User]()
	demoID := shared.NewUUID[user.UserDemographic]()

	tests := []struct {
		name     string
		demo     user.UserDemographic
		wantErr  bool
		validate func(*testing.T, *model.UserDemographic)
	}{
		{
			name: "全フィールドが設定されている場合",
			demo: user.NewUserDemographic(
				ctx,
				demoID,
				lo.ToPtr(1990),
				lo.ToPtr("正社員"),
				lo.ToPtr("男性"),
				lo.ToPtr("世田谷区"),
				lo.ToPtr(1),
				lo.ToPtr("東京都"),
			),
			wantErr: false,
			validate: func(t *testing.T, got *model.UserDemographic) {
				assert.True(t, got.City.Valid)
				assert.True(t, got.Prefecture.Valid)
				assert.True(t, got.YearOfBirth.Valid)
				assert.True(t, got.Gender.Valid)
				assert.True(t, got.HouseholdSize.Valid)
				assert.Equal(t, int16(1), got.HouseholdSize.Int16)
			},
		},
		{
			name: "一部フィールドがnilの場合",
			demo: user.NewUserDemographic(
				ctx,
				demoID,
				nil,              // 生年
				nil,              // 職業
				lo.ToPtr(""),     // 性別
				lo.ToPtr("世田谷区"), // 都市
				nil,              // 世帯人数
				nil,              // 都道府県
			),
			wantErr: false,
			validate: func(t *testing.T, got *model.UserDemographic) {
				assert.True(t, got.City.Valid)
				assert.False(t, got.Prefecture.Valid)
				assert.False(t, got.YearOfBirth.Valid)
				assert.True(t, got.Gender.Valid)
				assert.False(t, got.HouseholdSize.Valid)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncryptUserDemographics(ctx, encryptor, userID, &tt.demo)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, got)
			tt.validate(t, got)
		})
	}
}

func TestDecryptUserDemographics(t *testing.T) {
	ctx := context.Background()
	encryptor, err := NewEncryptor(&config.Config{
		ENCRYPTION_VERSION: Version1,
		ENCRYPTION_SECRET:  "12345678901234567890123456789012",
	})
	require.NoError(t, err)

	// テストデータの準備
	original := user.NewUserDemographic(
		ctx,
		shared.NewUUID[user.UserDemographic](),
		lo.ToPtr(1990),
		lo.ToPtr("正社員"),
		lo.ToPtr("男性"),
		lo.ToPtr("世田谷区"),
		lo.ToPtr(4),
		lo.ToPtr("回答しない"),
	)

	// 暗号化
	encrypted, err := EncryptUserDemographics(ctx, encryptor, shared.NewUUID[user.User](), &original)
	require.NoError(t, err)

	// 復号化のテスト
	t.Run("正常系: 全フィールドの復号化", func(t *testing.T) {
		decrypted, err := DecryptUserDemographics(ctx, encryptor, encrypted)
		require.NoError(t, err)
		assert.NotNil(t, decrypted)

		// 元のデータと一致することを確認
		assert.Equal(t, original.YearOfBirth(), decrypted.YearOfBirth())
		assert.Equal(t, original.Gender(), decrypted.Gender())
		assert.Equal(t, original.City().String(), decrypted.City().String())
		assert.Equal(t, original.Prefecture(), decrypted.Prefecture())
		assert.Equal(t, original.HouseholdSize(), decrypted.HouseholdSize())
	})

	t.Run("異常系: 不正な暗号文", func(t *testing.T) {
		invalidEncrypted := *encrypted
		invalidEncrypted.City.String = "invalid-cipher-text"
		_, err := DecryptUserDemographics(ctx, encryptor, &invalidEncrypted)
		assert.Error(t, err)
	})
}
