package policy_command

import (
	"context"
	"time"

	"github.com/neko-dream/server/internal/domain/model/consent"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"go.opentelemetry.io/otel"
)

type (
	AcceptPolicy interface {
		Execute(context.Context, AcceptPolicyInput) (*AcceptPolicyOutput, error)
	}

	AcceptPolicyInput struct {
		UserID    shared.UUID[user.User]
		Version   string
		IPAddress string
		UserAgent string
	}

	AcceptPolicyOutput struct {
		Success     bool
		ConsentedAt *time.Time
	}

	acceptPolicyInteractor struct {
		consentService consent.ConsentService
	}
)

func NewAcceptPolicy(
	consentService consent.ConsentService,
) AcceptPolicy {
	return &acceptPolicyInteractor{
		consentService: consentService,
	}
}

func (a *acceptPolicyInteractor) Execute(ctx context.Context, input AcceptPolicyInput) (*AcceptPolicyOutput, error) {
	ctx, span := otel.Tracer("policy_query").Start(ctx, "acceptPolicyInteractor.Execute")
	defer span.End()

	rec, err := a.consentService.RecordConsent(
		ctx,
		input.UserID,
		input.Version,
		input.IPAddress,
		input.UserAgent,
	)
	if err != nil {
		return nil, err
	}

	return &AcceptPolicyOutput{
		Success:     true,
		ConsentedAt: &rec.ConsentedAt,
	}, nil
}
