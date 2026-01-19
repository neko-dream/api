package session

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClaim(t *testing.T) {
	now := time.Now()
	ctx := clock.SetNow(context.Background(), now)
	sessionID := shared.NewUUID[Session]()

	t.Run("通常ユーザーのClaimを作成できる", func(t *testing.T) {
		u := user.NewUser(
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			lo.ToPtr("testuser"),
			lo.ToPtr("Test User"),
			"google|123456",
			shared.ProviderGoogle,
			lo.ToPtr("https://example.com/icon.png"),
		)
		u.ChangeEmail("test@example.com")
		u.SetEmailVerified(true)

		claim := NewClaim(ctx, u, sessionID, false, nil)

		assert.Equal(t, u.UserID().String(), claim.Sub)
		assert.Equal(t, sessionID.String(), claim.Jti)
		assert.Equal(t, now.Unix(), claim.Iat)
		// 通常ユーザーは7日間有効
		expectedExp := now.Add(7 * 24 * time.Hour).Unix()
		assert.Equal(t, expectedExp, claim.Exp)
		assert.Equal(t, "Test User", *claim.DisplayName)
		assert.Equal(t, "testuser", *claim.DisplayID)
		assert.Equal(t, "https://example.com/icon.png", *claim.IconURL)
		assert.True(t, claim.IsEmailVerified)
		assert.False(t, claim.IsGuest)
		assert.False(t, claim.RequiredPasswordChange)
	})

	t.Run("ゲストユーザーは30日間有効", func(t *testing.T) {
		u := user.NewUser(
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			lo.ToPtr("12345678"),
			lo.ToPtr("ゲストabc123"),
			"guest|abc123",
			shared.ProviderGuest,
			nil,
		)

		claim := NewClaim(ctx, u, sessionID, false, nil)

		// ゲストユーザーは30日間有効
		expectedExp := now.Add(30 * 24 * time.Hour).Unix()
		assert.Equal(t, expectedExp, claim.Exp)
		assert.True(t, claim.IsGuest)
		assert.True(t, claim.IsRegistered, "ゲストユーザーも登録済み扱い")
	})

	t.Run("パスワード変更要求フラグを設定できる", func(t *testing.T) {
		u := user.NewUser(
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			lo.ToPtr("testuser"),
			lo.ToPtr("Test User"),
			"google|123456",
			shared.ProviderGoogle,
			nil,
		)

		claim := NewClaim(ctx, u, sessionID, true, nil)

		assert.True(t, claim.RequiredPasswordChange)
	})
}

func TestNewClaimWithOrganization(t *testing.T) {
	now := time.Now()
	ctx := clock.SetNow(context.Background(), now)
	sessionID := shared.NewUUID[Session]()

	t.Run("組織情報を含むClaimを作成できる", func(t *testing.T) {
		u := user.NewUser(
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			lo.ToPtr("testuser"),
			lo.ToPtr("Test User"),
			"google|123456",
			shared.ProviderGoogle,
			nil,
		)

		orgID := "00000000-0000-0000-0000-000000000002"
		orgCode := "ORG001"
		orgRole := "管理者"
		orgType := 1

		claim := NewClaimWithOrganization(ctx, u, sessionID, false, &orgType, &orgID, &orgCode, &orgRole)

		assert.Equal(t, u.UserID().String(), claim.Sub)
		require.NotNil(t, claim.OrganizationID)
		assert.Equal(t, orgID, *claim.OrganizationID)
		require.NotNil(t, claim.OrganizationCode)
		assert.Equal(t, orgCode, *claim.OrganizationCode)
		require.NotNil(t, claim.OrganizationRole)
		assert.Equal(t, orgRole, *claim.OrganizationRole)
		require.NotNil(t, claim.OrgType)
		assert.Equal(t, orgType, *claim.OrgType)
	})

	t.Run("ゲストユーザーの組織Claimは30日有効", func(t *testing.T) {
		u := user.NewUser(
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001"),
			lo.ToPtr("12345678"),
			lo.ToPtr("ゲストabc123"),
			"guest|abc123",
			shared.ProviderGuest,
			nil,
		)

		orgID := "00000000-0000-0000-0000-000000000002"
		orgCode := "ORG001"
		orgRole := "一般"
		orgType := 1

		claim := NewClaimWithOrganization(ctx, u, sessionID, false, &orgType, &orgID, &orgCode, &orgRole)

		// ゲストユーザーは30日間有効
		expectedExp := now.Add(30 * 24 * time.Hour).Unix()
		assert.Equal(t, expectedExp, claim.Exp)
		assert.True(t, claim.IsGuest)
		assert.True(t, claim.IsRegistered, "ゲストユーザーも登録済み扱い")
		// 組織情報も設定されている
		require.NotNil(t, claim.OrganizationID)
		assert.Equal(t, orgID, *claim.OrganizationID)
	})
}

