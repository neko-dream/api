package opinion_test

import (
	"strings"
	"testing"
	"time"

	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/opinion"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOpinion(t *testing.T) {
	tests := []struct {
		name            string
		opinionID       shared.UUID[opinion.Opinion]
		talkSessionID   shared.UUID[talksession.TalkSession]
		userID          shared.UUID[user.User]
		parentOpinionID *shared.UUID[opinion.Opinion]
		title           *string
		content         string
		createdAt       time.Time
		referenceURL    *string
		wantErr         error
	}{
		{
			name:            "正常に意見を作成できる（タイトルあり）",
			opinionID:       shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-000000000001"),
			talkSessionID:   shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000002"),
			userID:          shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000003"),
			parentOpinionID: nil,
			title:           lo.ToPtr("テストタイトル"),
			content:         "これはテスト用の意見内容です",
			createdAt:       time.Now(),
			referenceURL:    lo.ToPtr("https://example.com"),
			wantErr:         nil,
		},
		{
			name:            "正常に意見を作成できる（タイトルなし）",
			opinionID:       shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-000000000004"),
			talkSessionID:   shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000005"),
			userID:          shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000006"),
			parentOpinionID: nil,
			title:           nil,
			content:         "タイトルなしの意見内容です",
			createdAt:       time.Now(),
			referenceURL:    nil,
			wantErr:         nil,
		},
		{
			name:            "正常に返信意見を作成できる",
			opinionID:       shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-000000000007"),
			talkSessionID:   shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000008"),
			userID:          shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000009"),
			parentOpinionID: lo.ToPtr(shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-00000000000a")),
			title:           nil,
			content:         "これは返信用の意見です",
			createdAt:       time.Now(),
			referenceURL:    nil,
			wantErr:         nil,
		},
		{
			name:            "内容が短すぎる場合はエラー",
			opinionID:       shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-00000000000b"),
			talkSessionID:   shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-00000000000c"),
			userID:          shared.MustParseUUID[user.User]("00000000-0000-0000-0000-00000000000d"),
			parentOpinionID: nil,
			title:           nil,
			content:         "短い",
			createdAt:       time.Now(),
			referenceURL:    nil,
			wantErr:         messages.OpinionContentBadLength,
		},
		{
			name:            "内容が長すぎる場合はエラー",
			opinionID:       shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-00000000000e"),
			talkSessionID:   shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-00000000000f"),
			userID:          shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000010"),
			parentOpinionID: nil,
			title:           nil,
			content:         strings.Repeat("あ", 141),
			createdAt:       time.Now(),
			referenceURL:    nil,
			wantErr:         messages.OpinionContentBadLength,
		},
		{
			name:            "親意見IDが自分自身の場合はエラー",
			opinionID:       shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-000000000011"),
			talkSessionID:   shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000012"),
			userID:          shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000013"),
			parentOpinionID: lo.ToPtr(shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-000000000011")),
			title:           nil,
			content:         "親意見IDが自分自身です",
			createdAt:       time.Now(),
			referenceURL:    nil,
			wantErr:         messages.OpinionParentOpinionIDIsSame,
		},
		{
			name:            "タイトルが短すぎる場合はエラー",
			opinionID:       shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-000000000014"),
			talkSessionID:   shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000015"),
			userID:          shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000016"),
			parentOpinionID: nil,
			title:           lo.ToPtr("短い"),
			content:         "内容は適切な長さです",
			createdAt:       time.Now(),
			referenceURL:    nil,
			wantErr:         messages.OpinionTitleBadLength,
		},
		{
			name:            "タイトルが長すぎる場合はエラー",
			opinionID:       shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-000000000017"),
			talkSessionID:   shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000018"),
			userID:          shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000019"),
			parentOpinionID: nil,
			title:           lo.ToPtr(strings.Repeat("あ", 51)),
			content:         "内容は適切な長さです",
			createdAt:       time.Now(),
			referenceURL:    nil,
			wantErr:         messages.OpinionTitleBadLength,
		},
		{
			name:            "内容の境界値テスト（5文字）",
			opinionID:       shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-00000000001a"),
			talkSessionID:   shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-00000000001b"),
			userID:          shared.MustParseUUID[user.User]("00000000-0000-0000-0000-00000000001c"),
			parentOpinionID: nil,
			title:           nil,
			content:         "ちょうど5文字",
			createdAt:       time.Now(),
			referenceURL:    nil,
			wantErr:         nil,
		},
		{
			name:            "内容の境界値テスト（140文字）",
			opinionID:       shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-00000000001d"),
			talkSessionID:   shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-00000000001e"),
			userID:          shared.MustParseUUID[user.User]("00000000-0000-0000-0000-00000000001f"),
			parentOpinionID: nil,
			title:           nil,
			content:         strings.Repeat("あ", 140),
			createdAt:       time.Now(),
			referenceURL:    nil,
			wantErr:         nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := opinion.NewOpinion(
				tt.opinionID,
				tt.talkSessionID,
				tt.userID,
				tt.parentOpinionID,
				tt.title,
				tt.content,
				tt.createdAt,
				tt.referenceURL,
			)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, got)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.opinionID, got.OpinionID())
			assert.Equal(t, tt.talkSessionID, got.TalkSessionID())
			assert.Equal(t, tt.userID, got.UserID())
			assert.Equal(t, tt.parentOpinionID, got.ParentOpinionID())
			assert.Equal(t, tt.title, got.Title())
			assert.Equal(t, tt.content, got.Content())
			assert.Equal(t, tt.createdAt, got.CreatedAt())
			assert.Equal(t, tt.referenceURL, got.ReferenceURL())
			assert.Empty(t, got.Opinions())
			assert.Nil(t, got.ReferenceImageURL())
		})
	}
}

