package vote_command

import (
	"context"
	"database/sql"
	"errors"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/analysis"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/domain/model/vote"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
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
		talksession.TalkSessionRepository
		vote.VoteRepository
		analysis.AnalysisService
		analysis.AnalysisRepository
		service.TalkSessionAccessControl
		*db.DBManager
	}
)

func NewVoteHandler(
	opinionService opinion.OpinionService,
	opinionRepository opinion.OpinionRepository,
	talkSessionRepository talksession.TalkSessionRepository,
	voteRepository vote.VoteRepository,
	analysisService analysis.AnalysisService,
	analysisRepository analysis.AnalysisRepository,
	talkSessionAccessControl service.TalkSessionAccessControl,
	DBManager *db.DBManager,
) Vote {
	return &voteHandler{
		OpinionService:           opinionService,
		OpinionRepository:        opinionRepository,
		VoteRepository:           voteRepository,
		TalkSessionRepository:    talkSessionRepository,
		AnalysisService:          analysisService,
		AnalysisRepository:       analysisRepository,
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
		return messages.OpinionNotFound
	}

	// セッションを探す
	session, err := i.TalkSessionRepository.FindByID(ctx, op.TalkSessionID())
	if err != nil {
		utils.HandleError(ctx, err, "TalkSessionRepository.FindByID")
		return messages.TalkSessionNotFound
	}

	// 終了していればエラー
	if session.IsFinished(ctx) {
		return messages.TalkSessionIsFinished
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

		return nil
	}); err != nil {
		return errtrace.Wrap(err)
	}

	// 非同期で分析を開始
	i.StartAnalysisIfNeeded(ctx, op.TalkSessionID())

	return nil
}

func (i *voteHandler) StartAnalysisIfNeeded(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) {
	ctx, span := otel.Tracer("vote_command").Start(ctx, "voteHandler.StartAnalysisIfNeeded")
	defer span.End()

	bg := context.Background()
	span = trace.SpanFromContext(ctx)
	bg = trace.ContextWithSpan(bg, span)
	go func() {
		_ = i.AnalysisService.StartAnalysis(bg, talkSessionID)

		// 分析レポートを取得
		report, err := i.AnalysisRepository.FindByTalkSessionID(bg, talkSessionID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			utils.HandleError(ctx, err, "AnalysisRepository.FindByTalkSessionID")
			return
		}
		if report.ShouldReGenerateReport() {
			_ = i.AnalysisService.GenerateReport(bg, talkSessionID)
		}
	}()
}
