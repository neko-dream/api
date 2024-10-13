package handler

import (
	"context"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/presentation/oas"
	vote_usecase "github.com/neko-dream/server/internal/usecase/vote"
	"github.com/neko-dream/server/pkg/utils"
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
		TalkSessionID:   shared.MustParseUUID[talksession.TalkSession](params.TalkSessionID),
		TargetOpinionID: shared.MustParseUUID[opinion.Opinion](params.OpinionID),
		UserID:          userID,
		VoteType:        string(value.VoteStatus.Value),
	})
	if err != nil {
		utils.HandleError(ctx, err, "postVoteUseCase.Execute")
		return nil, err
	}

	res := &oas.VoteOKApplicationJSON{}
	return res, nil
}
