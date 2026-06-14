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

	deletedQIDs, err := parseUUIDList[survey.Question](req.GetDeletedQuestionIDs())
	if err != nil {
		return nil, messages.BadRequestError
	}
	deletedCIDs, err := parseUUIDList[survey.Choice](req.GetDeletedChoiceIDs())
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
			MinValue:     optNilInt32ToIntPtr(qi.GetMinValue()),
			MaxValue:     optNilInt32ToIntPtr(qi.GetMaxValue()),
			MinLabel:     optNilStringToPtr(qi.GetMinLabel()),
			MaxLabel:     optNilStringToPtr(qi.GetMaxLabel()),
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
			MinValue:     optNilInt32ToIntPtr(qi.GetMinValue()),
			MaxValue:     optNilInt32ToIntPtr(qi.GetMaxValue()),
			MinLabel:     optNilStringToPtr(qi.GetMinLabel()),
			MaxLabel:     optNilStringToPtr(qi.GetMaxLabel()),
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
		selected, err := parseUUIDList[survey.Choice](ai.GetSelectedChoiceIDs())
		if err != nil {
			return nil, err
		}
		out = append(out, talksession_usecase.SurveyAnswerInput{
			QuestionID:        qid,
			Text:              optNilStringToPtr(ai.GetText()),
			SelectedChoiceIDs: selected,
			Other:             optNilStringToPtr(ai.GetOther()),
			Rating:            optNilInt32ToIntPtr(ai.GetRating()),
			Date:              optNilStringToPtr(ai.GetDate()),
		})
	}
	return out, nil
}

func parseUUIDList[T any](in []string) ([]shared.UUID[T], error) {
	out := make([]shared.UUID[T], 0, len(in))
	for _, s := range in {
		u, err := shared.ParseUUID[T](s)
		if err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, nil
}

func surveyToResponse(s *survey.Survey) *oas.TalkSessionSurvey {
	res := &oas.TalkSessionSurvey{
		SurveyID:      s.SurveyID().String(),
		TalkSessionID: s.TalkSessionID().String(),
	}
	questions := make([]oas.TalkSessionSurveyQuestion, 0, len(s.ActiveQuestions()))
	for _, q := range s.ActiveQuestions() {
		questions = append(questions, questionToResponse(q))
	}
	res.Questions = questions
	return res
}

func questionToResponse(q *survey.Question) oas.TalkSessionSurveyQuestion {
	choices := make([]oas.TalkSessionSurveyChoice, 0, len(q.ActiveChoices()))
	for _, c := range q.ActiveChoices() {
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
		MinValue:     intPtrToOptNilInt32(q.MinValue()),
		MaxValue:     intPtrToOptNilInt32(q.MaxValue()),
		MinLabel:     utils.ToOptNil[oas.OptNilString](q.MinLabel()),
		MaxLabel:     utils.ToOptNil[oas.OptNilString](q.MaxLabel()),
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
			Text:              utils.ToOptNil[oas.OptNilString](v.Text),
			SelectedChoiceIDs: selected,
			Other:             utils.ToOptNil[oas.OptNilString](v.Other),
			Rating:            intPtrToOptNilInt32(v.Rating),
			Date:              utils.ToOptNil[oas.OptNilString](v.Date),
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
