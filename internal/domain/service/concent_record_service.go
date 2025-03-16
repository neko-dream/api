package service

import (
	"context"
	"database/sql"
	"errors"
	"os/user"

	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/consent"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/pkg/utils"
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

	policy, err := c.policyRepository.FetchLatestPolicy(ctx)
	if err != nil {
		return false, err
	}

	record, err := c.consentRecordRepository.FindByUserAndVersion(ctx, userID.String(), policy.Version)
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

	// ポリシーが存在するかどうかを確認
	_, err := c.policyRepository.FindByVersion(ctx, version)
	if err != nil {
		utils.HandleError(ctx, err, "バージョンを取得できませんでした。")
		return nil, err
	}

	// ユーザーが同意済みかどうかを確認
	_, err = c.consentRecordRepository.FindByUserAndVersion(ctx, userID.String(), version)
	// ユーザーがすでに同意しているのなら、エラーを返す
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
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

	if err := c.consentRecordRepository.Save(ctx, record); err != nil {
		utils.HandleError(ctx, err, "同意を記録できませんでした。")
		return nil, err
	}

	return record, nil
}
