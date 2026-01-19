package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/auth"
	"github.com/neko-dream/api/internal/domain/model/organization"
	"github.com/neko-dream/api/internal/domain/model/session"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	mock_service "github.com/neko-dream/api/internal/domain/service/mock/service"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthorizationService_GetAuthContext(t *testing.T) {
	t.Run("セッションがない場合はnilを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthSvc := mock_service.NewMockAuthenticationService(ctrl)
		svc := NewAuthorizationService(mockAuthSvc)

		// セッションなしのコンテキスト
		ctx := context.Background()

		result, err := svc.GetAuthContext(ctx)
		require.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("セッションがある場合はAuthContextを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthSvc := mock_service.NewMockAuthenticationService(ctrl)
		svc := NewAuthorizationService(mockAuthSvc)

		now := time.Now()
		claim := &session.Claim{
			Sub:          "00000000-0000-0000-0000-000000000001",
			Iat:          now.Unix(),
			Exp:          now.Add(time.Hour).Unix(),
			Jti:          "00000000-0000-0000-0000-000000000002",
			DisplayName:  lo.ToPtr("Test User"),
			IsRegistered: true,
		}
		ctx := session.SetSession(context.Background(), claim)

		expectedAuthCtx := &auth.AuthenticationContext{
			UserID:       shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			SessionID:    shared.MustParseUUID[session.Session]("00000000-0000-0000-0000-000000000002"),
			DisplayName:  lo.ToPtr("Test User"),
			IsRegistered: true,
		}

		mockAuthSvc.EXPECT().
			BuildAuthContext(gomock.Any(), claim).
			Return(expectedAuthCtx, nil)

		result, err := svc.GetAuthContext(ctx)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedAuthCtx.UserID, result.UserID)
	})

	t.Run("BuildAuthContextがエラーを返した場合はエラーを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthSvc := mock_service.NewMockAuthenticationService(ctrl)
		svc := NewAuthorizationService(mockAuthSvc)

		claim := &session.Claim{
			Sub: "invalid-uuid",
			Jti: "00000000-0000-0000-0000-000000000002",
		}
		ctx := session.SetSession(context.Background(), claim)

		mockAuthSvc.EXPECT().
			BuildAuthContext(gomock.Any(), claim).
			Return(nil, errors.New("invalid user id"))

		result, err := svc.GetAuthContext(ctx)
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestAuthorizationService_GetUserID(t *testing.T) {
	t.Run("認証されていない場合はnilを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthSvc := mock_service.NewMockAuthenticationService(ctrl)
		svc := NewAuthorizationService(mockAuthSvc)

		ctx := context.Background()

		result, err := svc.GetUserID(ctx)
		require.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("認証されている場合はUserIDを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthSvc := mock_service.NewMockAuthenticationService(ctrl)
		svc := NewAuthorizationService(mockAuthSvc)

		claim := &session.Claim{
			Sub: "00000000-0000-0000-0000-000000000001",
			Jti: "00000000-0000-0000-0000-000000000002",
		}
		ctx := session.SetSession(context.Background(), claim)

		expectedUserID := shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001")
		expectedAuthCtx := &auth.AuthenticationContext{
			UserID: expectedUserID,
		}

		mockAuthSvc.EXPECT().
			BuildAuthContext(gomock.Any(), claim).
			Return(expectedAuthCtx, nil)

		result, err := svc.GetUserID(ctx)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedUserID, *result)
	})
}

func TestAuthorizationService_RequireAuthentication(t *testing.T) {
	t.Run("認証されていない場合はForbiddenErrorを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthSvc := mock_service.NewMockAuthenticationService(ctrl)
		svc := NewAuthorizationService(mockAuthSvc)

		ctx := context.Background()

		result, err := svc.RequireAuthentication(ctx)
		require.Error(t, err)
		assert.ErrorIs(t, err, messages.ForbiddenError)
		assert.Nil(t, result)
	})

	t.Run("認証されている場合はAuthContextを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthSvc := mock_service.NewMockAuthenticationService(ctrl)
		svc := NewAuthorizationService(mockAuthSvc)

		claim := &session.Claim{
			Sub: "00000000-0000-0000-0000-000000000001",
			Jti: "00000000-0000-0000-0000-000000000002",
		}
		ctx := session.SetSession(context.Background(), claim)

		expectedAuthCtx := &auth.AuthenticationContext{
			UserID:       shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			IsRegistered: true,
		}

		mockAuthSvc.EXPECT().
			BuildAuthContext(gomock.Any(), claim).
			Return(expectedAuthCtx, nil)

		result, err := svc.RequireAuthentication(ctx)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedAuthCtx.UserID, result.UserID)
	})
}

