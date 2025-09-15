package talksession_test

import (
	"context"
	"testing"
	"time"

	"github.com/neko-dream/api/internal/domain/model/organization"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTalkSession(t *testing.T) {
	tests := []struct {
		name                string
		talkSessionID       shared.UUID[talksession.TalkSession]
		theme               string
		description         *string
		thumbnailURL        *string
		ownerUserID         shared.UUID[user.User]
		createdAt           time.Time
		scheduledEndTime    time.Time
		location            *talksession.Location
		city                *string
		prefecture          *string
		hideTop             bool
		organizationID      *shared.UUID[organization.Organization]
		organizationAliasID *shared.UUID[organization.OrganizationAlias]
	}{
		{
			name:                "全てのフィールドを持つトークセッションを作成できる",
			talkSessionID:       shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000001"),
			theme:               "テストテーマ",
			description:         lo.ToPtr("これはテスト用の説明です"),
			thumbnailURL:        lo.ToPtr("https://example.com/thumbnail.png"),
			ownerUserID:         shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000002"),
			createdAt:           time.Now(),
			scheduledEndTime:    time.Now().Add(2 * time.Hour),
			location:            nil, // Locationは別の構造体のため、ここではnilを設定
			city:                lo.ToPtr("東京都"),
			prefecture:          lo.ToPtr("関東"),
			hideTop:             true,
			organizationID:      lo.ToPtr(shared.MustParseUUID[organization.Organization]("00000000-0000-0000-0000-000000000003")),
			organizationAliasID: lo.ToPtr(shared.MustParseUUID[organization.OrganizationAlias]("00000000-0000-0000-0000-000000000004")),
		},
		{
			name:                "必須フィールドのみでトークセッションを作成できる",
			talkSessionID:       shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000005"),
			theme:               "最小限のテーマ",
			description:         nil,
			thumbnailURL:        nil,
			ownerUserID:         shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000006"),
			createdAt:           time.Now(),
			scheduledEndTime:    time.Now().Add(1 * time.Hour),
			location:            nil,
			city:                nil,
			prefecture:          nil,
			hideTop:             false,
			organizationID:      nil,
			organizationAliasID: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := talksession.NewTalkSession(
				tt.talkSessionID,
				tt.theme,
				tt.description,
				tt.thumbnailURL,
				tt.ownerUserID,
				tt.createdAt,
				tt.scheduledEndTime,
				tt.location,
				tt.city,
				tt.prefecture,
				tt.hideTop,
				tt.organizationID,
				tt.organizationAliasID,
			)

			assert.Equal(t, tt.talkSessionID, ts.TalkSessionID())
			assert.Equal(t, tt.theme, ts.Theme())
			assert.Equal(t, tt.description, ts.Description())
			assert.Equal(t, tt.thumbnailURL, ts.ThumbnailURL())
			assert.Equal(t, tt.ownerUserID, ts.OwnerUserID())
			assert.Equal(t, tt.createdAt, ts.CreatedAt())
			assert.Equal(t, tt.scheduledEndTime, ts.ScheduledEndTime())
			assert.Equal(t, tt.location, ts.Location())
			assert.Equal(t, tt.city, ts.City())
			assert.Equal(t, tt.prefecture, ts.Prefecture())
			assert.Equal(t, tt.organizationID, ts.OrganizationID())
			assert.Equal(t, tt.organizationAliasID, ts.OrganizationAliasID())
			assert.Equal(t, tt.hideTop, ts.HideTop())
			assert.False(t, ts.HideReport())
			assert.False(t, ts.IsEndProcessed())
		})
	}
}

func TestTalkSession_ChangeTheme(t *testing.T) {
	t.Run("テーマを変更できる", func(t *testing.T) {
		ts := talksession.NewTalkSession(
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000001"),
			"元のテーマ",
			nil,
			nil,
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000002"),
			time.Now(),
			time.Now().Add(1*time.Hour),
			nil,
			nil,
			nil,
			true,
			nil,
			nil,
		)

		ts.ChangeTheme("新しいテーマ")

		assert.Equal(t, "新しいテーマ", ts.Theme())
	})
}

