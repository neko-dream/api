package repository

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/domain/model/vote"
	"github.com/neko-dream/server/internal/infrastructure/db"
	model "github.com/neko-dream/server/internal/infrastructure/db/sqlc"
)

type voteRepository struct {
	*db.DBManager
}

func NewVoteRepository(dbManager *db.DBManager) *voteRepository {
	return &voteRepository{dbManager}
}

func (o *voteRepository) FindByOpinionAndUserID(ctx context.Context, opinionID shared.UUID[opinion.Opinion], userID shared.UUID[user.User]) (*vote.Vote, error) {
	voteRow, err := o.GetQueries(ctx).GetVoteByUserIDAndOpinionID(ctx, model.GetVoteByUserIDAndOpinionIDParams{
		OpinionID: opinionID.UUID(),
		UserID:    userID.UUID(),
	})
	if err != nil {
		return nil, err
	}

	vote, err := vote.NewVote(
		shared.MustParseUUID[vote.Vote](voteRow.VoteID.String()),
		shared.MustParseUUID[opinion.Opinion](voteRow.OpinionID.String()),
		shared.MustParseUUID[user.User](voteRow.UserID.String()),
		opinion.VoteTypeFromInt(int(voteRow.VoteType)),
		voteRow.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return vote, nil
}