func TestAuthorizationService_RequireOrganizationRole(t *testing.T) {
	t.Run("認証されていない場合はForbiddenErrorを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthSvc := mock_service.NewMockAuthenticationService(ctrl)
		svc := NewAuthorizationService(mockAuthSvc)

		ctx := context.Background()

		result, err := svc.RequireOrganizationRole(ctx, organization.OrganizationUserRoleAdmin)
		require.Error(t, err)
		assert.ErrorIs(t, err, messages.ForbiddenError)
		assert.Nil(t, result)
	})

	t.Run("組織コンテキストがない場合はOrganizationContextRequiredを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthSvc := mock_service.NewMockAuthenticationService(ctrl)
		svc := NewAuthorizationService(mockAuthSvc)

		claim := &session.Claim{
			Sub: "00000000-0000-0000-0000-000000000001",
			Jti: "00000000-0000-0000-0000-000000000002",
		}
		ctx := session.SetSession(context.Background(), claim)

		// 組織情報なしのAuthContext
		authCtx := &auth.AuthenticationContext{
			UserID:           shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			OrganizationID:   nil, // 組織なし
			OrganizationRole: nil,
		}

		mockAuthSvc.EXPECT().
			BuildAuthContext(gomock.Any(), claim).
			Return(authCtx, nil)

		result, err := svc.RequireOrganizationRole(ctx, organization.OrganizationUserRoleAdmin)
		require.Error(t, err)
		assert.ErrorIs(t, err, messages.OrganizationContextRequired)
		assert.Nil(t, result)
	})

	t.Run("権限が不足している場合はOrganizationPermissionDeniedを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthSvc := mock_service.NewMockAuthenticationService(ctrl)
		svc := NewAuthorizationService(mockAuthSvc)

		claim := &session.Claim{
			Sub: "00000000-0000-0000-0000-000000000001",
			Jti: "00000000-0000-0000-0000-000000000002",
		}
		ctx := session.SetSession(context.Background(), claim)

		// メンバー権限のAuthContext
		orgID := shared.MustParseUUID[organization.Organization]("00000000-0000-0000-0000-000000000003")
		memberRole := organization.OrganizationUserRoleMember
		authCtx := &auth.AuthenticationContext{
			UserID:           shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			OrganizationID:   &orgID,
			OrganizationRole: &memberRole,
		}

		mockAuthSvc.EXPECT().
			BuildAuthContext(gomock.Any(), claim).
			Return(authCtx, nil)

		// Admin権限を要求
		result, err := svc.RequireOrganizationRole(ctx, organization.OrganizationUserRoleAdmin)
		require.Error(t, err)
		assert.ErrorIs(t, err, messages.OrganizationPermissionDenied)
		assert.Nil(t, result)
	})

	t.Run("十分な権限がある場合はAuthContextを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthSvc := mock_service.NewMockAuthenticationService(ctrl)
		svc := NewAuthorizationService(mockAuthSvc)

		claim := &session.Claim{
			Sub: "00000000-0000-0000-0000-000000000001",
			Jti: "00000000-0000-0000-0000-000000000002",
		}
		ctx := session.SetSession(context.Background(), claim)

		// Owner権限のAuthContext
		orgID := shared.MustParseUUID[organization.Organization]("00000000-0000-0000-0000-000000000003")
		ownerRole := organization.OrganizationUserRoleOwner
		authCtx := &auth.AuthenticationContext{
			UserID:           shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			OrganizationID:   &orgID,
			OrganizationRole: &ownerRole,
		}

		mockAuthSvc.EXPECT().
			BuildAuthContext(gomock.Any(), claim).
			Return(authCtx, nil)

		// Admin権限を要求（Ownerは上位なのでOK）
		result, err := svc.RequireOrganizationRole(ctx, organization.OrganizationUserRoleAdmin)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, authCtx.UserID, result.UserID)
	})
}

func TestAuthorizationService_RequireSuperAdmin(t *testing.T) {
	t.Run("SuperAdmin権限がある場合は成功する", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthSvc := mock_service.NewMockAuthenticationService(ctrl)
		svc := NewAuthorizationService(mockAuthSvc)

		claim := &session.Claim{
			Sub: "00000000-0000-0000-0000-000000000001",
			Jti: "00000000-0000-0000-0000-000000000002",
		}
		ctx := session.SetSession(context.Background(), claim)

		orgID := shared.MustParseUUID[organization.Organization]("00000000-0000-0000-0000-000000000003")
		superAdminRole := organization.OrganizationUserRoleSuperAdmin
		authCtx := &auth.AuthenticationContext{
			UserID:           shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			OrganizationID:   &orgID,
			OrganizationRole: &superAdminRole,
		}

		mockAuthSvc.EXPECT().
			BuildAuthContext(gomock.Any(), claim).
			Return(authCtx, nil)

		result, err := svc.RequireSuperAdmin(ctx)
		require.NoError(t, err)
		require.NotNil(t, result)
	})
}

