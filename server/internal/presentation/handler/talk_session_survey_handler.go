package handler

import (
	"context"
	"errors"

	"github.com/neko-dream/api/internal/application/usecase/talksession_usecase"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/talksession/survey"
	"github.com/neko-dream/api/internal/presentation/oas"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

func (t *talkSessionHandler) GetTalkSessionSurvey(ctx context.Context, params oas.GetTalkSessionSurveyParams) (oas.GetTalkSessionSurveyRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.GetTalkSessionSurvey")
	defer span.End()

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	s, err := t.getSurveyQuery.Execute(ctx, talkSessionID)
	if err != nil {
		if errors.Is(err, messages.TalkSessionSurveyNotFound) {
			return &oas.GetTalkSessionSurveyNotFound{}, nil
		}
		return nil, err
	}
	return surveyToResponse(s), nil
}

func (t *talkSessionHandler) CreateTalkSessionSurvey(ctx context.Context, req *oas.CreateTalkSessionSurveyReq, params oas.CreateTalkSessionSurveyParams) (oas.CreateTalkSessionSurveyRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.CreateTalkSessionSurvey")
	defer span.End()

	if req == nil {
		return nil, messages.BadRequestError
	}

	authCtx, err := t.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	questions, err := convertCreateQuestions(req.GetQuestions())
	if err != nil {
		return nil, err
	}

	s, err := t.createSurveyCommand.Execute(ctx, talksession_usecase.CreateSurveyInput{
		TalkSessionID: talkSessionID,
		ActorUserID:   authCtx.UserID,
		Questions:     questions,
	})
	if err != nil {
		return nil, err
	}
	return surveyToResponse(s), nil
}

func (t *talkSessionHandler) UpdateTalkSessionSurvey(ctx context.Context, req *oas.UpdateTalkSessionSurveyReq, params oas.UpdateTalkSessionSurveyParams) (oas.UpdateTalkSessionSurveyRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.UpdateTalkSessionSurvey")
	defer span.End()

	if req == nil {
		return nil, messages.BadRequestError
	}

	authCtx, err := t.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	questions, err := convertUpdateQuestions(req.GetQuestions())
	if err != nil {
		return nil, err
	}

	deletedQIDs, err := shared.ParseUUIDs[survey.Question](req.GetDeletedQuestionIDs())
	if err != nil {
		return nil, messages.BadRequestError
	}
	deletedCIDs, err := shared.ParseUUIDs[survey.Choice](req.GetDeletedChoiceIDs())
	if err != nil {
		return nil, messages.BadRequestError
	}

	s, err := t.updateSurveyCommand.Execute(ctx, talksession_usecase.UpdateSurveyInput{
		TalkSessionID:      talkSessionID,
		ActorUserID:        authCtx.UserID,
		Questions:          questions,
		DeletedQuestionIDs: deletedQIDs,
		DeletedChoiceIDs:   deletedCIDs,
	})
	if err != nil {
		return nil, err
	}
	return surveyToResponse(s), nil
}

func (t *talkSessionHandler) SubmitTalkSessionSurvey(ctx context.Context, req *oas.SubmitTalkSessionSurveyReq, params oas.SubmitTalkSessionSurveyParams) (oas.SubmitTalkSessionSurveyRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.SubmitTalkSessionSurvey")
	defer span.End()

	if req == nil {
		return nil, messages.BadRequestError
	}

	authCtx, err := t.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	answers, err := convertSubmitAnswers(req.GetAnswers())
	if err != nil {
		return nil, messages.BadRequestError
	}

	submissionID, err := t.submitSurveyCommand.Execute(ctx, talksession_usecase.SubmitSurveyInput{
		TalkSessionID: talkSessionID,
		UserID:        authCtx.UserID,
		Answers:       answers,
	})
	if err != nil {
		return nil, err
	}

	return &oas.SubmitTalkSessionSurveyOK{SubmissionID: submissionID.String()}, nil
}

func (t *talkSessionHandler) ListTalkSessionSurveySubmissions(ctx context.Context, params oas.ListTalkSessionSurveySubmissionsParams) (oas.ListTalkSessionSurveySubmissionsRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.ListTalkSessionSurveySubmissions")
	defer span.End()

	authCtx, err := t.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	limit := 20
	if params.Limit.IsSet() {
		limit = int(params.Limit.Value)
	}
	offset := 0
	if params.Offset.IsSet() {
		offset = int(params.Offset.Value)
	}

	out, err := t.listSurveySubmissionsQuery.Execute(ctx, talksession_usecase.ListSurveySubmissionsInput{
		TalkSessionID: talkSessionID,
		ActorUserID:   authCtx.UserID,
		Limit:         limit,
		Offset:        offset,
	})
	if err != nil {
		return nil, err
	}

	submissions := make([]oas.TalkSessionSurveySubmission, 0, len(out.Submissions))
	for _, sub := range out.Submissions {
		submissions = append(submissions, submissionToResponse(sub))
	}
	return &oas.ListTalkSessionSurveySubmissionsOK{
		Submissions: submissions,
		TotalCount:  out.TotalCount,
	}, nil
}

func (t *talkSessionHandler) GetTalkSessionSurveySummary(ctx context.Context, params oas.GetTalkSessionSurveySummaryParams) (oas.GetTalkSessionSurveySummaryRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "talkSessionHandler.GetTalkSessionSurveySummary")
	defer span.End()

	authCtx, err := t.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	talkSessionID, err := shared.ParseUUID[talksession.TalkSession](params.TalkSessionID)
	if err != nil {
		return nil, messages.BadRequestError
	}

	summary, err := t.getSurveySummaryQuery.Execute(ctx, talksession_usecase.GetSurveySummaryInput{
		TalkSessionID: talkSessionID,
		ActorUserID:   authCtx.UserID,
	})
	if err != nil {
		return nil, err
	}
	return summaryToResponse(summary), nil
}

