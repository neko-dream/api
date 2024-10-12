package vote

import (
	"context"
	"time"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type (
	VoteRepository interface {
		FindByOpinionAndUserID(ctx context.Context, opinionID shared.UUID[opinion.Opinion], userID shared.UUID[user.User]) (*Vote, error)
	}

	Vote struct {
		VoteID    shared.UUID[Vote]
		OpinionID shared.UUID[opinion.Opinion]
		UserID    shared.UUID[user.User]
		VoteType  opinion.VoteType
		CreatedAt time.Time
	}
)

func NewVote(
	voteID shared.UUID[Vote],
	opinionID shared.UUID[opinion.Opinion],
	userID shared.UUID[user.User],
	VoteType opinion.VoteType,
	createdAt time.Time,
) (*Vote, error) {
	if VoteType == opinion.UnVoted {
		return nil, messages.VoteUnvoteNotAllowed
	}

	return &Vote{
		VoteID:    voteID,
		OpinionID: opinionID,
		UserID:    userID,
		VoteType:  VoteType,
		CreatedAt: createdAt,
	}, nil
}
