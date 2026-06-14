package talksession_usecase

import (
	"context"
	"database/sql"
	"errors"

	"braces.dev/errtrace"
	"github.com/jackc/pgx/v5/pgconn"
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
		Text              *string
		SelectedChoiceIDs []shared.UUID[survey.Choice]
		Other             *string
		Rating            *int
		Date              *string
	}

	submitSurveyUseCase struct {
		surveyRepo survey.Repository
		dbm        *db.DBManager
	}
)

func NewSubmitSurveyUseCase(surveyRepo survey.Repository, dbm *db.DBManager) SubmitSurveyUseCase {
	return &submitSurveyUseCase{surveyRepo: surveyRepo, dbm: dbm}
}

func (uc *submitSurveyUseCase) Execute(ctx context.Context, input SubmitSurveyInput) (shared.UUID[survey.Submission], error) {
	ctx, span := otel.Tracer("talksession_command").Start(ctx, "submitSurveyUseCase.Execute")
	defer span.End()

	s, err := uc.surveyRepo.FindBySessionID(ctx, input.TalkSessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return shared.UUID[survey.Submission]{}, messages.TalkSessionSurveyNotFound
		}
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

	// 回答ドメインオブジェクトに変換
	submissionID := shared.NewUUID[survey.Submission]()
	answers := make([]*survey.Answer, 0, len(input.Answers))
	for _, ai := range input.Answers {
		v := survey.AnswerValue{
			Text:              ai.Text,
			SelectedChoiceIDs: ai.SelectedChoiceIDs,
			Other:             ai.Other,
			Rating:            ai.Rating,
			Date:              ai.Date,
		}
		answers = append(answers, survey.NewAnswer(
			shared.NewUUID[survey.Answer](),
			submissionID,
			ai.QuestionID,
			v,
		))
	}

	// バリデーション
	if err := survey.ValidateAnswers(s, answers); err != nil {
		e := *messages.TalkSessionSurveyAnswerInvalid
		e.Message = err.Error()
		return shared.UUID[survey.Submission]{}, &e
	}

	sub := survey.NewSubmission(submissionID, s.SurveyID(), input.UserID, answers)
	if err := uc.dbm.ExecTx(ctx, func(ctx context.Context) error {
		return uc.surveyRepo.CreateSubmission(ctx, sub)
	}); err != nil {
		// 事前チェックをすり抜けた並行リクエストの保険：UNIQUE 違反は AlreadySubmitted に丸める
		if isUniqueViolation(err) {
			return shared.UUID[survey.Submission]{}, messages.TalkSessionSurveyAlreadySubmitted
		}
		utils.HandleError(ctx, err, "surveyRepo.CreateSubmission")
		return shared.UUID[survey.Submission]{}, errtrace.Wrap(err)
	}

	return submissionID, nil
}

// isUniqueViolation は PostgreSQL unique_violation (23505) を検出する。
// pgx の標準エラー型を errors.As で取り出すため、otelsql / errtrace ラップを跨いでも安定。
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
