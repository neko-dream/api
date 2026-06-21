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
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

type (
	CreateSurveyUseCase interface {
		Execute(ctx context.Context, input CreateSurveyInput) (*survey.Survey, error)
	}

	CreateSurveyInput struct {
		TalkSessionID shared.UUID[talksession.TalkSession]
		ActorUserID   shared.UUID[user.User]
		Questions     []SurveyQuestionInput
	}

	SurveyQuestionInput struct {
		Text         string
		Type         survey.QuestionType
		IsRequired   bool
		AllowOther   bool
		MaxLength    *int
		MinValue     *int
		MaxValue     *int
		MinLabel     *string
		MaxLabel     *string
		DisplayOrder int
		ContainsPii  bool
		Choices      []SurveyChoiceInput
	}

	SurveyChoiceInput struct {
		Text         string
		DisplayOrder int
	}

	createSurveyUseCase struct {
		talkSessionRepo talksession.TalkSessionRepository
		surveyRepo      survey.Repository
		dbm             *db.DBManager
	}
)

func NewCreateSurveyUseCase(
	talkSessionRepo talksession.TalkSessionRepository,
	surveyRepo survey.Repository,
	dbm *db.DBManager,
) CreateSurveyUseCase {
	return &createSurveyUseCase{
		talkSessionRepo: talkSessionRepo,
		surveyRepo:      surveyRepo,
		dbm:             dbm,
	}
}

func (uc *createSurveyUseCase) Execute(ctx context.Context, input CreateSurveyInput) (*survey.Survey, error) {
	ctx, span := otel.Tracer("talksession_command").Start(ctx, "createSurveyUseCase.Execute")
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

	// オーナーのみ作成可能
	if ts.OwnerUserID() != input.ActorUserID {
		return nil, messages.ForbiddenError
	}

	// 既存チェック（1セッション1アンケート）
	existing, err := uc.surveyRepo.FindBySessionID(ctx, input.TalkSessionID)
	if err != nil {
		utils.HandleError(ctx, err, "surveyRepo.FindBySessionID")
		return nil, errtrace.Wrap(err)
	}
	if existing != nil {
		return nil, messages.TalkSessionSurveyAlreadyExists
	}

	// ドメインオブジェクト構築
	surveyID := shared.NewUUID[survey.Survey]()
	s := survey.NewSurvey(surveyID, input.TalkSessionID)
	for _, qi := range input.Questions {
		if !qi.Type.IsValid() {
			return nil, messages.TalkSessionSurveyAnswerInvalid
		}
		qid := shared.NewUUID[survey.Question]()
		q := survey.NewQuestion(qid, surveyID, qi.Text, qi.Type, survey.QuestionConfig{
			IsRequired:   qi.IsRequired,
			AllowOther:   qi.AllowOther,
			MaxLength:    qi.MaxLength,
			MinValue:     qi.MinValue,
			MaxValue:     qi.MaxValue,
			MinLabel:     qi.MinLabel,
			MaxLabel:     qi.MaxLabel,
			DisplayOrder: qi.DisplayOrder,
			ContainsPii:  qi.ContainsPii,
		})
		for _, ci := range qi.Choices {
			q.AddChoice(survey.NewChoice(shared.NewUUID[survey.Choice](), qid, ci.Text, ci.DisplayOrder))
		}
		s.AddQuestion(q)
	}

	// トランザクション内で survey 作成 と セッションの restrictions に survey.answered を追加
	if err := uc.dbm.ExecTx(ctx, func(ctx context.Context) error {
		if err := uc.surveyRepo.Create(ctx, s); err != nil {
			return err
		}

		// 既存制限に survey.answered を追加（既に含まれていれば no-op）
		if ts.EnsureRestriction(talksession.SurveyAnswered) {
			if err := uc.talkSessionRepo.Update(ctx, ts); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		utils.HandleError(ctx, err, "create survey transaction")
		return nil, errtrace.Wrap(err)
	}

	return s, nil
}
