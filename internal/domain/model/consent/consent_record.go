package consent

import (
	"context"
	"os/user"
	"time"

	"github.com/neko-dream/server/internal/domain/model/shared"
)

type ConsentRecordRepository interface {
	FindByID(ctx context.Context, id string) (*ConsentRecord, error)
	FindByUserAndVersion(ctx context.Context, userID string, version string) (*ConsentRecord, error)
	Save(ctx context.Context, record *ConsentRecord) error
}

type ConsentService interface {
	RecordConsent(ctx context.Context, userID shared.UUID[user.User], version string, ipAddress string, userAgent string) (*ConsentRecord, error)
	IsConsentValid(ctx context.Context, userID shared.UUID[user.User]) (bool, error)
}

type ConsentRecord struct {
	Id        shared.UUID[ConsentRecord]
	UserID    shared.UUID[user.User]
	Version   string
	IP        string
	UA        string
	CreatedAt time.Time
}

func NewConsentRecord(
	ctx context.Context,
	id shared.UUID[ConsentRecord],
	userID shared.UUID[user.User],
	version string,
	ip string,
	ua string,
	createdAt time.Time,
) *ConsentRecord {
	return &ConsentRecord{
		Id:        id,
		UserID:    userID,
		Version:   version,
		IP:        ip,
		UA:        ua,
		CreatedAt: createdAt,
	}
}
