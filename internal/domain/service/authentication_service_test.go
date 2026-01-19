package service

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/neko-dream/api/internal/domain/messages"
	mock_auth "github.com/neko-dream/api/internal/domain/mock/auth"
	mock_consent "github.com/neko-dream/api/internal/domain/mock/consent"
	mock_user "github.com/neko-dream/api/internal/domain/mock/user"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/consent"
	"github.com/neko-dream/api/internal/domain/model/organization"
	"github.com/neko-dream/api/internal/domain/model/session"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthenticationService_GenerateState(t *testing.T) {
	svc := &authenticationService{}

	t.Run("32文字の文字列を生成する", func(t *testing.T) {
		state, err := svc.GenerateState(context.Background())
		require.NoError(t, err)
		assert.Len(t, state, 32)
	})

	t.Run("英数字のみで構成される", func(t *testing.T) {
		state, err := svc.GenerateState(context.Background())
		require.NoError(t, err)

		matched, err := regexp.MatchString(`^[0-9A-Za-z]{32}$`, state)
		require.NoError(t, err)
		assert.True(t, matched, "期待: 英数字32文字, 実際: %s", state)
	})

	t.Run("毎回異なる値を生成する", func(t *testing.T) {
		states := make(map[string]bool)
		for i := 0; i < 100; i++ {
			state, err := svc.GenerateState(context.Background())
			require.NoError(t, err)
			states[state] = true
		}
		// 100回生成して、全て異なる値であることを期待
		assert.Equal(t, 100, len(states), "ランダム性が不足: 100回中%d種類のみ生成", len(states))
	})
}

func TestAuthenticationService_BuildAuthContext(t *testing.T) {
	svc := &authenticationService{}
	now := time.Now()

	t.Run("基本的なClaimから正しくAuthContextを構築する", func(t *testing.T) {
		claim := &session.Claim{
			Sub:                    "00000000-0000-0000-0000-000000000001",
			Iat:                    now.Unix(),
			Exp:                    now.Add(time.Hour).Unix(),
			Jti:                    "00000000-0000-0000-0000-000000000002",
			DisplayName:            lo.ToPtr("Test User"),
			DisplayID:              lo.ToPtr("testuser"),
			IconURL:                lo.ToPtr("https://example.com/icon.png"),
			IsRegistered:           true,
			IsEmailVerified:        true,
			RequiredPasswordChange: false,
			IsGuest:                false,
		}

		authCtx, err := svc.BuildAuthContext(context.Background(), claim)
		require.NoError(t, err)

		assert.Equal(t, "00000000-0000-0000-0000-000000000001", authCtx.UserID.String())
		assert.Equal(t, "00000000-0000-0000-0000-000000000002", authCtx.SessionID.String())
		assert.Equal(t, "Test User", *authCtx.DisplayName)
		assert.Equal(t, "testuser", *authCtx.DisplayID)
		assert.Equal(t, "https://example.com/icon.png", *authCtx.IconURL)
		assert.True(t, authCtx.IsRegistered)
		assert.True(t, authCtx.IsEmailVerified)
		assert.False(t, authCtx.RequiredPasswordChange)
		assert.False(t, authCtx.IsGuest)
		assert.Nil(t, authCtx.OrganizationID)
		assert.Nil(t, authCtx.OrganizationCode)
		assert.Nil(t, authCtx.OrganizationRole)
	})

	t.Run("ゲストユーザーのClaimを正しく処理する", func(t *testing.T) {
		claim := &session.Claim{
			Sub:          "00000000-0000-0000-0000-000000000001",
			Iat:          now.Unix(),
			Exp:          now.Add(time.Hour).Unix(),
			Jti:          "00000000-0000-0000-0000-000000000002",
			DisplayName:  lo.ToPtr("ゲストabc123"),
			DisplayID:    lo.ToPtr("12345678"),
			IsRegistered: true,
			IsGuest:      true,
		}

		authCtx, err := svc.BuildAuthContext(context.Background(), claim)
		require.NoError(t, err)

		assert.True(t, authCtx.IsGuest)
		assert.True(t, authCtx.IsRegistered)
	})

	t.Run("組織情報を含むClaimを正しく処理する", func(t *testing.T) {
		claim := &session.Claim{
			Sub:              "00000000-0000-0000-0000-000000000001",
			Iat:              now.Unix(),
			Exp:              now.Add(time.Hour).Unix(),
			Jti:              "00000000-0000-0000-0000-000000000002",
			IsRegistered:     true,
			OrganizationID:   lo.ToPtr("00000000-0000-0000-0000-000000000003"),
			OrganizationCode: lo.ToPtr("ORG001"),
			OrganizationRole: lo.ToPtr("管理者"), // NameToRoleは日本語を期待
		}

		authCtx, err := svc.BuildAuthContext(context.Background(), claim)
		require.NoError(t, err)

		require.NotNil(t, authCtx.OrganizationID)
		assert.Equal(t, "00000000-0000-0000-0000-000000000003", authCtx.OrganizationID.String())
		require.NotNil(t, authCtx.OrganizationCode)
		assert.Equal(t, "ORG001", *authCtx.OrganizationCode)
		require.NotNil(t, authCtx.OrganizationRole)
		assert.Equal(t, organization.OrganizationUserRoleAdmin, *authCtx.OrganizationRole)
	})

	t.Run("不正なUserIDの場合はエラーを返す", func(t *testing.T) {
		claim := &session.Claim{
			Sub: "invalid-uuid",
			Jti: "00000000-0000-0000-0000-000000000002",
		}

		_, err := svc.BuildAuthContext(context.Background(), claim)
		assert.Error(t, err)
	})

	t.Run("不正なSessionIDの場合はエラーを返す", func(t *testing.T) {
		claim := &session.Claim{
			Sub: "00000000-0000-0000-0000-000000000001",
			Jti: "invalid-uuid",
		}

		_, err := svc.BuildAuthContext(context.Background(), claim)
		assert.Error(t, err)
	})

	t.Run("パスワード変更要求フラグを正しく処理する", func(t *testing.T) {
		claim := &session.Claim{
			Sub:                    "00000000-0000-0000-0000-000000000001",
			Iat:                    now.Unix(),
			Exp:                    now.Add(time.Hour).Unix(),
			Jti:                    "00000000-0000-0000-0000-000000000002",
			RequiredPasswordChange: true,
		}

		authCtx, err := svc.BuildAuthContext(context.Background(), claim)
		require.NoError(t, err)

		assert.True(t, authCtx.RequiredPasswordChange)
	})
}

