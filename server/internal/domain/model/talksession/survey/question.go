package survey

import (
	"time"

	"github.com/neko-dream/api/internal/domain/model/shared"
)

type Question struct {
	questionID   shared.UUID[Question]
	surveyID     shared.UUID[Survey]
	text         string
	questionType QuestionType
	isRequired   bool
	allowOther   bool
	maxLength    *int
	minValue     *int
	maxValue     *int
	minLabel     *string
	maxLabel     *string
	displayOrder int
	containsPii  bool
	choices      []*Choice
	deletedAt    *time.Time
	createdAt    time.Time
	updatedAt    time.Time
}

type QuestionConfig struct {
	IsRequired   bool
	AllowOther   bool
	MaxLength    *int
	MinValue     *int
	MaxValue     *int
	MinLabel     *string
	MaxLabel     *string
	DisplayOrder int
	ContainsPii  bool
}

func NewQuestion(
	questionID shared.UUID[Question],
	surveyID shared.UUID[Survey],
	text string,
	questionType QuestionType,
	cfg QuestionConfig,
) *Question {
	now := time.Now()
	return &Question{
		questionID:   questionID,
		surveyID:     surveyID,
		text:         text,
		questionType: questionType,
		isRequired:   cfg.IsRequired,
		allowOther:   cfg.AllowOther,
		maxLength:    cfg.MaxLength,
		minValue:     cfg.MinValue,
		maxValue:     cfg.MaxValue,
		minLabel:     cfg.MinLabel,
		maxLabel:     cfg.MaxLabel,
		displayOrder: cfg.DisplayOrder,
		containsPii:  cfg.ContainsPii,
		createdAt:    now,
		updatedAt:    now,
	}
}

// ReconstructQuestion はDBから復元する。
func ReconstructQuestion(
	questionID shared.UUID[Question],
	surveyID shared.UUID[Survey],
	text string,
	questionType QuestionType,
	isRequired, allowOther bool,
	maxLength, minValue, maxValue *int,
	minLabel, maxLabel *string,
	displayOrder int,
	containsPii bool,
	choices []*Choice,
	deletedAt *time.Time,
	createdAt, updatedAt time.Time,
) *Question {
	return &Question{
		questionID:   questionID,
		surveyID:     surveyID,
		text:         text,
		questionType: questionType,
		isRequired:   isRequired,
		allowOther:   allowOther,
		maxLength:    maxLength,
		minValue:     minValue,
		maxValue:     maxValue,
		minLabel:     minLabel,
		maxLabel:     maxLabel,
		displayOrder: displayOrder,
		containsPii:  containsPii,
		choices:      choices,
		deletedAt:    deletedAt,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

func (q *Question) QuestionID() shared.UUID[Question] { return q.questionID }
func (q *Question) SurveyID() shared.UUID[Survey]     { return q.surveyID }
func (q *Question) Text() string                      { return q.text }
func (q *Question) Type() QuestionType                { return q.questionType }
func (q *Question) IsRequired() bool                  { return q.isRequired }
func (q *Question) AllowOther() bool                  { return q.allowOther }
func (q *Question) MaxLength() *int                   { return q.maxLength }
func (q *Question) MinValue() *int                    { return q.minValue }
func (q *Question) MaxValue() *int                    { return q.maxValue }
func (q *Question) MinLabel() *string                 { return q.minLabel }
func (q *Question) MaxLabel() *string                 { return q.maxLabel }
func (q *Question) DisplayOrder() int                 { return q.displayOrder }
func (q *Question) ContainsPii() bool                 { return q.containsPii }
func (q *Question) Choices() []*Choice                { return q.choices }
func (q *Question) DeletedAt() *time.Time             { return q.deletedAt }
func (q *Question) IsDeleted() bool                   { return q.deletedAt != nil }
func (q *Question) CreatedAt() time.Time              { return q.createdAt }
func (q *Question) UpdatedAt() time.Time              { return q.updatedAt }

// ActiveChoices は削除されていない選択肢のみ返す。
func (q *Question) ActiveChoices() []*Choice {
	out := make([]*Choice, 0, len(q.choices))
	for _, c := range q.choices {
		if !c.IsDeleted() {
			out = append(out, c)
		}
	}
	return out
}

func (q *Question) FindChoice(id shared.UUID[Choice]) *Choice {
	for _, c := range q.choices {
		if c.ChoiceID() == id {
			return c
		}
	}
	return nil
}

func (q *Question) AddChoice(c *Choice) {
	q.choices = append(q.choices, c)
	q.updatedAt = time.Now()
}

// RemoveChoice は in-memory から選択肢を削除する（hard delete 後の整合性維持用）。
func (q *Question) RemoveChoice(id shared.UUID[Choice]) {
	for i, c := range q.choices {
		if c.ChoiceID() == id {
			q.choices = append(q.choices[:i], q.choices[i+1:]...)
			q.updatedAt = time.Now()
			return
		}
	}
}

func (q *Question) SetText(text string) {
	q.text = text
	q.updatedAt = time.Now()
}

func (q *Question) SetType(t QuestionType) {
	q.questionType = t
	q.updatedAt = time.Now()
}

func (q *Question) SetConfig(cfg QuestionConfig) {
	q.isRequired = cfg.IsRequired
	q.allowOther = cfg.AllowOther
	q.maxLength = cfg.MaxLength
	q.minValue = cfg.MinValue
	q.maxValue = cfg.MaxValue
	q.minLabel = cfg.MinLabel
	q.maxLabel = cfg.MaxLabel
	q.displayOrder = cfg.DisplayOrder
	q.containsPii = cfg.ContainsPii
	q.updatedAt = time.Now()
}

func (q *Question) SoftDelete() {
	if q.deletedAt != nil {
		return
	}
	now := time.Now()
	q.deletedAt = &now
	q.updatedAt = now
}