func TestNewClaimFromMap(t *testing.T) {
	t.Run("jwt.MapClaimsからClaimを復元できる", func(t *testing.T) {
		now := time.Now()
		claims := jwt.MapClaims{
			"sub":                    "00000000-0000-0000-0000-000000000001",
			"iat":                    float64(now.Unix()),
			"exp":                    float64(now.Add(7 * 24 * time.Hour).Unix()),
			"jti":                    "00000000-0000-0000-0000-000000000002",
			"iconURL":                "https://example.com/icon.png",
			"displayName":            "Test User",
			"displayID":              "testuser",
			"isRegistered":           true,
			"isEmailVerified":        true,
			"requiredPasswordChange": false,
			"isGuest":                false,
		}

		claim := NewClaimFromMap(claims)

		assert.Equal(t, "00000000-0000-0000-0000-000000000001", claim.Sub)
		assert.Equal(t, "00000000-0000-0000-0000-000000000002", claim.Jti)
		assert.Equal(t, now.Unix(), claim.Iat)
		assert.Equal(t, "https://example.com/icon.png", *claim.IconURL)
		assert.Equal(t, "Test User", *claim.DisplayName)
		assert.Equal(t, "testuser", *claim.DisplayID)
		assert.True(t, claim.IsRegistered)
		assert.True(t, claim.IsEmailVerified)
		assert.False(t, claim.RequiredPasswordChange)
		assert.False(t, claim.IsGuest)
	})

	t.Run("組織情報を含むMapClaimsからClaimを復元できる", func(t *testing.T) {
		now := time.Now()
		claims := jwt.MapClaims{
			"sub":              "00000000-0000-0000-0000-000000000001",
			"iat":              float64(now.Unix()),
			"exp":              float64(now.Add(7 * 24 * time.Hour).Unix()),
			"jti":              "00000000-0000-0000-0000-000000000002",
			"isRegistered":     true,
			"orgType":          float64(1),
			"organizationID":   "00000000-0000-0000-0000-000000000003",
			"organizationCode": "ORG001",
			"organizationRole": "管理者",
		}

		claim := NewClaimFromMap(claims)

		require.NotNil(t, claim.OrgType)
		assert.Equal(t, 1, *claim.OrgType)
		require.NotNil(t, claim.OrganizationID)
		assert.Equal(t, "00000000-0000-0000-0000-000000000003", *claim.OrganizationID)
		require.NotNil(t, claim.OrganizationCode)
		assert.Equal(t, "ORG001", *claim.OrganizationCode)
		require.NotNil(t, claim.OrganizationRole)
		assert.Equal(t, "管理者", *claim.OrganizationRole)
	})

	t.Run("オプションフィールドがnilの場合も正しく処理される", func(t *testing.T) {
		now := time.Now()
		claims := jwt.MapClaims{
			"sub":          "00000000-0000-0000-0000-000000000001",
			"iat":          float64(now.Unix()),
			"exp":          float64(now.Add(time.Hour).Unix()),
			"jti":          "00000000-0000-0000-0000-000000000002",
			"isRegistered": true,
			// オプションフィールドは省略
		}

		claim := NewClaimFromMap(claims)

		assert.Nil(t, claim.IconURL)
		assert.Nil(t, claim.DisplayName)
		assert.Nil(t, claim.DisplayID)
		assert.Nil(t, claim.OrgType)
		assert.Nil(t, claim.OrganizationID)
		assert.False(t, claim.IsGuest)
	})
}

func TestClaim_UserID(t *testing.T) {
	t.Run("正しいUserIDを返す", func(t *testing.T) {
		claim := &Claim{
			Sub: "00000000-0000-0000-0000-000000000001",
		}

		userID, err := claim.UserID()
		require.NoError(t, err)
		assert.Equal(t, "00000000-0000-0000-0000-000000000001", userID.String())
	})

	t.Run("不正なUUIDの場合はエラーを返す", func(t *testing.T) {
		claim := &Claim{
			Sub: "invalid-uuid",
		}

		_, err := claim.UserID()
		require.Error(t, err)
	})
}

func TestClaim_SessionID(t *testing.T) {
	t.Run("正しいSessionIDを返す", func(t *testing.T) {
		claim := &Claim{
			Jti: "00000000-0000-0000-0000-000000000001",
		}

		sessionID, err := claim.SessionID()
		require.NoError(t, err)
		assert.Equal(t, "00000000-0000-0000-0000-000000000001", sessionID.String())
	})
}

func TestClaim_IsExpired(t *testing.T) {
	t.Run("有効期限内ならfalse", func(t *testing.T) {
		now := time.Now()
		ctx := clock.SetNow(context.Background(), now)

		claim := &Claim{
			Exp: now.Add(time.Hour).Unix(),
		}

		assert.False(t, claim.IsExpired(ctx))
	})

	t.Run("有効期限切れならtrue", func(t *testing.T) {
		now := time.Now()
		ctx := clock.SetNow(context.Background(), now)

		claim := &Claim{
			Exp: now.Add(-time.Hour).Unix(),
		}

		assert.True(t, claim.IsExpired(ctx))
	})
}

