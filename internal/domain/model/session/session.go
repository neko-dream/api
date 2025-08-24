package session

import (
	"context"
	"time"

	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/organization"
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
		DeactivateAllByUserID(context.Context, shared.UUID[user.User]) error
		FindBySessionID(context.Context, shared.UUID[Session]) (*Session, error)
		FindByUserID(context.Context, shared.UUID[user.User]) ([]Session, error)
	}

	SessionService interface {
		RefreshSession(context.Context, shared.UUID[user.User]) (*Session, error)
		DeactivateUserSessions(context.Context, shared.UUID[user.User]) error
		SwitchOrganization(context.Context, shared.UUID[user.User], shared.UUID[organization.Organization]) (*Session, error)
	}

	Session struct {
		sessionID      shared.UUID[Session]
		userID         shared.UUID[user.User]
		authProvider   shared.AuthProviderName
		status         status
		expires        time.Time
		lastActivity   time.Time
		organizationID *shared.UUID[any] // ログイン時に使用した組織ID（組織経由ログインの場合）
	}
)

func NewSession(
	sessionID shared.UUID[Session],
	userID shared.UUID[user.User],
	authProvider shared.AuthProviderName,
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

func NewSessionWithOrganization(
	sessionID shared.UUID[Session],
	userID shared.UUID[user.User],
	authProvider shared.AuthProviderName,
	status status,
	expires time.Time,
	lastActivity time.Time,
	organizationID *shared.UUID[any],
) *Session {
	return &Session{
		sessionID:      sessionID,
		userID:         userID,
		authProvider:   authProvider,
		status:         status,
		expires:        expires,
		lastActivity:   lastActivity,
		organizationID: organizationID,
	}
}

func (s *Session) UserID() shared.UUID[user.User] {
	return s.userID
}

func (s *Session) SessionID() shared.UUID[Session] {
	return s.sessionID
}

func (s *Session) Provider() shared.AuthProviderName {
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

func (s *Session) OrganizationID() *shared.UUID[any] {
	return s.organizationID
}

func SortByLastActivity(sessions []Session) []Session {
	sortedSession := make([]Session, len(sessions))
	copy(sortedSession, sessions)

	for i := range sortedSession {
		for j := range sortedSession[i+1:] {
			actualJ := i + 1 + j
			if sortedSession[i].lastActivity.Before(sortedSession[actualJ].lastActivity) {
				sortedSession[i], sortedSession[actualJ] = sortedSession[actualJ], sortedSession[i]
			}
		}
	}

	return sessions
}

func FilterActiveSessions(ctx context.Context, sessions []Session) []Session {
	ctx, span := otel.Tracer("session").Start(ctx, "FilterActiveSessions")
	defer span.End()

	_ = ctx

	activeSessions := make([]Session, 0)
	for _, sess := range sessions {
		if sess.IsActive(context.Background()) {
			activeSessions = append(activeSessions, sess)
		}
	}
	return activeSessions
}
