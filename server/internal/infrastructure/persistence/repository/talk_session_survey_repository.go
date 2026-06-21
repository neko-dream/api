package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/talksession/survey"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

type talkSessionSurveyRepository struct {
	*db.DBManager
}

func NewTalkSessionSurveyRepository(dbm *db.DBManager) survey.Repository {
	return &talkSessionSurveyRepository{DBManager: dbm}
}

func (r *talkSessionSurveyRepository) Create(ctx context.Context, s *survey.Survey) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "talkSessionSurveyRepository.Create")
	defer span.End()

	q := r.GetQueries(ctx)
	if err := q.CreateTalkSessionSurvey(ctx, model.CreateTalkSessionSurveyParams{
		SurveyID:      s.SurveyID().UUID(),
		TalkSessionID: s.TalkSessionID().UUID(),
		CreatedAt:     s.CreatedAt(),
		UpdatedAt:     s.UpdatedAt(),
	}); err != nil {
		utils.HandleError(ctx, err, "CreateTalkSessionSurvey")
		return err
	}

	for _, qn := range s.Questions() {
		if err := r.insertQuestion(ctx, qn); err != nil {
			return err
		}
		for _, c := range qn.Choices() {
			if err := r.insertChoice(ctx, c); err != nil {
				return err
			}
		}
	}
	return nil
}

// Save は質問・選択肢を upsert する。
// 既存 ID は UPDATE / soft-delete を反映、新規 ID は INSERT。
// hard delete は HardDeleteQuestion / HardDeleteChoice 経由。
func (r *talkSessionSurveyRepository) Save(ctx context.Context, s *survey.Survey) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "talkSessionSurveyRepository.Save")
	defer span.End()

	q := r.GetQueries(ctx)

	// 既存質問 ID を取得
	existingRows, err := q.ListSurveyQuestions(ctx, s.SurveyID().UUID())
	if err != nil {
		utils.HandleError(ctx, err, "ListSurveyQuestions")
		return err
	}
	existingQuestions := make(map[uuid.UUID]bool, len(existingRows))
	allQuestionIDs := make([]uuid.UUID, 0, len(existingRows))
	for _, row := range existingRows {
		existingQuestions[row.QuestionID] = true
		allQuestionIDs = append(allQuestionIDs, row.QuestionID)
	}

	// 既存選択肢を質問単位の N+1 にせず一括取得する
	existingChoices := make(map[uuid.UUID]bool)
	if len(allQuestionIDs) > 0 {
		choiceRows, err := q.ListSurveyChoicesByQuestionIDs(ctx, allQuestionIDs)
		if err != nil {
			utils.HandleError(ctx, err, "ListSurveyChoicesByQuestionIDs")
			return err
		}
		for _, row := range choiceRows {
			existingChoices[row.ChoiceID] = true
		}
	}

	for _, qn := range s.Questions() {
		qid := qn.QuestionID().UUID()
		if existingQuestions[qid] {
			// 更新 or soft delete
			if qn.IsDeleted() {
				if err := q.SoftDeleteSurveyQuestion(ctx, model.SoftDeleteSurveyQuestionParams{
					QuestionID: qid,
					DeletedAt:  sql.NullTime{Time: *qn.DeletedAt(), Valid: true},
				}); err != nil {
					utils.HandleError(ctx, err, "SoftDeleteSurveyQuestion")
					return err
				}
			} else {
				if err := r.updateQuestion(ctx, qn); err != nil {
					return err
				}
			}
		} else {
			if err := r.insertQuestion(ctx, qn); err != nil {
				return err
			}
		}

		for _, c := range qn.Choices() {
			cid := c.ChoiceID().UUID()
			if existingChoices[cid] {
				if c.IsDeleted() {
					if err := q.SoftDeleteSurveyChoice(ctx, model.SoftDeleteSurveyChoiceParams{
						ChoiceID:  cid,
						DeletedAt: sql.NullTime{Time: *c.DeletedAt(), Valid: true},
					}); err != nil {
						utils.HandleError(ctx, err, "SoftDeleteSurveyChoice")
						return err
					}
				} else {
					if err := q.UpdateSurveyChoice(ctx, model.UpdateSurveyChoiceParams{
						ChoiceID:     cid,
						ChoiceText:   c.Text(),
						DisplayOrder: int32(c.DisplayOrder()),
						UpdatedAt:    time.Now(),
					}); err != nil {
						utils.HandleError(ctx, err, "UpdateSurveyChoice")
						return err
					}
				}
			} else {
				if err := r.insertChoice(ctx, c); err != nil {
					return err
				}
			}
		}
	}

	if err := q.TouchSurvey(ctx, model.TouchSurveyParams{
		SurveyID:  s.SurveyID().UUID(),
		UpdatedAt: time.Now(),
	}); err != nil {
		utils.HandleError(ctx, err, "TouchSurvey")
		return err
	}
	return nil
}

