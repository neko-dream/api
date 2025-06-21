package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"braces.dev/errtrace"
	"github.com/google/uuid"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type sessionRepository struct {
	*db.DBManager
}

// DeactivateAllByUserID implements session.SessionRepository.
func (s *sessionRepository) DeactivateAllByUserID(ctx context.Context, userID shared.UUID[user.User]) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "sessionRepository.DeactivateAllByUserID")
	defer span.End()

	if err := s.GetQueries(ctx).DeactivateSessions(ctx, userID.UUID()); err != nil {
		return errtrace.Wrap(err)
	}

	return nil
}

// Create implements session.SessionRepository.
func (s *sessionRepository) Create(ctx context.Context, sess session.Session) (*session.Session, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "sessionRepository.Create")
	defer span.End()

	params := model.CreateSessionParams{
		SessionID:      sess.SessionID().UUID(),
		UserID:         sess.UserID().UUID(),
		Provider:       sess.Provider().String(),
		SessionStatus:  int32(sess.Status()),
		ExpiresAt:      sess.ExpiresAt(),
		LastActivityAt: sess.LastActivityAt(),
		CreatedAt:      clock.Now(ctx),
	}

	// organization_idの設定
	if sess.OrganizationID() != nil && !sess.OrganizationID().IsZero() {
		params.OrganizationID = uuid.NullUUID{
			UUID:  sess.OrganizationID().UUID(),
			Valid: true,
		}
	}

	if err := s.GetQueries(ctx).CreateSession(ctx, params); err != nil {
		return nil, errtrace.Wrap(err)
	}

	return &sess, nil
}

// FindBySessionID セッションIDを元にセッションを取得する
func (s *sessionRepository) FindBySessionID(ctx context.Context, sess shared.UUID[session.Session]) (*session.Session, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "sessionRepository.FindBySessionID")
	defer span.End()

	sessRow, err := s.GetQueries(ctx).FindSessionBySessionID(ctx, sess.UUID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errtrace.Wrap(err)
	}

	providerName, err := shared.NewAuthProviderName(sessRow.Provider)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	// organization_idの処理
	var orgID *shared.UUID[any]
	if sessRow.OrganizationID.Valid {
		id := shared.UUID[any](sessRow.OrganizationID.UUID)
		orgID = &id
	}

	if orgID != nil {
		return session.NewSessionWithOrganization(
			shared.UUID[session.Session](sessRow.SessionID),
			shared.UUID[user.User](sessRow.UserID),
			providerName,
			*session.NewSessionStatus(int(sessRow.SessionStatus)),
			sessRow.ExpiresAt,
			sessRow.LastActivityAt,
			orgID,
		), nil
	}

	return session.NewSession(
		shared.UUID[session.Session](sessRow.SessionID),
		shared.UUID[user.User](sessRow.UserID),
		providerName,
		*session.NewSessionStatus(int(sessRow.SessionStatus)),
		sessRow.ExpiresAt,
		sessRow.LastActivityAt,
	), nil
}

// FindByUserID implements session.SessionRepository.
func (s *sessionRepository) FindByUserID(ctx context.Context, userID shared.UUID[user.User]) ([]session.Session, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "sessionRepository.FindByUserID")
	defer span.End()

	sessionModels, err := s.GetQueries(ctx).FindActiveSessionsByUserID(ctx, userID.UUID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errtrace.Wrap(err)
	}

	sessions := make([]session.Session, 0, len(sessionModels))
	for _, sess := range sessionModels {
		providerName, err := shared.NewAuthProviderName(sess.Provider)
		if err != nil {
			utils.HandleError(ctx, err, fmt.Sprintf("NewAuthProviderName: %s", sess.Provider))
			continue
		}
		// organization_idの処理
		var orgID *shared.UUID[any]
		if sess.OrganizationID.Valid {
			id := shared.UUID[any](sess.OrganizationID.UUID)
			orgID = &id
		}

		if orgID != nil {
			sessions = append(sessions, *session.NewSessionWithOrganization(
				shared.UUID[session.Session](sess.SessionID),
				shared.UUID[user.User](sess.UserID),
				providerName,
				*session.NewSessionStatus(int(sess.SessionStatus)),
				sess.ExpiresAt,
				sess.LastActivityAt,
				orgID,
			))
		} else {
			sessions = append(sessions, *session.NewSession(
				shared.UUID[session.Session](sess.SessionID),
				shared.UUID[user.User](sess.UserID),
				providerName,
				*session.NewSessionStatus(int(sess.SessionStatus)),
				sess.ExpiresAt,
				sess.LastActivityAt,
			))
		}
	}

	return sessions, nil
}

// Update セッションの状態と最終アクティビティ時間を更新する
func (s *sessionRepository) Update(ctx context.Context, sess session.Session) (*session.Session, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "sessionRepository.Update")
	defer span.End()

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
