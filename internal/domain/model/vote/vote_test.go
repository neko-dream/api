package vote_test

import (
	"testing"
	"time"

	"github.com/neko-dream/api/internal/domain/model/opinion"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/domain/model/vote"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewVote(t *testing.T) {
	tests := []struct {
		name          string
		voteID        shared.UUID[vote.Vote]
		opinionID     shared.UUID[opinion.Opinion]
		talkSessionID shared.UUID[talksession.TalkSession]
		userID        shared.UUID[user.User]
		voteType      vote.VoteType
		createdAt     time.Time
		wantErr       bool
	}{
		{
			name:          "正常に賛成票を作成できる",
			voteID:        shared.MustParseUUID[vote.Vote]("00000000-0000-0000-0000-000000000001"),
			opinionID:     shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-000000000002"),
			talkSessionID: shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000003"),
			userID:        shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000004"),
			voteType:      vote.Agree,
			createdAt:     time.Now(),
			wantErr:       false,
		},
		{
			name:          "正常に反対票を作成できる",
			voteID:        shared.MustParseUUID[vote.Vote]("00000000-0000-0000-0000-000000000005"),
			opinionID:     shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-000000000006"),
			talkSessionID: shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000007"),
			userID:        shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000008"),
			voteType:      vote.Disagree,
			createdAt:     time.Now(),
			wantErr:       false,
		},
		{
			name:          "正常にパス票を作成できる",
			voteID:        shared.MustParseUUID[vote.Vote]("00000000-0000-0000-0000-000000000009"),
			opinionID:     shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-00000000000a"),
			talkSessionID: shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-00000000000b"),
			userID:        shared.MustParseUUID[user.User]("00000000-0000-0000-0000-00000000000c"),
			voteType:      vote.Pass,
			createdAt:     time.Now(),
			wantErr:       false,
		},
		{
			name:          "未投票状態でも作成できる",
			voteID:        shared.MustParseUUID[vote.Vote]("00000000-0000-0000-0000-00000000000d"),
			opinionID:     shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-00000000000e"),
			talkSessionID: shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-00000000000f"),
			userID:        shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000010"),
			voteType:      vote.UnVoted,
			createdAt:     time.Now(),
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vote.NewVote(
				tt.voteID,
				tt.opinionID,
				tt.talkSessionID,
				tt.userID,
				tt.voteType,
				tt.createdAt,
			)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.voteID, got.VoteID)
			assert.Equal(t, tt.opinionID, got.OpinionID)
			assert.Equal(t, tt.talkSessionID, got.TalkSessionID)
			assert.Equal(t, tt.userID, got.UserID)
			assert.Equal(t, tt.voteType, got.VoteType)
			assert.Equal(t, tt.createdAt, got.CreatedAt)
		})
	}
}

func TestVote_ChangeVoteType(t *testing.T) {
	tests := []struct {
		name         string
		initialType  vote.VoteType
		newType      vote.VoteType
		expectedType vote.VoteType
	}{
		{
			name:         "賛成から反対に変更できる",
			initialType:  vote.Agree,
			newType:      vote.Disagree,
			expectedType: vote.Disagree,
		},
		{
			name:         "反対からパスに変更できる",
			initialType:  vote.Disagree,
			newType:      vote.Pass,
			expectedType: vote.Pass,
		},
		{
			name:         "パスから賛成に変更できる",
			initialType:  vote.Pass,
			newType:      vote.Agree,
			expectedType: vote.Agree,
		},
		{
			name:         "未投票から賛成に変更できる",
			initialType:  vote.UnVoted,
			newType:      vote.Agree,
			expectedType: vote.Agree,
		},
		{
			name:         "同じ投票タイプに変更しても問題ない",
			initialType:  vote.Agree,
			newType:      vote.Agree,
			expectedType: vote.Agree,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := vote.NewVote(
				shared.MustParseUUID[vote.Vote]("00000000-0000-0000-0000-000000000001"),
				shared.MustParseUUID[opinion.Opinion]("00000000-0000-0000-0000-000000000002"),
				shared.MustParseUUID[talksession.TalkSession]("00000000-0000-0000-0000-000000000003"),
				shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000004"),
				tt.initialType,
				time.Now(),
			)
			require.NoError(t, err)

			v.ChangeVoteType(tt.newType)

			assert.Equal(t, tt.expectedType, v.VoteType)
		})
	}
}
