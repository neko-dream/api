package survey

import (
	"errors"
	"slices"
)

// QuestionType はアンケート質問の種類を表す。
type QuestionType string

const (
	QuestionTypeSingleChoice QuestionType = "single_choice"
	QuestionTypeDropdown     QuestionType = "dropdown"
)

var allQuestionTypes = []QuestionType{
	QuestionTypeSingleChoice,
	QuestionTypeDropdown,
}

func (t QuestionType) IsValid() bool {
	return slices.Contains(allQuestionTypes, t)
}

// HasChoices は選択肢を持つタイプか返す。
// 現在サポートする質問タイプ（single_choice / dropdown）はいずれも選択肢を持つ。
func (t QuestionType) HasChoices() bool {
	return t.IsValid()
}

func ParseQuestionType(s string) (QuestionType, error) {
	qt := QuestionType(s)
	if !qt.IsValid() {
		return "", errors.New("invalid question type: " + s)
	}
	return qt, nil
}
