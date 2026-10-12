package email

import (
	"context"

	email_template "github.com/neko-dream/api/internal/infrastructure/email/template"
)

type EmailSender interface {
	Send(context.Context, string, email_template.EmailTemplateType, map[string]any) error
}
