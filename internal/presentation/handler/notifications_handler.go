package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/neko-dream/server/internal/domain/model/notification"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/internal/presentation/oas"
	"go.opentelemetry.io/otel"
)

type notificationsHandler struct {
	dbManager                        *db.DBManager
	deviceRepository                 notification.DeviceRepository
	notificationPreferenceRepository user.NotificationPreferenceRepository
	authService                      service.AuthenticationService
	pushNotificationSender           notification.PushNotificationSender
	cfg                              *config.Config
	vapidKey                         string
	logger                           *slog.Logger
}

func NewNotificationsHandler(
	dbManager *db.DBManager,
	deviceRepository notification.DeviceRepository,
	notificationPreferenceRepository user.NotificationPreferenceRepository,
	authService service.AuthenticationService,
	pushNotificationSender notification.PushNotificationSender,
	cfg *config.Config,
) oas.NotificationsHandler {
	return &notificationsHandler{
		dbManager:                        dbManager,
		deviceRepository:                 deviceRepository,
		notificationPreferenceRepository: notificationPreferenceRepository,
		authService:                      authService,
		pushNotificationSender:           pushNotificationSender,
		cfg:                              cfg,
		vapidKey:                         cfg.FCM_VAPID_KEY,
		logger:                           slog.Default(),
	}
}

// RegisterDevice デバイス登録/更新
func (h *notificationsHandler) RegisterDevice(ctx context.Context, req *oas.RegisterDeviceReq) (oas.RegisterDeviceRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "notificationsHandler.RegisterDevice")
	defer span.End()

	// 認証チェック
	authCtx, err := requireAuthentication(h.authService, ctx)
	if err != nil {
		return &oas.RegisterDeviceUnauthorized{}, nil
	}

	exists, _ := h.checkDeviceTokenExists(ctx, req.DeviceToken)
	if exists != nil {
		var platformStr string
		switch req.Platform {
		case oas.RegisterDeviceReqPlatformIos:
			platformStr = "ios"
		case oas.RegisterDeviceReqPlatformAndroid:
			platformStr = "android"
		case oas.RegisterDeviceReqPlatformWeb:
			platformStr = "web"
		}
		response := &oas.Device{
			DeviceID:  exists.ID.String(),
			UserID:    exists.UserID.String(),
			Platform:  oas.DevicePlatform(platformStr),
			Enabled:   exists.Enabled,
			CreatedAt: exists.CreatedAt.Format(time.RFC3339),
			UpdatedAt: exists.UpdatedAt.Format(time.RFC3339),
		}
		if exists.DeviceName != nil {
			response.DeviceName = oas.NewOptString(*exists.DeviceName)
		}
		if exists.LastActiveAt != nil {
			response.LastActiveAt = oas.NewOptString(exists.LastActiveAt.Format(time.RFC3339))
		}
	}

	// プラットフォームを変換
	var platform notification.DevicePlatform
	switch req.Platform {
	case oas.RegisterDeviceReqPlatformIos:
		platform = notification.DevicePlatformAPNS
	case oas.RegisterDeviceReqPlatformAndroid:
		platform = notification.DevicePlatformGCM
	case oas.RegisterDeviceReqPlatformWeb:
		platform = notification.DevicePlatformWeb
	default:
		return &oas.RegisterDeviceBadRequest{}, nil
	}

	// デバイスを作成または更新
	device := notification.NewDevice(
		authCtx.UserID,
		req.DeviceToken,
		platform,
	)

	// オプションフィールドを設定
	if req.DeviceName.IsSet() {
		deviceName := req.DeviceName.Value
		device.DeviceName = &deviceName
	}
	if req.AppVersion.IsSet() {
		appVersion := req.AppVersion.Value
		device.AppVersion = &appVersion
	}
	if req.OsVersion.IsSet() {
		osVersion := req.OsVersion.Value
		device.OsVersion = &osVersion
	}

	// 保存
	if err := h.deviceRepository.Save(ctx, device); err != nil {
		h.logger.Error("デバイスの保存に失敗しました", slog.String("error", err.Error()))
		return &oas.RegisterDeviceBadRequest{}, nil
	}

	// レスポンスを作成
	// プラットフォームを文字列に変換してからDevicePlatformに変換
	var platformStr string
	switch req.Platform {
	case oas.RegisterDeviceReqPlatformIos:
		platformStr = "ios"
	case oas.RegisterDeviceReqPlatformAndroid:
		platformStr = "android"
	case oas.RegisterDeviceReqPlatformWeb:
		platformStr = "web"
	}
	response := &oas.Device{
		DeviceID:  device.ID.String(),
		UserID:    device.UserID.String(),
		Platform:  oas.DevicePlatform(platformStr),
		Enabled:   device.Enabled,
		CreatedAt: device.CreatedAt.Format(time.RFC3339),
		UpdatedAt: device.UpdatedAt.Format(time.RFC3339),
	}
	if device.DeviceName != nil {
		response.DeviceName = oas.NewOptString(*device.DeviceName)
	}
	if device.LastActiveAt != nil {
		response.LastActiveAt = oas.NewOptString(device.LastActiveAt.Format(time.RFC3339))
	}
	return response, nil
}

