package handler

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/neko-dream/api/internal/application/usecase/admin_usecase"
	"github.com/neko-dream/api/internal/presentation/oas"
	http_utils "github.com/neko-dream/api/pkg/http"
)

type failingAdminLogin struct{ err error }

func (f failingAdminLogin) Execute(context.Context, admin_usecase.AdminLoginInput) (*admin_usecase.AdminAuthOutput, error) {
	return nil, f.err
}

type failingAdminRefresh struct{ err error }

func (f failingAdminRefresh) Execute(context.Context, string) (*admin_usecase.AdminAuthOutput, error) {
	return nil, f.err
}

type failingAdminLogout struct{ err error }

func (f failingAdminLogout) Execute(context.Context, string) error { return f.err }

func TestAdminAuthInfrastructureErrorsAreNotCredentialFailures(t *testing.T) {
	dbErr := errors.New("database unavailable")
	h := adminHandler{adminLogin: failingAdminLogin{dbErr}, adminRefresh: failingAdminRefresh{dbErr}, adminLogout: failingAdminLogout{dbErr}}
	ctx := http_utils.WithHTTPResReqContext(context.Background(), httptest.NewRequest("POST", "/v1/admin/refresh", nil), httptest.NewRecorder())
	if res, err := h.AdminLogin(ctx, &oas.AdminLoginReq{}); res != nil || !errors.Is(err, dbErr) {
		t.Fatalf("login: %v %v", res, err)
	}
	if res, err := h.AdminRefresh(ctx); res != nil || !errors.Is(err, dbErr) {
		t.Fatalf("refresh: %v %v", res, err)
	}
	if res, err := h.AdminLogout(ctx); res != nil || !errors.Is(err, dbErr) {
		t.Fatalf("logout claimed success: %v %v", res, err)
	}
	h.adminLogin = failingAdminLogin{admin_usecase.ErrInvalidAdminCredentials}
	if res, err := h.AdminLogin(ctx, &oas.AdminLoginReq{}); err != nil {
		t.Fatal(err)
	} else if _, ok := res.(*oas.AdminLoginUnauthorized); !ok {
		t.Fatalf("unexpected login response %T", res)
	}
	h.adminRefresh = failingAdminRefresh{admin_usecase.ErrInvalidRefreshToken}
	if res, err := h.AdminRefresh(ctx); err != nil {
		t.Fatal(err)
	} else if _, ok := res.(*oas.AdminRefreshUnauthorized); !ok {
		t.Fatalf("unexpected refresh response %T", res)
	}
}
