package survey

import (
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
)

// Summary は survey の集計結果。
type Summary struct {
	SurveyID          shared.UUID[Survey]
	TalkSessionID     shared.UUID[talksession.TalkSession]
	SubmissionCount   int
	QuestionSummaries []*QuestionSummary
}

// QuestionSummary は質問ごとの集計結果。
// 質問タイプは single_choice / dropdown のみで、いずれも ChoiceCounts + OtherCount を持つ。
type QuestionSummary struct {
	QuestionID   shared.UUID[Question]
	Text         string
	Type         QuestionType
	AnswerCount  int
	ChoiceCounts []*ChoiceCount
	OtherCount   int
}

type ChoiceCount struct {
	ChoiceID shared.UUID[Choice]
	Text     string
	Count    int
}

// Aggregate は submissions を質問タイプに応じて集計する。
func Aggregate(s *Survey, submissions []*Submission) *Summary {
	out := &Summary{
		SurveyID:        s.SurveyID(),
		TalkSessionID:   s.TalkSessionID(),
		SubmissionCount: len(submissions),
	}
	active := s.ActiveQuestions()
	out.QuestionSummaries = make([]*QuestionSummary, 0, len(active))
	for _, q := range active {
		out.QuestionSummaries = append(out.QuestionSummaries, aggregateQuestion(q, submissions))
	}
	return out
}

func aggregateQuestion(q *Question, submissions []*Submission) *QuestionSummary {
	qs := &QuestionSummary{
		QuestionID: q.QuestionID(),
		Text:       q.Text(),
		Type:       q.Type(),
	}

	// この質問に対する全 Answer を抽出（呼び出し元で復号済みであること）
	var answers []AnswerValue
	for _, sub := range submissions {
		for _, a := range sub.Answers() {
			if a.QuestionID() == q.QuestionID() {
				answers = append(answers, a.Value())
				break
			}
		}
	}

	// single_choice / dropdown はいずれも選択肢系として集計する。
	qs.ChoiceCounts = buildChoiceCounts(q)
	idx := make(map[shared.UUID[Choice]]*ChoiceCount, len(qs.ChoiceCounts))
	for _, cc := range qs.ChoiceCounts {
		idx[cc.ChoiceID] = cc
	}
	for _, v := range answers {
		counted := false
		for _, cid := range v.SelectedChoiceIDs {
			if cc, ok := idx[cid]; ok {
				cc.Count++
				counted = true
			}
		}
		if v.Other != nil && *v.Other != "" {
			qs.OtherCount++
			counted = true
		}
		if counted {
			qs.AnswerCount++
		}
	}

	return qs
}

func buildChoiceCounts(q *Question) []*ChoiceCount {
	active := q.ActiveChoices()
	out := make([]*ChoiceCount, 0, len(active))
	for _, c := range active {
		out = append(out, &ChoiceCount{
			ChoiceID: c.ChoiceID(),
			Text:     c.Text(),
			Count:    0,
		})
	}
	return out
}