func (r *talkSessionSurveyRepository) insertQuestion(ctx context.Context, qn *survey.Question) error {
	return r.GetQueries(ctx).CreateSurveyQuestion(ctx, model.CreateSurveyQuestionParams{
		QuestionID:   qn.QuestionID().UUID(),
		SurveyID:     qn.SurveyID().UUID(),
		QuestionText: qn.Text(),
		QuestionType: model.TalkSessionSurveyQuestionType(qn.Type()),
		IsRequired:   qn.IsRequired(),
		AllowOther:   qn.AllowOther(),
		MaxLength:    utils.ToNullableSQL[sql.NullInt32](qn.MaxLength()),
		MinValue:     utils.ToNullableSQL[sql.NullInt32](qn.MinValue()),
		MaxValue:     utils.ToNullableSQL[sql.NullInt32](qn.MaxValue()),
		MinLabel:     utils.ToNullableSQL[sql.NullString](qn.MinLabel()),
		MaxLabel:     utils.ToNullableSQL[sql.NullString](qn.MaxLabel()),
		DisplayOrder: int32(qn.DisplayOrder()),
		ContainsPii:  qn.ContainsPii(),
		CreatedAt:    qn.CreatedAt(),
		UpdatedAt:    qn.UpdatedAt(),
	})
}

func (r *talkSessionSurveyRepository) updateQuestion(ctx context.Context, qn *survey.Question) error {
	return r.GetQueries(ctx).UpdateSurveyQuestion(ctx, model.UpdateSurveyQuestionParams{
		QuestionID:   qn.QuestionID().UUID(),
		QuestionText: qn.Text(),
		QuestionType: model.TalkSessionSurveyQuestionType(qn.Type()),
		IsRequired:   qn.IsRequired(),
		AllowOther:   qn.AllowOther(),
		MaxLength:    utils.ToNullableSQL[sql.NullInt32](qn.MaxLength()),
		MinValue:     utils.ToNullableSQL[sql.NullInt32](qn.MinValue()),
		MaxValue:     utils.ToNullableSQL[sql.NullInt32](qn.MaxValue()),
		MinLabel:     utils.ToNullableSQL[sql.NullString](qn.MinLabel()),
		MaxLabel:     utils.ToNullableSQL[sql.NullString](qn.MaxLabel()),
		DisplayOrder: int32(qn.DisplayOrder()),
		ContainsPii:  qn.ContainsPii(),
		UpdatedAt:    time.Now(),
	})
}

func (r *talkSessionSurveyRepository) insertChoice(ctx context.Context, c *survey.Choice) error {
	return r.GetQueries(ctx).CreateSurveyChoice(ctx, model.CreateSurveyChoiceParams{
		ChoiceID:     c.ChoiceID().UUID(),
		QuestionID:   c.QuestionID().UUID(),
		ChoiceText:   c.Text(),
		DisplayOrder: int32(c.DisplayOrder()),
		CreatedAt:    c.CreatedAt(),
		UpdatedAt:    c.UpdatedAt(),
	})
}

func (r *talkSessionSurveyRepository) FindBySurveyID(ctx context.Context, surveyID shared.UUID[survey.Survey]) (*survey.Survey, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "talkSessionSurveyRepository.FindBySurveyID")
	defer span.End()

	q := r.GetQueries(ctx)
	row, err := q.GetSurveyByID(ctx, surveyID.UUID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.reconstruct(ctx, row.SurveyID, row.TalkSessionID, row.CreatedAt, row.UpdatedAt)
}

func (r *talkSessionSurveyRepository) FindBySessionID(ctx context.Context, sessionID shared.UUID[talksession.TalkSession]) (*survey.Survey, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "talkSessionSurveyRepository.FindBySessionID")
	defer span.End()

	q := r.GetQueries(ctx)
	row, err := q.GetSurveyByTalkSessionID(ctx, sessionID.UUID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.reconstruct(ctx, row.SurveyID, row.TalkSessionID, row.CreatedAt, row.UpdatedAt)
}

