package session

import (
	"context"
	"testing"
	"time"

	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSessionStatus(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected *status
	}{
		{
			name:     "0はSESSION_ACTIVEを返す",
			input:    0,
			expected: func() *status { s := SESSION_ACTIVE; return &s }(),
		},
		{
			name:     "1はSESSION_INACTIVEを返す",
			input:    1,
			expected: func() *status { s := SESSION_INACTIVE; return &s }(),
		},
		{
			name:     "2はnilを返す",
			input:    2,
			expected: nil,
		},
		{
			name:     "-1はnilを返す",
			input:    -1,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewSessionStatus(tt.input)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, *tt.expected, *result)
			}
		})
	}
}

func TestSortByLastActivity(t *testing.T) {
	userID := shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001")
	now := time.Now()

	t.Run("lastActivityが新しい順にソートされる", func(t *testing.T) {
		sessions := []Session{
			*NewSession(
				shared.MustParseUUID[Session]("00000000-0000-0000-0000-000000000001"),
				userID,
				shared.ProviderGoogle,
				SESSION_ACTIVE,
				now.Add(time.Hour),
				now.Add(-3*time.Hour), // 3時間前
			),
			*NewSession(
				shared.MustParseUUID[Session]("00000000-0000-0000-0000-000000000002"),
				userID,
				shared.ProviderGoogle,
				SESSION_ACTIVE,
				now.Add(time.Hour),
				now.Add(-1*time.Hour), // 1時間前（最新）
			),
			*NewSession(
				shared.MustParseUUID[Session]("00000000-0000-0000-0000-000000000003"),
				userID,
				shared.ProviderGoogle,
				SESSION_ACTIVE,
				now.Add(time.Hour),
				now.Add(-2*time.Hour), // 2時間前
			),
		}

		sorted := SortByLastActivity(sessions)

		// lastActivityが新しい順: session2, session3, session1
		assert.Equal(t, "00000000-0000-0000-0000-000000000002", sorted[0].SessionID().String())
		assert.Equal(t, "00000000-0000-0000-0000-000000000003", sorted[1].SessionID().String())
		assert.Equal(t, "00000000-0000-0000-0000-000000000001", sorted[2].SessionID().String())
	})

	t.Run("元のスライスは変更されない", func(t *testing.T) {
		sessions := []Session{
			*NewSession(
				shared.MustParseUUID[Session]("00000000-0000-0000-0000-000000000001"),
				userID,
				shared.ProviderGoogle,
				SESSION_ACTIVE,
				now.Add(time.Hour),
				now.Add(-2*time.Hour),
			),
			*NewSession(
				shared.MustParseUUID[Session]("00000000-0000-0000-0000-000000000002"),
				userID,
				shared.ProviderGoogle,
				SESSION_ACTIVE,
				now.Add(time.Hour),
				now.Add(-1*time.Hour),
			),
		}

		originalFirstID := sessions[0].SessionID().String()
		_ = SortByLastActivity(sessions)

		// 元のスライスの順番は変わらないはず
		assert.Equal(t, originalFirstID, sessions[0].SessionID().String())
	})
}

func TestNewExpiresAt(t *testing.T) {
	t.Run("14日後の有効期限を返す", func(t *testing.T) {
		now := time.Now()
		ctx := clock.SetNow(context.Background(), now)

		expires := NewExpiresAt(ctx)
		require.NotNil(t, expires)

		// 2 * 24 * 7 = 14日
		expected := now.Add(14 * 24 * time.Hour)
		assert.Equal(t, expected.Unix(), expires.Unix())
	})
}

func TestNewGuestExpiresAt(t *testing.T) {
	t.Run("30日後の有効期限を返す", func(t *testing.T) {
		now := time.Now()
		ctx := clock.SetNow(context.Background(), now)

		expires := NewGuestExpiresAt(ctx)
		require.NotNil(t, expires)

		expected := now.Add(30 * 24 * time.Hour)
		assert.Equal(t, expected.Unix(), expires.Unix())
	})
}

func TestSession_IsActive(t *testing.T) {
	userID := shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001")
	now := time.Now()
	ctx := clock.SetNow(context.Background(), now)

	t.Run("有効期限内でACTIVEステータスならtrue", func(t *testing.T) {
		sess := NewSession(
			shared.NewUUID[Session](),
			userID,
			shared.ProviderGoogle,
			SESSION_ACTIVE,
			now.Add(time.Hour), // 1時間後に期限切れ
			now,
		)

		assert.True(t, sess.IsActive(ctx))
	})

	t.Run("有効期限切れならfalse", func(t *testing.T) {
		sess := NewSession(
			shared.NewUUID[Session](),
			userID,
			shared.ProviderGoogle,
			SESSION_ACTIVE,
			now.Add(-time.Hour), // 1時間前に期限切れ
			now,
		)

		assert.False(t, sess.IsActive(ctx))
	})

	t.Run("INACTIVEステータスならfalse", func(t *testing.T) {
		sess := NewSession(
			shared.NewUUID[Session](),
			userID,
			shared.ProviderGoogle,
			SESSION_INACTIVE,
			now.Add(time.Hour), // 有効期限内だが
			now,
		)

		assert.False(t, sess.IsActive(ctx))
	})
}

