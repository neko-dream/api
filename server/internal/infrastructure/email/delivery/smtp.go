package delivery

import (
	"context"
	"fmt"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"time"

	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/email"
	email_template "github.com/neko-dream/api/internal/infrastructure/email/template"
)

type smtpSender struct{ cfg *config.Config }

func (s *smtpSender) Send(ctx context.Context, to string, tmpl email_template.EmailTemplateType, data map[string]any) error {
	message, err := email.RenderMessage(s.cfg, to, tmpl, data)
	if err != nil {
		return err
	}
	return deliverSMTP(ctx, s.cfg.SMTP_ADDR, message)
}

func deliverSMTP(ctx context.Context, address string, message email.Message) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	connection, err := (&net.Dialer{}).DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("connect to local SMTP: %w", err)
	}
	defer connection.Close()
	stop := context.AfterFunc(ctx, func() { _ = connection.Close() })
	defer stop()
	deadline, _ := ctx.Deadline()
	if err := connection.SetDeadline(deadline); err != nil {
		return err
	}
	host, _, _ := net.SplitHostPort(address)
	client, err := smtp.NewClient(connection, host)
	if err != nil {
		return fmt.Errorf("initialize local SMTP: %w", err)
	}
	defer client.Close()
	from, err := mail.ParseAddress(message.From)
	if err != nil {
		return err
	}
	to, err := mail.ParseAddress(message.To)
	if err != nil {
		return err
	}
	if err := client.Mail(from.Address); err != nil {
		return fmt.Errorf("SMTP MAIL FROM: %w", err)
	}
	if err := client.Rcpt(to.Address); err != nil {
		return fmt.Errorf("SMTP RCPT TO: %w", err)
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP DATA: %w", err)
	}
	_, err = fmt.Fprintf(w, "From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\nContent-Transfer-Encoding: 8bit\r\n\r\n%s", message.From, message.To, mime.QEncoding.Encode("UTF-8", message.Subject), message.HTML)
	if err != nil {
		_ = w.Close()
		return fmt.Errorf("write SMTP message: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("accept SMTP message: %w", err)
	}
	return client.Quit()
}
