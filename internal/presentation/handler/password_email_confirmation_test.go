package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/neko-dream/api/internal/application/usecase/auth_usecase"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/http/cookie"
	"github.com/neko-dream/api/internal/presentation/oas"
	http_utils "github.com/neko-dream/api/pkg/http"
)

type emailLoginHandlerFake struct {
	begins, verifies int
	browser, token   string
}

func (f *emailLoginHandlerFake) Execute(context.Context, auth_usecase.PasswordLoginInput) (*auth_usecase.PasswordLoginOutput, error) {
	f.begins++
	return &auth_usecase.PasswordLoginOutput{BrowserSecret: "pending-secret", ExpiresAt: time.Now().Add(15 * time.Minute)}, nil
}
func (f *emailLoginHandlerFake) Verify(_ context.Context, token, browser string) (string, error) {
	f.verifies++
	f.token = token
	f.browser = browser
	return "signed-session", nil
}
func (f *emailLoginHandlerFake) VerifyPIN(ctx context.Context, pin, browser string) (string, error) {
	return f.Verify(ctx, pin, browser)
}
func emailHandlerHarness(env config.ENV) (*authHandler, *emailLoginHandlerFake) {
	cfg := &config.Config{Env: env, DOMAIN: "kotohiro.com", WEBSITE_URL: "https://localhost:3000", EMAIL_VERIFICATION_BASE_URL: "https://localhost:3000/api"}
	fake := &emailLoginHandlerFake{}
	return &authHandler{passwordLogin: fake, cfg: cfg, CookieManager: cookie.NewCookieManager(cfg)}, fake
}
func emailHandlerContext(method, origin, pending string) (context.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, "https://localhost:3000/api/auth/password/verify", nil)
	if origin != "" {
		req.Header.Set("Origin", origin)
	}
	if pending != "" {
		req.AddCookie(&http.Cookie{Name: passwordLoginPendingCookie, Value: pending})
	}
	w := httptest.NewRecorder()
	return http_utils.WithHTTPResReqContext(context.Background(), req, w), w
}

func TestPasswordLoginIssuesOnlyPendingCookie(t *testing.T) {
	a, fake := emailHandlerHarness(config.LOCAL)
	ctx, w := emailHandlerContext("POST", "https://localhost:3000", "")
	out, err := a.PasswordLogin(ctx, &oas.PasswordLoginReq{IdOrEmail: "city@example.com", Password: "password"})
	if err != nil {
		t.Fatal(err)
	}
	accepted := out.(*oas.PasswordLoginAcceptedHeaders)
	if fake.begins != 1 || len(accepted.SetCookie) != 1 || !strings.HasPrefix(accepted.SetCookie[0], passwordLoginPendingCookie+"=") || strings.Contains(accepted.SetCookie[0], "SessionId=") {
		t.Fatal("password alone issued session cookie")
	}
	if !strings.Contains(accepted.SetCookie[0], "Secure") || !strings.Contains(accepted.SetCookie[0], "HttpOnly") || w.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("pending cookie/security headers missing")
	}
	for _, env := range []config.ENV{config.DEV, config.PROD} {
		a, _ := emailHandlerHarness(env)
		ctx, _ := emailHandlerContext("POST", "https://localhost:3000", "")
		res, err := a.PasswordLogin(ctx, &oas.PasswordLoginReq{})
		if err != nil {
			t.Fatal(err)
		}
		value := res.(*oas.PasswordLoginAcceptedHeaders).SetCookie[0]
		if strings.Contains(strings.ToLower(value), "domain=") || !strings.Contains(value, "Path=/") || !strings.HasPrefix(value, "__Host-") {
			t.Fatal("pending cookie is not restricted to its issuing host")
		}
	}
}

func TestEmailConfirmationGETDoesNotConsume(t *testing.T) {
	for _, env := range []config.ENV{config.LOCAL, config.DEV, config.PROD, "", "prod", "dev"} {
		t.Run(string(env), func(t *testing.T) {
			a, fake := emailHandlerHarness(env)
			ctx, w := emailHandlerContext("GET", "", "")
			res, err := a.PasswordEmailConfirmation(ctx)
			if err != nil {
				t.Fatal(err)
			}
			body, err := io.ReadAll(res.Data)
			if err != nil {
				t.Fatal(err)
			}
			if fake.verifies != 0 || fake.begins != 0 {
				t.Fatal("GET consumed a login")
			}
			if strings.Contains(string(body), "{{NONCE}}") || !strings.Contains(w.Header().Get("Content-Security-Policy"), "script-src 'nonce-") {
				t.Fatal("CSP nonce missing")
			}
			if (env == config.LOCAL) != strings.Contains(string(body), `id="login"`) {
				t.Fatal("local login form leaked into deployed environment")
			}
			if w.Header().Get("Referrer-Policy") != "no-referrer" || w.Header().Get("Cache-Control") != "no-store" {
				t.Fatal("unsafe verification page headers")
			}
		})
	}
}