func (r *talkSessionSurveyRepository) reconstruct(
	ctx context.Context,
	surveyID, talkSessionID uuid.UUID,
	createdAt, updatedAt time.Time,
) (*survey.Survey, error) {
	q := r.GetQueries(ctx)

	qRows, err := q.ListSurveyQuestions(ctx, surveyID)
	if err != nil {
		return nil, err
	}

	questionIDs := make([]uuid.UUID, 0, len(qRows))
	for _, row := range qRows {
		questionIDs = append(questionIDs, row.QuestionID)
	}

	cRowsByQuestion := make(map[uuid.UUID][]*survey.Choice, len(questionIDs))
	if len(questionIDs) > 0 {
		cRows, err := q.ListSurveyChoicesByQuestionIDs(ctx, questionIDs)
		if err != nil {
			return nil, err
		}
		for _, row := range cRows {
			c := survey.ReconstructChoice(
				shared.UUID[survey.Choice](row.ChoiceID),
				shared.UUID[survey.Question](row.QuestionID),
				row.ChoiceText,
				int(row.DisplayOrder),
				utils.ToPtrIf(row.DeletedAt.Valid, row.DeletedAt.Time),
				row.CreatedAt,
				row.UpdatedAt,
			)
			cRowsByQuestion[row.QuestionID] = append(cRowsByQuestion[row.QuestionID], c)
		}
	}

	questions := make([]*survey.Question, 0, len(qRows))
	for _, row := range qRows {
		qt, err := survey.ParseQuestionType(row.QuestionType)
		if err != nil {
			return nil, err
		}
		q := survey.ReconstructQuestion(
			shared.UUID[survey.Question](row.QuestionID),
			shared.UUID[survey.Survey](row.SurveyID),
			row.QuestionText,
			qt,
			row.IsRequired,
			row.AllowOther,
			utils.ToPtrIf(row.MaxLength.Valid, int(row.MaxLength.Int32)),
			utils.ToPtrIf(row.MinValue.Valid, int(row.MinValue.Int32)),
			utils.ToPtrIf(row.MaxValue.Valid, int(row.MaxValue.Int32)),
			utils.ToPtrIf(row.MinLabel.Valid, row.MinLabel.String),
			utils.ToPtrIf(row.MaxLabel.Valid, row.MaxLabel.String),
			int(row.DisplayOrder),
			row.ContainsPii,
			cRowsByQuestion[row.QuestionID],
			utils.ToPtrIf(row.DeletedAt.Valid, row.DeletedAt.Time),
			row.CreatedAt,
			row.UpdatedAt,
		)
		questions = append(questions, q)
	}

	return survey.ReconstructSurvey(
		shared.UUID[survey.Survey](surveyID),
		shared.UUID[talksession.TalkSession](talkSessionID),
		questions,
		createdAt,
		updatedAt,
	), nil
}

func (r *talkSessionSurveyRepository) HardDeleteQuestion(ctx context.Context, questionID shared.UUID[survey.Question]) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "talkSessionSurveyRepository.HardDeleteQuestion")
	defer span.End()

	q := r.GetQueries(ctx)
	count, err := q.CountAnswersForQuestion(ctx, questionID.UUID())
	if err != nil {
		return err
	}
	if count > 0 {
		return messages.TalkSessionSurveyCannotHardDelete
	}
	return q.HardDeleteSurveyQuestion(ctx, questionID.UUID())
}

func (r *talkSessionSurveyRepository) HardDeleteChoice(ctx context.Context, choiceID shared.UUID[survey.Choice]) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "talkSessionSurveyRepository.HardDeleteChoice")
	defer span.End()

	q := r.GetQueries(ctx)
	count, err := q.CountAnswersForChoice(ctx, choiceID.UUID().String())
	if err != nil {
		return err
	}
	if count > 0 {
		return messages.TalkSessionSurveyCannotHardDelete
	}
	return q.HardDeleteSurveyChoice(ctx, choiceID.UUID())
}

