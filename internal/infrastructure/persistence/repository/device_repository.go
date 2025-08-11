package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/neko-dream/server/internal/domain/model/crypto"
	"github.com/neko-dream/server/internal/domain/model/notification"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	crypto_impl "github.com/neko-dream/server/internal/infrastructure/crypto"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"go.opentelemetry.io/otel"
)

// DeviceRepository デバイスリポジトリの実装
type DeviceRepository struct {
	*db.DBManager
	crypto.Encryptor
}

// NewDeviceRepository デバイスリポジトリのコンストラクタ
func NewDeviceRepository(
	DBManager *db.DBManager,
	Encryptor crypto.Encryptor,
) notification.DeviceRepository {
	return &DeviceRepository{
		DBManager: DBManager,
		Encryptor: Encryptor,
	}
}

// Save デバイス情報を保存（既存の場合は更新）
func (r *DeviceRepository) Save(ctx context.Context, device *notification.Device) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "DeviceRepository.Save")
	defer span.End()

	encryptedModel, err := crypto_impl.EncryptDevice(ctx, r.Encryptor, *device)
	if err != nil {
		return fmt.Errorf("failed to encrypt device: %w", err)
	}

	_, err = r.GetQueries(ctx).UpsertDevice(ctx, model.UpsertDeviceParams{
		DeviceID:    encryptedModel.DeviceID,
		UserID:      encryptedModel.UserID,
		DeviceToken: encryptedModel.DeviceToken,
		Platform:    string(encryptedModel.Platform),
		Enabled:     encryptedModel.Enabled,
		CreatedAt:   encryptedModel.CreatedAt,
		UpdatedAt:   encryptedModel.UpdatedAt,
		DeviceName:  encryptedModel.DeviceName,
		AppVersion:  encryptedModel.AppVersion,
		OsVersion:   encryptedModel.OsVersion,
	})
	if err != nil {
		return fmt.Errorf("failed to save device: %w", err)
	}

	return nil
}

// FindByUserID ユーザーIDでデバイスを検索
func (r *DeviceRepository) FindByUserID(ctx context.Context, userID shared.UUID[user.User]) ([]*notification.Device, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "DeviceRepository.FindByUserID")
	defer span.End()

	rows, err := r.GetQueries(ctx).GetDevicesByUserID(ctx, userID.UUID())
	if err != nil {
		return nil, fmt.Errorf("failed to find devices by user ID: %w", err)
	}

	devices := make([]*notification.Device, 0, len(rows))
	for _, row := range rows {
		device, err := crypto_impl.DecryptDevice(ctx, r.Encryptor, row)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt device: %w", err)
		}
		devices = append(devices, device)
	}

	return devices, nil
}

// FindByID デバイスIDで検索
func (r *DeviceRepository) FindByID(ctx context.Context, deviceID shared.UUID[notification.Device]) (*notification.Device, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "DeviceRepository.FindByID")
	defer span.End()

	row, err := r.GetQueries(ctx).GetDeviceByID(ctx, deviceID.UUID())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("device not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find device by ID: %w", err)
	}

	return crypto_impl.DecryptDevice(ctx, r.Encryptor, row)
}

// Delete デバイスを削除
func (r *DeviceRepository) Delete(ctx context.Context, deviceID shared.UUID[notification.Device]) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "DeviceRepository.Delete")
	defer span.End()

	err := r.GetQueries(ctx).DeleteDevice(ctx, deviceID.UUID())
	if err != nil {
		return fmt.Errorf("failed to delete device: %w", err)
	}
	return nil
}

// GetActiveDevicesByUserIDs 複数のユーザーIDからアクティブなデバイスを取得
func (r *DeviceRepository) GetActiveDevicesByUserIDs(ctx context.Context, userIDs []shared.UUID[user.User]) ([]*notification.Device, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "DeviceRepository.GetActiveDevicesByUserIDs")
	defer span.End()

	if len(userIDs) == 0 {
		return []*notification.Device{}, nil
	}

	uuidArray := make([]uuid.UUID, len(userIDs))
	for i, id := range userIDs {
		uuidArray[i] = id.UUID()
	}

	rows, err := r.GetQueries(ctx).GetActiveDevicesByUserIDs(ctx, uuidArray)
	if err != nil {
		return nil, fmt.Errorf("failed to get active devices by user IDs: %w", err)
	}

	devices := make([]*notification.Device, 0, len(rows))
	for _, row := range rows {
		device, err := crypto_impl.DecryptDevice(ctx, r.Encryptor, row)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt device: %w", err)
		}
		devices = append(devices, device)
	}

	return devices, nil
}

// InvalidateDevice デバイスを無効化
func (r *DeviceRepository) InvalidateDevice(ctx context.Context, deviceID shared.UUID[notification.Device]) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "DeviceRepository.InvalidateDevice")
	defer span.End()

	err := r.GetQueries(ctx).InvalidateDevice(ctx, deviceID.UUID())
	if err != nil {
		return fmt.Errorf("failed to invalidate device: %w", err)
	}
	return nil
}

// GetAllActiveDevices 全てのアクティブなデバイスを取得
func (r *DeviceRepository) GetAllActiveDevices(ctx context.Context) ([]*notification.Device, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "DeviceRepository.GetAllActiveDevices")
	defer span.End()

	rows, err := r.GetQueries(ctx).GetAllActiveDevices(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all active devices: %w", err)
	}

	devices := make([]*notification.Device, 0, len(rows))
	for _, row := range rows {
		device, err := crypto_impl.DecryptDevice(ctx, r.Encryptor, row)
		if err != nil {
			// デバイストークンの復号に失敗した場合はスキップ
			// （古い暗号化方式のデータなど）
			continue
		}
		devices = append(devices, device)
	}

	return devices, nil
}
