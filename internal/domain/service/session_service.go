package service

import (
	"context"
	"errors"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/clock"
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
		return errtrace.Wrap(err)
	}

	for _, sess := range sessions {
		sess.Deactivate(ctx)
	}

	return nil
}

// RefreshSession implements session.SessionService.
func (s *sessionService) RefreshSession(
	ctx context.Context,
	userID shared.UUID[user.User],
) (*session.Session, error) {
	// sessionを取得
	sessList, err := s.sessionRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	if len(sessList) == 0 {
		return nil, errtrace.Wrap(FailedToDeactivateSessionStatus)
	}

	sessList = session.SortByLastActivity(sessList)
	for _, sess := range sessList {
		// sessionが有効期限内であることを確認
		if !sess.IsActive(ctx) {
			return nil, errtrace.Wrap(SessionIsExpired)
		}

		// 最終アクティビティを更新
		sess.Deactivate(ctx)
		if _, err := s.sessionRepository.Update(ctx, sess); err != nil {
			return nil, errtrace.Wrap(err)
		}
	}

	// sessionを更新
	newSess := session.NewSession(
		shared.NewUUID[session.Session](),
		sessList[0].UserID(),
		sessList[0].Provider(),
		session.SESSION_ACTIVE,
		*session.NewExpiresAt(ctx),
		clock.Now(ctx),
	)
	updatedSess, err := s.sessionRepository.Create(ctx, *newSess)
	if err != nil {
		return nil, errtrace.Wrap(SessionRefreshFailed)
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