// GetDevices デバイス一覧取得
func (h *notificationsHandler) GetDevices(ctx context.Context) (oas.GetDevicesRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "notificationsHandler.GetDevices")
	defer span.End()

	// 認証チェック
	authCtx, err := requireAuthentication(h.authService, ctx)
	if err != nil {
		return &oas.GetDevicesUnauthorized{}, nil
	}
	userID := authCtx.UserID

	// デバイス一覧を取得
	devices, err := h.deviceRepository.FindByUserID(ctx, userID)
	if err != nil {
		h.logger.Error("デバイスの取得に失敗しました", slog.String("error", err.Error()))
		return &oas.GetDevicesUnauthorized{}, nil
	}

	// レスポンスを作成
	deviceList := make([]oas.Device, len(devices))
	for i, device := range devices {
		// プラットフォームを文字列に変換
		platform := "web"
		switch device.Platform {
		case notification.DevicePlatformAPNS:
			platform = "ios"
		case notification.DevicePlatformGCM:
			if device.DeviceName != nil && *device.DeviceName != "" {
				// デバイス名から判断（簡易的）
				platform = "android"
			}
		}

		deviceList[i] = oas.Device{
			DeviceID:  device.ID.String(),
			UserID:    device.UserID.String(),
			Platform:  oas.DevicePlatform(platform),
			Enabled:   device.Enabled,
			CreatedAt: device.CreatedAt.Format(time.RFC3339),
			UpdatedAt: device.UpdatedAt.Format(time.RFC3339),
		}
		if device.DeviceName != nil {
			deviceList[i].DeviceName = oas.NewOptString(*device.DeviceName)
		}
		if device.LastActiveAt != nil {
			deviceList[i].LastActiveAt = oas.NewOptString(device.LastActiveAt.Format(time.RFC3339))
		}
	}

	return &oas.GetDevicesOK{
		Devices: deviceList,
	}, nil
}

// DeleteDevice デバイス削除
func (h *notificationsHandler) DeleteDevice(ctx context.Context, params oas.DeleteDeviceParams) (oas.DeleteDeviceRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "notificationsHandler.DeleteDevice")
	defer span.End()

	authCtx, err := requireAuthentication(h.authService, ctx)
	if err != nil {
		return &oas.DeleteDeviceUnauthorized{}, nil
	}
	userID := authCtx.UserID

	deviceID, err := shared.ParseUUID[notification.Device](params.DeviceId)
	if err != nil {
		h.logger.Error("デバイスIDのパースに失敗しました", slog.String("error", err.Error()))
		return &oas.DeleteDeviceNotFound{}, nil
	}

	// デバイスの所有者確認と削除
	device, err := h.deviceRepository.FindByID(ctx, deviceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &oas.DeleteDeviceNotFound{}, nil
		}
		h.logger.Error("デバイスの取得に失敗しました", slog.String("error", err.Error()))
		return &oas.DeleteDeviceNotFound{}, nil
	}

	if device.UserID != userID {
		return &oas.DeleteDeviceNotFound{}, nil
	}

	if err := h.deviceRepository.Delete(ctx, deviceID); err != nil {
		h.logger.Error("デバイスの削除に失敗しました", slog.String("error", err.Error()))
		return &oas.DeleteDeviceNotFound{}, nil
	}

	return &oas.DeleteDeviceNoContent{}, nil
}

