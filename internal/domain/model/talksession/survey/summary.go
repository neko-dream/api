package survey

import (
	"sort"

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
// 質問タイプによって意味のあるフィールドが異なる：
//   - free_text / date:                    AnswerCount のみ
//   - single_choice / multi_choice / dropdown: ChoiceCounts + OtherCount
//   - rating:                              RatingAverage + RatingDistribution
type QuestionSummary struct {
	QuestionID         shared.UUID[Question]
	Text               string
	Type               QuestionType
	ContainsPii        bool
	AnswerCount        int
	ChoiceCounts       []*ChoiceCount
	OtherCount         int
	RatingAverage      *float64
	RatingDistribution []*RatingBucket
}

type ChoiceCount struct {
	ChoiceID shared.UUID[Choice]
	Text     string
	Count    int
}

type RatingBucket struct {
	Value int
	Count int
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
		QuestionID:  q.QuestionID(),
		Text:        q.Text(),
		Type:        q.Type(),
		ContainsPii: q.ContainsPii(),
	}

	// この質問に対する全 Answer を抽出
	var answers []AnswerValue
	for _, sub := range submissions {
		for _, a := range sub.Answers() {
			if a.QuestionID() == q.QuestionID() {
				answers = append(answers, a.Value())
				break
			}
		}
	}

	switch q.Type() {
	case QuestionTypeFreeText, QuestionTypeDate:
		for _, v := range answers {
			if !v.IsEmpty() {
				qs.AnswerCount++
			}
		}

	case QuestionTypeSingleChoice, QuestionTypeMultiChoice, QuestionTypeDropdown:
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

	case QuestionTypeRating:
		bucket := make(map[int]int)
		sum := 0
		for _, v := range answers {
			if v.Rating == nil {
				continue
			}
			qs.AnswerCount++
			sum += *v.Rating
			bucket[*v.Rating]++
		}
		if qs.AnswerCount > 0 {
			avg := float64(sum) / float64(qs.AnswerCount)
			qs.RatingAverage = &avg
		}
		qs.RatingDistribution = buildRatingDistribution(q, bucket)
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

// buildRatingDistribution は minValue〜maxValue 範囲（指定があれば）と
// 実際に出現した値の和集合を作り、value 昇順で返す。
func buildRatingDistribution(q *Question, bucket map[int]int) []*RatingBucket {
	values := make(map[int]struct{})
	if q.MinValue() != nil && q.MaxValue() != nil && *q.MinValue() <= *q.MaxValue() {
		for v := *q.MinValue(); v <= *q.MaxValue(); v++ {
			values[v] = struct{}{}
		}
	}
	for v := range bucket {
		values[v] = struct{}{}
	}
	sorted := make([]int, 0, len(values))
	for v := range values {
		sorted = append(sorted, v)
	}
	sort.Ints(sorted)
	out := make([]*RatingBucket, 0, len(sorted))
	for _, v := range sorted {
		out = append(out, &RatingBucket{Value: v, Count: bucket[v]})
	}
	return out
}