func summaryToResponse(s *survey.Summary) *oas.TalkSessionSurveySummary {
	qs := make([]oas.TalkSessionSurveyQuestionSummary, 0, len(s.QuestionSummaries))
	for _, q := range s.QuestionSummaries {
		choices := make([]oas.TalkSessionSurveyChoiceCount, 0, len(q.ChoiceCounts))
		for _, cc := range q.ChoiceCounts {
			choices = append(choices, oas.TalkSessionSurveyChoiceCount{
				ChoiceID: cc.ChoiceID.String(),
				Text:     cc.Text,
				Count:    int32(cc.Count),
			})
		}
		qs = append(qs, oas.TalkSessionSurveyQuestionSummary{
			QuestionID:   q.QuestionID.String(),
			Text:         q.Text,
			Type:         oas.TalkSessionSurveyQuestionType(q.Type),
			AnswerCount:  int32(q.AnswerCount),
			ChoiceCounts: choices,
			OtherCount:   int32(q.OtherCount),
		})
	}
	return &oas.TalkSessionSurveySummary{
		SurveyID:          s.SurveyID.String(),
		TalkSessionID:     s.TalkSessionID.String(),
		SubmissionCount:   int32(s.SubmissionCount),
		QuestionSummaries: qs,
	}
}

// --- 変換ヘルパ ---

func convertCreateQuestions(in []oas.TalkSessionSurveyQuestionInput) ([]talksession_usecase.SurveyQuestionInput, error) {
	out := make([]talksession_usecase.SurveyQuestionInput, 0, len(in))
	for _, qi := range in {
		qt, err := survey.ParseQuestionType(string(qi.GetType()))
		if err != nil {
			return nil, messages.TalkSessionSurveyAnswerInvalid
		}
		choices, err := convertChoiceInputsForCreate(qi.GetChoices())
		if err != nil {
			return nil, err
		}
		out = append(out, talksession_usecase.SurveyQuestionInput{
			Text:         qi.GetText(),
			Type:         qt,
			IsRequired:   qi.GetIsRequired(),
			AllowOther:   qi.GetAllowOther(),
			MaxLength:    optNilInt32ToIntPtr(qi.GetMaxLength()),
			DisplayOrder: int(qi.GetDisplayOrder()),
			Choices:      choices,
		})
	}
	return out, nil
}

func convertChoiceInputsForCreate(in []oas.TalkSessionSurveyChoiceInput) ([]talksession_usecase.SurveyChoiceInput, error) {
	out := make([]talksession_usecase.SurveyChoiceInput, 0, len(in))
	for _, ci := range in {
		out = append(out, talksession_usecase.SurveyChoiceInput{
			Text:         ci.GetText(),
			DisplayOrder: int(ci.GetDisplayOrder()),
		})
	}
	return out, nil
}

func convertUpdateQuestions(in []oas.TalkSessionSurveyQuestionInput) ([]talksession_usecase.UpdateSurveyQuestionInput, error) {
	out := make([]talksession_usecase.UpdateSurveyQuestionInput, 0, len(in))
	for _, qi := range in {
		qt, err := survey.ParseQuestionType(string(qi.GetType()))
		if err != nil {
			return nil, messages.TalkSessionSurveyAnswerInvalid
		}
		var qid *shared.UUID[survey.Question]
		if v, ok := optNilStringValue(qi.GetQuestionID()); ok {
			parsed, err := shared.ParseUUID[survey.Question](v)
			if err != nil {
				return nil, messages.BadRequestError
			}
			qid = &parsed
		}

		choices, err := convertChoiceInputsForUpdate(qi.GetChoices())
		if err != nil {
			return nil, err
		}

		out = append(out, talksession_usecase.UpdateSurveyQuestionInput{
			QuestionID:   qid,
			Text:         qi.GetText(),
			Type:         qt,
			IsRequired:   qi.GetIsRequired(),
			AllowOther:   qi.GetAllowOther(),
			MaxLength:    optNilInt32ToIntPtr(qi.GetMaxLength()),
			DisplayOrder: int(qi.GetDisplayOrder()),
			Choices:      choices,
		})
	}
	return out, nil
}

