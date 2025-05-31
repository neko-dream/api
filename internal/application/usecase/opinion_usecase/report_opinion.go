package opinion_usecase

import (
	"context"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type ReportOpinion interface {
	Execute(context.Context, ReportOpinionInput) error
}

type ReportOpinionInput struct {
	ReporterID shared.UUID[user.User]
	OpinionID  shared.UUID[opinion.Opinion]
	Reason     int32
	ReasonText *string
}

type reportOpinionInteractor struct {
	opinionRep opinion.OpinionRepository
	reportRep  opinion.ReportRepository
}

func NewReportOpinion(
	opinionRep opinion.OpinionRepository,
	reportRep opinion.ReportRepository,
) ReportOpinion {
	return &reportOpinionInteractor{
		opinionRep: opinionRep,
		reportRep:  reportRep,
	}
}

func (r *reportOpinionInteractor) Execute(ctx context.Context, input ReportOpinionInput) error {
	ctx, span := otel.Tracer("opinion_command").Start(ctx, "reportOpinionInteractor.Execute")
	defer span.End()

	opinion, err := r.opinionRep.FindByID(ctx, input.OpinionID)
	if err != nil {
		utils.HandleError(ctx, err, "opinionRep.FindByID")
		return messages.OpinionNotFound
	}

	report, err := opinion.Report(ctx, input.ReporterID, int(input.Reason), input.ReasonText)
	if err != nil {
		utils.HandleError(ctx, err, "opinion.Report")
		return messages.OpinionReportFailed
	}

	if err := r.reportRep.Create(ctx, *report); err != nil {
		utils.HandleError(ctx, err, "reportRep.Create")
		return messages.OpinionReportFailed
	}

	return nil
}
