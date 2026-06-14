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
	for _, row := range existingRows {
		existingQuestions[row.QuestionID] = true
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

		// 既存の選択肢取得
		choiceRows, err := q.ListSurveyChoicesByQuestionIDs(ctx, []uuid.UUID{qid})
		if err != nil {
			utils.HandleError(ctx, err, "ListSurveyChoicesByQuestionIDs")
			return err
		}
		existingChoices := make(map[uuid.UUID]bool, len(choiceRows))
		for _, row := range choiceRows {
			existingChoices[row.ChoiceID] = true
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
		MaxLength:    nullInt32FromPtr(qn.MaxLength()),
		MinValue:     nullInt32FromPtr(qn.MinValue()),
		MaxValue:     nullInt32FromPtr(qn.MaxValue()),
		MinLabel:     nullStringFromPtr(qn.MinLabel()),
		MaxLabel:     nullStringFromPtr(qn.MaxLabel()),
		DisplayOrder: int32(qn.DisplayOrder()),
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
		MaxLength:    nullInt32FromPtr(qn.MaxLength()),
		MinValue:     nullInt32FromPtr(qn.MinValue()),
		MaxValue:     nullInt32FromPtr(qn.MaxValue()),
		MinLabel:     nullStringFromPtr(qn.MinLabel()),
		MaxLabel:     nullStringFromPtr(qn.MaxLabel()),
		DisplayOrder: int32(qn.DisplayOrder()),
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
				nullTimeToPtr(row.DeletedAt),
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
			nullInt32ToIntPtr(row.MaxLength),
			nullInt32ToIntPtr(row.MinValue),
			nullInt32ToIntPtr(row.MaxValue),
			nullStringToPtr(row.MinLabel),
			nullStringToPtr(row.MaxLabel),
			int(row.DisplayOrder),
			cRowsByQuestion[row.QuestionID],
			nullTimeToPtr(row.DeletedAt),
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

	subIDs := make([]uuid.UUID, 0, len(subRows))
	for _, s := range subRows {
		subIDs = append(subIDs, s.SubmissionID)
	}

	answersBySubmission := make(map[uuid.UUID][]*survey.Answer, len(subIDs))
	if len(subIDs) > 0 {
		aRows, err := q.ListSurveyAnswersBySubmissionIDs(ctx, subIDs)
		if err != nil {
			return nil, 0, err
		}
		for _, row := range aRows {
			v, err := survey.UnmarshalAnswerValue(row.AnswerValue)
			if err != nil {
				return nil, 0, err
			}
			a := survey.ReconstructAnswer(
				shared.UUID[survey.Answer](row.AnswerID),
				shared.UUID[survey.Submission](row.SubmissionID),
				shared.UUID[survey.Question](row.QuestionID),
				v,
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
	return subs, int32(total), nil
}

// ヘルパ
func nullInt32FromPtr(p *int) sql.NullInt32 {
	if p == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: int32(*p), Valid: true}
}

func nullStringFromPtr(p *string) sql.NullString {
	if p == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *p, Valid: true}
}

func nullInt32ToIntPtr(n sql.NullInt32) *int {
	if !n.Valid {
		return nil
	}
	v := int(n.Int32)
	return &v
}

func nullStringToPtr(n sql.NullString) *string {
	if !n.Valid {
		return nil
	}
	return &n.String
}

func nullTimeToPtr(n sql.NullTime) *time.Time {
	if !n.Valid {
		return nil
	}
	return &n.Time
}
