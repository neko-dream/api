package policy_query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/consent"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/pkg/utils"
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
	}
)

func NewCheckConsent(
	consentService consent.ConsentService,
	consentRecordRep consent.ConsentRecordRepository,
	policyRep consent.PolicyRepository,
) CheckConsent {
	return &checkConsentInteractor{
		consentService:   consentService,
		consentRecordRep: consentRecordRep,
		policyRep:        policyRep,
	}
}

func (c *checkConsentInteractor) Execute(ctx context.Context, input CheckConsentInput) (*CheckConsentOutput, error) {
	ctx, span := otel.Tracer("policy_query").Start(ctx, "checkConsentInteractor.Execute")
	defer span.End()

	// 最新のポリシーのバージョンを取得
	policy, err := c.policyRep.FetchLatestPolicy(ctx)
	if err != nil {
		utils.HandleError(ctx, err, "ポリシーを取得できませんでした。")
		return nil, messages.PolicyFetchFailed
	}

	rec, err := c.consentRecordRep.FindByUserAndVersion(ctx, input.UserID, policy.Version)
	if errors.Is(err, sql.ErrNoRows) {
		return &CheckConsentOutput{
			ConsentGiven:  false,
			PolicyVersion: policy.Version,
		}, nil
	} else if err != nil {
		utils.HandleError(ctx, err, "同意情報を取得できませんでした。")
		return nil, messages.PolicyFetchFailed
	}

	return &CheckConsentOutput{
		PolicyVersion: policy.Version,
		ConsentGiven:  rec != nil,
		ConsentedAt:   lo.ToPtr(rec.ConsentedAt),
	}, nil
}