func (r *talkSessionSurveyRepository) HasAnswersForQuestion(ctx context.Context, questionID shared.UUID[survey.Question]) (bool, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "talkSessionSurveyRepository.HasAnswersForQuestion")
	defer span.End()

	count, err := r.GetQueries(ctx).CountAnswersForQuestion(ctx, questionID.UUID())
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *talkSessionSurveyRepository) HasAnswersForChoice(ctx context.Context, choiceID shared.UUID[survey.Choice]) (bool, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "talkSessionSurveyRepository.HasAnswersForChoice")
	defer span.End()

	count, err := r.GetQueries(ctx).CountAnswersForChoice(ctx, choiceID.UUID().String())
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *talkSessionSurveyRepository) CreateSubmission(ctx context.Context, sub *survey.Submission) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "talkSessionSurveyRepository.CreateSubmission")
	defer span.End()

	q := r.GetQueries(ctx)
	if err := q.CreateSurveySubmission(ctx, model.CreateSurveySubmissionParams{
		SubmissionID: sub.SubmissionID().UUID(),
		SurveyID:     sub.SurveyID().UUID(),
		UserID:       sub.UserID().UUID(),
		SubmittedAt:  sub.SubmittedAt(),
	}); err != nil {
		utils.HandleError(ctx, err, "CreateSurveySubmission")
		return err
	}

	for _, a := range sub.Answers() {
		payload, err := a.Value().Marshal()
		if err != nil {
			return err
		}
		if err := q.CreateSurveyAnswer(ctx, model.CreateSurveyAnswerParams{
			AnswerID:     a.AnswerID().UUID(),
			SubmissionID: a.SubmissionID().UUID(),
			QuestionID:   a.QuestionID().UUID(),
			AnswerValue:  payload,
			IsEncrypted:  a.IsEncrypted(),
			CreatedAt:    a.CreatedAt(),
		}); err != nil {
			utils.HandleError(ctx, err, "CreateSurveyAnswer")
			return err
		}
	}
	return nil
}

func (r *talkSessionSurveyRepository) HasSubmittedBySessionID(
	ctx context.Context,
	sessionID shared.UUID[talksession.TalkSession],
	userID shared.UUID[user.User],
) (bool, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "talkSessionSurveyRepository.HasSubmittedBySessionID")
	defer span.End()

	return r.GetQueries(ctx).HasSubmittedSurveyBySessionID(ctx, model.HasSubmittedSurveyBySessionIDParams{
		TalkSessionID: sessionID.UUID(),
		UserID:        userID.UUID(),
	})
}

func (r *talkSessionSurveyRepository) ListSubmissions(
	ctx context.Context,
	surveyID shared.UUID[survey.Survey],
	limit, offset int,
) ([]*survey.Submission, int32, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "talkSessionSurveyRepository.ListSubmissions")
	defer span.End()

	q := r.GetQueries(ctx)
	total, err := q.CountSurveySubmissions(ctx, surveyID.UUID())
	if err != nil {
		return nil, 0, err
	}

	subRows, err := q.ListSurveySubmissions(ctx, model.ListSurveySubmissionsParams{
		SurveyID: surveyID.UUID(),
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	subs, err := r.buildSubmissions(ctx, subRows)
	if err != nil {
		return nil, 0, err
	}
	return subs, int32(total), nil
}

func (r *talkSessionSurveyRepository) ListAllSubmissions(
	ctx context.Context,
	surveyID shared.UUID[survey.Survey],
) ([]*survey.Submission, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "talkSessionSurveyRepository.ListAllSubmissions")
	defer span.End()

	subRows, err := r.GetQueries(ctx).ListAllSurveySubmissions(ctx, surveyID.UUID())
	if err != nil {
		return nil, err
	}
	return r.buildSubmissions(ctx, subRows)
}

// buildSubmissions は submissions 行に紐付く answers をまとめて取得し、ドメインオブジェクトを組み立てる。
func (r *talkSessionSurveyRepository) buildSubmissions(
	ctx context.Context,
	subRows []model.TalkSessionSurveySubmission,
) ([]*survey.Submission, error) {
	subIDs := make([]uuid.UUID, 0, len(subRows))
	for _, s := range subRows {
		subIDs = append(subIDs, s.SubmissionID)
	}

	answersBySubmission := make(map[uuid.UUID][]*survey.Answer, len(subIDs))
	if len(subIDs) > 0 {
		aRows, err := r.GetQueries(ctx).ListSurveyAnswersBySubmissionIDs(ctx, subIDs)
		if err != nil {
			return nil, err
		}
		for _, row := range aRows {
			v, err := survey.UnmarshalAnswerValue(row.AnswerValue)
			if err != nil {
				return nil, err
			}
			a := survey.ReconstructAnswer(
				shared.UUID[survey.Answer](row.AnswerID),
				shared.UUID[survey.Submission](row.SubmissionID),
				shared.UUID[survey.Question](row.QuestionID),
				v,
				row.IsEncrypted,
				row.CreatedAt,
			)
			answersBySubmission[row.SubmissionID] = append(answersBySubmission[row.SubmissionID], a)
		}
	}

	subs := make([]*survey.Submission, 0, len(subRows))
	for _, row := range subRows {
		subs = append(subs, survey.ReconstructSubmission(
			shared.UUID[survey.Submission](row.SubmissionID),
			shared.UUID[survey.Survey](row.SurveyID),
			shared.UUID[user.User](row.UserID),
			row.SubmittedAt,
			answersBySubmission[row.SubmissionID],
		))
	}
	return subs, nil
}

