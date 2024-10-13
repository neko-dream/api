package talksession

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/shared/time"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type (
	TalkSessionRepository interface {
		Create(ctx context.Context, talkSession *TalkSession) error
		Update(ctx context.Context, talkSession *TalkSession) error
	}

	TalkSession struct {
		talkSessionID    shared.UUID[TalkSession]
		theme            string
		ownerUserID      shared.UUID[user.User]
		scheduledEndTime time.Time // 予定終了時間
		createdAt        time.Time // 作成日時
		location         *Location
		city             *string
		prefecture       *string
	}
)

func NewTalkSession(
	talkSessionID shared.UUID[TalkSession],
	theme string,
	ownerUserID shared.UUID[user.User],
	createdAt time.Time,
	scheduledEndTime time.Time,
	location *Location,
	city *string,
	prefecture *string,
) *TalkSession {
	return &TalkSession{
		talkSessionID:    talkSessionID,
		theme:            theme,
		ownerUserID:      ownerUserID,
		createdAt:        createdAt,
		scheduledEndTime: scheduledEndTime,
		location:         location,
		city:             city,
		prefecture:       prefecture,
	}
}

func (t *TalkSession) TalkSessionID() shared.UUID[TalkSession] {
	return t.talkSessionID
}

func (t *TalkSession) OwnerUserID() shared.UUID[user.User] {
	return t.ownerUserID
}

func (t *TalkSession) Theme() string {
	return t.theme
}

func (t *TalkSession) ScheduledEndTime() time.Time {
	return t.scheduledEndTime
}

func (t *TalkSession) CreatedAt() time.Time {
	return t.createdAt
}

func (t *TalkSession) Location() *Location {
	return t.location
}

func (t *TalkSession) City() *string {
	return t.city
}

func (t *TalkSession) Prefecture() *string {
	return t.prefecture
}

func (t *TalkSession) ChangeTheme(theme string) {
	t.theme = theme
}

// 終了しているかを調べる
func (t *TalkSession) IsFinished(ctx context.Context) bool {
	return t.scheduledEndTime.Before(time.Now(ctx).Time)
}
