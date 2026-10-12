package notification

import (
	"context"
	"time"

	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
)

// PushNotificationType プッシュ通知のタイプ
type PushNotificationType string

const (
	// PushNotificationTypeNewTalkSession 新規セッション
	PushNotificationTypeNewTalkSession PushNotificationType = "new_talk_session"
	// PushNotificationTypeTalkSessionEnd セッション終了
	PushNotificationTypeTalkSessionEnd PushNotificationType = "talk_session_end"
)

type PushNotification struct {
	RecipientID shared.UUID[user.User]
	Type        PushNotificationType
	Title       string
	Body        string
	Data        map[string]string
	Priority    string
	Sound       string
	CreatedAt   time.Time
}

// NewPushNotification プッシュ通知を生成
func NewPushNotification(
	recipientID shared.UUID[user.User],
	notificationType PushNotificationType,
	title string,
	body string,
) *PushNotification {
	return &PushNotification{
		RecipientID: recipientID,
		Type:        notificationType,
		Title:       title,
		Body:        body,
		Data:        make(map[string]string),
		Priority:    "high",
		Sound:       "default",
		CreatedAt:   time.Now(),
	}
}

// AddData データを追加
func (n *PushNotification) AddData(key, value string) {
	n.Data[key] = value
}

// SetPriority 優先度を設定
func (n *PushNotification) SetPriority(priority string) {
	n.Priority = priority
}

// SetSound 通知音を設定
func (n *PushNotification) SetSound(sound string) {
	n.Sound = sound
}

// PushNotificationSender プッシュ通知送信インターフェース
type PushNotificationSender interface {
	// Send 単一のプッシュ通知を送信
	Send(ctx context.Context, notification *PushNotification) error
	// SendBatch 複数のプッシュ通知をバッチ送信
	SendBatch(ctx context.Context, notifications []*PushNotification) error
}

// DevicePlatform デバイスプラットフォーム
type DevicePlatform string

const (
	DevicePlatformAPNS DevicePlatform = "APNS"
	DevicePlatformGCM  DevicePlatform = "GCM"
	DevicePlatformWeb  DevicePlatform = "web"
)

// Device デバイス情報
type Device struct {
	ID           shared.UUID[Device]
	UserID       shared.UUID[user.User]
	DeviceToken  string
	Platform     DevicePlatform
	DeviceName   *string
	AppVersion   *string
	OsVersion    *string
	Enabled      bool
	LastActiveAt *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewDevice デバイス情報を生成
func NewDevice(
	userID shared.UUID[user.User],
	deviceToken string,
	platform DevicePlatform,
) *Device {
	now := time.Now()
	return &Device{
		ID:          shared.NewUUID[Device](),
		UserID:      userID,
		DeviceToken: deviceToken,
		Platform:    platform,
		Enabled:     true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Disable デバイスを無効化
func (d *Device) Disable() {
	d.Enabled = false
	d.UpdatedAt = time.Now()
}

// Enable デバイスを有効化
func (d *Device) Enable() {
	d.Enabled = true
	d.UpdatedAt = time.Now()
}

type DeviceRepository interface {
	// Save デバイス情報を保存
	Save(ctx context.Context, device *Device) error
	// FindByUserID ユーザーIDでデバイスを検索
	FindByUserID(ctx context.Context, userID shared.UUID[user.User]) ([]*Device, error)
	// FindByID デバイスIDで検索
	FindByID(ctx context.Context, deviceID shared.UUID[Device]) (*Device, error)
	// Delete デバイスを削除
	Delete(ctx context.Context, deviceID shared.UUID[Device]) error
	// GetActiveDevicesByUserIDs 複数のユーザーIDからアクティブなデバイスを取得
	GetActiveDevicesByUserIDs(ctx context.Context, userIDs []shared.UUID[user.User]) ([]*Device, error)
	// InvalidateDevice デバイスを無効化
	InvalidateDevice(ctx context.Context, deviceID shared.UUID[Device]) error
	// GetAllActiveDevices 全てのアクティブなデバイスを取得
	GetAllActiveDevices(ctx context.Context) ([]*Device, error)
}