// GetNotificationPreferences 通知設定取得
func (h *notificationsHandler) GetNotificationPreferences(ctx context.Context) (oas.GetNotificationPreferencesRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "notificationsHandler.GetNotificationPreferences")
	defer span.End()

	// 認証チェック
	authCtx, err := requireAuthentication(h.authService, ctx)
	if err != nil {
		return &oas.GetNotificationPreferencesUnauthorized{}, nil
	}
	userID := authCtx.UserID

	// 通知設定を取得
	pref, err := h.notificationPreferenceRepository.FindByUserID(ctx, userID)
	if err != nil {
		h.logger.Error("通知設定の取得に失敗しました", slog.String("error", err.Error()))
		return &oas.GetNotificationPreferencesUnauthorized{}, nil
	}

	return &oas.NotificationPreferences{
		PushNotificationEnabled: pref.PushNotificationEnabled,
	}, nil
}

// UpdateNotificationPreferences 通知設定更新
func (h *notificationsHandler) UpdateNotificationPreferences(ctx context.Context, req *oas.UpdateNotificationPreferencesReq) (oas.UpdateNotificationPreferencesRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "notificationsHandler.UpdateNotificationPreferences")
	defer span.End()

	authCtx, err := requireAuthentication(h.authService, ctx)
	if err != nil {
		return &oas.UpdateNotificationPreferencesUnauthorized{}, nil
	}
	userID := authCtx.UserID

	// 現在の設定を取得
	pref, err := h.notificationPreferenceRepository.FindByUserID(ctx, userID)
	if err != nil {
		h.logger.Error("通知設定の取得に失敗しました", slog.String("error", err.Error()))
		return &oas.UpdateNotificationPreferencesBadRequest{}, nil
	}

	if req.PushNotificationEnabled.IsSet() {
		pref.PushNotificationEnabled = req.PushNotificationEnabled.Value
	}

	if err := h.notificationPreferenceRepository.Save(ctx, pref); err != nil {
		h.logger.Error("通知設定の保存に失敗しました", slog.String("error", err.Error()))
		return &oas.UpdateNotificationPreferencesBadRequest{}, nil
	}

	return &oas.NotificationPreferences{
		PushNotificationEnabled: pref.PushNotificationEnabled,
	}, nil
}

func (h *notificationsHandler) CheckDeviceExists(ctx context.Context, params oas.CheckDeviceExistsParams) (oas.CheckDeviceExistsRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "notificationsHandler.CheckDeviceExists")
	defer span.End()

	if params.DeviceToken == "" {
		return &oas.CheckDeviceExistsNotFound{}, nil
	}

	// 全てのアクティブなデバイスを取得してトークンをマッチング
	// 認証不要でデバイストークンの存在を確認する
	dev, err := h.checkDeviceTokenExists(ctx, params.DeviceToken)
	if err != nil {
		h.logger.Error("デバイストークンの存在確認に失敗しました",
			slog.String("error", err.Error()),
		)
		return &oas.CheckDeviceExistsOK{
			Exists: false,
		}, nil
	}

	// レスポンスを返す
	return &oas.CheckDeviceExistsOK{
		Exists: dev != nil,
	}, nil
}

// checkDeviceTokenExists デバイストークンが存在するか確認するヘルパーメソッド
func (h *notificationsHandler) checkDeviceTokenExists(ctx context.Context, deviceToken string) (*notification.Device, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "notificationsHandler.checkDeviceTokenExists")
	defer span.End()

	// 全てのアクティブなデバイスを取得
	devices, err := h.deviceRepository.GetAllActiveDevices(ctx)
	if err != nil {
		return nil, fmt.Errorf("アクティブデバイスの取得に失敗しました: %w", err)
	}

	// 各デバイスのトークンと比較
	for _, device := range devices {
		if device.DeviceToken == deviceToken {
			return device, nil
		}
	}

	return nil, nil
}

