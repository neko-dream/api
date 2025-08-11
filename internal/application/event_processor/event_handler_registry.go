package event_processor

import (
	"sort"
	"sync"

	"github.com/neko-dream/server/internal/application/event_processor/handlers"
	"github.com/neko-dream/server/internal/domain/model/event"
)

// EventHandlerRegistry イベントハンドラーを管理するレジストリ
type EventHandlerRegistry struct {
	handlers map[event.EventType][]handlers.EventHandler
	mu       sync.RWMutex
}

// NewEventHandlerRegistry 新しいEventHandlerRegistryを作成
func NewEventHandlerRegistry() *EventHandlerRegistry {
	return &EventHandlerRegistry{
		handlers: make(map[event.EventType][]handlers.EventHandler),
	}
}

// Register ハンドラーを登録
func (r *EventHandlerRegistry) Register(eventType event.EventType, handler handlers.EventHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.handlers[eventType] = append(r.handlers[eventType], handler)
	// 優先度順にソート（高い優先度が先）
	sort.Slice(r.handlers[eventType], func(i, j int) bool {
		return r.handlers[eventType][i].Priority() > r.handlers[eventType][j].Priority()
	})
}

// GetHandlers イベントタイプに対応するハンドラーを取得
func (r *EventHandlerRegistry) GetHandlers(eventType event.EventType) []handlers.EventHandler {
	r.mu.RLock()
	defer r.mu.RUnlock()

	handlerList := r.handlers[eventType]
	result := make([]handlers.EventHandler, len(handlerList))
	copy(result, handlerList)
	return result
}

// GetRegisteredEventTypes 登録されているすべてのイベントタイプを取得
func (r *EventHandlerRegistry) GetRegisteredEventTypes() []event.EventType {
	r.mu.RLock()
	defer r.mu.RUnlock()

	eventTypes := make([]event.EventType, 0, len(r.handlers))
	for eventType := range r.handlers {
		eventTypes = append(eventTypes, eventType)
	}
	return eventTypes
}

// Clear すべてのハンドラーをクリア（主にテスト用）
func (r *EventHandlerRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.handlers = make(map[event.EventType][]handlers.EventHandler)
}
