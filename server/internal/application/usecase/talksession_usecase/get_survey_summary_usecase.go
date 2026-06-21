package talksession_usecase

import (
	"context"
	"database/sql"
	"errors"

	"braces.dev/errtrace"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/talksession/survey"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

type (
	GetSurveySummaryUseCase interface {
		Execute(ctx context.Context, input GetSurveySummaryInput) (*survey.Summary, error)
	}

	GetSurveySummaryInput struct {
		TalkSessionID shared.UUID[talksession.TalkSession]
		ActorUserID   shared.UUID[user.User]
	}

	getSurveySummaryUseCase struct {
		talkSessionRepo talksession.TalkSessionRepository
		surveyRepo      survey.Repository
	}
)

func NewGetSurveySummaryUseCase(
	talkSessionRepo talksession.TalkSessionRepository,
	surveyRepo survey.Repository,
) GetSurveySummaryUseCase {
	return &getSurveySummaryUseCase{
		talkSessionRepo: talkSessionRepo,
		surveyRepo:      surveyRepo,
	}
}

func (uc *getSurveySummaryUseCase) Execute(ctx context.Context, input GetSurveySummaryInput) (*survey.Summary, error) {
	ctx, span := otel.Tracer("talksession_query").Start(ctx, "getSurveySummaryUseCase.Execute")
	defer span.End()

	ts, err := uc.talkSessionRepo.FindByID(ctx, input.TalkSessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, messages.TalkSessionNotFound
		}
		utils.HandleError(ctx, err, "TalkSessionRepository.FindByID")
		return nil, errtrace.Wrap(err)
	}
	if ts == nil {
		return nil, messages.TalkSessionNotFound
	}
	if ts.OwnerUserID() != input.ActorUserID {
		return nil, messages.ForbiddenError
	}

	s, err := uc.surveyRepo.FindBySessionID(ctx, input.TalkSessionID)
	if err != nil {
		utils.HandleError(ctx, err, "surveyRepo.FindBySessionID")
		return nil, errtrace.Wrap(err)
	}
	if s == nil {
		return nil, messages.TalkSessionSurveyNotFound
	}

	subs, err := uc.surveyRepo.ListAllSubmissions(ctx, s.SurveyID())
	if err != nil {
		utils.HandleError(ctx, err, "surveyRepo.ListAllSubmissions")
		return nil, errtrace.Wrap(err)
	}

	return survey.Aggregate(s, subs), nil
}
