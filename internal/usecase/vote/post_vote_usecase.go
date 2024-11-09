package vote_usecase

import (
	"context"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/analysis"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/domain/model/vote"
	"github.com/neko-dream/server/internal/infrastructure/db"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
)

type (
	PostVoteUseCase interface {
		Execute(context.Context, PostVoteInput) (*PostVoteOutput, error)
	}

	PostVoteInput struct {
		TalkSessionID   shared.UUID[talksession.TalkSession]
		TargetOpinionID shared.UUID[opinion.Opinion]
		UserID          shared.UUID[user.User]
		VoteType        string
	}

	PostVoteOutput struct {
	}

	postVoteInteractor struct {
		opinion.OpinionService
		vote.VoteRepository
		analysis.AnalysisService
		*db.DBManager
	}
)

func NewPostVoteUseCase(
	opinionService opinion.OpinionService,
	voteRepository vote.VoteRepository,
	analysisService analysis.AnalysisService,
	DBManager *db.DBManager,
) PostVoteUseCase {
	return &postVoteInteractor{
		OpinionService:  opinionService,
		VoteRepository:  voteRepository,
		AnalysisService: analysisService,
		DBManager:       DBManager,
	}
}

func (i *postVoteInteractor) Execute(ctx context.Context, input PostVoteInput) (*PostVoteOutput, error) {
	output := PostVoteOutput{}
	// Opinionに対して投票を行っているか確認
	voted, err := i.OpinionService.IsVoted(ctx, input.TargetOpinionID, input.UserID)
	if err != nil {
		utils.HandleError(ctx, err, "IsVoted")
		return nil, errtrace.Wrap(err)
	}
	// 投票を行っている場合、エラーを返す
	if voted {
		return nil, messages.OpinionAlreadyVoted
	}

	if err := i.ExecTx(ctx, func(ctx context.Context) error {
		// 投票を行っていない場合、投票を行う
		vote, err := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			input.TargetOpinionID,
			input.TalkSessionID,
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

		if err := i.AnalysisService.StartAnalysis(ctx, input.TalkSessionID); err != nil {
			utils.HandleError(ctx, err, "StartAnalysis")
			return err
		}

		return nil
	}); err != nil {
		return nil, errtrace.Wrap(err)
	}

	return &output, nil
}
