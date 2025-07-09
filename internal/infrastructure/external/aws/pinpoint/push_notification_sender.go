package pinpoint

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/pinpoint"
	"github.com/aws/aws-sdk-go-v2/service/pinpoint/types"
	"github.com/neko-dream/server/internal/domain/model/notification"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"go.opentelemetry.io/otel"
)

// PushNotificationSender AWS Pinpointを使用したプッシュ通知送信実装
type PushNotificationSender struct {
	client                     *pinpoint.Client
	applicationID              string
	deviceRepository           notification.DeviceRepository
	notificationPrefRepository user.NotificationPreferenceRepository
	logger                     *slog.Logger
}

// NewPushNotificationSender コンストラクタ
func NewPushNotificationSender(
	client *pinpoint.Client,
	cfg *config.Config,
	deviceRepository notification.DeviceRepository,
	notificationPrefRepository user.NotificationPreferenceRepository,
) notification.PushNotificationSender {
	return &PushNotificationSender{
		client:                     client,
		applicationID:              cfg.PINPOINT_APPLICATION_ID,
		deviceRepository:           deviceRepository,
		notificationPrefRepository: notificationPrefRepository,
		logger:                     slog.Default(),
	}
}

// Send 単一のプッシュ通知を送信
func (s *PushNotificationSender) Send(ctx context.Context, notification *notification.PushNotification) error {
	ctx, span := otel.Tracer("pinpoint").Start(ctx, "PushNotificationSender.Send")
	defer span.End()

	if shouldSkip, reason := s.shouldSkipNotification(ctx, notification); shouldSkip {
		s.logger.Info("通知をスキップします",
			slog.String("user_id", notification.RecipientID.String()),
			slog.String("reason", reason),
		)
		return nil
	}

	devices, err := s.deviceRepository.FindByUserID(ctx, notification.RecipientID)
	if err != nil {
		return fmt.Errorf("デバイス情報の取得に失敗しました: %w", err)
	}

	if len(devices) == 0 {
		s.logger.Info("送信先デバイスが見つかりません",
			slog.String("user_id", notification.RecipientID.String()),
		)
		return nil
	}

	// 有効なデバイスのみフィルタリング
	activeDevices := s.filterActiveDevices(devices)
	if len(activeDevices) == 0 {
		s.logger.Info("有効なデバイスが見つかりません",
			slog.String("user_id", notification.RecipientID.String()),
		)
		return nil
	}

	return s.sendToPinpoint(ctx, notification, activeDevices)
}

// SendBatch 複数のプッシュ通知をバッチ送信
func (s *PushNotificationSender) SendBatch(ctx context.Context, notifications []*notification.PushNotification) error {
	ctx, span := otel.Tracer("pinpoint").Start(ctx, "PushNotificationSender.SendBatch")
	defer span.End()

	for _, notif := range notifications {
		if err := s.Send(ctx, notif); err != nil {
			s.logger.Error("通知送信に失敗しました",
				slog.String("user_id", notif.RecipientID.String()),
				slog.String("error", err.Error()),
			)
			// エラーが発生しても他の通知は続行
			continue
		}
	}
	return nil
}

// shouldSkipNotification 通知をスキップすべきか判定
func (s *PushNotificationSender) shouldSkipNotification(
	ctx context.Context,
	notification *notification.PushNotification,
) (bool, string) {
	// 通知設定を取得
	preference, err := s.notificationPrefRepository.FindByUserID(ctx, notification.RecipientID)
	if err != nil {
		s.logger.Error("通知設定の取得に失敗しました",
			slog.String("user_id", notification.RecipientID.String()),
			slog.String("error", err.Error()),
		)
		return true, "設定取得エラー"
	}

	// 設定が見つからない場合はデフォルトで送信
	if preference == nil {
		return false, ""
	}

	// プッシュ通知が無効化されている場合
	if !preference.IsPushNotificationEnabled() {
		return true, "プッシュ通知が無効"
	}

	return false, ""
}

