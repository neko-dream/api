package analysis_usecase

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/neko-dream/api/internal/domain/model/analysis"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const (
	defaultDebounceInterval = 5 * time.Second
	defaultMaxConcurrency   = 3
)

// sessionState はセッション単位の分析スケジュール状態を管理する
type sessionState struct {
	timer    *time.Timer       // デバウンスタイマー（ScheduleDebounced時のみ使用）
	inFlight bool              // 分析が実行中かどうか
	pending  bool              // 実行中に新しいリクエストが来たかどうか
	spanCtx  trace.SpanContext // トレーシング用のSpanContext
}

type analysisScheduler struct {
	analysisService    analysis.AnalysisService
	analysisRepository analysis.AnalysisRepository

	mu       sync.Mutex
	sessions map[string]*sessionState

	semaphore chan struct{}
	wg        sync.WaitGroup

	debounceInterval time.Duration
	stopped          bool
}

// NewAnalysisScheduler は新しいAnalysisSchedulerを生成する
func NewAnalysisScheduler(
	analysisService analysis.AnalysisService,
	analysisRepository analysis.AnalysisRepository,
) AnalysisScheduler {
	return &analysisScheduler{
		analysisService:    analysisService,
		analysisRepository: analysisRepository,
		sessions:           make(map[string]*sessionState),
		semaphore:          make(chan struct{}, defaultMaxConcurrency),
		debounceInterval:   defaultDebounceInterval,
	}
}

func (s *analysisScheduler) Start(_ context.Context) {
	slog.Info("Analysis scheduler started",
		"debounceInterval", s.debounceInterval,
		"maxConcurrency", defaultMaxConcurrency,
	)
}

func (s *analysisScheduler) Stop() {
	s.mu.Lock()
	s.stopped = true
	// 全てのペンディングタイマーを停止
	for _, state := range s.sessions {
		if state.timer != nil {
			state.timer.Stop()
		}
	}
	s.mu.Unlock()

	// 実行中のgoroutineの完了を待つ
	s.wg.Wait()
	slog.Info("Analysis scheduler stopped")
}

func (s *analysisScheduler) ScheduleAnalysis(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession], mode ScheduleMode) {
	_, span := otel.Tracer("analysis_scheduler").Start(ctx, "AnalysisScheduler.ScheduleAnalysis")
	defer span.End()

	key := talkSessionID.String()
	spanCtx := span.SpanContext()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stopped {
		return
	}

	state, exists := s.sessions[key]

	switch mode {
	case ScheduleDebounced:
		if exists {
			state.spanCtx = spanCtx
			if state.timer != nil {
				// タイマーが動いている場合はリセット
				state.timer.Reset(s.debounceInterval)
				return
			}
			if state.inFlight {
				// 実行中の場合はpendingにする
				state.pending = true
				return
			}
		}
		// 新規または完了済み: デバウンスタイマーを開始
		if !exists {
			state = &sessionState{}
			s.sessions[key] = state
		}
		state.spanCtx = spanCtx
		state.timer = time.AfterFunc(s.debounceInterval, func() {
			s.onTimerFired(talkSessionID, key)
		})

	case ScheduleImmediate:
		if exists && state.inFlight {
			// 実行中の場合はpendingにする
			state.pending = true
			state.spanCtx = spanCtx
			return
		}
		if exists && state.timer != nil {
			// デバウンスタイマーが動いている場合は停止して即実行
			state.timer.Stop()
			state.timer = nil
		}
		if !exists {
			state = &sessionState{}
			s.sessions[key] = state
		}
		state.spanCtx = spanCtx
		state.inFlight = true
		s.dispatchLocked(talkSessionID, key, state)
	}
}

// onTimerFired はデバウンスタイマーが発火した時に呼ばれる
func (s *analysisScheduler) onTimerFired(talkSessionID shared.UUID[talksession.TalkSession], key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stopped {
		return
	}

	state, exists := s.sessions[key]
	if !exists {
		return
	}

	state.timer = nil
	state.inFlight = true
	s.dispatchLocked(talkSessionID, key, state)
}

// dispatchLocked はsemaphore取得とgoroutine起動を行う。mu.Lockを保持した状態で呼ぶこと。
func (s *analysisScheduler) dispatchLocked(talkSessionID shared.UUID[talksession.TalkSession], key string, state *sessionState) {
	spanCtx := state.spanCtx
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		// セマフォ取得（同時実行制限）
		s.semaphore <- struct{}{}
		defer func() { <-s.semaphore }()

		s.executeAnalysis(talkSessionID, key, spanCtx)
	}()
}

// executeAnalysis は実際の分析処理を実行する
func (s *analysisScheduler) executeAnalysis(talkSessionID shared.UUID[talksession.TalkSession], key string, spanCtx trace.SpanContext) {
	bg := context.Background()
	bg = trace.ContextWithSpanContext(bg, spanCtx)
	ctx, span := otel.Tracer("analysis_scheduler").Start(bg, "AnalysisScheduler.executeAnalysis")
	defer span.End()

	defer func() {
		if r := recover(); r != nil {
			utils.HandleError(ctx, nil, "AnalysisScheduler.executeAnalysis panicked")
		}
		s.onAnalysisComplete(talkSessionID, key)
	}()

	_ = s.analysisService.StartAnalysis(ctx, talkSessionID)

	report, err := s.analysisRepository.FindByTalkSessionID(ctx, talkSessionID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		utils.HandleError(ctx, err, "AnalysisRepository.FindByTalkSessionID")
		return
	}
	if report == nil {
		return
	}
	if report.ShouldReGenerateReport() {
		_ = s.analysisService.GenerateReport(ctx, talkSessionID)
	}
}

// onAnalysisComplete は分析完了後にpendingチェックとクリーンアップを行う
func (s *analysisScheduler) onAnalysisComplete(talkSessionID shared.UUID[talksession.TalkSession], key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, exists := s.sessions[key]
	if !exists {
		return
	}

	if state.pending && !s.stopped {
		// pendingリクエストがある場合は再ディスパッチ
		state.pending = false
		s.dispatchLocked(talkSessionID, key, state)
	} else {
		// クリーンアップ
		state.inFlight = false
		delete(s.sessions, key)
	}
}
