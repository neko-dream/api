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
	UpdateSurveyUseCase interface {
		Execute(ctx context.Context, input UpdateSurveyInput) (*survey.Survey, error)
	}

	UpdateSurveyInput struct {
		TalkSessionID      shared.UUID[talksession.TalkSession]
		ActorUserID        shared.UUID[user.User]
		Questions          []UpdateSurveyQuestionInput
		DeletedQuestionIDs []shared.UUID[survey.Question]
		DeletedChoiceIDs   []shared.UUID[survey.Choice]
	}

	UpdateSurveyQuestionInput struct {
		QuestionID   *shared.UUID[survey.Question] // nil なら新規追加
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
		Choices      []UpdateSurveyChoiceInput
	}

	UpdateSurveyChoiceInput struct {
		ChoiceID     *shared.UUID[survey.Choice] // nil なら新規追加
		Text         string
		DisplayOrder int
	}

	updateSurveyUseCase struct {
		talkSessionRepo talksession.TalkSessionRepository
		surveyRepo      survey.Repository
		dbm             *db.DBManager
	}
)

func NewUpdateSurveyUseCase(
	talkSessionRepo talksession.TalkSessionRepository,
	surveyRepo survey.Repository,
	dbm *db.DBManager,
) UpdateSurveyUseCase {
	return &updateSurveyUseCase{
		talkSessionRepo: talkSessionRepo,
		surveyRepo:      surveyRepo,
		dbm:             dbm,
	}
}

func (uc *updateSurveyUseCase) Execute(ctx context.Context, input UpdateSurveyInput) (*survey.Survey, error) {
	ctx, span := otel.Tracer("talksession_command").Start(ctx, "updateSurveyUseCase.Execute")
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

	// 既存質問・選択肢への更新適用
	for _, qi := range input.Questions {
		if !qi.Type.IsValid() {
			return nil, messages.TalkSessionSurveyAnswerInvalid
		}
		cfg := survey.QuestionConfig{
			IsRequired:   qi.IsRequired,
			AllowOther:   qi.AllowOther,
			MaxLength:    qi.MaxLength,
			MinValue:     qi.MinValue,
			MaxValue:     qi.MaxValue,
			MinLabel:     qi.MinLabel,
			MaxLabel:     qi.MaxLabel,
			DisplayOrder: qi.DisplayOrder,
		}
		if qi.QuestionID != nil {
			q := s.FindQuestion(*qi.QuestionID)
			if q == nil {
				return nil, messages.TalkSessionSurveyAnswerInvalid
			}
			// type 変更は回答ゼロのときのみ許可（既存回答のスキーマ整合性のため）
			if q.Type() != qi.Type {
				has, err := uc.surveyRepo.HasAnswersForQuestion(ctx, *qi.QuestionID)
				if err != nil {
					utils.HandleError(ctx, err, "surveyRepo.HasAnswersForQuestion")
					return nil, errtrace.Wrap(err)
				}
				if has {
					return nil, messages.TalkSessionSurveyQuestionTypeImmutable
				}
				q.SetType(qi.Type)
			}
			q.SetText(qi.Text)
			q.SetConfig(cfg)
			if err := applyChoices(q, qi.Choices); err != nil {
				return nil, err
			}
		} else {
			qid := shared.NewUUID[survey.Question]()
			q := survey.NewQuestion(qid, s.SurveyID(), qi.Text, qi.Type, cfg)
			for _, ci := range qi.Choices {
				q.AddChoice(survey.NewChoice(shared.NewUUID[survey.Choice](), qid, ci.Text, ci.DisplayOrder))
			}
			s.AddQuestion(q)
		}
	}

	// 削除指示の処理（回答があれば soft delete、なければ hard delete）
	if err := uc.dbm.ExecTx(ctx, func(ctx context.Context) error {
		// 質問
		for _, qid := range input.DeletedQuestionIDs {
			has, err := uc.surveyRepo.HasAnswersForQuestion(ctx, qid)
			if err != nil {
				return err
			}
			if has {
				if q := s.FindQuestion(qid); q != nil {
					q.SoftDelete()
				}
				continue
			}
			if err := uc.surveyRepo.HardDeleteQuestion(ctx, qid); err != nil {
				return err
			}
			// 後続の Save が「未知の質問＝新規」と誤判定して再挿入しないように in-memory からも削除
			s.RemoveQuestion(qid)
		}

		// 選択肢
		for _, cid := range input.DeletedChoiceIDs {
			has, err := uc.surveyRepo.HasAnswersForChoice(ctx, cid)
			if err != nil {
				return err
			}
			if has {
				for _, q := range s.Questions() {
					if c := q.FindChoice(cid); c != nil {
						c.SoftDelete()
						break
					}
				}
				continue
			}
			if err := uc.surveyRepo.HardDeleteChoice(ctx, cid); err != nil {
				return err
			}
			// 同様に in-memory からも削除
			for _, q := range s.Questions() {
				if q.FindChoice(cid) != nil {
					q.RemoveChoice(cid)
					break
				}
			}
		}

		// 永続化（既存は更新、新規は追加、soft delete も反映）
		if err := uc.surveyRepo.Save(ctx, s); err != nil {
			return err
		}
		return nil
	}); err != nil {
		utils.HandleError(ctx, err, "update survey transaction")
		return nil, errtrace.Wrap(err)
	}

	// 再読込で正確な状態を返す
	updated, err := uc.surveyRepo.FindBySessionID(ctx, input.TalkSessionID)
	if err != nil {
		utils.HandleError(ctx, err, "surveyRepo.FindBySessionID reload")
		return nil, errtrace.Wrap(err)
	}
	return updated, nil
}

func applyChoices(q *survey.Question, choices []UpdateSurveyChoiceInput) error {
	for _, ci := range choices {
		if ci.ChoiceID != nil {
			c := q.FindChoice(*ci.ChoiceID)
			if c == nil {
				return messages.TalkSessionSurveyAnswerInvalid
			}
			c.SetText(ci.Text)
			c.SetDisplayOrder(ci.DisplayOrder)
			continue
		}
		q.AddChoice(survey.NewChoice(shared.NewUUID[survey.Choice](), q.QuestionID(), ci.Text, ci.DisplayOrder))
	}
	return nil
}