func TestAuthorizationService_RequireOwner(t *testing.T) {
	t.Run("Owner権限がある場合は成功する", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthSvc := mock_service.NewMockAuthenticationService(ctrl)
		svc := NewAuthorizationService(mockAuthSvc)

		claim := &session.Claim{
			Sub: "00000000-0000-0000-0000-000000000001",
			Jti: "00000000-0000-0000-0000-000000000002",
		}
		ctx := session.SetSession(context.Background(), claim)

		orgID := shared.MustParseUUID[organization.Organization]("00000000-0000-0000-0000-000000000003")
		ownerRole := organization.OrganizationUserRoleOwner
		authCtx := &auth.AuthenticationContext{
			UserID:           shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			OrganizationID:   &orgID,
			OrganizationRole: &ownerRole,
		}

		mockAuthSvc.EXPECT().
			BuildAuthContext(gomock.Any(), claim).
			Return(authCtx, nil)

		result, err := svc.RequireOwner(ctx)
		require.NoError(t, err)
		require.NotNil(t, result)
	})
}

func TestAuthorizationService_RequireAdmin(t *testing.T) {
	t.Run("Admin権限がある場合は成功する", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthSvc := mock_service.NewMockAuthenticationService(ctrl)
		svc := NewAuthorizationService(mockAuthSvc)

		claim := &session.Claim{
			Sub: "00000000-0000-0000-0000-000000000001",
			Jti: "00000000-0000-0000-0000-000000000002",
		}
		ctx := session.SetSession(context.Background(), claim)

		orgID := shared.MustParseUUID[organization.Organization]("00000000-0000-0000-0000-000000000003")
		adminRole := organization.OrganizationUserRoleAdmin
		authCtx := &auth.AuthenticationContext{
			UserID:           shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			OrganizationID:   &orgID,
			OrganizationRole: &adminRole,
		}

		mockAuthSvc.EXPECT().
			BuildAuthContext(gomock.Any(), claim).
			Return(authCtx, nil)

		result, err := svc.RequireAdmin(ctx)
		require.NoError(t, err)
		require.NotNil(t, result)
	})
}

func TestAuthorizationService_IsAuthenticated(t *testing.T) {
	t.Run("セッションがない場合はfalseを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthSvc := mock_service.NewMockAuthenticationService(ctrl)
		svc := NewAuthorizationService(mockAuthSvc)

		ctx := context.Background()

		// セッションがない場合、GetAuthContextは(nil, nil)を返す
		// 認証されていないのでfalseを返すべき
		result := svc.IsAuthenticated(ctx)
		assert.False(t, result, "セッションがない場合は認証されていない")
	})

	t.Run("認証されている場合はtrueを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthSvc := mock_service.NewMockAuthenticationService(ctrl)
		svc := NewAuthorizationService(mockAuthSvc)

		claim := &session.Claim{
			Sub: "00000000-0000-0000-0000-000000000001",
			Jti: "00000000-0000-0000-0000-000000000002",
		}
		ctx := session.SetSession(context.Background(), claim)

		authCtx := &auth.AuthenticationContext{
			UserID: shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
		}

		mockAuthSvc.EXPECT().
			BuildAuthContext(gomock.Any(), claim).
			Return(authCtx, nil)

		result := svc.IsAuthenticated(ctx)
		assert.True(t, result)
	})

	t.Run("BuildAuthContextがエラーを返した場合はfalseを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthSvc := mock_service.NewMockAuthenticationService(ctrl)
		svc := NewAuthorizationService(mockAuthSvc)

		claim := &session.Claim{
			Sub: "invalid-uuid",
			Jti: "00000000-0000-0000-0000-000000000002",
		}
		ctx := session.SetSession(context.Background(), claim)

		mockAuthSvc.EXPECT().
			BuildAuthContext(gomock.Any(), claim).
			Return(nil, errors.New("invalid"))

		result := svc.IsAuthenticated(ctx)
		assert.False(t, result)
	})
}

