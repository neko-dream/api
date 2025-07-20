package vote_test

import (
	"testing"

	"github.com/neko-dream/server/internal/domain/model/vote"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestVoteType_Int(t *testing.T) {
	tests := []struct {
		name     string
		voteType vote.VoteType
		want     int
	}{
		{
			name:     "UnVotedは0を返す",
			voteType: vote.UnVoted,
			want:     0,
		},
		{
			name:     "Agreeは1を返す",
			voteType: vote.Agree,
			want:     1,
		},
		{
			name:     "Disagreeは2を返す",
			voteType: vote.Disagree,
			want:     2,
		},
		{
			name:     "Passは3を返す",
			voteType: vote.Pass,
			want:     3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.voteType.Int())
		})
	}
}

func TestVoteType_String(t *testing.T) {
	tests := []struct {
		name     string
		voteType vote.VoteType
		want     string
	}{
		{
			name:     "Agreeはagreeを返す",
			voteType: vote.Agree,
			want:     "agree",
		},
		{
			name:     "Disagreeはdisagreeを返す",
			voteType: vote.Disagree,
			want:     "disagree",
		},
		{
			name:     "Passはpassを返す",
			voteType: vote.Pass,
			want:     "pass",
		},
		{
			name:     "UnVotedは空文字を返す",
			voteType: vote.UnVoted,
			want:     "",
		},
		{
			name:     "不正な値は空文字を返す",
			voteType: vote.VoteType(99),
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.voteType.String())
		})
	}
}

func TestVoteTypeFromInt(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  vote.VoteType
	}{
		{
			name:  "1はAgreeを返す",
			input: 1,
			want:  vote.Agree,
		},
		{
			name:  "2はDisagreeを返す",
			input: 2,
			want:  vote.Disagree,
		},
		{
			name:  "3はPassを返す",
			input: 3,
			want:  vote.Pass,
		},
		{
			name:  "0はUnVotedを返す",
			input: 0,
			want:  vote.UnVoted,
		},
		{
			name:  "不正な値はUnVotedを返す",
			input: 99,
			want:  vote.UnVoted,
		},
		{
			name:  "負の値はUnVotedを返す",
			input: -1,
			want:  vote.UnVoted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, vote.VoteTypeFromInt(tt.input))
		})
	}
}

func TestVoteFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   *string
		want    *vote.VoteType
		wantErr bool
	}{
		{
			name:    "agreeはAgreeを返す",
			input:   lo.ToPtr("agree"),
			want:    lo.ToPtr(vote.Agree),
			wantErr: false,
		},
		{
			name:    "disagreeはDisagreeを返す",
			input:   lo.ToPtr("disagree"),
			want:    lo.ToPtr(vote.Disagree),
			wantErr: false,
		},
		{
			name:    "passはPassを返す",
			input:   lo.ToPtr("pass"),
			want:    lo.ToPtr(vote.Pass),
			wantErr: false,
		},
		{
			name:    "nilはnilを返す",
			input:   nil,
			want:    nil,
			wantErr: false,
		},
		{
			name:    "不正な文字列はエラーを返す",
			input:   lo.ToPtr("invalid"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "空文字はエラーを返す",
			input:   lo.ToPtr(""),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "大文字はエラーを返す",
			input:   lo.ToPtr("AGREE"),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vote.VoteFromString(tt.input)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}
			
			assert.NoError(t, err)
			if tt.want == nil {
				assert.Nil(t, got)
			} else {
				assert.NotNil(t, got)
				assert.Equal(t, *tt.want, *got)
			}
		})
	}
}