// SendTestNotification テスト通知送信
func (h *notificationsHandler) SendTestNotification(ctx context.Context, req *oas.SendTestNotificationReq) (oas.SendTestNotificationRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "notificationsHandler.SendTestNotification")
	defer span.End()

	// 環境チェック（dev/localのみ許可）
	if h.cfg.Env != config.LOCAL && h.cfg.Env != config.DEV {
		return &oas.SendTestNotificationBadRequest{}, nil
	}

	// 認証チェック
	authCtx, err := requireAuthentication(h.authService, ctx)
	if err != nil {
		return &oas.SendTestNotificationUnauthorized{}, nil
	}
	userID := authCtx.UserID

	// タイトルとボディの設定
	title := "Kotohiro テスト通知"
	body := fmt.Sprintf("これはテスト通知です。時刻: %s", time.Now().Format("15:04:05"))

	if req.Title.IsSet() && req.Title.Value != "" {
		title = req.Title.Value
	}
	if req.Body.IsSet() && req.Body.Value != "" {
		body = req.Body.Value
	}

	var devicesToSend []notification.Device

	// 特定のデバイスIDが指定されている場合
	if req.DeviceID.IsSet() && req.DeviceID.Value != "" {
		deviceID, err := shared.ParseUUID[notification.Device](req.DeviceID.Value)
		if err != nil {
			h.logger.Error("デバイスIDのパースに失敗しました", slog.String("error", err.Error()))
			return &oas.SendTestNotificationBadRequest{}, nil
		}

		device, err := h.deviceRepository.FindByID(ctx, deviceID)
		if err != nil {
			h.logger.Error("デバイスの取得に失敗しました", slog.String("error", err.Error()))
			return &oas.SendTestNotificationBadRequest{}, nil
		}

		// デバイスの所有者確認
		if device.UserID != userID {
			return &oas.SendTestNotificationBadRequest{}, nil
		}

		devicesToSend = append(devicesToSend, *device)
	} else {
		// ユーザーの全デバイスに送信
		devices, err := h.deviceRepository.FindByUserID(ctx, userID)
		if err != nil {
			h.logger.Error("デバイスの取得に失敗しました", slog.String("error", err.Error()))
			return &oas.SendTestNotificationBadRequest{}, nil
		}
		for _, device := range devices {
			devicesToSend = append(devicesToSend, *device)
		}
	}

	// 有効なデバイスのみフィルタリング
	var activeDevices []notification.Device
	for _, device := range devicesToSend {
		if device.Enabled {
			activeDevices = append(activeDevices, device)
		}
	}

	if len(activeDevices) == 0 {
		return &oas.SendTestNotificationOK{
			DevicesCount: 0,
			SuccessCount: 0,
		}, nil
	}

	// 通知を送信
	pushNotification := &notification.PushNotification{
		RecipientID: userID,
		Title:       title,
		Body:        body,
		Data: map[string]string{
			"type":      "test",
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}

	err = h.pushNotificationSender.Send(ctx, pushNotification)
	if err != nil {
		h.logger.Error("プッシュ通知の送信に失敗しました",
			slog.String("user_id", userID.String()),
			slog.String("error", err.Error()),
		)
		return &oas.SendTestNotificationOK{
			DevicesCount: int32(len(activeDevices)),
			SuccessCount: 0,
		}, nil
	}

	return &oas.SendTestNotificationOK{
		DevicesCount: int32(len(activeDevices)),
		SuccessCount: int32(len(activeDevices)),
	}, nil
}

// GetVapidKey VAPID公開鍵を取得
func (h *notificationsHandler) GetVapidKey(ctx context.Context) (*oas.GetVapidKeyOK, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "notificationsHandler.GetVapidKey")
	defer span.End()

	return &oas.GetVapidKeyOK{
		VapidKey: h.vapidKey,
	}, nil
}
