package repository

import (
	"context"
	"database/sql"
	"errors"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
)

type sessionRepository struct {
	*db.DBManager
}

// Create implements session.SessionRepository.
func (s *sessionRepository) Create(ctx context.Context, sess session.Session) (*session.Session, error) {
	if err := s.GetQueries(ctx).CreateSession(ctx, model.CreateSessionParams{
		SessionID:      sess.SessionID().UUID(),
		UserID:         sess.UserID().UUID(),
		Provider:       sess.Provider().String(),
		SessionStatus:  int32(sess.Status()),
		ExpiresAt:      sess.ExpiresAt(),
		LastActivityAt: sess.LastActivityAt(),
		CreatedAt:      clock.Now(ctx),
	}); err != nil {
		return nil, errtrace.Wrap(err)
	}

	return &sess, nil
}

// FindBySessionID セッションIDを元にセッションを取得する
func (s *sessionRepository) FindBySessionID(ctx context.Context, sess shared.UUID[session.Session]) (*session.Session, error) {
	sessRow, err := s.GetQueries(ctx).FindSessionBySessionID(ctx, sess.UUID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errtrace.Wrap(err)
	}

	return session.NewSession(
		shared.UUID[session.Session](sessRow.SessionID),
		shared.UUID[user.User](sessRow.UserID),
		auth.AuthProviderName(sessRow.Provider),
		*session.NewSessionStatus(int(sessRow.SessionStatus)),
		sessRow.ExpiresAt,
		sessRow.LastActivityAt,
	), nil
}

// FindByUserID implements session.SessionRepository.
func (s *sessionRepository) FindByUserID(ctx context.Context, userID shared.UUID[user.User]) ([]session.Session, error) {
	sessionModels, err := s.GetQueries(ctx).FindActiveSessionsByUserID(ctx, userID.UUID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errtrace.Wrap(err)
	}

	sessions := make([]session.Session, 0, len(sessionModels))
	for _, sess := range sessionModels {
		sessions = append(sessions, *session.NewSession(
			shared.UUID[session.Session](sess.SessionID),
			shared.UUID[user.User](sess.UserID),
			auth.AuthProviderName(sess.Provider),
			*session.NewSessionStatus(int(sess.SessionStatus)),
			sess.ExpiresAt,
			sess.LastActivityAt,
		))
	}

	return sessions, nil
}

// Update セッションの状態と最終アクティビティ時間を更新する
func (s *sessionRepository) Update(ctx context.Context, sess session.Session) (*session.Session, error) {
	if err := s.GetQueries(ctx).UpdateSession(ctx, model.UpdateSessionParams{
		SessionID:      sess.SessionID().UUID(),
		SessionStatus:  int32(sess.Status()),
		LastActivityAt: sess.LastActivityAt(),
	}); err != nil {
		return nil, errtrace.Wrap(err)
	}

	return &sess, nil
}

func NewSessionRepository(
	tm *db.DBManager,
) session.SessionRepository {
	return &sessionRepository{
		DBManager: tm,
	}
}
