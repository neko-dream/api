package session

import (
	"context"
	"time"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/pkg/oauth"
)

type status string

const (
	SESSION_ACTIVE   status = "ACTIVE"
	SESSION_INACTIVE status = "INACTIVE"
)

func NewSessionStatusStr(str string) *status {
	s := status(str)
	return &s
}

type expiresAt = time.Time

func NewExpiresAt() *expiresAt {
	t := time.Now().Add(2 * 24 * time.Hour * 7)
	e := expiresAt(t)
	return &e
}

type (
	SessionRepository interface {
		Create(context.Context, Session) (*Session, error)
		Update(context.Context, Session) (*Session, error)
		FindBySessionID(context.Context, shared.UUID[Session]) (*Session, error)
		FindByUserID(context.Context, shared.UUID[user.User]) ([]Session, error)
	}

	SessionService interface {
		RefreshSession(context.Context, Session) (*Session, error)
		DeactivateUserSessions(context.Context, shared.UUID[user.User]) error
	}

	Session struct {
		sessionID    shared.UUID[Session]
		userID       shared.UUID[user.User]
		authProvider oauth.AuthProviderName
		status       status
		expires      time.Time
		lastActivity time.Time
	}
)

func NewSession(
	sessionID shared.UUID[Session],
	userID shared.UUID[user.User],
	authProvider oauth.AuthProviderName,
	status status,
	expires time.Time,
) *Session {
	return &Session{
		sessionID:    sessionID,
		userID:       userID,
		authProvider: authProvider,
		status:       status,
		expires:      expires,
		lastActivity: time.Now(),
	}
}

func (s *Session) UserID() shared.UUID[user.User] {
	return s.userID
}

func (s *Session) SessionID() shared.UUID[Session] {
	return s.sessionID
}

func (s *Session) Provider() oauth.AuthProviderName {
	return s.authProvider
}

func (s *Session) IsActive() bool {
	return s.expires.After(time.Now())
}

func (s *Session) Deactivate(ctx context.Context) error {
	s.status = *NewSessionStatusStr("INACTIVE")
	return nil
}

func (s *Session) UpdateLastActivity() {
	s.lastActivity = time.Now()
}
