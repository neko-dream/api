package report_query

import (
	"context"

	"github.com/neko-dream/api/internal/application/query/dto"
	"github.com/neko-dream/api/internal/domain/model/opinion"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
)

type GetOpinionReportQuery interface {
	Execute(ctx context.Context, input GetOpinionReportInput) (*GetOpinionReportOutput, error)
}

type GetOpinionReportInput struct {
	OpinionID shared.UUID[opinion.Opinion]
	UserID    shared.UUID[user.User]
}

type GetOpinionReportOutput struct {
	Report dto.ReportDetail
}
