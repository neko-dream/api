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

// 空のトークンは本物のユースケースと同じく無効として扱う
func (c *capturingAdminRefresh) Execute(_ context.Context, token string) (*admin_usecase.AdminAuthOutput, error) {
	c.got = token
	if token == "" {
		return nil, admin_usecase.ErrInvalidRefreshToken
	}
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

func responseCookies(res *http.Response) map[string]string {
	got := map[string]string{}
	for _, c := range res.Cookies() {
		got[c.Name] = c.Value
	}
	return got
}

// 以前は cookie_utils.EncodeCookies が値を URL エンコードしていたため、base64（+ / =）のリフレッシュトークンが
// %2B などに変わった値でブラウザに保存されていた。保存済みの Cookie も使えるよう、読み取り時にデコードする。
func TestAdminRefreshDecodesURLEncodedCookie(t *testing.T) {
	const next = "Nx+Rot/ated=="
	refresh := &capturingAdminRefresh{out: &admin_usecase.AdminAuthOutput{AccessToken: "access", RefreshToken: next}}
	h := &adminHandler{cookieManager: cookie.NewCookieManager(&config.Config{Env: config.LOCAL}), adminRefresh: refresh}

	res := serveAdminRefresh(t, h, cookie.AdminRefreshCookieName+"=Q4Whk%2BUbhG%2FpBen4%3D")
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", res.StatusCode)
	}
	if refresh.got != "Q4Whk+UbhG/pBen4=" {
		t.Fatalf("token passed to usecase = %q", refresh.got)
	}
	// 新しいリフレッシュトークンは、+ / = を含んでもそのままの値でブラウザに保存される
	got := responseCookies(res)
	if v, ok := got[cookie.AdminRefreshCookieName]; !ok || v != next {
		t.Fatalf("refresh cookie = %q (present=%v), want %q; Set-Cookie = %q", v, ok, next, res.Header.Values("Set-Cookie"))
	}
	if got[cookie.AdminSessionCookieName] != "access" {
		t.Fatalf("access cookie = %q", got[cookie.AdminSessionCookieName])
	}
}

func TestAdminRefreshRejectsMalformedEncodedCookie(t *testing.T) {
	refresh := &capturingAdminRefresh{out: &admin_usecase.AdminAuthOutput{AccessToken: "access", RefreshToken: "next"}}
	h := &adminHandler{cookieManager: cookie.NewCookieManager(&config.Config{Env: config.LOCAL}), adminRefresh: refresh}

	res := serveAdminRefresh(t, h, cookie.AdminRefreshCookieName+"=broken%ZZ")
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", res.StatusCode)
	}
	if refresh.got != "" {
		t.Fatalf("malformed cookie passed to usecase as %q", refresh.got)
	}
	if _, ok := responseCookies(res)[cookie.AdminRefreshCookieName]; ok {
		t.Fatal("refresh cookie must not be issued for a malformed cookie")
	}
}
