package event_processor

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/neko-dream/api/internal/domain/model/event"
	"go.opentelemetry.io/otel"
)

type EventProcessor struct {
	eventStore event.EventStore
	registry   *EventHandlerRegistry
	logger     *slog.Logger
	interval   time.Duration
	batchSize  int
}

func NewEventProcessor(
	eventStore event.EventStore,
	registry *EventHandlerRegistry,
) *EventProcessor {
	return &EventProcessor{
		eventStore: eventStore,
		registry:   registry,
		logger:     slog.Default(),
		interval:   10 * time.Second,
		batchSize:  100,
	}
}

func (p *EventProcessor) WithInterval(interval time.Duration) *EventProcessor {
	p.interval = interval
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

	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	p.logger.Info("イベントプロセッサーを開始しました",
		slog.Duration("interval", p.interval),
		slog.Int("batch_size", p.batchSize),
	)

	p.processBatch(ctx)

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("イベントプロセッサーを停止します")
			return
		case <-ticker.C:
			p.processBatch(ctx)
		}
	}
}

func (p *EventProcessor) processBatch(ctx context.Context) {
	ctx, span := otel.Tracer("event_processor").Start(ctx, "EventProcessor.processBatch")
	defer span.End()

	eventTypes := p.registry.GetRegisteredEventTypes()
	if len(eventTypes) == 0 {
		p.logger.Debug("登録されているイベントタイプがありません")
		return
	}

	events, err := p.eventStore.GetUnprocessedEvents(ctx, eventTypes, p.batchSize)
	if err != nil {
		p.logger.Error("未処理イベントの取得に失敗しました",
			slog.String("error", err.Error()),
		)
		return
	}

	if len(events) == 0 {
		return
	}

	p.logger.Info("イベントを処理します",
		slog.Int("count", len(events)),
	)

	for _, storedEvent := range events {
		if err := p.processEvent(ctx, storedEvent); err != nil {
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
