package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/neko-dream/server/internal/domain/model/event"
	"github.com/neko-dream/server/internal/domain/model/notification"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"go.opentelemetry.io/otel"
)

type TalkSessionPushNotificationHandler struct {
	pushNotificationSender notification.PushNotificationSender
	userRepository         user.UserRepository
	talkSessionRepository  talksession.TalkSessionRepository
	logger                 *slog.Logger
}

func NewTalkSessionPushNotificationHandler(
	pushNotificationSender notification.PushNotificationSender,
	userRepository user.UserRepository,
	talkSessionRepository talksession.TalkSessionRepository,
) *TalkSessionPushNotificationHandler {
	return &TalkSessionPushNotificationHandler{
		pushNotificationSender: pushNotificationSender,
		userRepository:         userRepository,
		talkSessionRepository:  talkSessionRepository,
		logger:                 slog.Default(),
	}
}

// CanHandle このハンドラーがイベントを処理できるかチェック
func (h *TalkSessionPushNotificationHandler) CanHandle(eventType event.EventType) bool {
	return eventType == talksession.EventTypeTalkSessionStarted ||
		eventType == talksession.EventTypeTalkSessionEnded
}

func (h *TalkSessionPushNotificationHandler) Handle(ctx context.Context, storedEvent event.StoredEvent) error {
	ctx, span := otel.Tracer("handlers").Start(ctx, "TalkSessionPushNotificationHandler.Handle")
	defer span.End()

	switch storedEvent.EventType {
	case talksession.EventTypeTalkSessionStarted:
		return h.handleTalkSessionStarted(ctx, storedEvent)
	case talksession.EventTypeTalkSessionEnded:
		return h.handleTalkSessionEnded(ctx, storedEvent)
	default:
		return fmt.Errorf("未対応のイベントタイプ: %s", storedEvent.EventType)
	}
}

func (h *TalkSessionPushNotificationHandler) Priority() int {
	return 100 // 高優先度
}

// handleTalkSessionStarted セッション開始イベントを処理
func (h *TalkSessionPushNotificationHandler) handleTalkSessionStarted(ctx context.Context, storedEvent event.StoredEvent) error {
	var evt talksession.TalkSessionStartedEvent
	if err := json.Unmarshal(storedEvent.EventData, &evt); err != nil {
		return fmt.Errorf("イベントのデシリアライズに失敗しました: %w", err)
	}

	// セッション情報を取得
	session, err := h.talkSessionRepository.FindByID(ctx, evt.TalkSessionID)
	if err != nil {
		return fmt.Errorf("セッションの取得に失敗しました: %w", err)
	}
	if session == nil {
		return fmt.Errorf("セッションが見つかりません: %s", evt.TalkSessionID.String())
	}

	// 通知対象ユーザーを取得
	targetUsers, err := h.getNotificationTargetsForNewSession(ctx, session)
	if err != nil {
		return fmt.Errorf("通知対象ユーザーの取得に失敗しました: %w", err)
	}

	if len(targetUsers) == 0 {
		h.logger.Info("通知対象ユーザーがいません",
			slog.String("session_id", session.ID().String()),
		)
		return nil
	}

	// 通知を作成して送信
	notifications := h.createNewSessionNotifications(session, targetUsers)
	return h.pushNotificationSender.SendBatch(ctx, notifications)
}

// handleTalkSessionEnded トークセッション終了イベントを処理
func (h *TalkSessionPushNotificationHandler) handleTalkSessionEnded(ctx context.Context, storedEvent event.StoredEvent) error {
	// イベントデータをデシリアライズ
	var evt talksession.TalkSessionEndedEvent
	if err := json.Unmarshal(storedEvent.EventData, &evt); err != nil {
		return fmt.Errorf("イベントのデシリアライズに失敗しました: %w", err)
	}

	// トークセッション情報を取得
	session, err := h.talkSessionRepository.FindByID(ctx, evt.TalkSessionID)
	if err != nil {
		return fmt.Errorf("トークセッションの取得に失敗しました: %w", err)
	}
	if session == nil {
		return fmt.Errorf("トークセッションが見つかりません: %s", evt.TalkSessionID.String())
	}

	// 参加者全員に通知
	if len(evt.ParticipantIDs) == 0 {
		h.logger.Info("参加者がいません",
			slog.String("session_id", session.ID().String()),
		)
		return nil
	}

	// 通知を作成して送信
	notifications := h.createSessionEndNotifications(session, evt.ParticipantIDs)
	return h.pushNotificationSender.SendBatch(ctx, notifications)
}

// getNotificationTargetsForNewSession 新規セッションの通知対象を取得
func (h *TalkSessionPushNotificationHandler) getNotificationTargetsForNewSession(
	ctx context.Context,
	session *talksession.TalkSession,
) ([]shared.UUID[user.User], error) {
	// 運営（Neko Dream）のセッションかチェック
	if h.isKotohiroSession(session) {
		// TODO: 全アクティブユーザーを取得する機能を実装
		// 現在は空のリストを返す
		return []shared.UUID[user.User]{}, nil
	}

	// 運営以外のセッションは通知しない（将来的に拡張可能）
	return []shared.UUID[user.User]{}, nil
}

// isKotohiroSession 運営のセッションかチェック
func (h *TalkSessionPushNotificationHandler) isKotohiroSession(session *talksession.TalkSession) bool {
	if session.OrganizationID() == nil {
		return false
	}

	return session.OrganizationID().String() == organization.KotohiroOrganizationID.String()
}

// createNewSessionNotifications 新規セッション通知を作成
func (h *TalkSessionPushNotificationHandler) createNewSessionNotifications(
	session *talksession.TalkSession,
	targetUsers []shared.UUID[user.User],
) []*notification.PushNotification {
	notifications := make([]*notification.PushNotification, 0, len(targetUsers))
	for _, userID := range targetUsers {
		notif := notification.NewPushNotification(
			userID,
			notification.PushNotificationTypeNewTalkSession,
			"新しいセッションが開始しました",
			fmt.Sprintf("「%s」が開始しました", session.Theme()),
		)

		notif.AddData("talk_session_id", session.ID().String())
		notif.AddData("action", "open_talk_session")
		notifications = append(notifications, notif)
	}
	return notifications
}

// createSessionEndNotifications セッション終了通知を作成
func (h *TalkSessionPushNotificationHandler) createSessionEndNotifications(
	session *talksession.TalkSession,
	participantIDs []shared.UUID[user.User],
) []*notification.PushNotification {
	notifications := make([]*notification.PushNotification, 0, len(participantIDs))
	for _, userID := range participantIDs {
		notif := notification.NewPushNotification(
			userID,
			notification.PushNotificationTypeTalkSessionEnd,
			"セッションが終了しました！",
			fmt.Sprintf("「%s」が終了しました", session.Theme()),
		)
		// データを追加
		notif.AddData("talk_session_id", session.ID().String())
		notif.AddData("action", "open_talk_session_results")
		notifications = append(notifications, notif)
	}
	return notifications
}