func TestAuthenticationService_Authenticate(t *testing.T) {
	t.Run("既存ユーザーが見つかった場合はそのユーザーを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mock_user.NewMockUserRepository(ctrl)
		mockAuthFactory := mock_auth.NewMockAuthProviderFactory(ctrl)
		mockAuthProvider := mock_auth.NewMockAuthProvider(ctrl)
		mockConsentSvc := mock_consent.NewMockConsentService(ctrl)
		mockPolicyRepo := mock_consent.NewMockPolicyRepository(ctrl)

		subject := "google|123456"
		email := "test@example.com"
		existingUser := user.NewUser(
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			lo.ToPtr("testuser"),
			lo.ToPtr("Test User"),
			subject,
			shared.ProviderGoogle,
			nil,
		)

		mockAuthFactory.EXPECT().
			NewAuthProvider(gomock.Any(), "GOOGLE").
			Return(mockAuthProvider, nil)

		mockAuthProvider.EXPECT().
			VerifyAndIdentify(gomock.Any(), "auth_code").
			Return(&subject, &email, nil)

		mockUserRepo.EXPECT().
			FindBySubject(gomock.Any(), user.UserSubject(subject)).
			Return(&existingUser, nil)

		svc := &authenticationService{
			userRepository:      mockUserRepo,
			authProviderFactory: mockAuthFactory,
			consentService:      mockConsentSvc,
			policyRepository:    mockPolicyRepo,
		}

		result, err := svc.Authenticate(context.Background(), "GOOGLE", "auth_code")
		require.NoError(t, err)
		assert.Equal(t, existingUser.UserID(), result.UserID())
	})

	t.Run("新規ユーザーの場合は作成して返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mock_user.NewMockUserRepository(ctrl)
		mockAuthFactory := mock_auth.NewMockAuthProviderFactory(ctrl)
		mockAuthProvider := mock_auth.NewMockAuthProvider(ctrl)
		mockConsentSvc := mock_consent.NewMockConsentService(ctrl)
		mockPolicyRepo := mock_consent.NewMockPolicyRepository(ctrl)

		subject := "google|newuser"
		email := "newuser@example.com"

		mockAuthFactory.EXPECT().
			NewAuthProvider(gomock.Any(), "GOOGLE").
			Return(mockAuthProvider, nil)

		mockAuthProvider.EXPECT().
			VerifyAndIdentify(gomock.Any(), "auth_code").
			Return(&subject, &email, nil)

		mockUserRepo.EXPECT().
			FindBySubject(gomock.Any(), user.UserSubject(subject)).
			Return(nil, sql.ErrNoRows)

		mockPolicyRepo.EXPECT().
			FetchLatestPolicy(gomock.Any()).
			Return(&consent.Policy{Version: "1.0.0"}, nil)

		mockConsentSvc.EXPECT().
			RecordConsent(gomock.Any(), gomock.Any(), "1.0.0", "", "").
			Return(&consent.ConsentRecord{}, nil)

		// Createに渡されるUserの中身を検証
		mockUserRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, u user.User) error {
				// 新規ユーザーの検証
				assert.Equal(t, subject, u.Subject())
				assert.Equal(t, shared.ProviderGoogle, u.Provider())
				assert.Equal(t, email, *u.Email())
				assert.True(t, u.IsEmailVerified(), "OAuth経由のemailは検証済みであるべき")
				assert.False(t, u.IsGuest(), "新規ユーザーはゲストではない")
				return nil
			})

		svc := &authenticationService{
			userRepository:      mockUserRepo,
			authProviderFactory: mockAuthFactory,
			consentService:      mockConsentSvc,
			policyRepository:    mockPolicyRepo,
		}

		result, err := svc.Authenticate(context.Background(), "GOOGLE", "auth_code")
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, subject, result.Subject())
		assert.Equal(t, shared.ProviderGoogle, result.Provider())
		assert.Equal(t, email, *result.Email())
		assert.True(t, result.IsEmailVerified())
	})

	t.Run("AuthProviderの作成に失敗した場合はエラーを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mock_user.NewMockUserRepository(ctrl)
		mockAuthFactory := mock_auth.NewMockAuthProviderFactory(ctrl)
		mockConsentSvc := mock_consent.NewMockConsentService(ctrl)
		mockPolicyRepo := mock_consent.NewMockPolicyRepository(ctrl)

		mockAuthFactory.EXPECT().
			NewAuthProvider(gomock.Any(), "INVALID").
			Return(nil, errors.New("invalid provider"))

		svc := &authenticationService{
			userRepository:      mockUserRepo,
			authProviderFactory: mockAuthFactory,
			consentService:      mockConsentSvc,
			policyRepository:    mockPolicyRepo,
		}

		_, err := svc.Authenticate(context.Background(), "INVALID", "auth_code")
		assert.Error(t, err)
	})

	t.Run("認証に失敗した場合はエラーを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mock_user.NewMockUserRepository(ctrl)
		mockAuthFactory := mock_auth.NewMockAuthProviderFactory(ctrl)
		mockAuthProvider := mock_auth.NewMockAuthProvider(ctrl)
		mockConsentSvc := mock_consent.NewMockConsentService(ctrl)
		mockPolicyRepo := mock_consent.NewMockPolicyRepository(ctrl)

		mockAuthFactory.EXPECT().
			NewAuthProvider(gomock.Any(), "GOOGLE").
			Return(mockAuthProvider, nil)

		mockAuthProvider.EXPECT().
			VerifyAndIdentify(gomock.Any(), "invalid_code").
			Return(nil, nil, errors.New("authentication failed"))

		svc := &authenticationService{
			userRepository:      mockUserRepo,
			authProviderFactory: mockAuthFactory,
			consentService:      mockConsentSvc,
			policyRepository:    mockPolicyRepo,
		}

		_, err := svc.Authenticate(context.Background(), "GOOGLE", "invalid_code")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "authentication failed")
	})

	t.Run("subjectがnilの場合はForbiddenErrorを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mock_user.NewMockUserRepository(ctrl)
		mockAuthFactory := mock_auth.NewMockAuthProviderFactory(ctrl)
		mockAuthProvider := mock_auth.NewMockAuthProvider(ctrl)
		mockConsentSvc := mock_consent.NewMockConsentService(ctrl)
		mockPolicyRepo := mock_consent.NewMockPolicyRepository(ctrl)

		mockAuthFactory.EXPECT().
			NewAuthProvider(gomock.Any(), "GOOGLE").
			Return(mockAuthProvider, nil)

		// subjectがnilを返す
		mockAuthProvider.EXPECT().
			VerifyAndIdentify(gomock.Any(), "auth_code").
			Return(nil, nil, nil)

		svc := &authenticationService{
			userRepository:      mockUserRepo,
			authProviderFactory: mockAuthFactory,
			consentService:      mockConsentSvc,
			policyRepository:    mockPolicyRepo,
		}

		_, err := svc.Authenticate(context.Background(), "GOOGLE", "auth_code")
		require.Error(t, err)
		assert.ErrorIs(t, err, messages.ForbiddenError)
	})

	t.Run("ゲストユーザーがOAuth認証した場合はアップグレードされる", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mock_user.NewMockUserRepository(ctrl)
		mockAuthFactory := mock_auth.NewMockAuthProviderFactory(ctrl)
		mockAuthProvider := mock_auth.NewMockAuthProvider(ctrl)
		mockConsentSvc := mock_consent.NewMockConsentService(ctrl)
		mockPolicyRepo := mock_consent.NewMockPolicyRepository(ctrl)

		oauthSubject := "google|oauth_user"
		email := "upgraded@example.com"
		guestUserID := shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001")

		// ゲストユーザーを作成
		guestUser := user.NewUser(
			guestUserID,
			lo.ToPtr("12345678"),
			lo.ToPtr("ゲストabc123"),
			"guest|abc123",
			shared.ProviderGuest,
			nil,
		)

		mockAuthFactory.EXPECT().
			NewAuthProvider(gomock.Any(), "GOOGLE").
			Return(mockAuthProvider, nil)

		mockAuthProvider.EXPECT().
			VerifyAndIdentify(gomock.Any(), "auth_code").
			Return(&oauthSubject, &email, nil)

		// OAuthのsubjectで既存ユーザーは見つからない
		mockUserRepo.EXPECT().
			FindBySubject(gomock.Any(), user.UserSubject(oauthSubject)).
			Return(nil, sql.ErrNoRows)

		// 現在のセッションからゲストユーザーを取得
		mockUserRepo.EXPECT().
			FindByID(gomock.Any(), guestUserID).
			Return(&guestUser, nil)

		// ゲストユーザーに新しいAuthを追加
		mockUserRepo.EXPECT().
			AddUserAuthForMerge(gomock.Any(), guestUserID, shared.ProviderGoogle, oauthSubject).
			Return(nil)

		// プロバイダーを更新
		mockUserRepo.EXPECT().
			UpdateUserAuthProvider(gomock.Any(), guestUserID, shared.ProviderGoogle).
			Return(nil)

		// ユーザー情報を更新（email設定）
		mockUserRepo.EXPECT().
			Update(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, u user.User) error {
				assert.Equal(t, email, *u.Email())
				assert.True(t, u.IsEmailVerified())
				return nil
			})

		svc := &authenticationService{
			userRepository:      mockUserRepo,
			authProviderFactory: mockAuthFactory,
			consentService:      mockConsentSvc,
			policyRepository:    mockPolicyRepo,
		}

		// セッションにゲストユーザーのClaimを設定
		claim := &session.Claim{
			Sub:     guestUserID.String(),
			Jti:     "00000000-0000-0000-0000-000000000002",
			IsGuest: true,
		}
		ctx := session.SetSession(context.Background(), claim)

		result, err := svc.Authenticate(ctx, "GOOGLE", "auth_code")
		require.NoError(t, err)
		require.NotNil(t, result)
		// ゲストユーザーがアップグレードされて返される
		assert.Equal(t, guestUserID, result.UserID())
	})

	t.Run("退会ユーザーで31日以上経過している場合は新規ユーザーとして作成する", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mock_user.NewMockUserRepository(ctrl)
		mockAuthFactory := mock_auth.NewMockAuthProviderFactory(ctrl)
		mockAuthProvider := mock_auth.NewMockAuthProvider(ctrl)
		mockConsentSvc := mock_consent.NewMockConsentService(ctrl)
		mockPolicyRepo := mock_consent.NewMockPolicyRepository(ctrl)

		subject := "google|123456"
		email := "test@example.com"

		// 退会日を31日以上前に設定（ReactivationPeriodDays = 30）
		withdrawalDate := time.Now().AddDate(0, 0, -32)
		withdrawnUser := user.NewUserWithWithdrawalDate(
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			lo.ToPtr("olduser"),
			lo.ToPtr("Old User"),
			subject,
			shared.ProviderGoogle,
			nil,
			sql.NullTime{Time: withdrawalDate, Valid: true},
		)

		// Authenticate内でclock.Now(ctx)が使われるので、現在時刻を注入
		now := time.Now()
		ctx := clock.SetNow(context.Background(), now)

		mockAuthFactory.EXPECT().
			NewAuthProvider(gomock.Any(), "GOOGLE").
			Return(mockAuthProvider, nil)

		mockAuthProvider.EXPECT().
			VerifyAndIdentify(gomock.Any(), "auth_code").
			Return(&subject, &email, nil)

		// 退会ユーザーが見つかる
		mockUserRepo.EXPECT().
			FindBySubject(gomock.Any(), user.UserSubject(subject)).
			Return(&withdrawnUser, nil)

		// PrepareForDeleteUser呼び出し後、Updateで古いユーザーを更新
		mockUserRepo.EXPECT().
			Update(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, u user.User) error {
				// PrepareForDeleteUserによりsubjectが変更されている
				assert.True(t, strings.Contains(u.Subject(), "_withdrawn_"),
					"subjectに_withdrawn_サフィックスが付与されるべき")
				// emailがクリアされている
				assert.Nil(t, u.Email(), "退会処理によりemailはnilになるべき")
				return nil
			})

		// ChangeSubjectが呼ばれる
		mockUserRepo.EXPECT().
			ChangeSubject(gomock.Any(), withdrawnUser.UserID(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, userID shared.UUID[user.User], newSubject string) error {
				assert.True(t, strings.Contains(newSubject, "_withdrawn_"),
					"新しいsubjectには_withdrawn_が含まれるべき")
				return nil
			})

		// 同意ポリシー取得
		mockPolicyRepo.EXPECT().
			FetchLatestPolicy(gomock.Any()).
			Return(&consent.Policy{Version: "1.0.0"}, nil)

		// 同意記録
		mockConsentSvc.EXPECT().
			RecordConsent(gomock.Any(), gomock.Any(), "1.0.0", "", "").
			Return(&consent.ConsentRecord{}, nil)

		// 新規ユーザーとして作成される
		mockUserRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, u user.User) error {
				// 新規ユーザーの検証
				assert.Equal(t, subject, u.Subject())
				assert.Equal(t, shared.ProviderGoogle, u.Provider())
				assert.Equal(t, email, *u.Email())
				assert.True(t, u.IsEmailVerified())
				return nil
			})

		svc := &authenticationService{
			userRepository:      mockUserRepo,
			authProviderFactory: mockAuthFactory,
			consentService:      mockConsentSvc,
			policyRepository:    mockPolicyRepo,
		}

		result, err := svc.Authenticate(ctx, "GOOGLE", "auth_code")
		require.NoError(t, err)
		require.NotNil(t, result)
		// 新規ユーザーが作成されたので、元のsubjectで作成される
		assert.Equal(t, subject, result.Subject())
		assert.Equal(t, email, *result.Email())
	})

	t.Run("退会ユーザーで31日未満の場合は既存ユーザーとして返す（復活可能）", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mock_user.NewMockUserRepository(ctrl)
		mockAuthFactory := mock_auth.NewMockAuthProviderFactory(ctrl)
		mockAuthProvider := mock_auth.NewMockAuthProvider(ctrl)
		mockConsentSvc := mock_consent.NewMockConsentService(ctrl)
		mockPolicyRepo := mock_consent.NewMockPolicyRepository(ctrl)

		subject := "google|123456"
		email := "test@example.com"

		// 退会日を29日前に設定（31日未満）
		withdrawalDate := time.Now().AddDate(0, 0, -29)
		withdrawnUser := user.NewUserWithWithdrawalDate(
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			lo.ToPtr("stillvalid"),
			lo.ToPtr("Still Valid User"),
			subject,
			shared.ProviderGoogle,
			nil,
			sql.NullTime{Time: withdrawalDate, Valid: true},
		)

		now := time.Now()
		ctx := clock.SetNow(context.Background(), now)

		mockAuthFactory.EXPECT().
			NewAuthProvider(gomock.Any(), "GOOGLE").
			Return(mockAuthProvider, nil)

		mockAuthProvider.EXPECT().
			VerifyAndIdentify(gomock.Any(), "auth_code").
			Return(&subject, &email, nil)

		mockUserRepo.EXPECT().
			FindBySubject(gomock.Any(), user.UserSubject(subject)).
			Return(&withdrawnUser, nil)

		// Update/ChangeSubject/Createは呼ばれない（既存ユーザーとして返却）

		svc := &authenticationService{
			userRepository:      mockUserRepo,
			authProviderFactory: mockAuthFactory,
			consentService:      mockConsentSvc,
			policyRepository:    mockPolicyRepo,
		}

		result, err := svc.Authenticate(ctx, "GOOGLE", "auth_code")
		require.NoError(t, err)
		require.NotNil(t, result)
		// 退会済みだが復活可能期間内なので、そのまま返される
		assert.Equal(t, withdrawnUser.UserID(), result.UserID())
		assert.True(t, result.IsWithdrawn(), "退会状態は維持される")
	})
}
