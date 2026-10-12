package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"braces.dev/errtrace"
	"github.com/neko-dream/api/internal/domain/model/event"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
	"go.opentelemetry.io/otel"
)

type eventStore struct {
	*db.DBManager
}

func NewEventStore(dbManager *db.DBManager) event.EventStore {
	return &eventStore{
		DBManager: dbManager,
	}
}

func (s *eventStore) Store(ctx context.Context, evt event.DomainEvent) error {
	ctx, span := otel.Tracer("persistence").Start(ctx, "eventStore.Store")
	defer span.End()

	eventData, err := json.Marshal(evt)
	if err != nil {
		return errtrace.Wrap(fmt.Errorf("failed to marshal event: %w", err))
	}

	params := model.CreateDomainEventParams{
		ID:            evt.EventID().UUID(),
		EventType:     string(evt.EventType()),
		EventData:     eventData,
		AggregateID:   evt.AggregateID(),
		AggregateType: evt.AggregateType(),
		Status:        string(event.EventStatusPending),
		OccurredAt:    evt.OccurredAt(),
		RetryCount:    0,
	}

	_, err = s.GetQueries(ctx).CreateDomainEvent(ctx, params)
	return err
}

func (s *eventStore) StoreBatch(ctx context.Context, events []event.DomainEvent) error {
	ctx, span := otel.Tracer("persistence").Start(ctx, "eventStore.StoreBatch")
	defer span.End()

	if len(events) == 0 {
		return nil
	}

	return s.ExecTx(ctx, func(ctx context.Context) error {
		for _, evt := range events {
			eventData, err := json.Marshal(evt)
			if err != nil {
				return errtrace.Wrap(fmt.Errorf("failed to marshal event: %w", err))
			}

			params := model.CreateDomainEventParams{
				ID:            evt.EventID().UUID(),
				EventType:     string(evt.EventType()),
				EventData:     eventData,
				AggregateID:   evt.AggregateID(),
				AggregateType: evt.AggregateType(),
				Status:        string(event.EventStatusPending),
				OccurredAt:    evt.OccurredAt(),
				RetryCount:    0,
			}

			if _, err := s.GetQueries(ctx).CreateDomainEvent(ctx, params); err != nil {
				return errtrace.Wrap(fmt.Errorf("failed to insert event: %w", err))
			}
		}
		return nil
	})
}

func (s *eventStore) GetUnprocessedEvents(ctx context.Context, eventTypes []event.EventType, limit int) ([]event.StoredEvent, error) {
	ctx, span := otel.Tracer("persistence").Start(ctx, "eventStore.GetUnprocessedEvents")
	defer span.End()

	eventTypeStrings := make([]string, len(eventTypes))
	for i, et := range eventTypes {
		eventTypeStrings[i] = string(et)
	}

	params := model.GetUnprocessedEventsParams{
		RetryCount: 3, // Max retry count
		Column2:    eventTypeStrings,
		Limit:      int32(limit),
	}

	rows, err := s.GetQueries(ctx).GetUnprocessedEvents(ctx, params)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	events := make([]event.StoredEvent, len(rows))
	for i, row := range rows {
		var processedAt *time.Time
		if row.ProcessedAt.Valid {
			processedAt = &row.ProcessedAt.Time
		}
		var failedAt *time.Time
		if row.FailedAt.Valid {
			failedAt = &row.FailedAt.Time
		}
		var failureReason *string
		if row.FailureReason.Valid {
			failureReason = &row.FailureReason.String
		}

		events[i] = event.StoredEvent{
			ID:            shared.UUID[event.StoredEvent](row.ID),
			EventType:     event.EventType(row.EventType),
			EventData:     row.EventData,
			AggregateID:   row.AggregateID,
			AggregateType: row.AggregateType,
			Status:        event.EventStatus(row.Status),
			OccurredAt:    row.OccurredAt,
			ProcessedAt:   processedAt,
			FailedAt:      failedAt,
			FailureReason: failureReason,
			RetryCount:    int(row.RetryCount),
		}
	}

	return events, nil
}

// MarkAsProcessed marks an event as processed
func (s *eventStore) MarkAsProcessed(ctx context.Context, eventID shared.UUID[event.StoredEvent]) error {
	ctx, span := otel.Tracer("persistence").Start(ctx, "eventStore.MarkAsProcessed")
	defer span.End()

	_, err := s.GetQueries(ctx).MarkEventAsProcessed(ctx, eventID.UUID())
	return err
}

// MarkAsFailed marks an event as failed
func (s *eventStore) MarkAsFailed(ctx context.Context, eventID shared.UUID[event.StoredEvent], reason string) error {
	ctx, span := otel.Tracer("persistence").Start(ctx, "eventStore.MarkAsFailed")
	defer span.End()

	params := model.MarkEventAsFailedParams{
		ID:            eventID.UUID(),
		FailureReason: sql.NullString{String: reason, Valid: true},
	}

	_, err := s.GetQueries(ctx).MarkEventAsFailed(ctx, params)
	return err
}
