package session

import (
	"testing"
	"time"

	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/stretchr/testify/assert"
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
