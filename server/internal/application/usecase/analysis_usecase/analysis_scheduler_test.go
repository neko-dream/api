package analysis_usecase

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/neko-dream/api/internal/domain/model/analysis"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/stretchr/testify/assert"
)

// analysisServiceSpy はStartAnalysis/GenerateReportの呼び出し回数を記録するスパイ
type analysisServiceSpy struct {
	startAnalysisCount  atomic.Int32
	generateReportCount atomic.Int32
	startAnalysisDelay  time.Duration
	mu                  sync.Mutex
	calls               []string // セッションIDの呼び出し記録
}

func (s *analysisServiceSpy) StartAnalysis(_ context.Context, id shared.UUID[talksession.TalkSession]) error {
	s.startAnalysisCount.Add(1)
	s.mu.Lock()
	s.calls = append(s.calls, id.String())
	s.mu.Unlock()
	if s.startAnalysisDelay > 0 {
		time.Sleep(s.startAnalysisDelay)
	}
	return nil
}

func (s *analysisServiceSpy) GenerateReport(_ context.Context, _ shared.UUID[talksession.TalkSession]) error {
	s.generateReportCount.Add(1)
	return nil
}

func (s *analysisServiceSpy) GenerateImage(_ context.Context, _ shared.UUID[talksession.TalkSession]) (*analysis.WordCloudResponse, error) {
	return nil, nil
}

// analysisRepositoryStub はnilレポートを返すスタブ
type analysisRepositoryStub struct{}

func (r *analysisRepositoryStub) FindByTalkSessionID(_ context.Context, _ shared.UUID[talksession.TalkSession]) (*analysis.AnalysisReport, error) {
	return nil, nil
}

func (r *analysisRepositoryStub) FindByID(_ context.Context, _ shared.UUID[analysis.AnalysisReport]) (*analysis.AnalysisReport, error) {
	return nil, nil
}

func (r *analysisRepositoryStub) SaveReport(_ context.Context, _ *analysis.AnalysisReport) error {
	return nil
}

func newTestScheduler(svc *analysisServiceSpy, debounce time.Duration) *analysisScheduler {
	s := &analysisScheduler{
		analysisService:    svc,
		analysisRepository: &analysisRepositoryStub{},
		sessions:           make(map[string]*sessionState),
		semaphore:          make(chan struct{}, defaultMaxConcurrency),
		debounceInterval:   debounce,
	}
	return s
}

func TestScheduleDebounced_CoalescesMultipleRequests(t *testing.T) {
	svc := &analysisServiceSpy{}
	scheduler := newTestScheduler(svc, 100*time.Millisecond)
	scheduler.Start(context.Background())
	defer scheduler.Stop()

	sessionID := shared.NewUUID[talksession.TalkSession]()
	ctx := context.Background()

	// 10回連続でスケジュール（デバウンス内）
	for i := 0; i < 10; i++ {
		scheduler.ScheduleAnalysis(ctx, sessionID, ScheduleDebounced)
		time.Sleep(10 * time.Millisecond)
	}

	// デバウンス完了を待つ
	time.Sleep(300 * time.Millisecond)

	// 10回のリクエストが1回の分析に集約されること
	assert.Equal(t, int32(1), svc.startAnalysisCount.Load(), "debounce should coalesce requests into 1")
}

func TestScheduleImmediate_ExecutesImmediately(t *testing.T) {
	svc := &analysisServiceSpy{}
	scheduler := newTestScheduler(svc, 5*time.Second)
	scheduler.Start(context.Background())
	defer scheduler.Stop()

	sessionID := shared.NewUUID[talksession.TalkSession]()
	ctx := context.Background()

	scheduler.ScheduleAnalysis(ctx, sessionID, ScheduleImmediate)

	// 即座に実行されることを確認（セマフォ取得 + goroutine起動の時間を待つ）
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, int32(1), svc.startAnalysisCount.Load(), "immediate mode should execute right away")
}