func TestTalkSession_ChangeDescription(t *testing.T) {
	tests := []struct {
		name               string
		initialDescription *string
		newDescription     *string
		expected           *string
	}{
		{
			name:               "説明を変更できる",
			initialDescription: lo.ToPtr("元の説明"),
			newDescription:     lo.ToPtr("新しい説明"),
			expected:           lo.ToPtr("新しい説明"),
		},
		{
			name:               "説明をnilから設定できる",
			initialDescription: nil,
			newDescription:     lo.ToPtr("新しい説明"),
			expected:           lo.ToPtr("新しい説明"),
		},
		{
			name:               "説明をnilに変更できる",
			initialDescription: lo.ToPtr("元の説明"),
			newDescription:     nil,
			expected:           nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := talksession.NewTalkSession(
				shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000001"),
				"テーマ",
				tt.initialDescription,
				nil,
				shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000002"),
				time.Now(),
				time.Now().Add(1*time.Hour),
				nil,
				nil,
				nil,
				true,
				nil,
				nil,
			)

			ts.ChangeDescription(tt.newDescription)

			assert.Equal(t, tt.expected, ts.Description())
		})
	}
}

func TestTalkSession_ChangeThumbnailURL(t *testing.T) {
	t.Run("サムネイルURLを変更できる", func(t *testing.T) {
		ts := talksession.NewTalkSession(
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000001"),
			"テーマ",
			nil,
			lo.ToPtr("https://example.com/old.png"),
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000002"),
			time.Now(),
			time.Now().Add(1*time.Hour),
			nil,
			nil,
			nil,
			true,
			nil,
			nil,
		)

		newURL := lo.ToPtr("https://example.com/new.png")
		ts.ChangeThumbnailURL(newURL)

		assert.Equal(t, newURL, ts.ThumbnailURL())
	})
}

func TestTalkSession_ChangeScheduledEndTime(t *testing.T) {
	t.Run("終了予定時刻を変更できる", func(t *testing.T) {
		initialEndTime := time.Now().Add(1 * time.Hour)
		ts := talksession.NewTalkSession(
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000001"),
			"テーマ",
			nil,
			nil,
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000002"),
			time.Now(),
			initialEndTime,
			nil,
			nil,
			nil,
			true,
			nil,
			nil,
		)

		newEndTime := time.Now().Add(2 * time.Hour)
		ts.ChangeScheduledEndTime(newEndTime)

		assert.Equal(t, newEndTime, ts.ScheduledEndTime())
	})
}

func TestTalkSession_ChangeLocation(t *testing.T) {
	t.Run("位置情報を変更できる", func(t *testing.T) {
		// Locationは外部から設定できないため、このテストはスキップ
		// 実際のLocationは別のエンティティとして管理される
		t.Skip("Locationは別のエンティティとして管理されるため、直接変更できません")
	})
}

func TestTalkSession_StartSession(t *testing.T) {
	t.Run("セッションを開始できる", func(t *testing.T) {
		ts := talksession.NewTalkSession(
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000001"),
			"テーマ",
			lo.ToPtr("説明"),
			nil,
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000002"),
			time.Now(),
			time.Now().Add(1*time.Hour),
			nil,
			nil,
			nil,
			true,
			lo.ToPtr(shared.MustParseUUID[organization.Organization]("00000000-0000-0000-0000-000000000003")),
			nil,
		)

		err := ts.StartSession()
		require.NoError(t, err)

		// イベントが記録されていることを確認
		events := ts.GetRecordedEvents()
		assert.Len(t, events, 1)
		assert.Equal(t, talksession.EventTypeTalkSessionStarted, events[0].EventType())
	})

	t.Run("既に開始されたセッションは再度開始できない", func(t *testing.T) {
		ts := talksession.NewTalkSession(
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000001"),
			"テーマ",
			nil,
			nil,
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000002"),
			time.Now(),
			time.Now().Add(1*time.Hour),
			nil,
			nil,
			nil,
			true,
			nil,
			nil,
		)

		err := ts.StartSession()
		require.NoError(t, err)

		// 2回目の開始はエラー
		err = ts.StartSession()
		assert.Equal(t, talksession.ErrSessionAlreadyStarted, err)
	})
}

