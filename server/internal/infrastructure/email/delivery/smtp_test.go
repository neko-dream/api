package delivery

import (
	"bufio"
	"context"
	"fmt"
	"mime"
	"net"
	"net/mail"
	"net/textproto"
	"strings"
	"testing"
	"time"

	"github.com/neko-dream/api/internal/infrastructure/config"
	email_template "github.com/neko-dream/api/internal/infrastructure/email/template"
	"github.com/neko-dream/api/internal/infrastructure/external/aws/ses"
)

func startSMTP(t *testing.T, rejectRecipient bool) (string, <-chan string) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = listener.Close() })
	result := make(chan string, 1)
	go func() {
		connection, err := listener.Accept()
		if err != nil {
			return
		}
		defer connection.Close()
		_ = connection.SetDeadline(time.Now().Add(5 * time.Second))
		reader := textproto.NewReader(bufio.NewReader(connection))
		_, _ = fmt.Fprint(connection, "220 localhost SMTP\r\n")
		for {
			line, err := reader.ReadLine()
			if err != nil {
				return
			}
			switch {
			case strings.HasPrefix(line, "EHLO"), strings.HasPrefix(line, "HELO"), strings.HasPrefix(line, "MAIL"):
				_, _ = fmt.Fprint(connection, "250 OK\r\n")
			case strings.HasPrefix(line, "RCPT"):
				if rejectRecipient {
					_, _ = fmt.Fprint(connection, "550 recipient rejected\r\n")
				} else {
					_, _ = fmt.Fprint(connection, "250 OK\r\n")
				}
			case line == "DATA":
				_, _ = fmt.Fprint(connection, "354 send body\r\n")
				body, err := reader.ReadDotBytes()
				if err != nil {
					return
				}
				result <- string(body)
				_, _ = fmt.Fprint(connection, "250 accepted\r\n")
			case line == "QUIT":
				_, _ = fmt.Fprint(connection, "221 bye\r\n")
				return
			}
		}
	}()
	return listener.Addr().String(), result
}

func TestLocalSMTPDelivery(t *testing.T) {
	address, received := startSMTP(t, false)
	sender, err := NewSender(&config.Config{Env: config.LOCAL, SMTP_ADDR: address, EMAIL_FROM: "ことひろ <noreply@example.com>", APP_NAME: "ことひろ", WEBSITE_URL: "https://example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if err := sender.Send(context.Background(), "user@example.com", email_template.VerificationEmailTemplate, map[string]any{"Title": "確認メール", "VerificationURL": "https://example.com/verify", "ExpiryHours": 1}); err != nil {
		t.Fatal(err)
	}
	select {
	case raw := <-received:
		message, err := mail.ReadMessage(strings.NewReader(raw))
		if err != nil {
			t.Fatal(err)
		}
		subject, err := (&mime.WordDecoder{}).DecodeHeader(message.Header.Get("Subject"))
		if err != nil || subject != "確認メール" {
			t.Fatalf("invalid UTF-8 subject: %q %v", subject, err)
		}
		if !strings.Contains(raw, "https://example.com/verify") || !strings.Contains(message.Header.Get("Content-Type"), "text/html") {
			t.Fatal("invalid SMTP message")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("SMTP message not received")
	}
}

func TestLocalSMTPFailure(t *testing.T) {
	address, received := startSMTP(t, true)
	sender, err := NewSender(&config.Config{Env: config.LOCAL, SMTP_ADDR: address, EMAIL_FROM: "noreply@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	data := map[string]any{"Title": "確認", "VerificationURL": "https://example.com/verify", "ExpiryHours": 1}
	if err := sender.Send(context.Background(), "user@example.com", email_template.VerificationEmailTemplate, data); err == nil {
		t.Fatal("rejected recipient reported success")
	}
	select {
	case <-received:
		t.Fatal("rejected email delivered")
	default:
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := sender.Send(ctx, "user@example.com", email_template.VerificationEmailTemplate, data); err == nil {
		t.Fatal("canceled send reported success")
	}
}

func TestLocalSMTPTimeout(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	go func() {
		connection, err := listener.Accept()
		if err != nil {
			return
		}
		defer connection.Close()
		_ = connection.SetDeadline(time.Now().Add(time.Second))
		_, _ = connection.Read(make([]byte, 1))
	}()
	sender, err := NewSender(&config.Config{Env: config.LOCAL, SMTP_ADDR: listener.Addr().String(), EMAIL_FROM: "noreply@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	started := time.Now()
	err = sender.Send(ctx, "user@example.com", email_template.VerificationEmailTemplate, map[string]any{"Title": "確認", "VerificationURL": "https://example.com/verify", "ExpiryHours": 1})
	if err == nil || time.Since(started) > time.Second {
		t.Fatalf("context deadline ignored: %v", err)
	}
}

func TestLocalSenderRejectsExternalAddresses(t *testing.T) {
	for _, address := range []string{"", "smtp.example.com:25", "0.0.0.0:1025", "192.168.1.1:1025", "127.0.0.1", "localhost:", "localhost:abc", "localhost:0", "localhost:65536"} {
		t.Run(address, func(t *testing.T) {
			if _, err := NewSender(&config.Config{Env: config.LOCAL, SMTP_ADDR: address}); err == nil {
				t.Fatal("invalid local address accepted")
			}
		})
	}
	if _, err := NewSender(nil); err == nil {
		t.Fatal("nil config accepted")
	}
}

func TestDeployedSenderUsesSES(t *testing.T) {
	t.Setenv("ENV", "production")
	t.Setenv("AWS_EC2_METADATA_DISABLED", "true")
	for _, env := range []config.ENV{config.DEV, config.PROD} {
		sender, err := NewSender(&config.Config{Env: env, SMTP_ADDR: "invalid-local-address"})
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := sender.(*ses.SESEmailSender); !ok {
			t.Fatal("deployed environments must use SES")
		}
	}
}
