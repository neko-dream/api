package ses

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"

	"github.com/neko-dream/api/internal/infrastructure/config"
	email_template "github.com/neko-dream/api/internal/infrastructure/email/template"
)

func TestSESEmailSender_ProviderFailureDoesNotLogPrivateContent(t *testing.T) {
	var logs bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&logs, nil)))
	t.Cleanup(func() { slog.SetDefault(previous) })
	providerErr := errors.New("delivery failed for user@example.com: https://example.com/auth/password/setup#private-secret")
	client := &fakeSESClient{err: providerErr}
	sender, err := NewSESEmailSenderWithClient(&config.Config{EMAIL_FROM: "noreply@example.com", APP_NAME: "ことひろ", WEBSITE_URL: "https://example.com"}, client)
	if err != nil {
		t.Fatal(err)
	}
	err = sender.Send(context.Background(), "user@example.com", email_template.PasswordSetupEmailTemplate, map[string]any{"Title": "パスワード設定", "RecipientName": "ユーザー", "SetupURL": "https://example.com/auth/password/setup#private-secret", "ExpiryHours": 24})
	if !errors.Is(err, providerErr) {
		t.Fatalf("provider error was lost: %v", err)
	}
	if strings.Contains(logs.String(), "user@example.com") || strings.Contains(logs.String(), "private-secret") {
		t.Fatal("private provider error was logged")
	}
}

type fakeSESClient struct {
	input *sesv2.SendEmailInput
	ctx   context.Context
	err   error
	calls int
}

type contextKey struct{}

func (f *fakeSESClient) SendEmail(ctx context.Context, input *sesv2.SendEmailInput, _ ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error) {
	f.ctx, f.input = ctx, input
	f.calls++
	return &sesv2.SendEmailOutput{}, f.err
}

func TestSESEmailSender_Send(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(map[bool]string{false: "delivery", true: "provider failure"}[fail], func(t *testing.T) {
			providerErr := errors.New("SES rejected delivery")
			client := &fakeSESClient{}
			if fail {
				client.err = providerErr
			}
			cfg := &config.Config{EMAIL_FROM: "noreply@example.com", APP_NAME: "ことひろ", WEBSITE_URL: "https://example.com"}
			sender, err := NewSESEmailSenderWithClient(cfg, client)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.WithValue(context.Background(), contextKey{}, "request-context")
			err = sender.Send(ctx, "user@example.com", email_template.VerificationEmailTemplate, map[string]any{"Title": "メールの確認", "VerificationURL": "https://example.com/verify", "ExpiryHours": 24})
			if fail && !errors.Is(err, providerErr) {
				t.Fatalf("provider error was lost: %v", err)
			}
			if !fail && err != nil {
				t.Fatal(err)
			}
			if client.calls != 1 || client.ctx.Value(contextKey{}) != "request-context" {
				t.Fatal("unexpected delivery or context")
			}
			input := client.input
			if len(input.Destination.ToAddresses) != 1 || !strings.Contains(input.Destination.ToAddresses[0], "user@example.com") || !strings.Contains(*input.FromEmailAddress, "noreply@example.com") {
				t.Fatal("incorrect envelope")
			}
			if *input.Content.Simple.Subject.Data != "メールの確認" || *input.Content.Simple.Subject.Charset != "UTF-8" || *input.Content.Simple.Body.Html.Charset != "UTF-8" || !strings.Contains(*input.Content.Simple.Body.Html.Data, "https://example.com/verify") {
				t.Fatal("incorrect content")
			}
		})
	}
}

func TestSESEmailSender_ValidationDoesNotSend(t *testing.T) {
	client := &fakeSESClient{}
	sender, err := NewSESEmailSenderWithClient(&config.Config{EMAIL_FROM: "noreply@example.com"}, client)
	if err != nil {
		t.Fatal(err)
	}
	if err := sender.Send(context.Background(), "user@example.com", email_template.VerificationEmailTemplate, map[string]any{"Title": 42}); err == nil {
		t.Fatal("invalid input accepted")
	}
	if client.calls != 0 {
		t.Fatal("invalid email was sent")
	}
	if _, err := NewSESEmailSenderWithClient(nil, client); err == nil {
		t.Fatal("nil config accepted")
	}
	if _, err := NewSESEmailSenderWithClient(&config.Config{}, nil); err == nil {
		t.Fatal("nil client accepted")
	}
	if _, err := NewSESEmailSender(nil); err == nil {
		t.Fatal("nil config accepted by default constructor")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := sender.Send(ctx, "user@example.com", email_template.VerificationEmailTemplate, nil); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation not preserved: %v", err)
	}
	if client.calls != 0 {
		t.Fatal("canceled send called SES")
	}
}
