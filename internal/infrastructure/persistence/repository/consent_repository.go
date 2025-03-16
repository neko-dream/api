package repository

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/consent"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type consentRecordRepository struct {
	*db.DBManager
}

func NewConsentRecordRepository(
	DBManager *db.DBManager,
) consent.ConsentRecordRepository {
	return &consentRecordRepository{
		DBManager: DBManager,
	}
}

// FindByUserAndVersion ユーザーとバージョンからConsentRecordを取得する
func (c *consentRecordRepository) FindByUserAndVersion(ctx context.Context, userID shared.UUID[user.User], version string) (*consent.ConsentRecord, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "consentRecordRepository.FindByUserAndVersion")
	defer span.End()

	resRow, err := c.DBManager.GetQueries(ctx).FindConsentByUserAndVersion(ctx, model.FindConsentByUserAndVersionParams{
		UserID:        userID.UUID(),
		PolicyVersion: version,
	})
	if err != nil {
		utils.HandleError(ctx, err, "同意情報を取得できませんでした。")
		return nil, err
	}

	return &consent.ConsentRecord{
		ID:          shared.UUID[consent.ConsentRecord](resRow.PolicyConsent.PolicyConsentID),
		UserID:      shared.UUID[user.User](resRow.PolicyConsent.UserID),
		Version:     resRow.PolicyConsent.PolicyVersion,
		IP:          resRow.PolicyConsent.IpAddress,
		UA:          resRow.PolicyConsent.UserAgent,
		ConsentedAt: resRow.PolicyConsent.ConsentedAt,
	}, nil
}

// Create  ConsentRecordを作成する
func (c *consentRecordRepository) Create(ctx context.Context, record *consent.ConsentRecord) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "consentRecordRepository.Create")
	defer span.End()

	err := c.DBManager.GetQueries(ctx).CreatePolicyConsent(ctx, model.CreatePolicyConsentParams{
		PolicyConsentID: record.ID.UUID(),
		UserID:          record.UserID.UUID(),
		PolicyVersion:   record.Version,
		IpAddress:       record.IP,
		UserAgent:       record.UA,
		ConsentedAt:     record.ConsentedAt,
	})
	if err != nil {
		utils.HandleError(ctx, err, "同意情報を保存できませんでした。")
		return err
	}

	return nil
}
