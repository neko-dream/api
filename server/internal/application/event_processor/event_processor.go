package event_processor

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/neko-dream/api/internal/domain/model/event"
	"go.opentelemetry.io/otel"
)

// EventProcessor 保存されたイベントを処理する。
// DB を定期的に見に行くと Aurora Serverless v2 が自動一時停止できないため、
// 起動時と、イベントが保存されたとき(Notifier)だけ処理し、失敗したイベントがある間だけ retryInterval ごとに再試行する。
type EventProcessor struct {
	eventStore    event.EventStore
	notifier      *event.Notifier
	registry      *EventHandlerRegistry
	logger        *slog.Logger
	retryInterval time.Duration
	batchSize     int
}

func NewEventProcessor(
	eventStore event.EventStore,
	notifier *event.Notifier,
	registry *EventHandlerRegistry,
) *EventProcessor {
	return &EventProcessor{
		eventStore:    eventStore,
		notifier:      notifier,
		registry:      registry,
		logger:        slog.Default(),
		retryInterval: time.Minute,
		batchSize:     100,
	}
}

func (p *EventProcessor) WithRetryInterval(retryInterval time.Duration) *EventProcessor {
	p.retryInterval = retryInterval
	return p
}

func (p *EventProcessor) WithBatchSize(batchSize int) *EventProcessor {
	p.batchSize = batchSize
	return p
}

func (p *EventProcessor) WithLogger(logger *slog.Logger) *EventProcessor {
	p.logger = logger
	return p
}

func (p *EventProcessor) Start(ctx context.Context) {
	ctx, span := otel.Tracer("event_processor").Start(ctx, "EventProcessor.Start")
	defer span.End()

	p.logger.Info("イベントプロセッサーを開始しました",
		slog.Duration("retry_interval", p.retryInterval),
		slog.Int("batch_size", p.batchSize),
	)

	// 停止中に保存されたイベントを拾うため、起動時に一度処理する
	retry := p.processBatch(ctx)

	for {
		var retryC <-chan time.Time
		if retry {
			retryC = time.After(p.retryInterval)
		}

		select {
		case <-ctx.Done():
			p.logger.Info("イベントプロセッサーを停止します")
			return
		case <-p.notifier.C():
			retry = p.processBatch(ctx)
		case <-retryC:
			retry = p.processBatch(ctx)
		}
	}
}

// processBatch 未処理のイベントを処理する。再試行が必要なイベントが残っていれば true を返す
func (p *EventProcessor) processBatch(ctx context.Context) bool {
	ctx, span := otel.Tracer("event_processor").Start(ctx, "EventProcessor.processBatch")
	defer span.End()

	eventTypes := p.registry.GetRegisteredEventTypes()
	if len(eventTypes) == 0 {
		p.logger.Debug("登録されているイベントタイプがありません")
		return false
	}

	events, err := p.eventStore.GetUnprocessedEvents(ctx, eventTypes, p.batchSize)
	if err != nil {
		p.logger.Error("未処理イベントの取得に失敗しました",
			slog.String("error", err.Error()),
		)
		return true
	}

	if len(events) == 0 {
		return false
	}

	p.logger.Info("イベントを処理します",
		slog.Int("count", len(events)),
	)

	failed := false
	for _, storedEvent := range events {
		if err := p.processEvent(ctx, storedEvent); err != nil {
			failed = true
			p.logger.Error("イベント処理に失敗しました",
				slog.String("event_id", storedEvent.ID.String()),
				slog.String("event_type", string(storedEvent.EventType)),
				slog.String("error", err.Error()),
			)

			if markErr := p.eventStore.MarkAsFailed(ctx, storedEvent.ID, err.Error()); markErr != nil {
				p.logger.Error("イベントの失敗マークに失敗しました",
					slog.String("event_id", storedEvent.ID.String()),
					slog.String("error", markErr.Error()),
				)
			}
			continue
		}

		if err := p.eventStore.MarkAsProcessed(ctx, storedEvent.ID); err != nil {
			p.logger.Error("イベントの処理済みマークに失敗しました",
				slog.String("event_id", storedEvent.ID.String()),
				slog.String("error", err.Error()),
			)
		}
	}

	// バッチに収まらなかったイベントは続けて処理する。失敗があれば再試行の間隔を空ける
	if failed {
		return true
	}
	if len(events) >= p.batchSize {
		p.notifier.Notify()
	}
	return false
}

// processEvent 個別のイベントを処理
func (p *EventProcessor) processEvent(ctx context.Context, storedEvent event.StoredEvent) error {
	ctx, span := otel.Tracer("event_processor").Start(ctx, "EventProcessor.processEvent")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	handlers := p.registry.GetHandlers(storedEvent.EventType)

	if len(handlers) == 0 {
		p.logger.Warn("ハンドラーが登録されていません",
			slog.String("event_type", string(storedEvent.EventType)),
		)
		return nil
	}

	var errs []error
	for _, handler := range handlers {
		if !handler.CanHandle(storedEvent.EventType) {
			continue
		}

		p.logger.Debug("ハンドラーでイベントを処理します",
			slog.String("event_type", string(storedEvent.EventType)),
			slog.Int("priority", handler.Priority()),
		)

		if err := handler.Handle(ctx, storedEvent); err != nil {
			errs = append(errs, fmt.Errorf("handler (priority=%d): %w", handler.Priority(), err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("一部のハンドラーでエラーが発生: %v", errs)
	}

	return nil
}
