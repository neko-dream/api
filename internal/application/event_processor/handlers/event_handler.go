package handlers

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/event"
)

// EventHandler イベントハンドラーのインターフェース
type EventHandler interface {
	// CanHandle このハンドラーがイベントを処理できるかチェック
	CanHandle(eventType event.EventType) bool

	// Handle イベントを処理
	Handle(ctx context.Context, storedEvent event.StoredEvent) error

	// Priority 処理の優先度（同じイベントタイプに複数ハンドラーがある場合）
	// 大きい数値ほど高優先度
	Priority() int
}
