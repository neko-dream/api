package ses

import (
	"context"
	"testing"

	"github.com/neko-dream/api/internal/infrastructure/config"
	email_template "github.com/neko-dream/api/internal/infrastructure/email/template"
)

func TestSESEmailSender_Send(t *testing.T) {
	t.Skip("送信を伴うテストなのでスキップ")
	ctx := context.Background()
	cfg := &config.Config{
		EMAIL_FROM:  "noreply@kotohiro.com",
		APP_NAME:    "ことひろ",
		WEBSITE_URL: "https://kotohiro.com",
	}
	emailSender := NewSESEmailSender(cfg)

	err := emailSender.Send(
		ctx,
		"info@kotohiro.com",
		email_template.VerificationEmailTemplate,
		map[string]any{
			"CompanyLogo":     "https://github.com/neko-dream/api/raw/develop/docs/public/assets/icon.png",
			"Title":           "タイトル",
			"RecipientName":   "ユーザー名",
			"VerificationURL": "https://example.com/verify?code=123456",
			"ExpiryHours":     24,
			"ContactEmail":    "info@kotohiro.com",
		},
	)
	if err != nil {
		t.Fatalf("failed to send email: %v", err)
	}
}
