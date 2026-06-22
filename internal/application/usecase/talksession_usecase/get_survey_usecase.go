package talksession_usecase

import (
	"context"

	"braces.dev/errtrace"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/talksession/survey"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

type (
	GetSurveyUseCase interface {
		Execute(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) (*survey.Survey, error)
	}

	getSurveyUseCase struct {
		surveyRepo survey.Repository
	}
)

func NewGetSurveyUseCase(surveyRepo survey.Repository) GetSurveyUseCase {
	return &getSurveyUseCase{surveyRepo: surveyRepo}
}

func (uc *getSurveyUseCase) Execute(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) (*survey.Survey, error) {
	ctx, span := otel.Tracer("talksession_query").Start(ctx, "getSurveyUseCase.Execute")
	defer span.End()

	// FindBySessionID は存在しない場合 (nil, nil) を返す契約
	s, err := uc.surveyRepo.FindBySessionID(ctx, talkSessionID)
	if err != nil {
		utils.HandleError(ctx, err, "surveyRepo.FindBySessionID")
		return nil, errtrace.Wrap(err)
	}
	if s == nil {
		return nil, messages.TalkSessionSurveyNotFound
	}
	return s, nil
}
