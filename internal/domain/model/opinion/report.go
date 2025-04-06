package opinion

import (
	"context"
	"time"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type ReportRepository interface {
	Create(context.Context, Report) error
	UpdateStatus(context.Context, shared.UUID[Report], Status) error
	FindByOpinionID(context.Context, shared.UUID[Opinion]) ([]Report, error)
}

type Report struct {
	OpinionReportID shared.UUID[Report]
	OpinionID       shared.UUID[Opinion]
	TalkSessionID   shared.UUID[talksession.TalkSession]
	ReporterID      shared.UUID[user.User]
	Reason          Reason
	ReasonText      *string
	Status          Status
	CreatedAt       time.Time
}

func NewReport(
	opinionReportID shared.UUID[Report],
	opinionID shared.UUID[Opinion],
	talkSessionID shared.UUID[talksession.TalkSession],
	reporterID shared.UUID[user.User],
	reason int,
	reasonText *string,
	status string,
	createdAt time.Time,
) (*Report, error) {
	statusType, err := NewStatus(status)
	if err != nil {
		return nil, err
	}

	return &Report{
		OpinionReportID: opinionReportID,
		OpinionID:       opinionID,
		TalkSessionID:   talkSessionID,
		ReporterID:      reporterID,
		Reason:          Reason(reason),
		ReasonText:      reasonText,
		Status:          statusType,
		CreatedAt:       createdAt,
	}, nil
}