// filterActiveDevices 有効なデバイスのみフィルタリング
func (s *PushNotificationSender) filterActiveDevices(devices []*notification.Device) []*notification.Device {
	var activeDevices []*notification.Device
	for _, device := range devices {
		if device.Enabled {
			activeDevices = append(activeDevices, device)
		}
	}
	return activeDevices
}

// sendToPinpoint AWS Pinpoint経由で通知を送信
func (s *PushNotificationSender) sendToPinpoint(
	ctx context.Context,
	notification *notification.PushNotification,
	devices []*notification.Device,
) error {
	addresses := make(map[string]types.AddressConfiguration)
	for _, device := range devices {
		channelType := s.getChannelType(device.Platform)
		addresses[device.DeviceToken] = types.AddressConfiguration{
			ChannelType: channelType,
		}
	}

	dataPayload := make(map[string]string)
	for k, v := range notification.Data {
		dataPayload[k] = v
	}
	dataPayload["title"] = notification.Title
	dataPayload["body"] = notification.Body

	messageRequest := &pinpoint.SendMessagesInput{
		ApplicationId: aws.String(s.applicationID),
		MessageRequest: &types.MessageRequest{
			Addresses: addresses,
			MessageConfiguration: &types.DirectMessageConfiguration{
				GCMMessage: &types.GCMMessage{
					Title:    aws.String(notification.Title),
					Body:     aws.String(notification.Body),
					Data:     dataPayload,
					Priority: aws.String(notification.Priority),
				},
			},
		},
	}

	// 送信実行
	output, err := s.client.SendMessages(ctx, messageRequest)
	if err != nil {
		return fmt.Errorf("Pinpoint送信エラー: %w", err)
	}

	s.processSendResults(ctx, output, devices)

	return nil
}

// getChannelType プラットフォームからチャンネルタイプを取得
func (s *PushNotificationSender) getChannelType(platform notification.DevicePlatform) types.ChannelType {
	switch platform {
	case notification.DevicePlatformAPNS:
		return types.ChannelTypeApns
	case notification.DevicePlatformGCM:
		return types.ChannelTypeGcm
	default:
		return types.ChannelTypeGcm
	}
}

// processSendResults 送信結果を処理
func (s *PushNotificationSender) processSendResults(
	ctx context.Context,
	output *pinpoint.SendMessagesOutput,
	devices []*notification.Device,
) {
	if output.MessageResponse == nil || output.MessageResponse.Result == nil {
		return
	}

	// デバイストークンとデバイスIDのマップを作成
	tokenToDevice := make(map[string]*notification.Device)
	for _, device := range devices {
		tokenToDevice[device.DeviceToken] = device
	}

	// 各送信結果を処理
	for address, result := range output.MessageResponse.Result {
		if result.StatusCode == nil {
			continue
		}

		statusCode := *result.StatusCode
		if statusCode >= 400 {
			s.logger.Warn("通知送信失敗",
				slog.String("address", address),
				slog.Int("status_code", int(statusCode)),
				slog.String("message", aws.ToString(result.StatusMessage)),
			)

			// 410 Gone の場合はデバイスを無効化
			if statusCode == 410 {
				if device, ok := tokenToDevice[address]; ok {
					s.handleInvalidDevice(ctx, device)
				}
			}
		} else {
			s.logger.Debug("通知送信成功",
				slog.String("address", address),
				slog.Int("status_code", int(statusCode)),
			)
		}
	}
}

// handleInvalidDevice 無効なデバイスを処理
func (s *PushNotificationSender) handleInvalidDevice(ctx context.Context, device *notification.Device) {
	device.Disable()
	if err := s.deviceRepository.Save(ctx, device); err != nil {
		s.logger.Error("デバイス無効化の保存に失敗しました",
			slog.String("device_id", device.ID.String()),
			slog.String("error", err.Error()),
		)
	}
}
