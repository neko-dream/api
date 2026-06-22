package survey

import (
	"time"

	"github.com/neko-dream/api/internal/domain/model/shared"
)

type Choice struct {
	choiceID     shared.UUID[Choice]
	questionID   shared.UUID[Question]
	text         string
	displayOrder int
	deletedAt    *time.Time
	createdAt    time.Time
	updatedAt    time.Time
}

func NewChoice(
	choiceID shared.UUID[Choice],
	questionID shared.UUID[Question],
	text string,
	displayOrder int,
) *Choice {
	now := time.Now()
	return &Choice{
		choiceID:     choiceID,
		questionID:   questionID,
		text:         text,
		displayOrder: displayOrder,
		createdAt:    now,
		updatedAt:    now,
	}
}

// ReconstructChoice はDBから復元する。
func ReconstructChoice(
	choiceID shared.UUID[Choice],
	questionID shared.UUID[Question],
	text string,
	displayOrder int,
	deletedAt *time.Time,
	createdAt, updatedAt time.Time,
) *Choice {
	return &Choice{
		choiceID:     choiceID,
		questionID:   questionID,
		text:         text,
		displayOrder: displayOrder,
		deletedAt:    deletedAt,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

func (c *Choice) ChoiceID() shared.UUID[Choice]     { return c.choiceID }
func (c *Choice) QuestionID() shared.UUID[Question] { return c.questionID }
func (c *Choice) Text() string                      { return c.text }
func (c *Choice) DisplayOrder() int                  { return c.displayOrder }
func (c *Choice) DeletedAt() *time.Time             { return c.deletedAt }
func (c *Choice) IsDeleted() bool                   { return c.deletedAt != nil }
func (c *Choice) CreatedAt() time.Time              { return c.createdAt }
func (c *Choice) UpdatedAt() time.Time              { return c.updatedAt }

func (c *Choice) SetText(text string) {
	c.text = text
	c.updatedAt = time.Now()
}

func (c *Choice) SetDisplayOrder(order int) {
	c.displayOrder = order
	c.updatedAt = time.Now()
}

// SoftDelete は論理削除する。
func (c *Choice) SoftDelete() {
	if c.deletedAt != nil {
		return
	}
	now := time.Now()
	c.deletedAt = &now
	c.updatedAt = now
}
