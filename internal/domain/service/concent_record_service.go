package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/consent"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

type consentService struct {
	consentRecordRepository consent.ConsentRecordRepository
	policyRepository        consent.PolicyRepository
}

func NewConsentService(
	consentRecordRepository consent.ConsentRecordRepository,
	policyRepository consent.PolicyRepository,
) consent.ConsentService {
	return &consentService{
		consentRecordRepository: consentRecordRepository,
		policyRepository:        policyRepository,
	}
}

// IsConsentValid ユーザーが最新のポリシーに同意済みかどうかを取得する
func (c *consentService) IsConsentValid(ctx context.Context, userID shared.UUID[user.User]) (bool, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "consentService.IsConsentValid")
	defer span.End()

	policy, err := c.policyRepository.FetchLatestPolicy(ctx)
	if err != nil {
		return false, err
	}

	record, err := c.consentRecordRepository.FindByUserAndVersion(ctx, userID, policy.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return record != nil, nil
}

// RecordConsent ユーザーの同意を記録する
func (c *consentService) RecordConsent(ctx context.Context, userID shared.UUID[user.User], version string, ipAddress string, userAgent string) (*consent.ConsentRecord, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "consentService.RecordConsent")
	defer span.End()

	// ポリシーが存在するかどうかを確認
	_, err := c.policyRepository.FindByVersion(ctx, version)
	if err != nil {
		utils.HandleError(ctx, err, "バージョンを取得できませんでした。")
		return nil, messages.PolicyNotFound
	}

	// ユーザーが同意済みかどうかを確認
	rec, err := c.consentRecordRepository.FindByUserAndVersion(ctx, userID, version)
	// ユーザーがすでに同意しているのなら、エラーを返す
	if !errors.Is(err, sql.ErrNoRows) || rec != nil {
		return nil, messages.PolicyAlreadyConsented
	}

	// 同意を記録
	record := consent.NewConsentRecord(
		ctx,
		shared.NewUUID[consent.ConsentRecord](),
		userID,
		version,
		ipAddress,
		userAgent,
		clock.Now(ctx),
	)

	if err := c.consentRecordRepository.Create(ctx, record); err != nil {
		utils.HandleError(ctx, err, "同意を記録できませんでした。")
		return nil, err
	}

	return record, nil
}
