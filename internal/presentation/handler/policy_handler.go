package handler

import (
	"context"
	"time"

	"github.com/neko-dream/api/internal/application/query/policy_query"
	"github.com/neko-dream/api/internal/application/usecase/policy_usecase"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/service"
	"github.com/neko-dream/api/internal/presentation/oas"
	http_utils "github.com/neko-dream/api/pkg/http"
	"github.com/neko-dream/api/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type policyHandler struct {
	checkConsentQuery    policy_query.CheckConsent
	acceptPolicy         policy_usecase.AcceptPolicy
	authorizationService service.AuthorizationService
}

func NewPolicyHandler(
	checkConsentQuery policy_query.CheckConsent,
	acceptPolicy policy_usecase.AcceptPolicy,
	authorizationService service.AuthorizationService,
) oas.PolicyHandler {
	return &policyHandler{
		checkConsentQuery:    checkConsentQuery,
		acceptPolicy:         acceptPolicy,
		authorizationService: authorizationService,
	}
}

func (h *policyHandler) GetPolicyConsentStatus(ctx context.Context) (oas.GetPolicyConsentStatusRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "PolicyHandler.CheckConsent")
	defer span.End()

	authCtx, err := h.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	output, err := h.checkConsentQuery.Execute(ctx, policy_query.CheckConsentInput{
		UserID: authCtx.UserID,
	})
	if err != nil {
		return nil, err
	}
	var consentedAt *string
	if output.ConsentedAt != nil {
		consentedAt = lo.ToPtr(output.ConsentedAt.Format(time.RFC3339))
	}

	return &oas.PolicyConsentStatus{
		PolicyVersion: output.PolicyVersion,
		ConsentGiven:  output.ConsentGiven,
		ConsentedAt:   utils.ToOptNil[oas.OptNilString](consentedAt),
	}, nil
}

func (h *policyHandler) PolicyConsent(ctx context.Context, req *oas.PolicyConsentReq) (oas.PolicyConsentRes, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "PolicyHandler.AcceptPolicy")
	defer span.End()

	authCtx, err := h.authorizationService.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, messages.BadRequestError
	}

	request := http_utils.GetHTTPRequest(ctx)
	ipAddress := request.Header.Get("X-Forwarded-For")
	userAgent := request.Header.Get("User-Agent")

	output, err := h.acceptPolicy.Execute(ctx, policy_usecase.AcceptPolicyInput{
		UserID:    authCtx.UserID,
		Version:   req.PolicyVersion,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	})
	if err != nil {
		return nil, err
	}

	var consentedAt *string
	if output.ConsentedAt != nil {
		consentedAt = lo.ToPtr(output.ConsentedAt.Format(time.RFC3339))
	}

	return &oas.PolicyConsentStatus{
		PolicyVersion: req.PolicyVersion,
		ConsentedAt:   utils.ToOptNil[oas.OptNilString](consentedAt),
		ConsentGiven:  output.Success,
	}, nil
}
