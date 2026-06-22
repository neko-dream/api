package survey

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/neko-dream/api/internal/domain/model/shared"
)

// ValidateAnswers は提出された回答が survey の質問と整合しているか検証する。
// 戻り値の error は最初の違反を表す。
func ValidateAnswers(s *Survey, answers []*Answer) error {
	if s == nil {
		return errors.New("survey is nil")
	}

	// 質問 ID → 質問 のマップ
	active := s.ActiveQuestions()
	qByID := make(map[shared.UUID[Question]]*Question, len(active))
	for _, q := range active {
		qByID[q.QuestionID()] = q
	}

	// 重複回答の検出
	seen := make(map[shared.UUID[Question]]bool, len(answers))
	for _, a := range answers {
		if seen[a.QuestionID()] {
			return fmt.Errorf("duplicate answer for question %s", a.QuestionID())
		}
		seen[a.QuestionID()] = true
	}

	// 未知 question_id の検出 と 個別バリデーション
	for _, a := range answers {
		q, ok := qByID[a.QuestionID()]
		if !ok {
			return fmt.Errorf("answer references unknown or deleted question %s", a.QuestionID())
		}
		if err := validateAnswerAgainstQuestion(q, a); err != nil {
			return err
		}
	}

	// 必須質問の未回答チェック
	for _, q := range active {
		if !q.IsRequired() {
			continue
		}
		if !seen[q.QuestionID()] {
			return fmt.Errorf("required question %s is not answered", q.QuestionID())
		}
	}

	return nil
}

func validateAnswerAgainstQuestion(q *Question, a *Answer) error {
	v := a.Value()

	// 必須質問は空回答を許さない
	if q.IsRequired() && v.IsEmpty() {
		return fmt.Errorf("required question %s has empty answer", q.QuestionID())
	}

	// 必須でない場合、空回答は許容（回答省略扱い）
	if v.IsEmpty() {
		return nil
	}

	switch q.Type() {
	case QuestionTypeSingleChoice, QuestionTypeDropdown:
		hasChoice := len(v.SelectedChoiceIDs) >= 1
		hasOther, err := validateOther(q, v)
		if err != nil {
			return err
		}
		if !hasChoice && !hasOther {
			return fmt.Errorf("single_choice/dropdown question %s requires exactly one choice or other", q.QuestionID())
		}
		// single_choice / dropdown は choice と other の同時送信を許さない（排他）
		if hasChoice && hasOther {
			return fmt.Errorf("single_choice/dropdown question %s allows either choice or other, not both", q.QuestionID())
		}
		if len(v.SelectedChoiceIDs) > 1 {
			return fmt.Errorf("single_choice/dropdown question %s allows only one choice", q.QuestionID())
		}
		if err := validateChoiceIDsBelongToQuestion(q, v.SelectedChoiceIDs); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unsupported question type %s", q.Type())
	}

	return nil
}

// validateOther は AllowOther=true の choice 系設問に対する Other フィールドを検証する。
// 戻り値 hasOther は「Other に意味のある値があるか」を示す（型/許可フラグ判定はここで吸収）。
func validateOther(q *Question, v AnswerValue) (hasOther bool, err error) {
	if !q.AllowOther() || v.Other == nil || *v.Other == "" {
		return false, nil
	}
	if q.MaxLength() != nil && utf8.RuneCountInString(*v.Other) > *q.MaxLength() {
		return true, fmt.Errorf("other text exceeds max_length for question %s", q.QuestionID())
	}
	return true, nil
}

func validateChoiceIDsBelongToQuestion(q *Question, ids []shared.UUID[Choice]) error {
	if len(ids) == 0 {
		return nil
	}
	activeIDs := make(map[shared.UUID[Choice]]bool, len(q.Choices()))
	for _, c := range q.ActiveChoices() {
		activeIDs[c.ChoiceID()] = true
	}
	for _, id := range ids {
		if !activeIDs[id] {
			return fmt.Errorf("choice %s does not belong to active choices of question %s", id, q.QuestionID())
		}
	}
	return nil
}
