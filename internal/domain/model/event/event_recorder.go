package event

// EventRecorder イベントを記録する集約のための基底型
type EventRecorder struct {
	events []DomainEvent
}

// RecordEvent イベントを記録
func (r *EventRecorder) RecordEvent(event DomainEvent) {
	r.events = append(r.events, event)
}

// GetRecordedEvents 記録されたイベントを取得
func (r *EventRecorder) GetRecordedEvents() []DomainEvent {
	return r.events
}

// ClearRecordedEvents 記録されたイベントをクリア
func (r *EventRecorder) ClearRecordedEvents() {
	r.events = r.events[:0]
}