func TestClaim_IssueAt_ExpiresAt(t *testing.T) {
	now := time.Now()
	exp := now.Add(7 * 24 * time.Hour)

	claim := &Claim{
		Iat: now.Unix(),
		Exp: exp.Unix(),
	}

	assert.Equal(t, now.Unix(), claim.IssueAt().Unix())
	assert.Equal(t, exp.Unix(), claim.ExpiresAt().Unix())
}

func TestClaim_AudienceIssuer(t *testing.T) {
	claim := &Claim{}

	assert.Equal(t, "https://api.kotohiro.com", claim.Audience())
	assert.Equal(t, "https://api.kotohiro.com", claim.Issuer())
}

func TestClaim_GenMapClaim(t *testing.T) {
	t.Run("基本的なClaimをMapClaimsに変換できる", func(t *testing.T) {
		claim := &Claim{
			Sub:                    "00000000-0000-0000-0000-000000000001",
			Iat:                    1000000,
			Exp:                    2000000,
			Jti:                    "00000000-0000-0000-0000-000000000002",
			IconURL:                lo.ToPtr("https://example.com/icon.png"),
			DisplayName:            lo.ToPtr("Test User"),
			DisplayID:              lo.ToPtr("testuser"),
			IsRegistered:           true,
			IsEmailVerified:        true,
			RequiredPasswordChange: false,
			IsGuest:                false,
		}

		mapClaims := claim.GenMapClaim()

		assert.Equal(t, claim.Sub, (*mapClaims)["sub"])
		assert.Equal(t, claim.Iat, (*mapClaims)["iat"])
		assert.Equal(t, claim.Exp, (*mapClaims)["exp"])
		assert.Equal(t, claim.Jti, (*mapClaims)["jti"])
		assert.Equal(t, Issuer, (*mapClaims)["iss"])
		assert.Equal(t, Audience, (*mapClaims)["aud"])
		assert.Equal(t, claim.IconURL, (*mapClaims)["iconURL"])
		assert.Equal(t, claim.DisplayName, (*mapClaims)["displayName"])
		assert.Equal(t, claim.DisplayID, (*mapClaims)["displayID"])
		assert.Equal(t, claim.IsRegistered, (*mapClaims)["isRegistered"])
		assert.Equal(t, claim.IsEmailVerified, (*mapClaims)["isEmailVerified"])
		assert.Equal(t, claim.RequiredPasswordChange, (*mapClaims)["requiredPasswordChange"])
		assert.Equal(t, claim.IsGuest, (*mapClaims)["isGuest"])
	})

	t.Run("組織情報を含むClaimをMapClaimsに変換できる", func(t *testing.T) {
		orgType := 1
		claim := &Claim{
			Sub:              "00000000-0000-0000-0000-000000000001",
			Iat:              1000000,
			Exp:              2000000,
			Jti:              "00000000-0000-0000-0000-000000000002",
			IsRegistered:     true,
			OrgType:          &orgType,
			OrganizationID:   lo.ToPtr("00000000-0000-0000-0000-000000000003"),
			OrganizationCode: lo.ToPtr("ORG001"),
			OrganizationRole: lo.ToPtr("管理者"),
		}

		mapClaims := claim.GenMapClaim()

		assert.Equal(t, orgType, (*mapClaims)["orgType"])
		assert.Equal(t, "00000000-0000-0000-0000-000000000003", (*mapClaims)["organizationID"])
		assert.Equal(t, "ORG001", (*mapClaims)["organizationCode"])
		assert.Equal(t, "管理者", (*mapClaims)["organizationRole"])
	})

	t.Run("OrgTypeがnilの場合はnilが設定される", func(t *testing.T) {
		claim := &Claim{
			Sub:          "00000000-0000-0000-0000-000000000001",
			Iat:          1000000,
			Exp:          2000000,
			Jti:          "00000000-0000-0000-0000-000000000002",
			IsRegistered: true,
			OrgType:      nil,
		}

		mapClaims := claim.GenMapClaim()

		assert.Nil(t, (*mapClaims)["orgType"])
	})
}

func TestSetSession_GetSession(t *testing.T) {
	t.Run("セッションをコンテキストに設定・取得できる", func(t *testing.T) {
		claim := &Claim{
			Sub: "00000000-0000-0000-0000-000000000001",
			Jti: "00000000-0000-0000-0000-000000000002",
		}

		ctx := SetSession(context.Background(), claim)
		result := GetSession(ctx)

		require.NotNil(t, result)
		assert.Equal(t, claim.Sub, result.Sub)
		assert.Equal(t, claim.Jti, result.Jti)
	})

	t.Run("セッションが設定されていない場合はnilを返す", func(t *testing.T) {
		ctx := context.Background()
		result := GetSession(ctx)

		assert.Nil(t, result)
	})
}
