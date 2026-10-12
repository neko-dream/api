package delivery

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/neko-dream/api/internal/infrastructure/config"
	email_template "github.com/neko-dream/api/internal/infrastructure/email/template"
)

func TestMailpitIntegration(t *testing.T) {
	address := os.Getenv("MAIL_TEST_SMTP_ADDR")
	if address == "" {
		t.Skip("set MAIL_TEST_SMTP_ADDR and MAIL_TEST_HTTP_URL to test Mailpit delivery")
	}
	baseURL := os.Getenv("MAIL_TEST_HTTP_URL")
	if baseURL == "" {
		t.Fatal("MAIL_TEST_HTTP_URL is required")
	}
	cfg := &config.Config{Env: config.LOCAL, SMTP_ADDR: address, EMAIL_FROM: "ことひろ <noreply@localhost.test>", APP_NAME: "ことひろ", WEBSITE_URL: "https://localhost:3000"}
	sender, err := NewSender(cfg)
	if err != nil {
		t.Fatal(err)
	}
	subject := "ことひろ・メール送信テスト " + uuid.NewString()
	if err := sender.Send(context.Background(), "mail-test@localhost.test", email_template.VerificationEmailTemplate, map[string]any{"Title": subject, "VerificationURL": "https://localhost:3000/auth/login", "ExpiryHours": 1}); err != nil {
		t.Fatal(err)
	}
	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Get(baseURL + "/api/v1/messages")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("Mailpit HTTP status: %d", response.StatusCode)
	}
	var mailbox struct {
		Messages []struct {
			Subject string
			To      []struct{ Address string }
		}
	}
	if err := json.NewDecoder(response.Body).Decode(&mailbox); err != nil {
		t.Fatal(err)
	}
	for _, message := range mailbox.Messages {
		if message.Subject == subject && len(message.To) == 1 && message.To[0].Address == "mail-test@localhost.test" {
			return
		}
	}
	t.Fatal("SMTP delivery was not found in the local mailbox")
}
