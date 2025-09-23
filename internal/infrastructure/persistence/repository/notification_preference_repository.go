package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
	"go.opentelemetry.io/otel"
)

type notificationPreferenceRepository struct {
	*db.DBManager
}

// NewNotificationPreferenceRepository creates a new notification preference repository
func NewNotificationPreferenceRepository(dbManager *db.DBManager) user.NotificationPreferenceRepository {
	return &notificationPreferenceRepository{
		DBManager: dbManager,
	}
}

// GetByUserIDs returns notification preferences for the given user IDs
func (r *notificationPreferenceRepository) GetByUserIDs(ctx context.Context, userIDs []shared.UUID[user.User]) (map[shared.UUID[user.User]]*user.NotificationPreference, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "notificationPreferenceRepository.GetByUserIDs")
	defer span.End()

	// Convert UUIDs to array
	uuidArray := make([]uuid.UUID, len(userIDs))
	for i, id := range userIDs {
		uuidArray[i] = id.UUID()
	}

	prefs, err := r.GetQueries(ctx).GetNotificationPreferencesByUserIDs(ctx, uuidArray)
	if err != nil {
		return nil, err
	}

	result := make(map[shared.UUID[user.User]]*user.NotificationPreference)
	// Map existing preferences
	for _, pref := range prefs {
		result[shared.UUID[user.User](pref.UserID)] = r.toDomainNotificationPreference(&pref)
	}

	for _, userID := range userIDs {
		if _, exists := result[userID]; !exists {
			result[userID] = r.getDefaultPreference(userID)
		}
	}

	return result, nil
}

// FindByUserID returns notification preference for the given user ID
func (r *notificationPreferenceRepository) FindByUserID(ctx context.Context, userID shared.UUID[user.User]) (*user.NotificationPreference, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "notificationPreferenceRepository.FindByUserID")
	defer span.End()

	pref, err := r.GetQueries(ctx).GetNotificationPreference(ctx, userID.UUID())
	if err != nil {
		if err == sql.ErrNoRows {
			// Return default preferences if not found
			return r.getDefaultPreference(userID), nil
		}
		return nil, err
	}

	return r.toDomainNotificationPreference(&pref), nil
}

func (r *notificationPreferenceRepository) Save(ctx context.Context, pref *user.NotificationPreference) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "notificationPreferenceRepository.Save")
	defer span.End()

	params := model.UpsertNotificationPreferenceParams{
		UserID:                  pref.UserID.UUID(),
		PushNotificationEnabled: pref.PushNotificationEnabled,
	}

	_, err := r.GetQueries(ctx).UpsertNotificationPreference(ctx, params)
	return err
}

// toDomainNotificationPreference converts from database model to domain model
func (r *notificationPreferenceRepository) toDomainNotificationPreference(pref *model.NotificationPreference) *user.NotificationPreference {
	return &user.NotificationPreference{
		UserID:                  shared.UUID[user.User](pref.UserID),
		PushNotificationEnabled: pref.PushNotificationEnabled,
	}
}

// getDefaultPreference returns default notification preferences
func (r *notificationPreferenceRepository) getDefaultPreference(userID shared.UUID[user.User]) *user.NotificationPreference {
	return &user.NotificationPreference{
		UserID:                  userID,
		PushNotificationEnabled: true,
	}
}
