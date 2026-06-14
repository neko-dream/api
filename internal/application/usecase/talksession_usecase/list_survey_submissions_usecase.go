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
	ListSurveySubmissionsUseCase interface {
		Execute(ctx context.Context, input ListSurveySubmissionsInput) (*ListSurveySubmissionsOutput, error)
	}

	ListSurveySubmissionsInput struct {
		TalkSessionID shared.UUID[talksession.TalkSession]
		ActorUserID   shared.UUID[user.User]
		Limit         int
		Offset        int
	}

	ListSurveySubmissionsOutput struct {
		Submissions []*survey.Submission
		TotalCount  int32
	}

	listSurveySubmissionsUseCase struct {
		talkSessionRepo talksession.TalkSessionRepository
		surveyRepo      survey.Repository
	}
)

func NewListSurveySubmissionsUseCase(
	talkSessionRepo talksession.TalkSessionRepository,
	surveyRepo survey.Repository,
) ListSurveySubmissionsUseCase {
	return &listSurveySubmissionsUseCase{
		talkSessionRepo: talkSessionRepo,
		surveyRepo:      surveyRepo,
	}
}

func (uc *listSurveySubmissionsUseCase) Execute(ctx context.Context, input ListSurveySubmissionsInput) (*ListSurveySubmissionsOutput, error) {
	ctx, span := otel.Tracer("talksession_query").Start(ctx, "listSurveySubmissionsUseCase.Execute")
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

	limit := input.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := input.Offset
	if offset < 0 {
		offset = 0
	}

	subs, total, err := uc.surveyRepo.ListSubmissions(ctx, s.SurveyID(), limit, offset)
	if err != nil {
		utils.HandleError(ctx, err, "surveyRepo.ListSubmissions")
		return nil, errtrace.Wrap(err)
	}

	return &ListSurveySubmissionsOutput{
		Submissions: subs,
		TotalCount:  total,
	}, nil
}
