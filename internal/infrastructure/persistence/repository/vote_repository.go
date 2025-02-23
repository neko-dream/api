package repository

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/domain/model/vote"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"go.opentelemetry.io/otel"
)

type voteRepository struct {
	*db.DBManager
}

func NewVoteRepository(dbManager *db.DBManager) vote.VoteRepository {
	return &voteRepository{dbManager}
}

func (o *voteRepository) Create(ctx context.Context, vote vote.Vote) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "voteRepository.Create")
	defer span.End()

	if err := o.GetQueries(ctx).CreateVote(ctx, model.CreateVoteParams{
		VoteID:        vote.VoteID.UUID(),
		OpinionID:     vote.OpinionID.UUID(),
		TalkSessionID: vote.TalkSessionID.UUID(),
		UserID:        vote.UserID.UUID(),
		VoteType:      int16(vote.VoteType.Int()),
		CreatedAt:     vote.CreatedAt,
	}); err != nil {
		return err
	}

	return nil
}

func (o *voteRepository) FindByOpinionAndUserID(ctx context.Context, opinionID shared.UUID[opinion.Opinion], userID shared.UUID[user.User]) (*vote.Vote, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "voteRepository.FindByOpinionAndUserID")
	defer span.End()

	voteRow, err := o.GetQueries(ctx).FindVoteByUserIDAndOpinionID(ctx, model.FindVoteByUserIDAndOpinionIDParams{
		OpinionID: opinionID.UUID(),
		UserID:    userID.UUID(),
	})
	if err != nil {
		return nil, err
	}

	voteID, err := shared.ParseUUID[vote.Vote](voteRow.VoteID.String())
	if err != nil {
		return nil, err
	}
	opID, err := shared.ParseUUID[opinion.Opinion](voteRow.OpinionID.String())
	if err != nil {
		return nil, err
	}
	tsID, err := shared.ParseUUID[talksession.TalkSession](voteRow.TalkSessionID.String())
	if err != nil {
		return nil, err
	}
	uID, err := shared.ParseUUID[user.User](voteRow.UserID.String())
	if err != nil {
		return nil, err
	}

	vote, err := vote.NewVote(
		voteID,
		opID,
		tsID,
		uID,
		vote.VoteTypeFromInt(int(voteRow.VoteType)),
		voteRow.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return vote, nil
}
