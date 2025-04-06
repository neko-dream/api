package report_command

import (
	"context"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type SolveReportCommand interface {
	Execute(ctx context.Context, input SolveReportInput) error
}

type SolveReportInput struct {
	OpinionID shared.UUID[opinion.Opinion]
	UserID    shared.UUID[user.User]
	Status    opinion.Status
}

type solveReportCommandInteractor struct {
	reportRep      opinion.ReportRepository
	opinionRep     opinion.OpinionRepository
	talkSessionRep talksession.TalkSessionRepository
	*db.DBManager
}

func NewSolveReportCommandInteractor(
	reportRepository opinion.ReportRepository,
	opinionRepository opinion.OpinionRepository,
	talkSessionRepository talksession.TalkSessionRepository,
	dbManager *db.DBManager,
) SolveReportCommand {
	return &solveReportCommandInteractor{
		reportRep:      reportRepository,
		opinionRep:     opinionRepository,
		talkSessionRep: talkSessionRepository,
		DBManager:      dbManager,
	}
}

// Execute implements SolveReportCommand.
func (s *solveReportCommandInteractor) Execute(ctx context.Context, input SolveReportInput) error {
	ctx, span := otel.Tracer("report_command").Start(ctx, "solveReportCommandInteractor.Execute")
	defer span.End()

	// opinion取得
	opr, err := s.opinionRep.FindByID(ctx, input.OpinionID)
	if err != nil {
		return err
	}

	// talkSession取得
	talkSession, err := s.talkSessionRep.FindByID(ctx, opr.TalkSessionID())
	if err != nil {
		return err
	}

	// talkSessionのオーナーと一致するか確認
	if talkSession.OwnerUserID() != input.UserID {
		return messages.TalkSessionNotFound
	}

	if err := s.DBManager.ExecTx(ctx, func(ctx context.Context) error {
		// report一覧取得
		reports, err := s.reportRep.FindByOpinionID(ctx, input.OpinionID)
		if err != nil {
			return err
		}

		// reportをいずれかの状態に更新
		for _, report := range reports {
			if err := s.reportRep.UpdateStatus(ctx, report.OpinionReportID, input.Status); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		utils.HandleError(ctx, err, "SolveReportCommandInteractor.Execute")
		return err
	}

	return nil
}
