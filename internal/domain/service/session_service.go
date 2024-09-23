package service

import (
	"context"
	"errors"

	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type sessionService struct {
	sessionRepository session.SessionRepository
}

var (
	SessionIsExpired                = errors.New("セッションの期限が切れています。再ログインしてください")
	FailedToDeactivateSessionStatus = errors.New("セッションステータスの無効化に失敗しました")
	SessionRefreshFailed            = errors.New("セッションの更新に失敗しました。再ログインしてください")
)

func (s *sessionService) DeactivateUserSessions(ctx context.Context, userID shared.UUID[user.User]) error {
	sessions, err := s.sessionRepository.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}

	for _, sess := range sessions {
		if err := sess.Deactivate(ctx); err != nil {
			return err
		}
	}

	return nil
}

// RefreshSession implements session.SessionService.
func (s *sessionService) RefreshSession(ctx context.Context, sess session.Session) (*session.Session, error) {
	// sessionが有効期限内であることを確認
	if !sess.IsActive() {
		return nil, SessionIsExpired
	}

	if err := sess.Deactivate(ctx); err != nil {
		return nil, FailedToDeactivateSessionStatus
	}

	// 最終アクティビティを更新
	sess.UpdateLastActivity()
	if _, err := s.sessionRepository.Update(ctx, sess); err != nil {
		return nil, err
	}

	// sessionを更新
	newSess := session.NewSession(
		shared.NewUUID[session.Session](),
		sess.UserID(),
		sess.Provider(),
		session.SESSION_ACTIVE,
		*session.NewExpiresAt(),
	)
	updatedSess, err := s.sessionRepository.Create(ctx, *newSess)
	if err != nil {
		return nil, SessionRefreshFailed
	}

	return updatedSess, nil
}

func NewSessionService(
	sessionRepository session.SessionRepository,
) session.SessionService {
	return &sessionService{
		sessionRepository: sessionRepository,
	}
}