func convertChoiceInputsForUpdate(in []oas.TalkSessionSurveyChoiceInput) ([]talksession_usecase.UpdateSurveyChoiceInput, error) {
	out := make([]talksession_usecase.UpdateSurveyChoiceInput, 0, len(in))
	for _, ci := range in {
		var cid *shared.UUID[survey.Choice]
		if v, ok := optNilStringValue(ci.GetChoiceID()); ok {
			parsed, err := shared.ParseUUID[survey.Choice](v)
			if err != nil {
				return nil, messages.BadRequestError
			}
			cid = &parsed
		}
		out = append(out, talksession_usecase.UpdateSurveyChoiceInput{
			ChoiceID:     cid,
			Text:         ci.GetText(),
			DisplayOrder: int(ci.GetDisplayOrder()),
		})
	}
	return out, nil
}

func convertSubmitAnswers(in []oas.TalkSessionSurveyAnswerInput) ([]talksession_usecase.SurveyAnswerInput, error) {
	out := make([]talksession_usecase.SurveyAnswerInput, 0, len(in))
	for _, ai := range in {
		qid, err := shared.ParseUUID[survey.Question](ai.GetQuestionID())
		if err != nil {
			return nil, err
		}
		selected, err := shared.ParseUUIDs[survey.Choice](ai.GetSelectedChoiceIDs())
		if err != nil {
			return nil, err
		}
		out = append(out, talksession_usecase.SurveyAnswerInput{
			QuestionID:        qid,
			SelectedChoiceIDs: selected,
			Other:             optNilStringToPtr(ai.GetOther()),
		})
	}
	return out, nil
}

func surveyToResponse(s *survey.Survey) *oas.TalkSessionSurvey {
	res := &oas.TalkSessionSurvey{
		SurveyID:      s.SurveyID().String(),
		TalkSessionID: s.TalkSessionID().String(),
	}
	active := s.ActiveQuestions()
	questions := make([]oas.TalkSessionSurveyQuestion, 0, len(active))
	for _, q := range active {
		questions = append(questions, questionToResponse(q))
	}
	res.Questions = questions
	return res
}

func questionToResponse(q *survey.Question) oas.TalkSessionSurveyQuestion {
	active := q.ActiveChoices()
	choices := make([]oas.TalkSessionSurveyChoice, 0, len(active))
	for _, c := range active {
		choices = append(choices, oas.TalkSessionSurveyChoice{
			ChoiceID:     c.ChoiceID().String(),
			Text:         c.Text(),
			DisplayOrder: int32(c.DisplayOrder()),
		})
	}
	return oas.TalkSessionSurveyQuestion{
		QuestionID:   q.QuestionID().String(),
		Text:         q.Text(),
		Type:         oas.TalkSessionSurveyQuestionType(q.Type()),
		IsRequired:   q.IsRequired(),
		AllowOther:   q.AllowOther(),
		MaxLength:    intPtrToOptNilInt32(q.MaxLength()),
		DisplayOrder: int32(q.DisplayOrder()),
		Choices:      choices,
	}
}

func submissionToResponse(sub *survey.Submission) oas.TalkSessionSurveySubmission {
	answers := make([]oas.TalkSessionSurveyAnswer, 0, len(sub.Answers()))
	for _, a := range sub.Answers() {
		v := a.Value()
		selected := make([]string, 0, len(v.SelectedChoiceIDs))
		for _, cid := range v.SelectedChoiceIDs {
			selected = append(selected, cid.String())
		}
		answers = append(answers, oas.TalkSessionSurveyAnswer{
			QuestionID:        a.QuestionID().String(),
			SelectedChoiceIDs: selected,
			Other:             utils.ToOptNil[oas.OptNilString](v.Other),
		})
	}
	return oas.TalkSessionSurveySubmission{
		SubmissionID: sub.SubmissionID().String(),
		UserID:       sub.UserID().String(),
		SubmittedAt:  sub.SubmittedAt(),
		Answers:      answers,
	}
}

// OptNil 系の入力変換（utils.ToOptNil の逆方向ヘルパが無いためここに置く）

func optNilStringValue(v oas.OptNilString) (string, bool) {
	if !v.IsSet() || v.IsNull() {
		return "", false
	}
	return v.Value, true
}

func optNilStringToPtr(v oas.OptNilString) *string {
	if !v.IsSet() || v.IsNull() {
		return nil
	}
	s := v.Value
	return &s
}

func optNilInt32ToIntPtr(v oas.OptNilInt32) *int {
	if !v.IsSet() || v.IsNull() {
		return nil
	}
	i := int(v.Value)
	return &i
}

// intPtrToOptNilInt32 はドメインの *int を OAS の OptNilInt32 に変換する。
// utils.ToOptNil は OptNilInt32 を扱わないため、ogen 生成の NewOptNilInt32 を使う。
func intPtrToOptNilInt32(p *int) oas.OptNilInt32 {
	if p == nil {
		return oas.OptNilInt32{}
	}
	return oas.NewOptNilInt32(int32(*p))
}
