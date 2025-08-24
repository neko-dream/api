package handler

import (
	"context"
	"time"

	"github.com/neko-dream/server/internal/application/query/policy_query"
	"github.com/neko-dream/server/internal/application/usecase/policy_usecase"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/presentation/oas"
	http_utils "github.com/neko-dream/server/pkg/http"
	"github.com/neko-dream/server/pkg/utils"
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

	authCtx, err := h.authorizationService.RequireAuth(ctx)
	if err != nil {
		return nil, err
	}
	userID := authCtx.UserID

	output, err := h.checkConsentQuery.Execute(ctx, policy_query.CheckConsentInput{
		UserID: shared.UUID[user.User](userID),
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
	authCtx, err := h.authorizationService.RequireAuth(ctx)
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
