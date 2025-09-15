package repository

import (
	"context"

	"github.com/neko-dream/api/internal/domain/model/event"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"go.opentelemetry.io/otel"
)

// EventStoreMock is a mock implementation of EventStore for testing
type EventStoreMock struct {
	events []event.DomainEvent
}

// NewEventStoreMock creates a new mock event store
func NewEventStoreMock() event.EventStore {
	return &EventStoreMock{
		events: make([]event.DomainEvent, 0),
	}
}

// Store stores a single event
func (m *EventStoreMock) Store(ctx context.Context, evt event.DomainEvent) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "EventStoreMock.Store")
	defer span.End()

	m.events = append(m.events, evt)
	return nil
}

func (m *EventStoreMock) StoreBatch(ctx context.Context, events []event.DomainEvent) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "EventStoreMock.StoreBatch")
	defer span.End()

	m.events = append(m.events, events...)
	return nil
}

func (m *EventStoreMock) GetUnprocessedEvents(ctx context.Context, eventTypes []event.EventType, limit int) ([]event.StoredEvent, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "EventStoreMock.GetUnprocessedEvents")
	defer span.End()

	return []event.StoredEvent{}, nil
}

// MarkAsProcessed marks an event as processed
func (m *EventStoreMock) MarkAsProcessed(ctx context.Context, eventID shared.UUID[event.StoredEvent]) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "EventStoreMock.MarkAsProcessed")
	defer span.End()

	return nil
}

// MarkAsFailed marks an event as failed
func (m *EventStoreMock) MarkAsFailed(ctx context.Context, eventID shared.UUID[event.StoredEvent], reason string) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "EventStoreMock.MarkAsFailed")
	defer span.End()

	return nil
}

// GetEvents returns stored events (for testing)
func (m *EventStoreMock) GetEvents() []event.DomainEvent {
	return m.events
}
