package talksession_usecase

import (
	"context"
	"database/sql"
	"errors"

	"braces.dev/errtrace"
	"github.com/neko-dream/api/internal/domain/messages"
	domain_crypto "github.com/neko-dream/api/internal/domain/model/crypto"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/talksession/survey"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/crypto"
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
		encryptor       domain_crypto.Encryptor
	}
)

func NewGetSurveySummaryUseCase(
	talkSessionRepo talksession.TalkSessionRepository,
	surveyRepo survey.Repository,
	encryptor domain_crypto.Encryptor,
) GetSurveySummaryUseCase {
	return &getSurveySummaryUseCase{
		talkSessionRepo: talkSessionRepo,
		surveyRepo:      surveyRepo,
		encryptor:       encryptor,
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

	// 回答はすべて暗号化されているため、集計前にサーバ側で復号する。
	for _, sub := range subs {
		for _, a := range sub.Answers() {
			if !a.IsEncrypted() {
				continue
			}
			dec, err := crypto.DecryptAnswerValue(ctx, uc.encryptor, a.Value())
			if err != nil {
				utils.HandleError(ctx, err, "crypto.DecryptAnswerValue")
				return nil, errtrace.Wrap(err)
			}
			a.SetDecryptedValue(dec)
		}
	}

	return survey.Aggregate(s, subs), nil
}
