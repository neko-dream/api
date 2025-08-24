package user_test

import (
	"testing"
	"time"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_Withdraw(t *testing.T) {
	tests := []struct {
		name        string
		setupUser   func() *user.User
		withdrawAt  time.Time
		wantErr     error
		checkResult func(t *testing.T, u *user.User)
	}{
		{
			name: "正常に退会できる",
			setupUser: func() *user.User {
				u := user.NewUser(
					shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
					lo.ToPtr("display_id"),
					lo.ToPtr("display_name"),
					"test|123",
					shared.AuthProviderName("test"),
					lo.ToPtr("https://example.com/icon.png"),
				)
				return &u
			},
			withdrawAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:    nil,
			checkResult: func(t *testing.T, u *user.User) {
				assert.True(t, u.IsWithdrawn())
				withdrawalDate := u.WithdrawalDate()
				require.NotNil(t, withdrawalDate)
				assert.Equal(t, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), *withdrawalDate)
			},
		},
		{
			name: "既に退会済みの場合はエラー",
			setupUser: func() *user.User {
				u := user.NewUser(
					shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
					lo.ToPtr("display_id"),
					lo.ToPtr("display_name"),
					"test|123",
					shared.AuthProviderName("test"),
					lo.ToPtr("https://example.com/icon.png"),
				)
				// 既に退会済みにする
				_ = u.Withdraw(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
				return &u
			},
			withdrawAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			wantErr:    user.ErrAlreadyWithdrawn,
			checkResult: func(t *testing.T, u *user.User) {
				assert.True(t, u.IsWithdrawn())
				withdrawalDate := u.WithdrawalDate()
				require.NotNil(t, withdrawalDate)
				// 最初の退会日時が保持されている
				assert.Equal(t, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), *withdrawalDate)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := tt.setupUser()
			err := u.Withdraw(tt.withdrawAt)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
			tt.checkResult(t, u)
		})
	}
}

func TestUser_Reactivate(t *testing.T) {
	tests := []struct {
		name         string
		setupUser    func() *user.User
		reactivateAt time.Time
		wantErr      error
		checkResult  func(t *testing.T, u *user.User)
	}{
		{
			name: "30日以内なら復活できる",
			setupUser: func() *user.User {
				u := user.NewUser(
					shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
					lo.ToPtr("display_id"),
					lo.ToPtr("display_name"),
					"test|123",
					shared.AuthProviderName("test"),
					lo.ToPtr("https://example.com/icon.png"),
				)
				// 退会済みにする
				_ = u.Withdraw(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
				return &u
			},
			reactivateAt: time.Date(2024, 1, 30, 23, 59, 59, 0, time.UTC), // 30日目
			wantErr:      nil,
			checkResult: func(t *testing.T, u *user.User) {
				assert.False(t, u.IsWithdrawn())
				assert.Nil(t, u.WithdrawalDate())
			},
		},
		{
			name: "31日経過したら復活できない",
			setupUser: func() *user.User {
				u := user.NewUser(
					shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
					lo.ToPtr("display_id"),
					lo.ToPtr("display_name"),
					"test|123",
					shared.AuthProviderName("test"),
					lo.ToPtr("https://example.com/icon.png"),
				)
				// 退会済みにする
				_ = u.Withdraw(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
				return &u
			},
			reactivateAt: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), // 31日経過
			wantErr:      user.ErrReactivationPeriodExpired,
			checkResult: func(t *testing.T, u *user.User) {
				assert.True(t, u.IsWithdrawn())
				withdrawalDate := u.WithdrawalDate()
				require.NotNil(t, withdrawalDate)
				assert.Equal(t, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), *withdrawalDate)
			},
		},
		{
			name: "退会していない場合はエラー",
			setupUser: func() *user.User {
				u := user.NewUser(
					shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
					lo.ToPtr("display_id"),
					lo.ToPtr("display_name"),
					"test|123",
					shared.AuthProviderName("test"),
					lo.ToPtr("https://example.com/icon.png"),
				)
				return &u
			},
			reactivateAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:      user.ErrNotWithdrawn,
			checkResult: func(t *testing.T, u *user.User) {
				assert.False(t, u.IsWithdrawn())
				assert.Nil(t, u.WithdrawalDate())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := tt.setupUser()
			err := u.Reactivate(tt.reactivateAt)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
			tt.checkResult(t, u)
		})
	}
}

