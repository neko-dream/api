package vote_command

import (
	"context"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/analysis"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/domain/model/vote"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type (
	Vote interface {
		Execute(context.Context, VoteInput) error
	}

	VoteInput struct {
		TargetOpinionID shared.UUID[opinion.Opinion]
		UserID          shared.UUID[user.User]
		VoteType        string
	}

	voteHandler struct {
		opinion.OpinionService
		opinion.OpinionRepository
		vote.VoteRepository
		analysis.AnalysisService
		service.TalkSessionAccessControl
		*db.DBManager
	}
)

func NewVoteHandler(
	opinionService opinion.OpinionService,
	opinionRepository opinion.OpinionRepository,
	voteRepository vote.VoteRepository,
	analysisService analysis.AnalysisService,
	talkSessionAccessControl service.TalkSessionAccessControl,
	DBManager *db.DBManager,
) Vote {
	return &voteHandler{
		OpinionService:           opinionService,
		OpinionRepository:        opinionRepository,
		VoteRepository:           voteRepository,
		AnalysisService:          analysisService,
		TalkSessionAccessControl: talkSessionAccessControl,
		DBManager:                DBManager,
	}
}

func (i *voteHandler) Execute(ctx context.Context, input VoteInput) error {
	ctx, span := otel.Tracer("vote_command").Start(ctx, "voteHandler.Execute")
	defer span.End()

	// opinionを探す
	op, err := i.OpinionRepository.FindByID(ctx, input.TargetOpinionID)
	if err != nil {
		utils.HandleError(ctx, err, "OpinionRepository.FindByID")
		return err
	}

	// 参加制限を満たしているか確認。満たしていない場合はエラーを返す
	if _, err := i.TalkSessionAccessControl.CanUserJoin(ctx, op.TalkSessionID(), lo.ToPtr(input.UserID)); err != nil {
		utils.HandleError(ctx, err, "TalkSessionAccessControl.CanUserJoin")
		return err
	}

	// Opinionに対して投票を行っているか確認
	voted, err := i.OpinionService.IsVoted(ctx, input.TargetOpinionID, input.UserID)
	if err != nil {
		utils.HandleError(ctx, err, "IsVoted")
		return errtrace.Wrap(err)
	}
	// 投票を行っている場合、エラーを返す
	if voted {
		return messages.OpinionAlreadyVoted
	}

	if err := i.ExecTx(ctx, func(ctx context.Context) error {
		// 投票を行っていない場合、投票を行う
		vote, err := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			input.TargetOpinionID,
			op.TalkSessionID(),
			input.UserID,
			vote.VoteFromString(lo.ToPtr(input.VoteType)),
			clock.Now(ctx),
		)
		if err != nil {
			utils.HandleError(ctx, err, "NewVote")
			return err
		}

		if err := i.VoteRepository.Create(ctx, *vote); err != nil {
			return messages.VoteFailed
		}

		if err := i.AnalysisService.StartAnalysis(ctx, op.TalkSessionID()); err != nil {
			utils.HandleError(ctx, err, "StartAnalysis")
			return err
		}

		return nil
	}); err != nil {
		return errtrace.Wrap(err)
	}

	return nil
}
