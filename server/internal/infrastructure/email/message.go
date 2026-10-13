package email

import (
	"bytes"
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"strings"

	"github.com/neko-dream/api/internal/infrastructure/config"
	email_template "github.com/neko-dream/api/internal/infrastructure/email/template"
)

// Message is the validated, rendered payload shared by delivery transports.
type Message struct {
	From, To, Subject, HTML string
}

var loginPINPattern = regexp.MustCompile(`^[0-9]{6}$`)

func RenderMessage(cfg *config.Config, to string, tmpl email_template.EmailTemplateType, data map[string]any) (Message, error) {
	if cfg == nil {
		return Message{}, fmt.Errorf("email configuration is missing")
	}
	from, err := parseAddress(cfg.EMAIL_FROM)
	if err != nil {
		return Message{}, fmt.Errorf("invalid sender: %w", err)
	}
	recipient, err := parseAddress(to)
	if err != nil {
		return Message{}, fmt.Errorf("invalid recipient: %w", err)
	}
	subject, ok := data["Title"].(string)
	if !ok || strings.TrimSpace(subject) == "" || strings.ContainsAny(subject, "\r\n") {
		return Message{}, fmt.Errorf("email title must be a non-empty single-line string")
	}
	fields := email_template.DataWithCommonFields(cfg, data)
	if err := validateTemplate(tmpl, fields); err != nil {
		return Message{}, err
	}
	// Optional template fields have explicit defaults so missing required fields
	// can be rejected without breaking conditional sections.
	for _, key := range []string{"CompanyLogo", "RecipientName", "Password"} {
		if _, ok := fields[key]; !ok {
			fields[key] = ""
		}
	}
	if _, ok := fields["ExistingAccount"]; !ok {
		fields["ExistingAccount"] = false
	}
	if _, ok := fields["ContactEmail"]; !ok {
		fields["ContactEmail"] = from.Address
	}
	t, err := email_template.LoadMailTemplate(tmpl)
	if err != nil {
		return Message{}, fmt.Errorf("load email template: %w", err)
	}
	var body bytes.Buffer
	if err := t.Option("missingkey=error").ExecuteTemplate(&body, string(tmpl), fields); err != nil {
		return Message{}, fmt.Errorf("render email template: %w", err)
	}
	return Message{From: from.String(), To: recipient.String(), Subject: subject, HTML: body.String()}, nil
}

func parseAddress(value string) (*mail.Address, error) {
	if strings.ContainsAny(value, "\r\n") {
		return nil, fmt.Errorf("email address must not contain line breaks")
	}
	return mail.ParseAddress(strings.TrimSpace(value))
}

func validateTemplate(tmpl email_template.EmailTemplateType, fields map[string]any) error {
	var required []string
	var linkKey string
	switch tmpl {
	case email_template.PasswordSetupEmailTemplate:
		linkKey = "SetupURL"
		required = []string{"RecipientName"}
		hours, ok := fields["ExpiryHours"].(int)
		if !ok || hours <= 0 {
			return fmt.Errorf("ExpiryHours must be a positive integer")
		}
	case email_template.PasswordLoginEmailTemplate:
		linkKey = "VerificationURL"
		pin, ok := fields["PIN"].(string)
		if !ok || !loginPINPattern.MatchString(pin) {
			return fmt.Errorf("PIN must contain exactly six ASCII digits")
		}
		minutes, ok := fields["ExpiryMinutes"].(int)
		if !ok || minutes <= 0 {
			return fmt.Errorf("ExpiryMinutes must be a positive integer")
		}
	case email_template.VerificationEmailTemplate:
		linkKey = "VerificationURL"
		hours, ok := fields["ExpiryHours"].(int)
		if !ok || hours <= 0 {
			return fmt.Errorf("ExpiryHours must be a positive integer")
		}
	case email_template.OrganizationInvitationEmailTemplate:
		required = []string{"OrganizationName", "Email"}
		linkKey = "InvitationURL"
		existing, ok := fields["ExistingAccount"].(bool)
		if fields["ExistingAccount"] != nil && !ok {
			return fmt.Errorf("ExistingAccount must be a boolean")
		}
		if !existing {
			required = append(required, "Password")
		}
	default:
		return fmt.Errorf("unsupported email template %q", tmpl)
	}
	for _, key := range required {
		value, ok := fields[key].(string)
		if !ok || strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s must be a non-empty string", key)
		}
	}
	link, ok := fields[linkKey].(string)
	parsed, err := url.Parse(link)
	if !ok || err != nil || parsed.Host == "" || (parsed.Scheme != "https" && parsed.Scheme != "http") {
		return fmt.Errorf("%s must be an absolute HTTP(S) URL", linkKey)
	}
	if tmpl == email_template.PasswordSetupEmailTemplate && (parsed.Scheme != "https" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment == "") {
		return fmt.Errorf("SetupURL must be an HTTPS URL with a fragment token")
	}
	return nil
}
