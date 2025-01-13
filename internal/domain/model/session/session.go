package session

import (
	"context"
	"time"

	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"go.opentelemetry.io/otel"
)

type status int

const (
	SESSION_ACTIVE status = iota
	SESSION_INACTIVE
)

func NewSessionStatus(num int) *status {
	switch num {
	case 0:
		s := SESSION_ACTIVE
		return &s
	case 1:
		s := SESSION_INACTIVE
		return &s
	default:
		return nil
	}
}

type expiresAt = time.Time

func NewExpiresAt(ctx context.Context) *expiresAt {
	ctx, span := otel.Tracer("session").Start(ctx, "NewExpiresAt")
	defer span.End()

	t := clock.Now(ctx).Add(2 * 24 * time.Hour * 7)
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
		RefreshSession(context.Context, shared.UUID[user.User]) (*Session, error)
		DeactivateUserSessions(context.Context, shared.UUID[user.User]) error
	}

	Session struct {
		sessionID    shared.UUID[Session]
		userID       shared.UUID[user.User]
		authProvider auth.AuthProviderName
		status       status
		expires      time.Time
		lastActivity time.Time
	}
)

func NewSession(
	sessionID shared.UUID[Session],
	userID shared.UUID[user.User],
	authProvider auth.AuthProviderName,
	status status,
	expires time.Time,
	lastActivity time.Time,
) *Session {
	return &Session{
		sessionID:    sessionID,
		userID:       userID,
		authProvider: authProvider,
		status:       status,
		expires:      expires,
		lastActivity: lastActivity,
	}
}

func (s *Session) UserID() shared.UUID[user.User] {
	return s.userID
}

func (s *Session) SessionID() shared.UUID[Session] {
	return s.sessionID
}

func (s *Session) Provider() auth.AuthProviderName {
	return s.authProvider
}

func (s *Session) Status() status {
	return s.status
}

func (s *Session) IsActive(ctx context.Context) bool {
	ctx, span := otel.Tracer("session").Start(ctx, "Session.IsActive")
	defer span.End()

	return s.expires.After(clock.Now(ctx)) && s.status == SESSION_ACTIVE
}

func (s *Session) Deactivate(ctx context.Context) {
	ctx, span := otel.Tracer("session").Start(ctx, "Session.Deactivate")
	defer span.End()

	s.status = SESSION_INACTIVE
	s.UpdateLastActivity(ctx)
}

func (s *Session) ExpiresAt() time.Time {
	return s.expires
}

func (s *Session) LastActivityAt() time.Time {
	return s.lastActivity
}

func (s *Session) UpdateLastActivity(ctx context.Context) {
	ctx, span := otel.Tracer("session").Start(ctx, "Session.UpdateLastActivity")
	defer span.End()

	s.lastActivity = clock.Now(ctx)
}

func SortByLastActivity(sessions []Session) []Session {
	sortedSession := make([]Session, len(sessions))
	copy(sortedSession, sessions)

	for i := 0; i < len(sortedSession); i++ {
		for j := i + 1; j < len(sortedSession); j++ {
			if sortedSession[i].lastActivity.Before(sortedSession[j].lastActivity) {
				sortedSession[i], sortedSession[j] = sortedSession[j], sortedSession[i]
			}
		}
	}

	return sessions
}
