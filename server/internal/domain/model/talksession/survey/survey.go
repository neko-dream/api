// Package survey はトークセッション参加前に必須回答となるアンケートのドメインモデル。
package survey

import (
	"time"

	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
)

// Survey は talk_session に紐付くアンケート定義（集約ルート）。
type Survey struct {
	surveyID      shared.UUID[Survey]
	talkSessionID shared.UUID[talksession.TalkSession]
	questions     []*Question
	createdAt     time.Time
	updatedAt     time.Time
}

func NewSurvey(
	surveyID shared.UUID[Survey],
	talkSessionID shared.UUID[talksession.TalkSession],
) *Survey {
	now := time.Now()
	return &Survey{
		surveyID:      surveyID,
		talkSessionID: talkSessionID,
		createdAt:     now,
		updatedAt:     now,
	}
}

// ReconstructSurvey はDBから復元する。
func ReconstructSurvey(
	surveyID shared.UUID[Survey],
	talkSessionID shared.UUID[talksession.TalkSession],
	questions []*Question,
	createdAt, updatedAt time.Time,
) *Survey {
	return &Survey{
		surveyID:      surveyID,
		talkSessionID: talkSessionID,
		questions:     questions,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}

func (s *Survey) SurveyID() shared.UUID[Survey]                          { return s.surveyID }
func (s *Survey) TalkSessionID() shared.UUID[talksession.TalkSession]    { return s.talkSessionID }
func (s *Survey) Questions() []*Question                                 { return s.questions }
func (s *Survey) CreatedAt() time.Time                                   { return s.createdAt }
func (s *Survey) UpdatedAt() time.Time                                   { return s.updatedAt }

// ActiveQuestions は削除されていない質問のみ display_order 順で返す。
func (s *Survey) ActiveQuestions() []*Question {
	out := make([]*Question, 0, len(s.questions))
	for _, q := range s.questions {
		if !q.IsDeleted() {
			out = append(out, q)
		}
	}
	return out
}

func (s *Survey) FindQuestion(id shared.UUID[Question]) *Question {
	for _, q := range s.questions {
		if q.QuestionID() == id {
			return q
		}
	}
	return nil
}

func (s *Survey) AddQuestion(q *Question) {
	s.questions = append(s.questions, q)
	s.updatedAt = time.Now()
}

// RemoveQuestion は in-memory アグリゲートから質問を削除する（hard delete 後の整合性維持用）。
// DB 側の削除は呼び出し元の責任。
func (s *Survey) RemoveQuestion(id shared.UUID[Question]) {
	for i, q := range s.questions {
		if q.QuestionID() == id {
			s.questions = append(s.questions[:i], s.questions[i+1:]...)
			s.updatedAt = time.Now()
			return
		}
	}
}

func (s *Survey) Touch() {
	s.updatedAt = time.Now()
}