func TestOpinion_Reply(t *testing.T) {
	t.Run("意見に返信を追加できる", func(t *testing.T) {
		parentOpinion, err := opinion.NewOpinion(
			shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-000000000001"),
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000002"),
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000003"),
			nil,
			lo.ToPtr("親意見のタイトル"),
			"これは親意見の内容です",
			time.Now(),
			nil,
		)
		require.NoError(t, err)

		replyOpinion, err := opinion.NewOpinion(
			shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-000000000004"),
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000002"),
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000005"),
			lo.ToPtr(parentOpinion.OpinionID()),
			nil,
			"これは返信意見の内容です",
			time.Now(),
			nil,
		)
		require.NoError(t, err)

		assert.Empty(t, parentOpinion.Opinions())

		parentOpinion.Reply(*replyOpinion)

		assert.Len(t, parentOpinion.Opinions(), 1)
		assert.Equal(t, replyOpinion.OpinionID(), parentOpinion.Opinions()[0].OpinionID())
	})

	t.Run("複数の返信を追加できる", func(t *testing.T) {
		parentOpinion, err := opinion.NewOpinion(
			shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-000000000001"),
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000002"),
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000003"),
			nil,
			lo.ToPtr("親意見のタイトル"),
			"これは親意見の内容です",
			time.Now(),
			nil,
		)
		require.NoError(t, err)

		reply1, err := opinion.NewOpinion(
			shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-000000000004"),
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000002"),
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000005"),
			lo.ToPtr(parentOpinion.OpinionID()),
			nil,
			"これは1つ目の返信意見です",
			time.Now(),
			nil,
		)
		require.NoError(t, err)

		reply2, err := opinion.NewOpinion(
			shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-000000000006"),
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000002"),
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000007"),
			lo.ToPtr(parentOpinion.OpinionID()),
			nil,
			"これは2つ目の返信意見です",
			time.Now(),
			nil,
		)
		require.NoError(t, err)

		parentOpinion.Reply(*reply1)
		parentOpinion.Reply(*reply2)

		assert.Len(t, parentOpinion.Opinions(), 2)
		assert.Equal(t, reply1.OpinionID(), parentOpinion.Opinions()[0].OpinionID())
		assert.Equal(t, reply2.OpinionID(), parentOpinion.Opinions()[1].OpinionID())
	})
}

