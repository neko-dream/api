package report_query

import (
	"context"

	"github.com/neko-dream/server/internal/application/query/dto"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
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
