package vote

import (
	"context"
	"time"

	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type (
	VoteRepository interface {
		Create(ctx context.Context, vote Vote) error
		FindByOpinionAndUserID(ctx context.Context, opinionID shared.UUID[opinion.Opinion], userID shared.UUID[user.User]) (*Vote, error)
	}

	Vote struct {
		VoteID    shared.UUID[Vote]
		OpinionID shared.UUID[opinion.Opinion]
		UserID    shared.UUID[user.User]
		VoteType  VoteType
		CreatedAt time.Time
	}
)

func NewVote(
	voteID shared.UUID[Vote],
	parentOpinionID shared.UUID[opinion.Opinion],
	userID shared.UUID[user.User],
	VoteType VoteType,
	createdAt time.Time,
) (*Vote, error) {
	return &Vote{
		VoteID:    voteID,
		OpinionID: parentOpinionID,
		UserID:    userID,
		VoteType:  VoteType,
		CreatedAt: createdAt,
	}, nil
}
