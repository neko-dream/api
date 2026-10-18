package event_processor

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/neko-dream/api/internal/domain/model/event"
	"github.com/neko-dream/api/internal/domain/model/shared"
)

const testEventType event.EventType = "test.event"

type eventStoreStub struct {
	event.EventStore
	mu      sync.Mutex
	pending []event.StoredEvent
	fetches int
}

func (s *eventStoreStub) GetUnprocessedEvents(_ context.Context, _ []event.EventType, limit int) ([]event.StoredEvent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.fetches++
	n := min(limit, len(s.pending))
	return append([]event.StoredEvent(nil), s.pending[:n]...), nil
}

func (s *eventStoreStub) MarkAsProcessed(_ context.Context, id shared.UUID[event.StoredEvent]) error {
	s.remove(id)
	return nil
}

func (s *eventStoreStub) MarkAsFailed(context.Context, shared.UUID[event.StoredEvent], string) error {
	return nil
}

func (s *eventStoreStub) add(n int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for range n {
		s.pending = append(s.pending, event.StoredEvent{ID: shared.NewUUID[event.StoredEvent](), EventType: testEventType})
	}
}

func (s *eventStoreStub) remove(id shared.UUID[event.StoredEvent]) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, e := range s.pending {
		if e.ID == id {
			s.pending = append(s.pending[:i], s.pending[i+1:]...)
			return
		}
	}
}

func (s *eventStoreStub) counts() (pending, fetches int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.pending), s.fetches
}

type handlerStub struct {
	mu     sync.Mutex
	err    error
	called int
}

func (h *handlerStub) CanHandle(eventType event.EventType) bool { return eventType == testEventType }
func (h *handlerStub) Priority() int                            { return 0 }
func (h *handlerStub) Handle(context.Context, event.StoredEvent) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.called++
	return h.err
}

func (h *handlerStub) calls() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.called
}

func startProcessor(t *testing.T, store *eventStoreStub, handler *handlerStub, retryInterval time.Duration, batchSize int) *event.Notifier {
	t.Helper()
	registry := NewEventHandlerRegistry()
	registry.Register(testEventType, handler)
	notifier := event.NewNotifier()
	p := NewEventProcessor(store, notifier, registry).WithRetryInterval(retryInterval).WithBatchSize(batchSize)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		p.Start(ctx)
		close(done)
	}()
	t.Cleanup(func() {
		cancel()
		<-done
	})
	return notifier
}

func waitFor(t *testing.T, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatal("条件を満たしませんでした")
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func TestEventProcessor_ProcessesOnStartAndNotify(t *testing.T) {
	store := &eventStoreStub{}
	store.add(1)
	handler := &handlerStub{}
	notifier := startProcessor(t, store, handler, time.Hour, 100)

	// 起動時に処理する
	waitFor(t, func() bool { p, _ := store.counts(); return p == 0 })

	// 通知されたら処理する
	store.add(1)
	notifier.Notify()
	waitFor(t, func() bool { p, _ := store.counts(); return p == 0 })
	if got := handler.calls(); got != 2 {
		t.Fatalf("handler calls = %d, want 2", got)
	}
}

func TestEventProcessor_DoesNotPollWhenIdle(t *testing.T) {
	store := &eventStoreStub{}
	startProcessor(t, store, &handlerStub{}, 10*time.Millisecond, 100)

	waitFor(t, func() bool { _, f := store.counts(); return f == 1 })
	time.Sleep(100 * time.Millisecond)
	if _, fetches := store.counts(); fetches != 1 {
		t.Fatalf("fetches = %d, want 1 (アイドル時に DB を見に行かない)", fetches)
	}
}

func TestEventProcessor_RetriesOnlyWhileFailing(t *testing.T) {
	store := &eventStoreStub{}
	store.add(1)
	handler := &handlerStub{err: errors.New("failed")}
	startProcessor(t, store, handler, 10*time.Millisecond, 100)

	// 失敗している間は再試行する
	waitFor(t, func() bool { return handler.calls() >= 3 })

	// 成功したら再試行をやめる
	handler.mu.Lock()
	handler.err = nil
	handler.mu.Unlock()
	waitFor(t, func() bool { p, _ := store.counts(); return p == 0 })
	_, fetches := store.counts()
	time.Sleep(100 * time.Millisecond)
	if _, got := store.counts(); got != fetches {
		t.Fatalf("fetches = %d, want %d (成功後は再試行しない)", got, fetches)
	}
}

func TestEventProcessor_ContinuesWhenBatchIsFull(t *testing.T) {
	store := &eventStoreStub{}
	store.add(5)
	startProcessor(t, store, &handlerStub{}, time.Hour, 2)

	waitFor(t, func() bool { p, _ := store.counts(); return p == 0 })
}
