package talksession

import (
	"context"
	"time"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type (
	TalkSessionRepository interface {
		Create(ctx context.Context, talkSession *TalkSession) error
		Update(ctx context.Context, talkSession *TalkSession) error
		FindByID(ctx context.Context, talkSessionID shared.UUID[TalkSession]) (*TalkSession, error)
	}

	TalkSession struct {
		talkSessionID shared.UUID[TalkSession]
		theme         string
		ownerUserID   shared.UUID[user.User]
		finishedAt    *time.Time
	}
)

func NewTalkSession(
	talkSessionID shared.UUID[TalkSession],
	theme string,
	ownerUserID shared.UUID[user.User],
	finishedAt *time.Time,
) *TalkSession {
	return &TalkSession{
		talkSessionID: talkSessionID,
		theme:         theme,
		ownerUserID:   ownerUserID,
		finishedAt:    finishedAt,
	}
}

func (t *TalkSession) TalkSessionID() shared.UUID[TalkSession] {
	return t.talkSessionID
}

func (t *TalkSession) Theme() string {
	return t.theme
}

func (t *TalkSession) ChangeTheme(theme string) {
	t.theme = theme
}

func (t *TalkSession) OwnerUserID() shared.UUID[user.User] {
	return t.ownerUserID
}

func (t *TalkSession) IsFinished() bool {
	return t.finishedAt != nil
}