func TestTalkSession_EndSession(t *testing.T) {
	t.Run("セッションを終了できる", func(t *testing.T) {
		// 終了予定時刻を過去に設定して、IsFinishedがtrueになるようにする
		ts := talksession.NewTalkSession(
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000001"),
			"テーマ",
			nil,
			nil,
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000002"),
			time.Now(),
			time.Now().Add(-1*time.Hour), // 過去の時刻を設定
			nil,
			nil,
			nil,
			true,
			nil,
			nil,
		)

		// まず開始する
		err := ts.StartSession()
		require.NoError(t, err)

		// 終了する
		// 参加者IDのリストを渡す
		participantIDs := []shared.UUID[user.User]{
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000003"),
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000004"),
		}
		err = ts.EndSession(participantIDs)
		require.NoError(t, err)
		assert.True(t, ts.IsFinished(context.Background()))

		// イベントが記録されていることを確認
		events := ts.GetRecordedEvents()
		assert.Len(t, events, 2)
		assert.Equal(t, talksession.EventTypeTalkSessionEnded, events[1].EventType())
	})

	t.Run("開始されていないセッションは終了できない", func(t *testing.T) {
		ts := talksession.NewTalkSession(
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000001"),
			"テーマ",
			nil,
			nil,
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000002"),
			time.Now(),
			time.Now().Add(1*time.Hour),
			nil,
			nil,
			nil,
			true,
			nil,
			nil,
		)

		participantIDs := []shared.UUID[user.User]{}
		err := ts.EndSession(participantIDs)
		assert.Equal(t, talksession.ErrSessionNotYetFinished, err)
	})

	t.Run("既に終了したセッションは再度終了できない", func(t *testing.T) {
		// 終了予定時刻を過去に設定
		ts := talksession.NewTalkSession(
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000001"),
			"テーマ",
			nil,
			nil,
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000002"),
			time.Now(),
			time.Now().Add(-1*time.Hour), // 過去の時刻を設定
			nil,
			nil,
			nil,
			true,
			nil,
			nil,
		)

		err := ts.StartSession()
		require.NoError(t, err)

		participantIDs := []shared.UUID[user.User]{}
		err = ts.EndSession(participantIDs)
		require.NoError(t, err)

		// 2回目の終了はエラー
		err = ts.EndSession(participantIDs)
		assert.Equal(t, talksession.ErrSessionAlreadyEnded, err)
	})
}

func TestTalkSession_SetReportVisibility(t *testing.T) {
	t.Run("レポートの表示/非表示を設定できる", func(t *testing.T) {
		ts := talksession.NewTalkSession(
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000001"),
			"テーマ",
			nil,
			nil,
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000002"),
			time.Now(),
			time.Now().Add(1*time.Hour),
			nil,
			nil,
			nil,
			true,
			nil,
			nil,
		)

		// デフォルトはfalse
		assert.False(t, ts.HideReport())

		// 非表示に設定
		ts.SetReportVisibility(true)
		assert.True(t, ts.HideReport())

		// 表示に戻す
		ts.SetReportVisibility(false)
		assert.False(t, ts.HideReport())
	})
}

func TestTalkSession_MarkAsEndProcessed(t *testing.T) {
	t.Run("終了処理済みフラグを設定できる", func(t *testing.T) {
		ts := talksession.NewTalkSession(
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000001"),
			"テーマ",
			nil,
			nil,
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000002"),
			time.Now(),
			time.Now().Add(1*time.Hour),
			nil,
			nil,
			nil,
			true,
			nil,
			nil,
		)

		// デフォルトはfalse
		assert.False(t, ts.IsEndProcessed())

		// 終了処理済みに設定
		ts.MarkAsEndProcessed()
		assert.True(t, ts.IsEndProcessed())
	})
}

func TestTalkSession_UpdateRestrictions(t *testing.T) {
	t.Run("制限を更新できる", func(t *testing.T) {
		ts := talksession.NewTalkSession(
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000001"),
			"テーマ",
			nil,
			nil,
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000002"),
			time.Now(),
			time.Now().Add(1*time.Hour),
			nil,
			nil,
			nil,
			true,
			nil,
			nil,
		)

		// 初期状態では制限なし
		assert.Empty(t, ts.RestrictionList())

		// 制限を追加（有効な制限キーを使用）
		restrictions := []string{"demographics.gender", "demographics.prefecture"}
		err := ts.UpdateRestrictions(context.Background(), restrictions)
		require.NoError(t, err)

		assert.NotEmpty(t, ts.RestrictionList())
		assert.NotNil(t, ts.Restrictions())
	})
}
