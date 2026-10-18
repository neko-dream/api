package admin_usecase

import (
	"context"
	"log/slog"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"braces.dev/errtrace"
	"github.com/neko-dream/api/internal/application/usecase/analysis_usecase"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/talksession/survey"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"go.opentelemetry.io/otel"
)

const (
	editTalkSessionThemeMaxLength       = 100
	editTalkSessionDescriptionMaxLength = 40000
)

// EditTalkSession 管理者がセッションを編集する。オーナー以外のセッションも編集できる。
type EditTalkSession interface {
	Execute(context.Context, EditTalkSessionInput) error
}

// EditTalkSessionInput は編集後の全項目を受け取る。変更しない項目には現在の値を渡す。
type EditTalkSessionInput struct {
	TalkSessionID    shared.UUID[talksession.TalkSession]
	Theme            string
	Description      string
	ScheduledEndTime time.Time
	Prefecture       string // 空文字は未設定
	City             string // 空文字は未設定
	HideTop          bool
	HideReport       bool
	DisableAnalysis  bool
	Restrictions     []string
}

type editTalkSessionTxRunner interface {
	ExecTx(ctx context.Context, fn func(context.Context) error) error
}

type editTalkSession struct {
	talkSessions      talksession.TalkSessionRepository
	surveys           survey.Repository
	analysisScheduler analysis_usecase.AnalysisScheduler
	txRunner          editTalkSessionTxRunner
}

func NewEditTalkSession(
	talkSessions talksession.TalkSessionRepository,
	surveys survey.Repository,
	analysisScheduler analysis_usecase.AnalysisScheduler,
	dbManager *db.DBManager,
) EditTalkSession {
	return &editTalkSession{
		talkSessions:      talkSessions,
		surveys:           surveys,
		analysisScheduler: analysisScheduler,
		txRunner:          dbManager,
	}
}

func (e *editTalkSession) Execute(ctx context.Context, input EditTalkSessionInput) error {
	ctx, span := otel.Tracer("admin_usecase").Start(ctx, "editTalkSession.Execute")
	defer span.End()

	claim := admin.GetAdminSession(ctx)
	if claim == nil {
		slog.WarnContext(ctx, "talk session edit denied", "talk_session_id", input.TalkSessionID.String())
		return messages.ForbiddenError
	}

	if err := e.edit(ctx, input, claim.Sub); err != nil {
		slog.WarnContext(ctx, "talk session edit by admin failed", "admin_user_id", claim.Sub, "talk_session_id", input.TalkSessionID.String())
		return err
	}
	return nil
}

func (e *editTalkSession) edit(ctx context.Context, input EditTalkSessionInput, adminUserID string) error {
	if strings.TrimSpace(input.Theme) == "" {
		return messages.BadRequestError
	}
	if utf8.RuneCountInString(input.Theme) > editTalkSessionThemeMaxLength {
		return messages.TalkSessionThemeTooLong
	}
	if utf8.RuneCountInString(input.Description) > editTalkSessionDescriptionMaxLength {
		return messages.TalkSessionDescriptionTooLong
	}

	ts, err := e.talkSessions.FindByID(ctx, input.TalkSessionID)
	if err != nil || ts == nil {
		return messages.TalkSessionNotFound
	}

	// 終了済みのセッションでも他の項目を直せるよう、終了日時は変えたときだけ検証する
	scheduledEndTimeChanged := !ts.ScheduledEndTime().Equal(input.ScheduledEndTime)
	if scheduledEndTimeChanged && input.ScheduledEndTime.Before(clock.Now(ctx)) {
		return messages.InvalidScheduledEndTime
	}

	// アンケートがないセッションに回答済み制限を付けると、誰も参加できなくなる
	if slices.Contains(input.Restrictions, string(talksession.SurveyAnswered)) {
		s, err := e.surveys.FindBySessionID(ctx, input.TalkSessionID)
		if err != nil {
			return errtrace.Wrap(err)
		}
		if s == nil {
			return messages.BadRequestError
		}
	}

	before := ts.RestrictionList()
	changed := changedTalkSessionFields(ts, input, scheduledEndTimeChanged)
	wasAnalysisDisabled := ts.DisableAnalysis()

	if err := ts.UpdateRestrictions(ctx, input.Restrictions); err != nil {
		return messages.BadRequestError
	}
	if !slices.Equal(restrictionKeys(before), restrictionKeys(ts.RestrictionList())) {
		changed = append(changed, "restrictions")
	}

	ts.ChangeTheme(input.Theme)
	ts.ChangeDescription(&input.Description)
	ts.ChangeScheduledEndTime(input.ScheduledEndTime)
	ts.ChangePrefecture(emptyToNil(input.Prefecture))
	ts.ChangeCity(emptyToNil(input.City))
	ts.ChangeHideTop(input.HideTop)
	ts.SetReportVisibility(input.HideReport)
	ts.SetDisableAnalysis(input.DisableAnalysis)

	if err := e.txRunner.ExecTx(ctx, func(ctx context.Context) error {
		return e.talkSessions.Update(ctx, ts)
	}); err != nil {
		return errtrace.Wrap(err)
	}

	// 監査ログには変更した項目名だけを残し、値（本文など）は残さない
	slog.InfoContext(ctx, "talk session edited by admin",
		"admin_user_id", adminUserID,
		"talk_session_id", input.TalkSessionID.String(),
		"changed_fields", changed,
	)

	if wasAnalysisDisabled && !input.DisableAnalysis {
		e.analysisScheduler.ScheduleAnalysis(ctx, input.TalkSessionID, analysis_usecase.ScheduleImmediate)
	}
	return nil
}

func changedTalkSessionFields(ts *talksession.TalkSession, input EditTalkSessionInput, scheduledEndTimeChanged bool) []string {
	var changed []string
	add := func(name string, isChanged bool) {
		if isChanged {
			changed = append(changed, name)
		}
	}
	add("theme", ts.Theme() != input.Theme)
	add("description", nilToEmpty(ts.Description()) != input.Description)
	add("scheduled_end_time", scheduledEndTimeChanged)
	add("prefecture", nilToEmpty(ts.Prefecture()) != input.Prefecture)
	add("city", nilToEmpty(ts.City()) != input.City)
	add("hide_top", ts.HideTop() != input.HideTop)
	add("hide_report", ts.HideReport() != input.HideReport)
	add("disable_analysis", ts.DisableAnalysis() != input.DisableAnalysis)
	return changed
}

func restrictionKeys(restrictions talksession.Restrictions) []string {
	keys := make([]string, 0, len(restrictions))
	for _, r := range restrictions {
		keys = append(keys, string(r))
	}
	slices.Sort(keys)
	return keys
}

func emptyToNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func nilToEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
