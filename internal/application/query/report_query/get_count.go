package report_query

import (
	"context"

	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/opinion"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/user"
	"go.opentelemetry.io/otel"
)

type GetCountQuery interface {
	Execute(ctx context.Context, input GetCountInput) (*GetCountOutput, error)
}

type GetCountInput struct {
	TalkSessionID shared.UUID[talksession.TalkSession]
	UserID        shared.UUID[user.User]
	Status        string
}

type GetCountOutput struct {
	Count int
}

type getCountQueryInteractor struct {
	reportRep      opinion.ReportRepository
	talkSessionRep talksession.TalkSessionRepository
}

func NewGetCountQueryInteractor(
	reportRepository opinion.ReportRepository,
	talkSessionRepository talksession.TalkSessionRepository,
) GetCountQuery {
	return &getCountQueryInteractor{
		reportRep:      reportRepository,
		talkSessionRep: talkSessionRepository,
	}
}

// Execute implements GetCountQuery.
func (g *getCountQueryInteractor) Execute(ctx context.Context, input GetCountInput) (*GetCountOutput, error) {
	ctx, span := otel.Tracer("report_query").Start(ctx, "getCountQueryInteractor.Execute")
	defer span.End()

	// talkSession取得
	talkSession, err := g.talkSessionRep.FindByID(ctx, input.TalkSessionID)
	if err != nil {
		return nil, err
	}

	// talkSessionのオーナーと一致するか確認
	if talkSession.OwnerUserID() != input.UserID {
		return nil, messages.TalkSessionNotFound
	}

	count, err := g.reportRep.CountByTalkSessionIDAndStatus(ctx, input.TalkSessionID, opinion.Status(input.Status))
	if err != nil {
		return nil, err
	}

	return &GetCountOutput{
		Count: count,
	}, nil
}
