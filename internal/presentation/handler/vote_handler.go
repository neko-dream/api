package handler

import (
	"context"

	"github.com/neko-dream/server/internal/application/usecase/vote_usecase"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type voteHandler struct {
	voteCommand vote_usecase.Vote
	authService service.AuthenticationService
}

func NewVoteHandler(
	voteCommand vote_usecase.Vote,
	authService service.AuthenticationService,
) oas.VoteHandler {
	return &voteHandler{
		voteCommand: voteCommand,
		authService: authService,
	}
}

// Vote2 implements oas.VoteHandler.
func (v *voteHandler) Vote2(ctx context.Context, req *oas.Vote2Req, params oas.Vote2Params) (oas.Vote2Res, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "voteHandler.Vote")
	defer span.End()

	authCtx, err := requireAuthentication(v.authService, ctx)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, messages.RequiredParameterError
	}

	value := req
	targetOpinionID, err := shared.ParseUUID[opinion.Opinion](params.OpinionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	err = v.voteCommand.Execute(ctx, vote_usecase.VoteInput{
		TargetOpinionID: targetOpinionID,
		UserID:          authCtx.UserID,
		VoteType:        string(value.VoteStatus),
	})
	if err != nil {
		utils.HandleError(ctx, err, "postVoteUseCase.Execute")
		return nil, err
	}

	res := &oas.Vote2OKApplicationJSON{}
	return res, nil
}
