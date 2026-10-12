package handler

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/presentation/oas"
)

type adminSecurityTokens struct {
	admin.AdminTokenManager
	claim *admin.AdminClaim
	err   error
}

type adminSecurityUsers struct{ admin.AdminUserRepository }

func (adminSecurityUsers) FindByID(context.Context, shared.UUID[admin.AdminUser]) (*admin.AdminUser, error) {
	return &admin.AdminUser{}, nil
}

func (m adminSecurityTokens) Parse(context.Context, string) (*admin.AdminClaim, error) {
	return m.claim, m.err
}

func TestAdminSecurityRefreshContract(t *testing.T) {
	id := shared.NewUUID[admin.AdminUser]()
	for _, tc := range []struct {
		name     string
		claim    *admin.AdminClaim
		parseErr error
		wantErr  error
	}{
		{"expired JWT", nil, fmt.Errorf("wrapped: %w", jwt.ErrTokenExpired), messages.TokenExpiredError},
		{"invalid signature", nil, jwt.ErrTokenSignatureInvalid, messages.ForbiddenError},
		{"nil claim", nil, nil, messages.ForbiddenError},
		{"expired claim", &admin.AdminClaim{Exp: time.Now().Add(-time.Minute).Unix()}, nil, messages.TokenExpiredError},
		{"valid", &admin.AdminClaim{Sub: id.String(), Exp: time.Now().Add(time.Hour).Unix()}, nil, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := securityHandler{adminTokenManager: adminSecurityTokens{claim: tc.claim, err: tc.parseErr}, adminUserRepo: adminSecurityUsers{}}
			ctx, err := h.HandleAdminCookieAuth(context.Background(), "GetAdminTokenInfo", oas.AdminCookieAuth{})
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("error %v want %v", err, tc.wantErr)
			}
			if tc.wantErr != nil && admin.GetAdminSession(ctx) != nil {
				t.Fatal("failed authentication populated context")
			}
			if tc.wantErr == nil && admin.GetAdminSession(ctx) != tc.claim {
				t.Fatal("authenticated context missing")
			}
		})
	}
}