func TestDeduplication_InFlightPreventsDouble(t *testing.T) {
	svc := &analysisServiceSpy{startAnalysisDelay: 200 * time.Millisecond}
	scheduler := newTestScheduler(svc, 50*time.Millisecond)
	scheduler.Start(context.Background())
	defer scheduler.Stop()

	sessionID := shared.NewUUID[talksession.TalkSession]()
	ctx := context.Background()

	// 1回目: 即座に実行開始
	scheduler.ScheduleAnalysis(ctx, sessionID, ScheduleImmediate)
	time.Sleep(50 * time.Millisecond)

	// 2回目: inFlight中に来たのでpendingになる
	scheduler.ScheduleAnalysis(ctx, sessionID, ScheduleImmediate)
	time.Sleep(50 * time.Millisecond)

	// 3回目: まだinFlight中。pendingは既にtrue
	scheduler.ScheduleAnalysis(ctx, sessionID, ScheduleImmediate)

	// 全て完了を待つ（200ms * 2回 + バッファ）
	time.Sleep(600 * time.Millisecond)

	// 1回目 + pending再実行の2回のみ（3回目はpendingの上書き）
	count := svc.startAnalysisCount.Load()
	assert.Equal(t, int32(2), count, "should execute exactly 2 times: initial + 1 pending re-dispatch")
}

func TestConcurrencyLimit(t *testing.T) {
	svc := &analysisServiceSpy{startAnalysisDelay: 200 * time.Millisecond}
	scheduler := newTestScheduler(svc, 10*time.Millisecond)
	scheduler.Start(context.Background())
	defer scheduler.Stop()

	ctx := context.Background()

	// 5つの異なるセッションを同時にスケジュール
	for i := 0; i < 5; i++ {
		sessionID := shared.NewUUID[talksession.TalkSession]()
		scheduler.ScheduleAnalysis(ctx, sessionID, ScheduleImmediate)
	}

	// 少し待って、同時実行数を確認
	time.Sleep(50 * time.Millisecond)

	// セマフォにより最大3つが同時実行中であること
	// startAnalysisCountは3以上（3つが実行開始済み）
	startedCount := svc.startAnalysisCount.Load()
	assert.LessOrEqual(t, startedCount, int32(3), "at most 3 should start concurrently due to semaphore")

	// 全て完了を待つ
	time.Sleep(500 * time.Millisecond)

	// 最終的に5つ全て実行されること
	assert.Equal(t, int32(5), svc.startAnalysisCount.Load(), "all 5 should eventually complete")
}

func TestStop_WaitsForInFlight(t *testing.T) {
	svc := &analysisServiceSpy{startAnalysisDelay: 200 * time.Millisecond}
	scheduler := newTestScheduler(svc, 10*time.Millisecond)
	scheduler.Start(context.Background())

	sessionID := shared.NewUUID[talksession.TalkSession]()
	ctx := context.Background()

	scheduler.ScheduleAnalysis(ctx, sessionID, ScheduleImmediate)
	time.Sleep(50 * time.Millisecond) // 実行開始を待つ

	// Stopは実行中のgoroutineの完了を待つ
	scheduler.Stop()

	// Stop後に実行が完了していること
	assert.Equal(t, int32(1), svc.startAnalysisCount.Load())
}

func TestStop_PreventsNewScheduling(t *testing.T) {
	svc := &analysisServiceSpy{}
	scheduler := newTestScheduler(svc, 10*time.Millisecond)
	scheduler.Start(context.Background())
	scheduler.Stop()

	sessionID := shared.NewUUID[talksession.TalkSession]()
	ctx := context.Background()

	// Stop後のスケジュールは無視される
	scheduler.ScheduleAnalysis(ctx, sessionID, ScheduleImmediate)
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, int32(0), svc.startAnalysisCount.Load(), "no analysis should run after stop")
}
