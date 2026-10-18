package admin_usecase

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/neko-dream/api/internal/application/usecase/analysis_usecase"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/talksession/survey"
	"github.com/neko-dream/api/internal/domain/model/user"
)

var editTestNow = time.Date(2026, 10, 15, 12, 0, 0, 0, time.UTC)

type editTestTalkSessions struct {
	talksession.TalkSessionRepository
	session *talksession.TalkSession
	saved   *talksession.TalkSession
	findErr error
}

func (r *editTestTalkSessions) FindByID(context.Context, shared.UUID[talksession.TalkSession]) (*talksession.TalkSession, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	return r.session, nil
}

func (r *editTestTalkSessions) Update(_ context.Context, ts *talksession.TalkSession) error {
	r.saved = ts
	return nil
}

type editTestSurveys struct {
	survey.Repository
	exists bool
}

func (r *editTestSurveys) FindBySessionID(context.Context, shared.UUID[talksession.TalkSession]) (*survey.Survey, error) {
	if !r.exists {
		return nil, nil
	}
	return &survey.Survey{}, nil
}

type editTestScheduler struct {
	analysis_usecase.AnalysisScheduler
	scheduled int
}

func (s *editTestScheduler) ScheduleAnalysis(context.Context, shared.UUID[talksession.TalkSession], analysis_usecase.ScheduleMode) {
	s.scheduled++
}

type editTestTx struct{}

func (editTestTx) ExecTx(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }

type editTestHarness struct {
	ctx       context.Context
	sessions  *editTestTalkSessions
	surveys   *editTestSurveys
	scheduler *editTestScheduler
	uc        EditTalkSession
	input     EditTalkSessionInput
	log       *bytes.Buffer
	actor     string
}

func newEditTestHarness(t *testing.T) *editTestHarness {
	t.Helper()
	var output bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&output, nil)))
	t.Cleanup(func() { slog.SetDefault(previous) })

	actor := shared.NewUUID[admin.AdminUser]().String()
	ctx := admin.SetAdminSession(clock.SetNow(context.Background(), editTestNow), &admin.AdminClaim{Sub: actor})
	id := shared.NewUUID[talksession.TalkSession]()
	city, prefecture, description := "千代田区", "東京都", "<p>編集前</p>"
	ts := talksession.NewTalkSession(
		id, "編集前のテーマ", &description, nil, shared.NewUUID[user.User](),
		editTestNow.Add(-48*time.Hour), editTestNow.Add(24*time.Hour),
		nil, &city, &prefecture, false, nil, nil, nil,
	)
	h := &editTestHarness{
		ctx:       ctx,
		sessions:  &editTestTalkSessions{session: ts},
		surveys:   &editTestSurveys{},
		scheduler: &editTestScheduler{},
		log:       &output,
		actor:     actor,
	}
	h.uc = &editTalkSession{talkSessions: h.sessions, surveys: h.surveys, analysisScheduler: h.scheduler, txRunner: editTestTx{}}
	// 何も変えない入力を基準にし、各テストで必要な項目だけ変える
	h.input = EditTalkSessionInput{
		TalkSessionID:    id,
		Theme:            ts.Theme(),
		Description:      description,
		ScheduledEndTime: ts.ScheduledEndTime(),
		Prefecture:       prefecture,
		City:             city,
	}
	return h
}