func TestAuthorizationService_IsInOrganization(t *testing.T) {
	t.Run("セッションがない場合はfalseを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthSvc := mock_service.NewMockAuthenticationService(ctrl)
		svc := NewAuthorizationService(mockAuthSvc)

		ctx := context.Background()

		result := svc.IsInOrganization(ctx)
		assert.False(t, result)
	})

	t.Run("組織に所属している場合はtrueを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthSvc := mock_service.NewMockAuthenticationService(ctrl)
		svc := NewAuthorizationService(mockAuthSvc)

		claim := &session.Claim{
			Sub: "00000000-0000-0000-0000-000000000001",
			Jti: "00000000-0000-0000-0000-000000000002",
		}
		ctx := session.SetSession(context.Background(), claim)

		orgID := shared.MustParseUUID[organization.Organization]("00000000-0000-0000-0000-000000000003")
		authCtx := &auth.AuthenticationContext{
			UserID:         shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			OrganizationID: &orgID,
		}

		mockAuthSvc.EXPECT().
			BuildAuthContext(gomock.Any(), claim).
			Return(authCtx, nil)

		result := svc.IsInOrganization(ctx)
		assert.True(t, result)
	})

	t.Run("組織に所属していない場合はfalseを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthSvc := mock_service.NewMockAuthenticationService(ctrl)
		svc := NewAuthorizationService(mockAuthSvc)

		claim := &session.Claim{
			Sub: "00000000-0000-0000-0000-000000000001",
			Jti: "00000000-0000-0000-0000-000000000002",
		}
		ctx := session.SetSession(context.Background(), claim)

		authCtx := &auth.AuthenticationContext{
			UserID:         shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			OrganizationID: nil,
		}

		mockAuthSvc.EXPECT().
			BuildAuthContext(gomock.Any(), claim).
			Return(authCtx, nil)

		result := svc.IsInOrganization(ctx)
		assert.False(t, result)
	})
}

func TestAuthorizationService_IsKotohiro(t *testing.T) {
	t.Run("認証されていない場合はfalseを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthSvc := mock_service.NewMockAuthenticationService(ctrl)
		svc := NewAuthorizationService(mockAuthSvc)

		ctx := context.Background()

		result := svc.IsKotohiro(ctx)
		assert.False(t, result)
	})

	t.Run("組織に所属していない場合はfalseを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthSvc := mock_service.NewMockAuthenticationService(ctrl)
		svc := NewAuthorizationService(mockAuthSvc)

		claim := &session.Claim{
			Sub: "00000000-0000-0000-0000-000000000001",
			Jti: "00000000-0000-0000-0000-000000000002",
		}
		ctx := session.SetSession(context.Background(), claim)

		authCtx := &auth.AuthenticationContext{
			UserID:         shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			OrganizationID: nil,
		}

		mockAuthSvc.EXPECT().
			BuildAuthContext(gomock.Any(), claim).
			Return(authCtx, nil)

		result := svc.IsKotohiro(ctx)
		assert.False(t, result)
	})

	t.Run("kotohiro組織に所属している場合はtrueを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthSvc := mock_service.NewMockAuthenticationService(ctrl)
		svc := NewAuthorizationService(mockAuthSvc)

		claim := &session.Claim{
			Sub: "00000000-0000-0000-0000-000000000001",
			Jti: "00000000-0000-0000-0000-000000000002",
		}
		ctx := session.SetSession(context.Background(), claim)

		// KotohiroOrganizationID は 00000000-0000-0000-0000-000000000001
		kotohiroOrgID := organization.KotohiroOrganizationID
		authCtx := &auth.AuthenticationContext{
			UserID:         shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			OrganizationID: &kotohiroOrgID,
		}

		mockAuthSvc.EXPECT().
			BuildAuthContext(gomock.Any(), claim).
			Return(authCtx, nil)

		result := svc.IsKotohiro(ctx)
		assert.True(t, result)
	})

	t.Run("kotohiro以外の組織に所属している場合はfalseを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthSvc := mock_service.NewMockAuthenticationService(ctrl)
		svc := NewAuthorizationService(mockAuthSvc)

		claim := &session.Claim{
			Sub: "00000000-0000-0000-0000-000000000001",
			Jti: "00000000-0000-0000-0000-000000000002",
		}
		ctx := session.SetSession(context.Background(), claim)

		orgID := shared.MustParseUUID[organization.Organization]("00000000-0000-0000-0000-000000000003")
		otherCode := "other-org"
		authCtx := &auth.AuthenticationContext{
			UserID:           shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			OrganizationID:   &orgID,
			OrganizationCode: &otherCode,
		}

		mockAuthSvc.EXPECT().
			BuildAuthContext(gomock.Any(), claim).
			Return(authCtx, nil)

		result := svc.IsKotohiro(ctx)
		assert.False(t, result)
	})
}
