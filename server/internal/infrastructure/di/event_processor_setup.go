package di

import (
	"github.com/neko-dream/api/internal/application/event_processor"
	"github.com/neko-dream/api/internal/application/event_processor/handlers"
	"github.com/neko-dream/api/internal/domain/model/event"
	"github.com/neko-dream/api/internal/domain/model/talksession"
)

func SetupEventProcessor(
	eventStore event.EventStore,
	registry *event_processor.EventHandlerRegistry,
	pushHandler *handlers.TalkSessionPushNotificationHandler,
) *event_processor.EventProcessor {

	registry.Register(talksession.EventTypeTalkSessionStarted, pushHandler)
	registry.Register(talksession.EventTypeTalkSessionEnded, pushHandler)

	return event_processor.NewEventProcessor(eventStore, registry)
}