func TestUser_CanReactivate(t *testing.T) {
	tests := []struct {
		name      string
		setupUser func() *user.User
		checkAt   time.Time
		want      bool
	}{
		{
			name: "退会していない場合はfalse",
			setupUser: func() *user.User {
				u := user.NewUser(
					shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
					lo.ToPtr("display_id"),
					lo.ToPtr("display_name"),
					"test|123",
					shared.AuthProviderName("test"),
					nil,
				)
				return &u
			},
			checkAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			want:    false,
		},
		{
			name: "退会後30日以内ならtrue",
			setupUser: func() *user.User {
				u := user.NewUser(
					shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
					lo.ToPtr("display_id"),
					lo.ToPtr("display_name"),
					"test|123",
					shared.AuthProviderName("test"),
					nil,
				)
				_ = u.Withdraw(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
				return &u
			},
			checkAt: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), // 15日後
			want:    true,
		},
		{
			name: "退会後30日ちょうどならtrue",
			setupUser: func() *user.User {
				u := user.NewUser(
					shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
					lo.ToPtr("display_id"),
					lo.ToPtr("display_name"),
					"test|123",
					shared.AuthProviderName("test"),
					nil,
				)
				_ = u.Withdraw(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
				return &u
			},
			checkAt: time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC), // 30日後
			want:    true,
		},
		{
			name: "退会後31日経過したらfalse",
			setupUser: func() *user.User {
				u := user.NewUser(
					shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
					lo.ToPtr("display_id"),
					lo.ToPtr("display_name"),
					"test|123",
					shared.AuthProviderName("test"),
					nil,
				)
				_ = u.Withdraw(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
				return &u
			},
			checkAt: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), // 31日後
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := tt.setupUser()
			got, err := u.CanReactivate(tt.checkAt)
			assert.Equal(t, tt.want, got)

			// エラーの検証
			if u.IsWithdrawn() && u.IsReactivationPeriodExpired(tt.checkAt) {
				// 退会済みで期限切れの場合
				assert.ErrorIs(t, err, user.ErrReactivationPeriodExpired)
			} else if !u.IsWithdrawn() {
				// 退会していない場合
				assert.ErrorIs(t, err, user.ErrNotWithdrawn)
			} else {
				// 退会済みで期限内の場合
				assert.NoError(t, err)
			}
		})
	}
}

func TestUser_IsReactivationPeriodExpired(t *testing.T) {
	tests := []struct {
		name      string
		setupUser func() *user.User
		checkAt   time.Time
		want      bool
	}{
		{
			name: "退会していない場合はfalse",
			setupUser: func() *user.User {
				u := user.NewUser(
					shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
					lo.ToPtr("display_id"),
					lo.ToPtr("display_name"),
					"test|123",
					shared.AuthProviderName("test"),
					nil,
				)
				return &u
			},
			checkAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			want:    false,
		},
		{
			name: "退会後30日以内ならfalse",
			setupUser: func() *user.User {
				u := user.NewUser(
					shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
					lo.ToPtr("display_id"),
					lo.ToPtr("display_name"),
					"test|123",
					shared.AuthProviderName("test"),
					nil,
				)
				_ = u.Withdraw(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
				return &u
			},
			checkAt: time.Date(2024, 1, 30, 0, 0, 0, 0, time.UTC), // 29日後
			want:    false,
		},
		{
			name: "退会後31日経過したらtrue",
			setupUser: func() *user.User {
				u := user.NewUser(
					shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
					lo.ToPtr("display_id"),
					lo.ToPtr("display_name"),
					"test|123",
					shared.AuthProviderName("test"),
					nil,
				)
				_ = u.Withdraw(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
				return &u
			},
			checkAt: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), // 31日後
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := tt.setupUser()
			got := u.IsReactivationPeriodExpired(tt.checkAt)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUser_PrepareForDuplication(t *testing.T) {
	tests := []struct {
		name           string
		setupUser      func() *user.User
		wantNewSubject string
	}{
		{
			name: "退会していない場合は元の値を返す",
			setupUser: func() *user.User {
				u := user.NewUser(
					shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
					lo.ToPtr("display_id"),
					lo.ToPtr("display_name"),
					"original_subject",
					shared.AuthProviderName("test"),
					nil,
				)
				u.ChangeEmail("user@example.com")
				return &u
			},
			wantNewSubject: "original_subject",
		},
		{
			name: "退会済みの場合は変更された値を返す",
			setupUser: func() *user.User {
				u := user.NewUser(
					shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
					lo.ToPtr("display_id"),
					lo.ToPtr("display_name"),
					"original_subject",
					shared.AuthProviderName("test"),
					nil,
				)
				u.ChangeEmail("user@example.com")
				_ = u.Withdraw(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
				return &u
			},
			wantNewSubject: "original_subject_withdrawn_1704067200",
		},
		{
			name: "emailがnilの場合でも正常に動作",
			setupUser: func() *user.User {
				u := user.NewUser(
					shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
					lo.ToPtr("display_id"),
					lo.ToPtr("display_name"),
					"original_subject",
					shared.AuthProviderName("test"),
					nil,
				)
				_ = u.Withdraw(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
				return &u
			},
			wantNewSubject: "original_subject_withdrawn_1704067200",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := tt.setupUser()
			gotSubject := u.PrepareForDeleteUser()
			assert.Equal(t, tt.wantNewSubject, gotSubject)
			// PrepareForDuplicationはemailを返さないのでその部分のテストを削除
			// 実際のemailの変更はAuthServiceで行われる
		})
	}
}
