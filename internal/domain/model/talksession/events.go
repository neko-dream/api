package talksession

import (
	"time"

	"github.com/neko-dream/server/internal/domain/model/event"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
)

const (
	EventTypeTalkSessionStarted event.EventType = "talksession.started"
	EventTypeTalkSessionEnded   event.EventType = "talksession.ended"
)

type TalkSessionStartedEvent struct {
	event.BaseEvent
	TalkSessionID    shared.UUID[TalkSession]                `json:"talk_session_id"`
	OwnerID          shared.UUID[user.User]                  `json:"owner_id"`
	Theme            string                                  `json:"theme"`
	Description      string                                  `json:"description"`
	OrganizationID   *shared.UUID[organization.Organization] `json:"organization_id,omitempty"`
	ScheduledEndTime time.Time                               `json:"scheduled_end_time"`
}

func NewTalkSessionStartedEvent(
	talkSessionID shared.UUID[TalkSession],
	ownerID shared.UUID[user.User],
	theme string,
	description string,
	organizationID *shared.UUID[organization.Organization],
	scheduledEndTime time.Time,
) *TalkSessionStartedEvent {
	return &TalkSessionStartedEvent{
		BaseEvent:        event.NewBaseEvent(EventTypeTalkSessionStarted, talkSessionID.String(), "TalkSession"),
		TalkSessionID:    talkSessionID,
		OwnerID:          ownerID,
		Theme:            theme,
		Description:      description,
		OrganizationID:   organizationID,
		ScheduledEndTime: scheduledEndTime,
	}
}

type TalkSessionEndedEvent struct {
	event.BaseEvent
	TalkSessionID  shared.UUID[TalkSession] `json:"talk_session_id"`
	OwnerID        shared.UUID[user.User]   `json:"owner_id"`
	Theme          string                   `json:"theme"`
	ParticipantIDs []shared.UUID[user.User] `json:"participant_ids"`
	EndedAt        time.Time                `json:"ended_at"`
}

func NewTalkSessionEndedEvent(
	talkSessionID shared.UUID[TalkSession],
	ownerID shared.UUID[user.User],
	theme string,
	participantIDs []shared.UUID[user.User],
) *TalkSessionEndedEvent {
	return &TalkSessionEndedEvent{
		BaseEvent:      event.NewBaseEvent(EventTypeTalkSessionEnded, talkSessionID.String(), "TalkSession"),
		TalkSessionID:  talkSessionID,
		OwnerID:        ownerID,
		Theme:          theme,
		ParticipantIDs: participantIDs,
		EndedAt:        time.Now(),
	}
}
