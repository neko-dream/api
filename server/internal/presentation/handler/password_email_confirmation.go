package handler

import (
	"context"
	"crypto/rand"
	"embed"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/presentation/oas"
	cookie_utils "github.com/neko-dream/api/pkg/cookie"
	http_utils "github.com/neko-dream/api/pkg/http"
)

const passwordLoginPendingCookie = "__Host-PasswordLoginPending"

//go:embed password_email_confirmation.html
var emailConfirmationFiles embed.FS

func setEmailLoginSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
}

func (a *authHandler) PasswordEmailConfirmation(ctx context.Context) (oas.PasswordEmailConfirmationOK, error) {
	body, err := emailConfirmationFiles.ReadFile("password_email_confirmation.html")
	if err != nil {
		return oas.PasswordEmailConfirmationOK{}, err
	}
	var secret [24]byte
	if _, err := rand.Read(secret[:]); err != nil {
		return oas.PasswordEmailConfirmationOK{}, err
	}
	nonce := base64.RawStdEncoding.EncodeToString(secret[:])
	w := http_utils.GetHTTPResponse(ctx)
	setEmailLoginSecurityHeaders(w)
	w.Header().Set("Content-Security-Policy", fmt.Sprintf("default-src 'none'; script-src 'nonce-%s'; style-src 'nonce-%s'; connect-src 'self'; form-action 'self'; base-uri 'none'; frame-ancestors 'none'", nonce, nonce))
	page := strings.ReplaceAll(string(body), "{{NONCE}}", nonce)
	if a.cfg.Env != config.LOCAL {
		start := strings.Index(page, "<!--LOCAL-START-->")
		end := strings.Index(page, "<!--LOCAL-END-->")
		if start < 0 || end <= start {
			return oas.PasswordEmailConfirmationOK{}, fmt.Errorf("local-only form boundaries are missing")
		}
		page = page[:start] + page[end+len("<!--LOCAL-END-->"):]
	}
	return oas.PasswordEmailConfirmationOK{Data: strings.NewReader(page)}, nil
}

func (a *authHandler) VerifyPasswordEmail(ctx context.Context, req *oas.VerifyPasswordEmailReq) (oas.VerifyPasswordEmailRes, error) {
	w := http_utils.GetHTTPResponse(ctx)
	setEmailLoginSecurityHeaders(w)
	if err := a.validateEmailLoginOrigin(ctx); err != nil {
		return nil, err
	}
	// Validate the fixed post-login destination before consuming the one-time link.
	u, err := url.Parse(a.cfg.WEBSITE_URL)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return nil, messages.InternalServerError
	}
	u.Path = strings.TrimRight(u.Path, "/") + "/home"
	pending, err := http_utils.GetHTTPRequest(ctx).Cookie(passwordLoginPendingCookie)
	if err != nil {
		return nil, messages.InvalidEmailLoginLinkError
	}
	token, err := a.passwordLogin.Verify(ctx, req.Token, pending.Value)
	if err != nil {
		return nil, err
	}
	expired := a.CookieManager.CreateSessionCookie("")
	expired.Name = passwordLoginPendingCookie
	expired.Domain = ""
	expired.MaxAge = -1
	expired.Expires = time.Unix(1, 0).UTC()
	return &oas.VerifyPasswordEmailOKHeaders{
		SetCookie: cookie_utils.EncodeCookies([]*http.Cookie{a.CookieManager.CreateSessionCookie(token), expired}),
		Response:  oas.VerifyPasswordEmailOK{RedirectURL: u.String()},
	}, nil
}

func (a *authHandler) VerifyPasswordEmailPIN(ctx context.Context, req *oas.VerifyPasswordEmailPINReq) (oas.VerifyPasswordEmailPINRes, error) {
	w := http_utils.GetHTTPResponse(ctx)
	setEmailLoginSecurityHeaders(w)
	if err := a.validateEmailLoginOrigin(ctx); err != nil {
		return nil, err
	}
	// Validate the fixed post-login destination before consuming the one-time link.
	u, err := url.Parse(a.cfg.WEBSITE_URL)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return nil, messages.InternalServerError
	}
	u.Path = strings.TrimRight(u.Path, "/") + "/home"
	pending, err := http_utils.GetHTTPRequest(ctx).Cookie(passwordLoginPendingCookie)
	if err != nil {
		return nil, messages.InvalidEmailLoginPINError
	}
	token, err := a.passwordLogin.VerifyPIN(ctx, req.Pin, pending.Value)
	if err != nil {
		return nil, err
	}
	expired := a.CookieManager.CreateSessionCookie("")
	expired.Name = passwordLoginPendingCookie
	expired.Domain = ""
	expired.MaxAge = -1
	expired.Expires = time.Unix(1, 0).UTC()
	return &oas.VerifyPasswordEmailPINOKHeaders{
		SetCookie: cookie_utils.EncodeCookies([]*http.Cookie{a.CookieManager.CreateSessionCookie(token), expired}),
		Response:  oas.VerifyPasswordEmailPINOK{RedirectURL: u.String()},
	}, nil
}

// Login CSRF must not let another site silently start a login in the victim's
// browser. Non-browser API clients may omit Origin; links still require both secrets.
func (a *authHandler) validateEmailLoginOrigin(ctx context.Context) error {
	origin := http_utils.GetHTTPRequest(ctx).Header.Get("Origin")
	if origin == "" {
		return nil
	}
	for _, allowed := range a.cfg.BrowserOrigins() {
		if origin == allowed {
			return nil
		}
	}
	return messages.ForbiddenError
}
