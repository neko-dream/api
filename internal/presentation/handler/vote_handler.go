package handler

import (
	"context"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/neko-dream/server/internal/usecase/command/vote_command"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type voteHandler struct {
	voteCommand vote_command.Vote
}

func NewVoteHandler(
	voteCommand vote_command.Vote,
) oas.VoteHandler {
	return &voteHandler{
		voteCommand: voteCommand,
	}
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
	err = v.voteCommand.Execute(ctx, vote_command.VoteInput{
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