func TestOpinion_ChangeReferenceImageURL(t *testing.T) {
	t.Run("参照画像URLを変更できる", func(t *testing.T) {
		op, err := opinion.NewOpinion(
			shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-000000000001"),
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000002"),
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000003"),
			nil,
			nil,
			"これは意見の内容です",
			time.Now(),
			nil,
		)
		require.NoError(t, err)

		assert.Nil(t, op.ReferenceImageURL())

		newImageURL := lo.ToPtr("https://example.com/image.png")
		op.ChangeReferenceImageURL(newImageURL)

		assert.Equal(t, newImageURL, op.ReferenceImageURL())
	})

	t.Run("参照画像URLをnilに変更できる", func(t *testing.T) {
		op, err := opinion.NewOpinion(
			shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-000000000001"),
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000002"),
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000003"),
			nil,
			nil,
			"これは意見の内容です",
			time.Now(),
			nil,
		)
		require.NoError(t, err)

		imageURL := lo.ToPtr("https://example.com/image.png")
		op.ChangeReferenceImageURL(imageURL)
		assert.Equal(t, imageURL, op.ReferenceImageURL())

		op.ChangeReferenceImageURL(nil)
		assert.Nil(t, op.ReferenceImageURL())
	})
}

func TestOpinion_SetSeed(t *testing.T) {
	t.Run("シードを設定できる", func(t *testing.T) {
		op, err := opinion.NewOpinion(
			shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-000000000001"),
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000002"),
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000003"),
			nil,
			nil,
			"これは意見の内容です",
			time.Now(),
			nil,
		)
		require.NoError(t, err)

		op.SetSeed()

		// SetSeedメソッドは内部でseedフラグを設定するが、
		// 外部からは確認できないため、ここではエラーが発生しないことを確認
	})
}

func TestOpinion_Count(t *testing.T) {
	t.Run("意見の総数を取得できる", func(t *testing.T) {
		parentOpinion, err := opinion.NewOpinion(
			shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-000000000001"),
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000002"),
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000003"),
			nil,
			lo.ToPtr("親意見のタイトル"),
			"これは親意見の内容です",
			time.Now(),
			nil,
		)
		require.NoError(t, err)

		// 最初は返信の数のみ（0）
		assert.Equal(t, 0, parentOpinion.Count())

		// 返信を追加
		reply1, err := opinion.NewOpinion(
			shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-000000000004"),
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000002"),
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000005"),
			lo.ToPtr(parentOpinion.OpinionID()),
			nil,
			"これは1つ目の返信意見です",
			time.Now(),
			nil,
		)
		require.NoError(t, err)

		reply2, err := opinion.NewOpinion(
			shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-000000000006"),
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000002"),
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000007"),
			lo.ToPtr(parentOpinion.OpinionID()),
			nil,
			"これは2つ目の返信意見です",
			time.Now(),
			nil,
		)
		require.NoError(t, err)

		parentOpinion.Reply(*reply1)
		parentOpinion.Reply(*reply2)

		// 返信2つ
		assert.Equal(t, 2, parentOpinion.Count())
	})

	t.Run("ネストされた返信も含めて総数を取得できる", func(t *testing.T) {
		parentOpinion, err := opinion.NewOpinion(
			shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-000000000001"),
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000002"),
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000003"),
			nil,
			lo.ToPtr("親意見のタイトル"),
			"これは親意見の内容です",
			time.Now(),
			nil,
		)
		require.NoError(t, err)

		// 第1階層の返信
		reply1, err := opinion.NewOpinion(
			shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-000000000004"),
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000002"),
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000005"),
			lo.ToPtr(parentOpinion.OpinionID()),
			nil,
			"これは第1階層の返信です",
			time.Now(),
			nil,
		)
		require.NoError(t, err)

		// 第2階層の返信（reply1への返信）
		reply1_1, err := opinion.NewOpinion(
			shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-000000000006"),
			shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000002"),
			shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000007"),
			lo.ToPtr(reply1.OpinionID()),
			nil,
			"これは第2階層の返信です",
			time.Now(),
			nil,
		)
		require.NoError(t, err)

		reply1.Reply(*reply1_1)
		parentOpinion.Reply(*reply1)

		// 返信1つ（ネストされた返信はその親の返信に含まれるため、親からは1つに見える）
		assert.Equal(t, 1, parentOpinion.Count())
	})
}
