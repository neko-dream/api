package opinion_query

import (
	"context"

	"github.com/neko-dream/api/internal/application/query/dto"
	"github.com/neko-dream/api/internal/domain/model/opinion"
	"go.opentelemetry.io/otel"
)

type GetReportReasons interface {
	Execute(ctx context.Context) ([]dto.ReportReason, error)
}

type getReportReasons struct {
}

func NewGetReportReasons() GetReportReasons {
	return &getReportReasons{}
}

// Execute implements GetReportReasons.
func (g *getReportReasons) Execute(ctx context.Context) ([]dto.ReportReason, error) {
	ctx, span := otel.Tracer("opinion_query").Start(ctx, "getReportReasons.Execute")
	defer span.End()

	_ = ctx

	var reasons []dto.ReportReason

	for _, reason := range opinion.ReasonValues() {
		reasons = append(reasons, dto.ReportReason{
			ReasonID: int(reason),
			Reason:   reason.StringJP(),
		})
	}

	return reasons, nil
}
