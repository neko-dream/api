package talksession_usecase

import (
	"context"

	"braces.dev/errtrace"
	"github.com/neko-dream/api/internal/domain/messages"
	domain_crypto "github.com/neko-dream/api/internal/domain/model/crypto"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/talksession/survey"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/crypto"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

type (
	SubmitSurveyUseCase interface {
		Execute(ctx context.Context, input SubmitSurveyInput) (shared.UUID[survey.Submission], error)
	}

	SubmitSurveyInput struct {
		TalkSessionID shared.UUID[talksession.TalkSession]
		UserID        shared.UUID[user.User]
		Answers       []SurveyAnswerInput
	}

	SurveyAnswerInput struct {
		QuestionID        shared.UUID[survey.Question]
		SelectedChoiceIDs []shared.UUID[survey.Choice]
		Other             *string
	}

	submitSurveyUseCase struct {
		surveyRepo survey.Repository
		dbm        *db.DBManager
		encryptor  domain_crypto.Encryptor
	}
)

func NewSubmitSurveyUseCase(
	surveyRepo survey.Repository,
	dbm *db.DBManager,
	encryptor domain_crypto.Encryptor,
) SubmitSurveyUseCase {
	return &submitSurveyUseCase{surveyRepo: surveyRepo, dbm: dbm, encryptor: encryptor}
}

func (uc *submitSurveyUseCase) Execute(ctx context.Context, input SubmitSurveyInput) (shared.UUID[survey.Submission], error) {
	ctx, span := otel.Tracer("talksession_command").Start(ctx, "submitSurveyUseCase.Execute")
	defer span.End()

	// FindBySessionID は存在しない場合 (nil, nil) を返す契約
	s, err := uc.surveyRepo.FindBySessionID(ctx, input.TalkSessionID)
	if err != nil {
		utils.HandleError(ctx, err, "surveyRepo.FindBySessionID")
		return shared.UUID[survey.Submission]{}, errtrace.Wrap(err)
	}
	if s == nil {
		return shared.UUID[survey.Submission]{}, messages.TalkSessionSurveyNotFound
	}

	// 重複送信を事前チェック（1ユーザー1送信、修正不可）
	submitted, err := uc.surveyRepo.HasSubmittedBySessionID(ctx, input.TalkSessionID, input.UserID)
	if err != nil {
		utils.HandleError(ctx, err, "surveyRepo.HasSubmittedBySessionID")
		return shared.UUID[survey.Submission]{}, errtrace.Wrap(err)
	}
	if submitted {
		return shared.UUID[survey.Submission]{}, messages.TalkSessionSurveyAlreadySubmitted
	}

	// 回答ドメインオブジェクトに変換（暗号化前のバリデーション用）
	submissionID := shared.NewUUID[survey.Submission]()
	answers := make([]*survey.Answer, 0, len(input.Answers))
	for _, ai := range input.Answers {
		v := survey.AnswerValue{
			SelectedChoiceIDs: ai.SelectedChoiceIDs,
			Other:             ai.Other,
		}
		answers = append(answers, survey.NewAnswer(
			shared.NewUUID[survey.Answer](),
			submissionID,
			ai.QuestionID,
			v,
		))
	}

	// バリデーション（平文に対して実施）
	if err := survey.ValidateAnswers(s, answers); err != nil {
		return shared.UUID[survey.Submission]{}, messages.TalkSessionSurveyAnswerInvalid.WithMessage(err.Error())
	}

	// すべての回答を暗号化して保存する（選択肢 ID / その他テキストを含む）。
	for _, a := range answers {
		enc, err := crypto.EncryptAnswerValue(ctx, uc.encryptor, a.Value())
		if err != nil {
			utils.HandleError(ctx, err, "crypto.EncryptAnswerValue")
			return shared.UUID[survey.Submission]{}, errtrace.Wrap(err)
		}
		a.SetEncryptedValue(enc)
	}

	sub := survey.NewSubmission(submissionID, s.SurveyID(), input.UserID, answers)
	if err := uc.dbm.ExecTx(ctx, func(ctx context.Context) error {
		return uc.surveyRepo.CreateSubmission(ctx, sub)
	}); err != nil {
		// 事前チェックをすり抜けた並行リクエストの保険：UNIQUE 違反は AlreadySubmitted に丸める
		if utils.IsUniqueViolation(err) {
			return shared.UUID[survey.Submission]{}, messages.TalkSessionSurveyAlreadySubmitted
		}
		utils.HandleError(ctx, err, "surveyRepo.CreateSubmission")
		return shared.UUID[survey.Submission]{}, errtrace.Wrap(err)
	}

	return submissionID, nil
}
