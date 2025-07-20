package user_test

import (
	"context"
	"strings"
	"testing"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name        string
		userID      shared.UUID[user.User]
		displayID   *string
		displayName *string
		subject     string
		provider    shared.AuthProviderName
		iconURL     *string
	}{
		{
			name:        "全てのフィールドを持つユーザーを作成できる",
			userID:      shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			displayID:   lo.ToPtr("testuser123"),
			displayName: lo.ToPtr("Test User"),
			subject:     "auth0|123456",
			provider:    shared.AuthProviderName("auth0"),
			iconURL:     lo.ToPtr("https://example.com/icon.png"),
		},
		{
			name:        "必須フィールドのみでユーザーを作成できる",
			userID:      shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000002"),
			displayID:   nil,
			displayName: nil,
			subject:     "google|789012",
			provider:    shared.AuthProviderName("google"),
			iconURL:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := user.NewUser(
				tt.userID,
				tt.displayID,
				tt.displayName,
				tt.subject,
				tt.provider,
				tt.iconURL,
			)

			assert.Equal(t, tt.userID, u.UserID())
			assert.Equal(t, tt.displayID, u.DisplayID())
			assert.Equal(t, tt.displayName, u.DisplayName())
			assert.Equal(t, tt.subject, u.Subject())
			assert.Equal(t, tt.provider, u.Provider())
			assert.Equal(t, tt.iconURL, u.IconURL())
			assert.Nil(t, u.Email())
			assert.False(t, u.IsEmailVerified())
			assert.Nil(t, u.Demographics())
		})
	}
}

func TestUser_SetDisplayID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "正常な表示IDを設定できる",
			input:   "user123",
			wantErr: nil,
		},
		{
			name:    "4文字の表示IDを設定できる",
			input:   "abcd",
			wantErr: nil,
		},
		{
			name:    "30文字の表示IDを設定できる",
			input:   strings.Repeat("a", 30),
			wantErr: nil,
		},
		{
			name:    "空文字はエラーになる",
			input:   "",
			wantErr: messages.UserDisplayIDInvalidError,
		},
		{
			name:    "3文字以下はエラーになる",
			input:   "abc",
			wantErr: messages.UserDisplayIDTooShort,
		},
		{
			name:    "31文字以上はエラーになる",
			input:   strings.Repeat("a", 31),
			wantErr: messages.UserDisplayIDTooLong,
		},
		{
			name:    "日本語でも文字数でカウントされる",
			input:   "あいうえお",
			wantErr: nil,
		},
		{
			name:    "絵文字も1文字としてカウントされる",
			input:   "user😀",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := user.NewUser(
				shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
				nil,
				nil,
				"test|123",
				shared.AuthProviderName("test"),
				nil,
			)

			err := u.SetDisplayID(tt.input)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, u.DisplayID())
			} else {
				require.NoError(t, err)
				assert.NotNil(t, u.DisplayID())
				assert.Equal(t, tt.input, *u.DisplayID())
			}
		})
	}
}

func TestUser_ChangeName(t *testing.T) {
	tests := []struct {
		name        string
		initialName *string
		newName     *string
		expected    *string
	}{
		{
			name:        "名前を変更できる",
			initialName: lo.ToPtr("Old Name"),
			newName:     lo.ToPtr("New Name"),
			expected:    lo.ToPtr("New Name"),
		},
		{
			name:        "nilから名前を設定できる",
			initialName: nil,
			newName:     lo.ToPtr("New Name"),
			expected:    lo.ToPtr("New Name"),
		},
		{
			name:        "nilを渡しても変更されない",
			initialName: lo.ToPtr("Existing Name"),
			newName:     nil,
			expected:    lo.ToPtr("Existing Name"),
		},
		{
			name:        "空文字を渡しても変更されない",
			initialName: lo.ToPtr("Existing Name"),
			newName:     lo.ToPtr(""),
			expected:    lo.ToPtr("Existing Name"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := user.NewUser(
				shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
				nil,
				tt.initialName,
				"test|123",
				shared.AuthProviderName("test"),
				nil,
			)

			u.ChangeName(context.Background(), tt.newName)

			assert.Equal(t, tt.expected, u.DisplayName())
		})
	}
}

func TestUser_IconURL(t *testing.T) {
	t.Run("アイコンURLを変更できる", func(t *testing.T) {
		u := user.NewUser(
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			nil,
			nil,
			"test|123",
			shared.AuthProviderName("test"),
			lo.ToPtr("https://example.com/old.png"),
		)

		newURL := lo.ToPtr("https://example.com/new.png")
		u.ChangeIconURL(newURL)

		assert.Equal(t, newURL, u.IconURL())
	})

	t.Run("アイコンを削除できる", func(t *testing.T) {
		u := user.NewUser(
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			nil,
			nil,
			"test|123",
			shared.AuthProviderName("test"),
			lo.ToPtr("https://example.com/icon.png"),
		)

		u.DeleteIcon()

		assert.Nil(t, u.IconURL())
	})
}

func TestUser_Email(t *testing.T) {
	t.Run("メールアドレスを設定できる", func(t *testing.T) {
		u := user.NewUser(
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			nil,
			nil,
			"test|123",
			shared.AuthProviderName("test"),
			nil,
		)

		u.ChangeEmail("test@example.com")

		assert.NotNil(t, u.Email())
		assert.Equal(t, "test@example.com", *u.Email())
	})

	t.Run("メール認証状態を設定できる", func(t *testing.T) {
		u := user.NewUser(
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			nil,
			nil,
			"test|123",
			shared.AuthProviderName("test"),
			nil,
		)

		assert.False(t, u.IsEmailVerified())

		u.SetEmailVerified(true)
		assert.True(t, u.IsEmailVerified())

		u.SetEmailVerified(false)
		assert.False(t, u.IsEmailVerified())
	})
}

func TestUser_Verify(t *testing.T) {
	tests := []struct {
		name        string
		displayID   *string
		displayName *string
		expected    bool
	}{
		{
			name:        "表示IDと表示名が両方設定されている場合はtrue",
			displayID:   lo.ToPtr("user123"),
			displayName: lo.ToPtr("Test User"),
			expected:    true,
		},
		{
			name:        "表示IDのみ設定されている場合はfalse",
			displayID:   lo.ToPtr("user123"),
			displayName: nil,
			expected:    false,
		},
		{
			name:        "表示名のみ設定されている場合はfalse",
			displayID:   nil,
			displayName: lo.ToPtr("Test User"),
			expected:    false,
		},
		{
			name:        "両方設定されていない場合はfalse",
			displayID:   nil,
			displayName: nil,
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := user.NewUser(
				shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
				tt.displayID,
				tt.displayName,
				"test|123",
				shared.AuthProviderName("test"),
				nil,
			)

			assert.Equal(t, tt.expected, u.Verify())
		})
	}
}

func TestUser_Demographics(t *testing.T) {
	t.Run("デモグラフィック情報を設定できる", func(t *testing.T) {
		u := user.NewUser(
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			nil,
			nil,
			"test|123",
			shared.AuthProviderName("test"),
			nil,
		)

		assert.Nil(t, u.Demographics())

		demographics := user.NewUserDemographic(
			context.Background(),
			shared.MustParseUUID[user.UserDemographic]("00000000-0000-0000-0000-000000000002"),
			lo.ToPtr(1990),
			lo.ToPtr("male"),
			lo.ToPtr("Tokyo"),
			lo.ToPtr("東京都"),
		)

		u.SetDemographics(demographics)

		assert.NotNil(t, u.Demographics())
		assert.Equal(t, demographics.ID(), u.Demographics().ID())
	})
}