package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/neko-dream/api/internal/application/usecase/admin_usecase"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/http/cookie"
	"github.com/neko-dream/api/internal/presentation/oas"
	httputils "github.com/neko-dream/api/pkg/http"
)

type capturingAdminRefresh struct {
	got string
	out *admin_usecase.AdminAuthOutput
}

func (c *capturingAdminRefresh) Execute(_ context.Context, token string) (*admin_usecase.AdminAuthOutput, error) {
	c.got = token
	return c.out, nil
}

func serveAdminRefresh(t *testing.T, h *adminHandler, cookieHeader string) *http.Response {
	t.Helper()
	srv, err := oas.NewServer(&handlers{AdminHandler: h}, NewSecurityHandler(nil, nil, nil, nil), oas.WithErrorHandler(CustomErrorHandler))
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest(http.MethodPost, "/v1/admin/refresh", nil)
	if cookieHeader != "" {
		r.Header.Set("Cookie", cookieHeader)
	}
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, r.WithContext(httputils.WithHTTPResReqContext(r.Context(), r, w)))
	return w.Result()
}

// ogen は Set-Cookie の値を URL エンコードするため、base64（+ / =）のリフレッシュトークンは
// %2B などに変わった値でブラウザに保存されていた。保存済みの Cookie も使えるよう、読み取り時にデコードする。
func TestAdminRefreshDecodesURLEncodedCookie(t *testing.T) {
	refresh := &capturingAdminRefresh{out: &admin_usecase.AdminAuthOutput{AccessToken: "access", RefreshToken: "next_token-1"}}
	h := &adminHandler{cookieManager: cookie.NewCookieManager(&config.Config{Env: config.LOCAL}), adminRefresh: refresh}

	res := serveAdminRefresh(t, h, cookie.AdminRefreshCookieName+"=Q4Whk%2BUbhG%2FpBen4%3D")
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", res.StatusCode)
	}
	if refresh.got != "Q4Whk+UbhG/pBen4=" {
		t.Fatalf("token passed to usecase = %q", refresh.got)
	}
	// URL で使える文字だけのトークンは、エンコードされずにそのまま保存される
	for _, c := range res.Cookies() {
		if c.Name == cookie.AdminRefreshCookieName && c.Value != "next_token-1" {
			t.Fatalf("refresh cookie value = %q", c.Value)
		}
	}
}

func TestAdminRefreshRejectsMalformedEncodedCookie(t *testing.T) {
	refresh := &capturingAdminRefresh{out: &admin_usecase.AdminAuthOutput{AccessToken: "access", RefreshToken: "next"}}
	h := &adminHandler{cookieManager: cookie.NewCookieManager(&config.Config{Env: config.LOCAL}), adminRefresh: refresh}
	serveAdminRefresh(t, h, cookie.AdminRefreshCookieName+"=broken%ZZ")
	if refresh.got != "" {
		t.Fatalf("malformed cookie passed to usecase as %q", refresh.got)
	}
}
