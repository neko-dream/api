package handler

import (
	"context"

	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/presentation/oas"
)

// アンケート機能は別 PR で実装予定。ここではビルドを通すためのスタブのみ。

func (t *talkSessionHandler) GetTalkSessionSurvey(ctx context.Context, params oas.GetTalkSessionSurveyParams) (oas.GetTalkSessionSurveyRes, error) {
	return nil, messages.InternalServerError
}

func (t *talkSessionHandler) CreateTalkSessionSurvey(ctx context.Context, req *oas.CreateTalkSessionSurveyReq, params oas.CreateTalkSessionSurveyParams) (oas.CreateTalkSessionSurveyRes, error) {
	return nil, messages.InternalServerError
}

func (t *talkSessionHandler) UpdateTalkSessionSurvey(ctx context.Context, req *oas.UpdateTalkSessionSurveyReq, params oas.UpdateTalkSessionSurveyParams) (oas.UpdateTalkSessionSurveyRes, error) {
	return nil, messages.InternalServerError
}

func (t *talkSessionHandler) SubmitTalkSessionSurvey(ctx context.Context, req *oas.SubmitTalkSessionSurveyReq, params oas.SubmitTalkSessionSurveyParams) (oas.SubmitTalkSessionSurveyRes, error) {
	return nil, messages.InternalServerError
}

func (t *talkSessionHandler) ListTalkSessionSurveySubmissions(ctx context.Context, params oas.ListTalkSessionSurveySubmissionsParams) (oas.ListTalkSessionSurveySubmissionsRes, error) {
	return nil, messages.InternalServerError
}
