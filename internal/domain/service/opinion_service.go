package service

import (
	"context"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/domain/model/vote"
)

type opinionService struct {
	opinionRepo opinion.OpinionRepository
	voteRepo    vote.VoteRepository
}

func NewOpinionService(
	opinionRepo opinion.OpinionRepository,
	voteRepo vote.VoteRepository,
) opinion.OpinionService {
	return &opinionService{
		opinionRepo: opinionRepo,
		voteRepo:    voteRepo,
	}
}

// IsVotedOrReplied implements opinion.OpinionService.
func (o *opinionService) IsVoted(ctx context.Context, opinionID shared.UUID[opinion.Opinion], userID shared.UUID[user.User]) (bool, error) {
	v, err := o.voteRepo.FindByOpinionAndUserID(ctx, opinionID, userID)
	if err != nil {
		return true, messages.VoteFailed
	}
	if v != nil {
		return true, nil
	}

	return false, nil
}
