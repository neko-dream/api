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
	logger                           *slog.Logger
}

func NewNotificationsHandler(
	dbManager *db.DBManager,
	deviceRepository notification.DeviceRepository,
	notificationPreferenceRepository user.NotificationPreferenceRepository,
	authService service.AuthenticationService,
	pushNotificationSender notification.PushNotificationSender,
) oas.NotificationsHandler {
	return &notificationsHandler{
		dbManager:                        dbManager,
		deviceRepository:                 deviceRepository,
		notificationPreferenceRepository: notificationPreferenceRepository,
		authService:                      authService,
		pushNotificationSender:           pushNotificationSender,
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

	// プラットフォームを変換
	var platform notification.DevicePlatform
	switch req.Platform {
	case oas.RegisterDeviceReqPlatformIos:
		platform = notification.DevicePlatformAPNS
	case oas.RegisterDeviceReqPlatformAndroid:
		platform = notification.DevicePlatformGCM
	case oas.RegisterDeviceReqPlatformWeb:
		platform = notification.DevicePlatformGCM // WebもFCMを使用
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
	exists, err := h.checkDeviceTokenExists(ctx, params.DeviceToken)
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
		Exists: exists,
	}, nil
}

// checkDeviceTokenExists デバイストークンが存在するか確認するヘルパーメソッド
func (h *notificationsHandler) checkDeviceTokenExists(ctx context.Context, deviceToken string) (bool, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "notificationsHandler.checkDeviceTokenExists")
	defer span.End()

	// 全てのアクティブなデバイスを取得
	devices, err := h.deviceRepository.GetAllActiveDevices(ctx)
	if err != nil {
		return false, fmt.Errorf("アクティブデバイスの取得に失敗しました: %w", err)
	}

	// 各デバイスのトークンと比較
	for _, device := range devices {
		if device.DeviceToken == deviceToken {
			return true, nil
		}
	}

	return false, nil
}
