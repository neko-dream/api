package event

import (
	"context"
	"time"

	"github.com/neko-dream/server/internal/domain/model/shared"
)

// EventStore イベントストアのインターフェース
type EventStore interface {
	// Store トランザクション内でイベントを保存
	Store(ctx context.Context, event DomainEvent) error

	// StoreBatch トランザクション内で複数のイベントを保存
	StoreBatch(ctx context.Context, events []DomainEvent) error

	// GetUnprocessedEvents 未処理のイベントを取得
	GetUnprocessedEvents(ctx context.Context, eventTypes []EventType, limit int) ([]StoredEvent, error)

	// MarkAsProcessed イベントを処理済みとしてマーク
	MarkAsProcessed(ctx context.Context, eventID shared.UUID[StoredEvent]) error

	// MarkAsFailed イベントを失敗としてマーク
	MarkAsFailed(ctx context.Context, eventID shared.UUID[StoredEvent], reason string) error
}

type StoredEvent struct {
	ID            shared.UUID[StoredEvent]
	EventType     EventType
	EventData     []byte
	AggregateID   string
	AggregateType string
	Status        EventStatus
	OccurredAt    time.Time
	ProcessedAt   *time.Time
	FailedAt      *time.Time
	FailureReason *string
	RetryCount    int
}

type EventStatus string

const (
	EventStatusPending    EventStatus = "pending"
	EventStatusProcessing EventStatus = "processing"
	EventStatusProcessed  EventStatus = "processed"
	EventStatusFailed     EventStatus = "failed"
)
