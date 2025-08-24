package service

import (
	"context"
	"errors"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
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
	ctx, span := otel.Tracer("service").Start(ctx, "sessionService.DeactivateUserSessions")
	defer span.End()

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
	ctx, span := otel.Tracer("service").Start(ctx, "sessionService.RefreshSession")
	defer span.End()

	// sessionを取得
	sessList, err := s.sessionRepository.FindByUserID(ctx, userID)
	if err != nil {
		utils.HandleError(ctx, err, "sessionRepository.FindByUserID")
		return nil, errtrace.Wrap(err)
	}
	if len(sessList) == 0 {
		return nil, errtrace.Wrap(FailedToDeactivateSessionStatus)
	}

	sessList = session.FilterActiveSessions(ctx, session.SortByLastActivity(sessList))
	for _, sess := range sessList {
		// 最終アクティビティを更新
		sess.Deactivate(ctx)
		sess.UpdateLastActivity(ctx)
	}
	if err := s.sessionRepository.DeactivateAllByUserID(ctx, userID); err != nil {
		utils.HandleError(ctx, err, "sessionRepository.Update")
		return nil, errtrace.Wrap(err)
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
		utils.HandleError(ctx, err, "sessionRepository.Create")
		return nil, errtrace.Wrap(SessionRefreshFailed)
	}

	return updatedSess, nil
}

// SwitchOrganization implements session.SessionService.
func (s *sessionService) SwitchOrganization(
	ctx context.Context,
	userID shared.UUID[user.User],
	organizationID shared.UUID[organization.Organization],
) (*session.Session, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "sessionService.SwitchOrganization")
	defer span.End()

	// 既存セッションを取得してプロバイダー情報を引き継ぐ
	sessions, err := s.sessionRepository.FindByUserID(ctx, userID)
	if err != nil {
		utils.HandleError(ctx, err, "sessionRepository.FindByUserID")
		return nil, errtrace.Wrap(err)
	}
	// 最新のアクティブなセッションからプロバイダー情報を取得
	var provider shared.AuthProviderName
	activeSessions := session.FilterActiveSessions(ctx, session.SortByLastActivity(sessions))
	if len(activeSessions) > 0 {
		provider = activeSessions[0].Provider()
	}

	// 既存のアクティブなセッションをすべて無効化
	if err := s.sessionRepository.DeactivateAllByUserID(ctx, userID); err != nil {
		utils.HandleError(ctx, err, "sessionRepository.DeactivateAllByUserID")
		return nil, errtrace.Wrap(err)
	}

	// 新しい組織付きセッションを作成
	// organizationIDをany型にキャスト
	orgID := shared.UUID[any](organizationID.UUID())
	newSess := session.NewSessionWithOrganization(
		shared.NewUUID[session.Session](),
		userID,
		provider,
		session.SESSION_ACTIVE,
		*session.NewExpiresAt(ctx),
		clock.Now(ctx),
		&orgID,
	)

	createdSess, err := s.sessionRepository.Create(ctx, *newSess)
	if err != nil {
		utils.HandleError(ctx, err, "sessionRepository.Create")
		return nil, errtrace.Wrap(err)
	}

	return createdSess, nil
}

func NewSessionService(
	sessionRepository session.SessionRepository,
) session.SessionService {
	return &sessionService{
		sessionRepository: sessionRepository,
	}
}
