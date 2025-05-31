package handler

import (
	"context"

	"github.com/neko-dream/server/internal/application/usecase/vote_usecase"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type voteHandler struct {
	voteCommand vote_usecase.Vote
}

func NewVoteHandler(
	voteCommand vote_usecase.Vote,
) oas.VoteHandler {
	return &voteHandler{
		voteCommand: voteCommand,
	}
}

// Vote2 implements oas.VoteHandler.
func (v *voteHandler) Vote2(ctx context.Context, req oas.OptVote2Req, params oas.Vote2Params) (oas.Vote2Res, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "voteHandler.Vote")
	defer span.End()

	claim := session.GetSession(ctx)
	userID, err := claim.UserID()
	if err != nil {
		return nil, messages.ForbiddenError
	}
	if !req.IsSet() {
		return nil, messages.RequiredParameterError
	}

	value := req.Value
	targetOpinionID, err := shared.ParseUUID[opinion.Opinion](params.OpinionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	err = v.voteCommand.Execute(ctx, vote_usecase.VoteInput{
		TargetOpinionID: targetOpinionID,
		UserID:          userID,
		VoteType:        string(value.VoteStatus.Value),
	})
	if err != nil {
		utils.HandleError(ctx, err, "postVoteUseCase.Execute")
		return nil, err
	}

	res := &oas.Vote2OKApplicationJSON{}
	return res, nil
}

// Vote implements oas.VoteHandler.
func (v *voteHandler) Vote(ctx context.Context, req oas.OptVoteReq, params oas.VoteParams) (oas.VoteRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "voteHandler.Vote")
	defer span.End()

	claim := session.GetSession(ctx)
	userID, err := claim.UserID()
	if err != nil {
		return nil, messages.ForbiddenError
	}
	if !req.IsSet() {
		return nil, messages.RequiredParameterError
	}

	value := req.Value
	targetOpinionID, err := shared.ParseUUID[opinion.Opinion](params.OpinionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	err = v.voteCommand.Execute(ctx, vote_usecase.VoteInput{
		TargetOpinionID: targetOpinionID,
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