func TestEditTalkSessionByAdmin(t *testing.T) {
	t.Run("管理者でなければ編集できない", func(t *testing.T) {
		h := newEditTestHarness(t)
		if err := h.uc.Execute(context.Background(), h.input); !errors.Is(err, messages.ForbiddenError) {
			t.Fatalf("err = %v, want %v", err, messages.ForbiddenError)
		}
		if h.sessions.saved != nil {
			t.Fatal("session was saved without admin")
		}
	})

	t.Run("存在しないセッションは編集できない", func(t *testing.T) {
		h := newEditTestHarness(t)
		h.sessions.findErr = errors.New("not found")
		if err := h.uc.Execute(h.ctx, h.input); !errors.Is(err, messages.TalkSessionNotFound) {
			t.Fatalf("err = %v, want %v", err, messages.TalkSessionNotFound)
		}
	})

	t.Run("各項目を更新し、変更した項目名を監査ログに残す", func(t *testing.T) {
		h := newEditTestHarness(t)
		h.input.Theme = "編集後のテーマ"
		h.input.Description = "<p>編集後</p>"
		h.input.ScheduledEndTime = editTestNow.Add(72 * time.Hour)
		h.input.Prefecture = ""
		h.input.City = ""
		h.input.HideTop = true
		h.input.HideReport = true
		h.input.DisableAnalysis = true
		h.input.Restrictions = []string{string(talksession.DemographicsPrefecture)}

		if err := h.uc.Execute(h.ctx, h.input); err != nil {
			t.Fatal(err)
		}
		saved := h.sessions.saved
		if saved == nil {
			t.Fatal("session was not saved")
		}
		if saved.Theme() != "編集後のテーマ" || *saved.Description() != "<p>編集後</p>" || !saved.ScheduledEndTime().Equal(h.input.ScheduledEndTime) {
			t.Fatal("content was not updated")
		}
		if saved.Prefecture() != nil || saved.City() != nil {
			t.Fatal("empty location was not cleared")
		}
		if !saved.HideTop() || !saved.HideReport() || !saved.DisableAnalysis() {
			t.Fatal("visibility settings were not updated")
		}
		if len(saved.Restrictions()) != 1 || saved.Restrictions()[0].Key != talksession.DemographicsPrefecture {
			t.Fatal("restrictions were not updated")
		}
		for _, expected := range []string{"talk session edited by admin", h.actor, h.input.TalkSessionID.String(), "theme", "description", "scheduled_end_time", "prefecture", "city", "hide_top", "hide_report", "disable_analysis", "restrictions"} {
			if !strings.Contains(h.log.String(), expected) {
				t.Fatalf("audit log is missing %q: %s", expected, h.log.String())
			}
		}
		if strings.Contains(h.log.String(), "編集後") {
			t.Fatal("audit log must not contain edited values")
		}
	})

	t.Run("終了済みのセッションでも、終了日時を変えなければ他の項目を編集できる", func(t *testing.T) {
		h := newEditTestHarness(t)
		h.sessions.session.ChangeScheduledEndTime(editTestNow.Add(-time.Hour))
		h.input.ScheduledEndTime = editTestNow.Add(-time.Hour)
		h.input.Theme = "終了後に直したテーマ"
		if err := h.uc.Execute(h.ctx, h.input); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("終了日時を過去に変えることはできない", func(t *testing.T) {
		h := newEditTestHarness(t)
		h.input.ScheduledEndTime = editTestNow.Add(-time.Hour)
		if err := h.uc.Execute(h.ctx, h.input); !errors.Is(err, messages.InvalidScheduledEndTime) {
			t.Fatalf("err = %v, want %v", err, messages.InvalidScheduledEndTime)
		}
		if h.sessions.saved != nil {
			t.Fatal("invalid edit was saved")
		}
	})

	t.Run("入力が不正なら保存しない", func(t *testing.T) {
		for name, tc := range map[string]struct {
			edit func(*EditTalkSessionInput)
			want error
		}{
			"テーマが空":          {func(in *EditTalkSessionInput) { in.Theme = " " }, messages.BadRequestError},
			"テーマが長すぎる":       {func(in *EditTalkSessionInput) { in.Theme = strings.Repeat("あ", 101) }, messages.TalkSessionThemeTooLong},
			"本文が長すぎる":        {func(in *EditTalkSessionInput) { in.Description = strings.Repeat("あ", 40001) }, messages.TalkSessionDescriptionTooLong},
			"不正な参加制限":        {func(in *EditTalkSessionInput) { in.Restrictions = []string{"unknown"} }, messages.BadRequestError},
			"アンケートなしで回答済み制限": {func(in *EditTalkSessionInput) { in.Restrictions = []string{string(talksession.SurveyAnswered)} }, messages.BadRequestError},
		} {
			t.Run(name, func(t *testing.T) {
				h := newEditTestHarness(t)
				tc.edit(&h.input)
				if err := h.uc.Execute(h.ctx, h.input); !errors.Is(err, tc.want) {
					t.Fatalf("err = %v, want %v", err, tc.want)
				}
				if h.sessions.saved != nil {
					t.Fatal("invalid edit was saved")
				}
			})
		}
	})

	t.Run("アンケートがあれば回答済みの参加制限を付けられる", func(t *testing.T) {
		h := newEditTestHarness(t)
		h.surveys.exists = true
		h.input.Restrictions = []string{string(talksession.SurveyAnswered)}
		if err := h.uc.Execute(h.ctx, h.input); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("分析を無効から有効に戻したときだけ分析を予約する", func(t *testing.T) {
		h := newEditTestHarness(t)
		h.sessions.session.SetDisableAnalysis(true)
		h.input.DisableAnalysis = false
		if err := h.uc.Execute(h.ctx, h.input); err != nil {
			t.Fatal(err)
		}
		if h.scheduler.scheduled != 1 {
			t.Fatalf("scheduled = %d, want 1", h.scheduler.scheduled)
		}

		h = newEditTestHarness(t)
		h.input.Theme = "分析設定は変えない"
		if err := h.uc.Execute(h.ctx, h.input); err != nil {
			t.Fatal(err)
		}
		if h.scheduler.scheduled != 0 {
			t.Fatalf("scheduled = %d, want 0", h.scheduler.scheduled)
		}
	})
}
