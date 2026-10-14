package handler

import (
	"context"

	"github.com/neko-dream/api/internal/application/usecase/auth_usecase"
	"github.com/neko-dream/api/internal/domain/messages"
	password_auth "github.com/neko-dream/api/internal/domain/model/auth/password"
	"github.com/neko-dream/api/internal/presentation/oas"
	http_utils "github.com/neko-dream/api/pkg/http"
)

func (a *authHandler) GetPasswordPolicy(ctx context.Context) (*oas.GetPasswordPolicyOK, error) {
	setEmailLoginSecurityHeaders(http_utils.GetHTTPResponse(ctx))
	policy, err := password_auth.NewPasswordPolicy(len(a.cfg.HASH_PEPPER))
	if err != nil {
		return nil, err
	}
	return &oas.GetPasswordPolicyOK{
		MinLength:                   int32(policy.MinLength),
		MaxBytes:                    int32(policy.MaxBytes),
		ForbidControlCharacters:     policy.ForbidControlCharacters,
		ForbidSurroundingWhitespace: policy.ForbidSurroundingWhitespace,
	}, nil
}

func (a *authHandler) SetupPassword(ctx context.Context, req *oas.SetupPasswordReq) error {
	setEmailLoginSecurityHeaders(http_utils.GetHTTPResponse(ctx))
	if err := a.validateEmailLoginOrigin(ctx); err != nil {
		return err
	}
	if req == nil {
		return messages.BadRequestError
	}
	return a.setupPassword.Execute(ctx, auth_usecase.SetupPasswordInput{Token: req.Token, Password: req.Password})
}