func TestSession_Deactivate(t *testing.T) {
	userID := shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001")
	now := time.Now()
	ctx := clock.SetNow(context.Background(), now)

	t.Run("セッションを無効化する", func(t *testing.T) {
		sess := NewSession(
			shared.NewUUID[Session](),
			userID,
			shared.ProviderGoogle,
			SESSION_ACTIVE,
			now.Add(time.Hour),
			now.Add(-time.Hour),
		)

		assert.True(t, sess.IsActive(ctx), "無効化前はアクティブ")

		sess.Deactivate(ctx)

		assert.False(t, sess.IsActive(ctx), "無効化後は非アクティブ")
		assert.Equal(t, SESSION_INACTIVE, sess.Status())
		// lastActivityも更新される
		assert.Equal(t, now.Unix(), sess.LastActivityAt().Unix())
	})
}

func TestSession_UpdateLastActivity(t *testing.T) {
	userID := shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001")
	oldTime := time.Now().Add(-time.Hour)
	newTime := time.Now()
	ctx := clock.SetNow(context.Background(), newTime)

	t.Run("lastActivityを現在時刻に更新する", func(t *testing.T) {
		sess := NewSession(
			shared.NewUUID[Session](),
			userID,
			shared.ProviderGoogle,
			SESSION_ACTIVE,
			newTime.Add(time.Hour),
			oldTime, // 1時間前
		)

		assert.Equal(t, oldTime.Unix(), sess.LastActivityAt().Unix())

		sess.UpdateLastActivity(ctx)

		assert.Equal(t, newTime.Unix(), sess.LastActivityAt().Unix())
	})
}

func TestNewSessionWithOrganization(t *testing.T) {
	userID := shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001")
	sessionID := shared.NewUUID[Session]()
	orgID := shared.NewUUID[any]()
	now := time.Now()

	t.Run("組織IDを含むセッションを作成できる", func(t *testing.T) {
		sess := NewSessionWithOrganization(
			sessionID,
			userID,
			shared.ProviderGoogle,
			SESSION_ACTIVE,
			now.Add(time.Hour),
			now,
			&orgID,
		)

		require.NotNil(t, sess)
		assert.Equal(t, sessionID, sess.SessionID())
		assert.Equal(t, userID, sess.UserID())
		assert.Equal(t, shared.ProviderGoogle, sess.Provider())
		require.NotNil(t, sess.OrganizationID())
		assert.Equal(t, orgID, *sess.OrganizationID())
	})

	t.Run("組織IDなしでも作成できる", func(t *testing.T) {
		sess := NewSessionWithOrganization(
			sessionID,
			userID,
			shared.ProviderGoogle,
			SESSION_ACTIVE,
			now.Add(time.Hour),
			now,
			nil,
		)

		require.NotNil(t, sess)
		assert.Nil(t, sess.OrganizationID())
	})
}

func TestSession_ExpiresAt(t *testing.T) {
	userID := shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001")
	now := time.Now()
	expiresTime := now.Add(24 * time.Hour)

	t.Run("正しい有効期限を返す", func(t *testing.T) {
		sess := NewSession(
			shared.NewUUID[Session](),
			userID,
			shared.ProviderGoogle,
			SESSION_ACTIVE,
			expiresTime,
			now,
		)

		result := sess.ExpiresAt()
		assert.Equal(t, expiresTime.Unix(), result.Unix())
	})
}

func TestFilterActiveSessions(t *testing.T) {
	userID := shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001")
	now := time.Now()
	ctx := clock.SetNow(context.Background(), now)

	t.Run("アクティブなセッションのみをフィルタする", func(t *testing.T) {
		sessions := []Session{
			*NewSession(
				shared.MustParseUUID[Session]("00000000-0000-0000-0000-000000000001"),
				userID,
				shared.ProviderGoogle,
				SESSION_ACTIVE,
				now.Add(time.Hour), // アクティブ
				now,
			),
			*NewSession(
				shared.MustParseUUID[Session]("00000000-0000-0000-0000-000000000002"),
				userID,
				shared.ProviderGoogle,
				SESSION_INACTIVE, // 非アクティブ
				now.Add(time.Hour),
				now,
			),
			*NewSession(
				shared.MustParseUUID[Session]("00000000-0000-0000-0000-000000000003"),
				userID,
				shared.ProviderGoogle,
				SESSION_ACTIVE,
				now.Add(-time.Hour), // 期限切れ
				now,
			),
			*NewSession(
				shared.MustParseUUID[Session]("00000000-0000-0000-0000-000000000004"),
				userID,
				shared.ProviderGoogle,
				SESSION_ACTIVE,
				now.Add(2*time.Hour), // アクティブ
				now,
			),
		}

		filtered := FilterActiveSessions(ctx, sessions)

		assert.Len(t, filtered, 2, "アクティブなセッションは2つ")
		assert.Equal(t, "00000000-0000-0000-0000-000000000001", filtered[0].SessionID().String())
		assert.Equal(t, "00000000-0000-0000-0000-000000000004", filtered[1].SessionID().String())
	})

	t.Run("全て非アクティブなら空リストを返す", func(t *testing.T) {
		sessions := []Session{
			*NewSession(
				shared.NewUUID[Session](),
				userID,
				shared.ProviderGoogle,
				SESSION_INACTIVE,
				now.Add(time.Hour),
				now,
			),
		}

		filtered := FilterActiveSessions(ctx, sessions)
		assert.Empty(t, filtered)
	})
}
