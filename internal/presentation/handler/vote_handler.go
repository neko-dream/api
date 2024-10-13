package handler

import (
	"context"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/presentation/oas"
	vote_usecase "github.com/neko-dream/server/internal/usecase/vote"
)

type voteHandler struct {
	postVoteUseCase vote_usecase.PostVoteUseCase
}

func NewVoteHandler(
	postVoteUseCase vote_usecase.PostVoteUseCase,
) oas.VoteHandler {
	return &voteHandler{
		postVoteUseCase: postVoteUseCase,
	}
}

// Vote implements oas.VoteHandler.
func (v *voteHandler) Vote(ctx context.Context, req oas.OptVoteReq, params oas.VoteParams) (oas.VoteRes, error) {
	claim := session.GetSession(ctx)
	userID, err := claim.UserID()
	if err != nil {
		return nil, messages.ForbiddenError
	}
	if !req.IsSet() {
		return nil, messages.RequiredParameterError
	}

	value := req.Value
	_, err = v.postVoteUseCase.Execute(ctx, vote_usecase.PostVoteInput{
		TargetOpinionID: shared.MustParseUUID[opinion.Opinion](params.OpinionID),
		UserID:          userID,
		VoteType:        string(value.VoteStatus.Value),
	})
	if err != nil {
		return nil, err
	}

	res := &oas.VoteOKApplicationJSON{}
	return res, nil
}
