package email

import (
	"maps"
	"strings"
	"testing"

	"github.com/neko-dream/api/internal/infrastructure/config"
	email_template "github.com/neko-dream/api/internal/infrastructure/email/template"
)

func mailConfig() *config.Config {
	return &config.Config{EMAIL_FROM: "ことひろ <noreply@example.com>", APP_NAME: "ことひろ", WEBSITE_URL: "https://example.com"}
}

func TestRenderPasswordLoginPIN(t *testing.T) {
	for _, pin := range []any{"000001", "", "12345", "1234567", "１２３４５６", 123456, nil} {
		data := map[string]any{"Title": "ログインの本人確認", "VerificationURL": "https://example.com/verify#token", "ExpiryMinutes": 15, "PIN": pin}
		m, err := RenderMessage(mailConfig(), "user@example.com", email_template.PasswordLoginEmailTemplate, data)
		if pin == "000001" {
			if err != nil || !strings.Contains(m.HTML, "000001") || !strings.Contains(m.HTML, "15") {
				t.Fatal("missing login PIN", err)
			}
		} else if err == nil {
			t.Fatalf("invalid PIN accepted: %v", pin)
		}
	}
}

func verificationData() map[string]any {
	return map[string]any{"Title": "メールの確認", "VerificationURL": "https://example.com/verify?token=one&next=two", "ExpiryHours": 1}
}

func TestRenderMessageVerification(t *testing.T) {
	data := verificationData()
	original := maps.Clone(data)
	m, err := RenderMessage(mailConfig(), "user@example.com", email_template.VerificationEmailTemplate, data)
	if err != nil {
		t.Fatal(err)
	}
	if m.Subject != "メールの確認" || !strings.Contains(m.HTML, "メールアドレスを確認") || !strings.Contains(m.HTML, "one&amp;next=two") || !strings.Contains(m.HTML, "mailto:noreply@example.com") {
		t.Fatalf("rendered message is missing expected content: %#v", m)
	}
	if !maps.Equal(data, original) {
		t.Fatal("render mutated the caller's data")
	}
}

func TestRenderMessageRejectsInvalidData(t *testing.T) {
	tests := []struct {
		name   string
		change func(*config.Config, *string, map[string]any)
	}{
		{"missing sender", func(c *config.Config, _ *string, _ map[string]any) { c.EMAIL_FROM = "" }},
		{"invalid sender", func(c *config.Config, _ *string, _ map[string]any) { c.EMAIL_FROM = "not-an-address" }},
		{"sender injection", func(c *config.Config, _ *string, _ map[string]any) {
			c.EMAIL_FROM = "a@example.com\r\nBcc: b@example.com"
		}},
		{"missing recipient", func(_ *config.Config, to *string, _ map[string]any) { *to = "" }},
		{"invalid recipient", func(_ *config.Config, to *string, _ map[string]any) { *to = "invalid" }},
		{"multiple recipients", func(_ *config.Config, to *string, _ map[string]any) { *to = "a@example.com,b@example.com" }},
		{"recipient injection", func(_ *config.Config, to *string, _ map[string]any) { *to = "a@example.com\nBcc: b@example.com" }},
		{"missing title", func(_ *config.Config, _ *string, d map[string]any) { delete(d, "Title") }},
		{"title wrong type", func(_ *config.Config, _ *string, d map[string]any) { d["Title"] = 10 }},
		{"blank title", func(_ *config.Config, _ *string, d map[string]any) { d["Title"] = "  " }},
		{"title injection", func(_ *config.Config, _ *string, d map[string]any) { d["Title"] = "title\r\nBcc: x@example.com" }},
		{"missing URL", func(_ *config.Config, _ *string, d map[string]any) { delete(d, "VerificationURL") }},
		{"URL wrong type", func(_ *config.Config, _ *string, d map[string]any) { d["VerificationURL"] = 10 }},
		{"relative URL", func(_ *config.Config, _ *string, d map[string]any) { d["VerificationURL"] = "/verify" }},
		{"unsafe URL", func(_ *config.Config, _ *string, d map[string]any) { d["VerificationURL"] = "javascript:alert(1)" }},
		{"invalid URL", func(_ *config.Config, _ *string, d map[string]any) { d["VerificationURL"] = "https://%zz" }},
		{"missing expiry", func(_ *config.Config, _ *string, d map[string]any) { delete(d, "ExpiryHours") }},
		{"zero expiry", func(_ *config.Config, _ *string, d map[string]any) { d["ExpiryHours"] = 0 }},
		{"negative expiry", func(_ *config.Config, _ *string, d map[string]any) { d["ExpiryHours"] = -1 }},
		{"expiry wrong type", func(_ *config.Config, _ *string, d map[string]any) { d["ExpiryHours"] = "1" }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, to, data := mailConfig(), "user@example.com", verificationData()
			tt.change(cfg, &to, data)
			if _, err := RenderMessage(cfg, to, email_template.VerificationEmailTemplate, data); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
	if _, err := RenderMessage(nil, "user@example.com", email_template.VerificationEmailTemplate, verificationData()); err == nil {
		t.Fatal("nil config accepted")
	}
	if _, err := RenderMessage(mailConfig(), "user@example.com", "missing.tpl", verificationData()); err == nil {
		t.Fatal("unsupported template accepted")
	}
}

func TestRenderInvitation(t *testing.T) {
	for _, existing := range []bool{false, true} {
		t.Run(map[bool]string{false: "new account", true: "existing account"}[existing], func(t *testing.T) {
			data := map[string]any{"Title": "招待", "OrganizationName": "<script>市役所</script>", "Email": "user@example.com", "InvitationURL": "https://example.com", "ExistingAccount": existing}
			if !existing {
				data["Password"] = "initial-password"
			}
			m, err := RenderMessage(mailConfig(), "user@example.com", email_template.OrganizationInvitationEmailTemplate, data)
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(m.HTML, "<script>") || !strings.Contains(m.HTML, "&lt;script&gt;") {
				t.Fatal("organization name was not escaped")
			}
			if existing && (strings.Contains(m.HTML, "初期パスワード") || !strings.Contains(m.HTML, "既存のアカウント")) {
				t.Fatal("existing account received new credentials")
			}
			if !existing && (!strings.Contains(m.HTML, "initial-password") || !strings.Contains(m.HTML, "初回ログイン")) {
				t.Fatal("new account instructions missing")
			}
			for _, key := range []string{"OrganizationName", "Email", "InvitationURL"} {
				invalid := maps.Clone(data)
				delete(invalid, key)
				if _, err := RenderMessage(mailConfig(), "user@example.com", email_template.OrganizationInvitationEmailTemplate, invalid); err == nil {
					t.Fatalf("missing %s accepted", key)
				}
			}
			if !existing {
				delete(data, "Password")
				if _, err := RenderMessage(mailConfig(), "user@example.com", email_template.OrganizationInvitationEmailTemplate, data); err == nil {
					t.Fatal("missing initial password accepted")
				}
			}
		})
	}
}
