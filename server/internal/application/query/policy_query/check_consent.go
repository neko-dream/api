package policy_query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/consent"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"github.com/neko-dream/api/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type (
	CheckConsent interface {
		Execute(context.Context, CheckConsentInput) (*CheckConsentOutput, error)
	}

	CheckConsentInput struct {
		UserID shared.UUID[user.User]
	}

	CheckConsentOutput struct {
		PolicyVersion string
		ConsentedAt   *time.Time
		ConsentGiven  bool
	}

	checkConsentInteractor struct {
		consentService   consent.ConsentService
		consentRecordRep consent.ConsentRecordRepository
		policyRep        consent.PolicyRepository
		dbm              *db.DBManager
	}
)

func NewCheckConsent(
	consentService consent.ConsentService,
	consentRecordRep consent.ConsentRecordRepository,
	policyRep consent.PolicyRepository,
	dbm *db.DBManager,
) CheckConsent {
	return &checkConsentInteractor{
		consentService:   consentService,
		consentRecordRep: consentRecordRep,
		policyRep:        policyRep,
		dbm:              dbm,
	}
}

func (c *checkConsentInteractor) Execute(ctx context.Context, input CheckConsentInput) (*CheckConsentOutput, error) {
	ctx, span := otel.Tracer("policy_query").Start(ctx, "checkConsentInteractor.Execute")
	defer span.End()

	var recOut *CheckConsentOutput

	err := c.dbm.ExecTx(ctx, func(ctx context.Context) error {
		// 最新のポリシーのバージョンを取得
		policy, err := c.policyRep.FetchLatestPolicy(ctx)
		if err != nil {
			utils.HandleError(ctx, err, "ポリシーを取得できませんでした。")
			return messages.PolicyFetchFailed
		}

		rec, err := c.consentRecordRep.FindByUserAndVersion(ctx, input.UserID, policy.Version)
		if errors.Is(err, sql.ErrNoRows) {
			recOut = lo.ToPtr(CheckConsentOutput{
				ConsentGiven:  false,
				PolicyVersion: policy.Version,
			})
			return nil
		} else if err != nil {
			utils.HandleError(ctx, err, "同意情報を取得できませんでした。")
			return messages.PolicyFetchFailed
		}

		recOut = lo.ToPtr(CheckConsentOutput{
			ConsentGiven:  true,
			PolicyVersion: policy.Version,
			ConsentedAt:   lo.ToPtr(rec.ConsentedAt),
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return recOut, nil
}
