package analysis_usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/analysis"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"go.opentelemetry.io/otel"
)

// ApplyFeedbackUseCase レポートにフィードバックを行う
type ApplyFeedbackUseCase interface {
	Execute(ctx context.Context, input ApplyFeedbackInput) error
}

type ApplyFeedbackInput struct {
	UserID           shared.UUID[user.User]
	AnalysisReportID shared.UUID[analysis.AnalysisReport]
	FeedbackType     analysis.FeedbackType
}

func NewApplyFeedbackInteractor(
	analysisRepository analysis.AnalysisRepository,
	dbm *db.DBManager,
) ApplyFeedbackUseCase {
	return &applyFeedbackInteractor{
		analysisRepository: analysisRepository,
		dbm:                dbm,
	}
}

type applyFeedbackInteractor struct {
	analysisRepository analysis.AnalysisRepository
	dbm                *db.DBManager
}

func (a *applyFeedbackInteractor) Execute(ctx context.Context, input ApplyFeedbackInput) error {
	ctx, span := otel.Tracer("analysis_usecase").Start(ctx, "applyFeedbackInteractor.Execute")
	defer span.End()

	if err := a.dbm.ExecTx(ctx, func(ctx context.Context) error {
		report, err := a.analysisRepository.FindByID(ctx, input.AnalysisReportID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return messages.AnalysisReportNotFound
			}
			return err
		}

		if report.HasReceivedFeedbackFrom(ctx, input.UserID) {
			return messages.AnalysisReportAlreadyFeedbacked
		}

		report.ApplyFeedback(ctx, input.FeedbackType, input.UserID)

		if err := a.analysisRepository.SaveReport(ctx, report); err != nil {
			return messages.AnalysisReportFeedbackFailed
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