func TestEmailConfirmationRequiresCookieAndTrustedOrigin(t *testing.T) {
	for _, tc := range []struct {
		name, origin, pending string
		want                  error
	}{{"no pending cookie", "https://localhost:3000", "", messages.InvalidEmailLoginLinkError}, {"cross-site", "https://evil.example", "pending", messages.ForbiddenError}, {"null origin", "null", "pending", messages.ForbiddenError}} {
		t.Run(tc.name, func(t *testing.T) {
			a, fake := emailHandlerHarness(config.LOCAL)
			ctx, _ := emailHandlerContext("POST", tc.origin, tc.pending)
			if _, err := a.VerifyPasswordEmail(ctx, &oas.VerifyPasswordEmailReq{Token: "email-token"}); !errors.Is(err, tc.want) {
				t.Fatal(err)
			}
			if fake.verifies != 0 {
				t.Fatal("invalid request reached verification")
			}
		})
	}
	a, fake := emailHandlerHarness(config.LOCAL)
	ctx, _ := emailHandlerContext("POST", "https://localhost:3000", "pending")
	res, err := a.VerifyPasswordEmail(ctx, &oas.VerifyPasswordEmailReq{Token: "email-token"})
	if err != nil {
		t.Fatal(err)
	}
	out := res.(*oas.VerifyPasswordEmailOKHeaders)
	if fake.token != "email-token" || fake.browser != "pending" || out.Response.RedirectURL != "https://localhost:3000/home" || len(out.SetCookie) != 2 {
		t.Fatal("verification result incorrect")
	}
	if !strings.HasPrefix(out.SetCookie[0], "SessionId=") || !strings.HasPrefix(out.SetCookie[1], passwordLoginPendingCookie+"=") || !strings.Contains(out.SetCookie[1], "Expires=") {
		t.Fatal("session establishment or pending cleanup missing")
	}
	ctx, _ = emailHandlerContext("POST", "https://evil.example", "")
	if _, err := a.PasswordLogin(ctx, &oas.PasswordLoginReq{}); !errors.Is(err, messages.ForbiddenError) {
		t.Fatal("login CSRF accepted")
	}
}

func TestPINConfirmationRequiresCookieAndTrustedOrigin(t *testing.T) {
	for _, tc := range []struct {
		origin, pending string
		want            error
	}{
		{"https://localhost:3000", "", messages.InvalidEmailLoginPINError},
		{"https://evil.example", "pending", messages.ForbiddenError},
		{"null", "pending", messages.ForbiddenError},
	} {
		a, fake := emailHandlerHarness(config.LOCAL)
		ctx, _ := emailHandlerContext("POST", tc.origin, tc.pending)
		if _, err := a.VerifyPasswordEmailPIN(ctx, &oas.VerifyPasswordEmailPINReq{Pin: "000001"}); !errors.Is(err, tc.want) {
			t.Fatal(err)
		}
		if fake.verifies != 0 {
			t.Fatal("invalid PIN request reached verifier")
		}
	}
	a, fake := emailHandlerHarness(config.LOCAL)
	ctx, w := emailHandlerContext("POST", "https://localhost:3000", "pending")
	res, err := a.VerifyPasswordEmailPIN(ctx, &oas.VerifyPasswordEmailPINReq{Pin: "000001"})
	if err != nil {
		t.Fatal(err)
	}
	out := res.(*oas.VerifyPasswordEmailPINOKHeaders)
	if fake.token != "000001" || fake.browser != "pending" || out.Response.RedirectURL != "https://localhost:3000/home" || len(out.SetCookie) != 2 {
		t.Fatal("incorrect PIN result")
	}
	if !strings.HasPrefix(out.SetCookie[0], "SessionId=") || !strings.Contains(out.SetCookie[1], "Max-Age=0") || w.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("PIN security cookies/headers missing")
	}
}
