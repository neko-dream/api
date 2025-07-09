package crypto

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/neko-dream/server/internal/domain/model/crypto"
	"github.com/neko-dream/server/internal/domain/model/notification"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"go.opentelemetry.io/otel"
)

func EncryptDevice(
	ctx context.Context,
	encryptor crypto.Encryptor,
	device notification.Device,
) (model.Device, error) {
	ctx, span := otel.Tracer("crypto").Start(ctx, "EncryptDevice")
	defer span.End()

	encryptedToken, err := encryptor.EncryptString(ctx, device.DeviceToken)
	if err != nil {
		return model.Device{}, err
	}
	var deviceName, appVersion, osVersion sql.NullString
	var LastActiveAt sql.NullTime
	if device.DeviceName != nil {
		deviceName = sql.NullString{
			String: *device.DeviceName,
			Valid:  true,
		}
	}
	if device.AppVersion != nil {
		appVersion = sql.NullString{
			String: *device.AppVersion,
			Valid:  true,
		}
	}
	if device.OsVersion != nil {
		osVersion = sql.NullString{
			String: *device.OsVersion,
			Valid:  true,
		}
	}
	if device.LastActiveAt != nil {
		LastActiveAt = sql.NullTime{
			Time:  *device.LastActiveAt,
			Valid: true,
		}
	}

	return model.Device{
		DeviceID:     device.ID.UUID(),
		UserID:       device.UserID.UUID(),
		DeviceToken:  encryptedToken,
		Platform:     string(device.Platform),
		Enabled:      device.Enabled,
		CreatedAt:    device.CreatedAt,
		UpdatedAt:    device.UpdatedAt,
		DeviceName:   deviceName,
		AppVersion:   appVersion,
		OsVersion:    osVersion,
		LastActiveAt: LastActiveAt,
	}, nil
}

func DecryptDevice(
	ctx context.Context,
	decryptor crypto.Encryptor,
	device model.Device,
) (*notification.Device, error) {
	ctx, span := otel.Tracer("crypto").Start(ctx, "DecryptDevice")
	defer span.End()

	decryptedToken, err := decryptor.DecryptString(ctx, device.DeviceToken)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt device token: %w", err)
	}

	var deviceName, appVersion, osVersion *string
	if device.DeviceName.Valid {
		deviceName = &device.DeviceName.String
	}
	if device.AppVersion.Valid {
		appVersion = &device.AppVersion.String
	}
	if device.OsVersion.Valid {
		osVersion = &device.OsVersion.String
	}

	var lastActiveAt *time.Time
	if device.LastActiveAt.Valid {
		lastActiveAt = &device.LastActiveAt.Time
	}

	return &notification.Device{
		ID:           shared.UUID[notification.Device](device.DeviceID),
		UserID:       shared.UUID[user.User](device.UserID),
		DeviceToken:  decryptedToken,
		Platform:     notification.DevicePlatform(device.Platform),
		Enabled:      device.Enabled,
		CreatedAt:    device.CreatedAt,
		UpdatedAt:    device.UpdatedAt,
		DeviceName:   deviceName,
		AppVersion:   appVersion,
		OsVersion:    osVersion,
		LastActiveAt: lastActiveAt,
	}, nil
}
