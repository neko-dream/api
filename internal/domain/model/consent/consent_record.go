package consent

import (
	"context"
	"time"

	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"go.opentelemetry.io/otel"
)

type ConsentRecordRepository interface {
	FindByUserAndVersion(ctx context.Context, userID shared.UUID[user.User], version string) (*ConsentRecord, error)
	Create(ctx context.Context, record *ConsentRecord) error
}

type ConsentService interface {
	RecordConsent(ctx context.Context, userID shared.UUID[user.User], version string, ipAddress string, userAgent string) (*ConsentRecord, error)
	IsConsentValid(ctx context.Context, userID shared.UUID[user.User]) (bool, error)
}

type ConsentRecord struct {
	ID          shared.UUID[ConsentRecord]
	UserID      shared.UUID[user.User]
	Version     string
	IP          string
	UA          string
	ConsentedAt time.Time
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
	ctx, span := otel.Tracer("consent").Start(ctx, "NewConsentRecord")
	defer span.End()

	_ = ctx

	return &ConsentRecord{
		ID:          id,
		UserID:      userID,
		Version:     version,
		IP:          ip,
		UA:          ua,
		ConsentedAt: createdAt,
	}
}
