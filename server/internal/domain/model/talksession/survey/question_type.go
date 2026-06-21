package survey

import (
	"errors"
	"slices"
)

// QuestionType はアンケート質問の種類を表す。
type QuestionType string

const (
	QuestionTypeFreeText     QuestionType = "free_text"
	QuestionTypeSingleChoice QuestionType = "single_choice"
	QuestionTypeMultiChoice  QuestionType = "multi_choice"
	QuestionTypeRating       QuestionType = "rating"
	QuestionTypeDate         QuestionType = "date"
	QuestionTypeDropdown     QuestionType = "dropdown"
)

var allQuestionTypes = []QuestionType{
	QuestionTypeFreeText,
	QuestionTypeSingleChoice,
	QuestionTypeMultiChoice,
	QuestionTypeRating,
	QuestionTypeDate,
	QuestionTypeDropdown,
}

func (t QuestionType) IsValid() bool {
	return slices.Contains(allQuestionTypes, t)
}

// HasChoices は選択肢を持つタイプか返す。
func (t QuestionType) HasChoices() bool {
	switch t {
	case QuestionTypeSingleChoice, QuestionTypeMultiChoice, QuestionTypeDropdown:
		return true
	}
	return false
}

func ParseQuestionType(s string) (QuestionType, error) {
	qt := QuestionType(s)
	if !qt.IsValid() {
		return "", errors.New("invalid question type: " + s)
	}
	return qt, nil
}
