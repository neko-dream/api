package user

import (
	"context"

	"github.com/neko-dream/api/internal/domain/model/shared"
)

type NotificationType string

const (
	NotificationTypeNewTalkSession NotificationType = "new_talk_session"
	NotificationTypeTalkSessionEnd NotificationType = "talk_session_end"
)

type NotificationPreference struct {
	UserID                  shared.UUID[User]
	PushNotificationEnabled bool // プッシュ通知の有効/無効
}

func (np *NotificationPreference) IsPushNotificationEnabled() bool {
	if np == nil {
		return true // デフォルトは有効
	}
	return np.PushNotificationEnabled
}

type NotificationPreferenceRepository interface {
	GetByUserIDs(ctx context.Context, userIDs []shared.UUID[User]) (map[shared.UUID[User]]*NotificationPreference, error)

	FindByUserID(ctx context.Context, userID shared.UUID[User]) (*NotificationPreference, error)

	Save(ctx context.Context, pref *NotificationPreference) error
}
