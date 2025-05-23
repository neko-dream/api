package vote

import (
	"context"
	"time"

	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type (
	VoteRepository interface {
		Create(ctx context.Context, vote Vote) error
		Update(ctx context.Context, vote Vote) error
		FindByOpinionAndUserID(ctx context.Context, opinionID shared.UUID[opinion.Opinion], userID shared.UUID[user.User]) (*Vote, error)
	}

	Vote struct {
		VoteID        shared.UUID[Vote]
		OpinionID     shared.UUID[opinion.Opinion]
		TalkSessionID shared.UUID[talksession.TalkSession]
		UserID        shared.UUID[user.User]
		VoteType      VoteType
		CreatedAt     time.Time
	}
)

func (v *Vote) ChangeVoteType(voteType VoteType) {
	v.VoteType = voteType
}

func NewVote(
	voteID shared.UUID[Vote],
	parentOpinionID shared.UUID[opinion.Opinion],
	talkSessionID shared.UUID[talksession.TalkSession],
	userID shared.UUID[user.User],
	VoteType VoteType,
	createdAt time.Time,
) (*Vote, error) {
	return &Vote{
		VoteID:        voteID,
		OpinionID:     parentOpinionID,
		TalkSessionID: talkSessionID,
		UserID:        userID,
		VoteType:      VoteType,
		CreatedAt:     createdAt,
	}, nil
}
