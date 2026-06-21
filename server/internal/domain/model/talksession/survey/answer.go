package survey

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/neko-dream/api/internal/domain/model/shared"
)

// AnswerValue は質問への回答値（JSONB として永続化）。
// 質問タイプによって有効なフィールドが異なる：
//   - free_text:                  Text
//   - single_choice / dropdown:   SelectedChoiceIDs（1件）、または Other（allowOther=true）
//   - multi_choice:               SelectedChoiceIDs（1件以上）、加えて Other
//   - rating:                     Rating
//   - date:                       Date (YYYY-MM-DD)
type AnswerValue struct {
	Text              *string              `json:"text,omitempty"`
	SelectedChoiceIDs []shared.UUID[Choice] `json:"selectedChoiceIDs,omitempty"`
	Other             *string              `json:"other,omitempty"`
	Rating            *int                 `json:"rating,omitempty"`
	Date              *string              `json:"date,omitempty"`
}

func (v AnswerValue) IsEmpty() bool {
	if v.Text != nil && *v.Text != "" {
		return false
	}
	if len(v.SelectedChoiceIDs) > 0 {
		return false
	}
	if v.Other != nil && *v.Other != "" {
		return false
	}
	if v.Rating != nil {
		return false
	}
	if v.Date != nil && *v.Date != "" {
		return false
	}
	return true
}

func (v AnswerValue) Marshal() ([]byte, error) {
	return json.Marshal(v)
}

func UnmarshalAnswerValue(data []byte) (AnswerValue, error) {
	var v AnswerValue
	if len(data) == 0 {
		return v, errors.New("empty answer value")
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return v, err
	}
	return v, nil
}

type Answer struct {
	answerID     shared.UUID[Answer]
	submissionID shared.UUID[Submission]
	questionID   shared.UUID[Question]
	value        AnswerValue
	isEncrypted  bool
	createdAt    time.Time
}

func NewAnswer(
	answerID shared.UUID[Answer],
	submissionID shared.UUID[Submission],
	questionID shared.UUID[Question],
	value AnswerValue,
) *Answer {
	return &Answer{
		answerID:     answerID,
		submissionID: submissionID,
		questionID:   questionID,
		value:        value,
		createdAt:    time.Now(),
	}
}

func ReconstructAnswer(
	answerID shared.UUID[Answer],
	submissionID shared.UUID[Submission],
	questionID shared.UUID[Question],
	value AnswerValue,
	isEncrypted bool,
	createdAt time.Time,
) *Answer {
	return &Answer{
		answerID:     answerID,
		submissionID: submissionID,
		questionID:   questionID,
		value:        value,
		isEncrypted:  isEncrypted,
		createdAt:    createdAt,
	}
}

func (a *Answer) AnswerID() shared.UUID[Answer]         { return a.answerID }
func (a *Answer) SubmissionID() shared.UUID[Submission] { return a.submissionID }
func (a *Answer) QuestionID() shared.UUID[Question]     { return a.questionID }
func (a *Answer) Value() AnswerValue                    { return a.value }
func (a *Answer) IsEncrypted() bool                     { return a.isEncrypted }
func (a *Answer) CreatedAt() time.Time                  { return a.createdAt }

// SetEncryptedValue は保存直前に呼び出される。
// PII フラグがある質問の場合、暗号文の値と is_encrypted=true をセットする。
func (a *Answer) SetEncryptedValue(v AnswerValue) {
	a.value = v
	a.isEncrypted = true
}

// SetDecryptedValue は読み出し後に呼び出される。
// is_encrypted=true の Answer に対し、復号した値をセットする。
// 復号後も isEncrypted フラグは true のまま（DB 上の状態を表す）。
func (a *Answer) SetDecryptedValue(v AnswerValue) {
	a.value = v
}
