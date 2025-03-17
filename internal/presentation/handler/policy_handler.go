package handler

import (
	"context"
	"time"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/neko-dream/server/internal/usecase/command/policy_command"
	"github.com/neko-dream/server/internal/usecase/query/policy_query"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type policyHandler struct {
	checkConsentQuery policy_query.CheckConsent
	acceptPolicy      policy_command.AcceptPolicy
}

func NewPolicyHandler(
	checkConsentQuery policy_query.CheckConsent,
	acceptPolicy policy_command.AcceptPolicy,
) oas.PolicyHandler {
	return &policyHandler{
		checkConsentQuery: checkConsentQuery,
		acceptPolicy:      acceptPolicy,
	}
}

func (h *policyHandler) GetPolicyConsentStatus(ctx context.Context) (oas.GetPolicyConsentStatusRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "PolicyHandler.CheckConsent")
	defer span.End()

	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}

	userID, err := claim.UserID()
	if err != nil {
		utils.HandleError(ctx, err, "claim.UserID")
		return nil, messages.ForbiddenError
	}

	output, err := h.checkConsentQuery.Execute(ctx, policy_query.CheckConsentInput{
		UserID: shared.UUID[user.User](userID),
	})
	if err != nil {
		return nil, err
	}

	return &oas.GetPolicyConsentStatusOK{
		PolicyVersion: output.PolicyVersion,
		ConsentGiven:  output.ConsentGiven,
		ConsentedAt:   utils.ToOptNil[oas.OptNilString](output.ConsentedAt.Format(time.RFC3339)),
	}, nil
}

func (h *policyHandler) PolicyConsent(ctx context.Context, params oas.PolicyConsentParams) (oas.PolicyConsentRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "PolicyHandler.AcceptPolicy")
	defer span.End()

	claim := session.GetSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}

	userID, err := claim.UserID()
	if err != nil {
		utils.HandleError(ctx, err, "claim.UserID")
		return nil, messages.ForbiddenError
	}

	var ip, ua string
	if v, ok := params.XFowarderdFor.Get(); ok {
		ip = v
	}
	if v, ok := params.UserAgent.Get(); ok {
		ua = v
	}

	output, err := h.acceptPolicy.Execute(ctx, policy_command.AcceptPolicyInput{
		UserID:    shared.UUID[user.User](userID),
		Version:   params.PolicyVersion,
		IPAddress: ip,
		UserAgent: ua,
	})
	if err != nil {
		return nil, err
	}

	return &oas.PolicyConsentOK{
		PolicyVersion: params.PolicyVersion,
		ConsentedAt:   utils.ToOptNil[oas.OptNilString](output.ConsentedAt.Format(time.RFC3339)),
		ConsentGiven:  output.Success,
	}, nil
}
