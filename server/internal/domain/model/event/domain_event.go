package event

import (
	"time"

	"github.com/neko-dream/api/internal/domain/model/shared"
)

// EventType represents the type of domain event
type EventType string

// DomainEvent ドメインイベントの基底インターフェース
type DomainEvent interface {
	EventID() shared.UUID[BaseEvent]
	EventType() EventType
	OccurredAt() time.Time
	AggregateID() string
	AggregateType() string
}

// BaseEvent イベントの基底実装
type BaseEvent struct {
	ID          shared.UUID[BaseEvent] `json:"id"`
	Type        EventType              `json:"type"`
	Timestamp   time.Time              `json:"timestamp"`
	AggregateId string                 `json:"aggregate_id"`
	Aggregate   string                 `json:"aggregate"`
}

// NewBaseEvent 新しい基底イベントを生成
func NewBaseEvent(eventType EventType, aggregateID, aggregateType string) BaseEvent {
	return BaseEvent{
		ID:          shared.NewUUID[BaseEvent](),
		Type:        eventType,
		Timestamp:   time.Now(),
		AggregateId: aggregateID,
		Aggregate:   aggregateType,
	}
}

func (e BaseEvent) EventID() shared.UUID[BaseEvent] { return e.ID }
func (e BaseEvent) EventType() EventType            { return e.Type }
func (e BaseEvent) OccurredAt() time.Time           { return e.Timestamp }
func (e BaseEvent) AggregateID() string             { return e.AggregateId }
func (e BaseEvent) AggregateType() string           { return e.Aggregate }
