package ses

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/email"
	email_template "github.com/neko-dream/api/internal/infrastructure/email/template"
	"github.com/neko-dream/api/internal/infrastructure/external/aws"
	"go.opentelemetry.io/otel"
)

type SESEmailSender struct {
	*config.Config
	sesClient SESClient
}

// SESClient allows delivery tests without AWS credentials or external mail.
type SESClient interface {
	SendEmail(context.Context, *sesv2.SendEmailInput, ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error)
}

func NewSESEmailSender(
	cfg *config.Config,
) (email.EmailSender, error) {
	if cfg == nil {
		return nil, fmt.Errorf("email configuration is missing")
	}
	awsCfg, err := aws.NewAWSConfig()
	if err != nil {
		return nil, fmt.Errorf("initialize SES: %w", err)
	}
	return NewSESEmailSenderWithClient(cfg, sesv2.NewFromConfig(awsCfg))
}

func NewSESEmailSenderWithClient(cfg *config.Config, client SESClient) (email.EmailSender, error) {
	if cfg == nil || client == nil {
		return nil, fmt.Errorf("email configuration and SES client are required")
	}
	return &SESEmailSender{Config: cfg, sesClient: client}, nil
}

// Send sends an email using AWS SES.
func (s *SESEmailSender) Send(
	ctx context.Context,
	to string,
	tmpl email_template.EmailTemplateType,
	data map[string]any,
) error {
	ctx, span := otel.Tracer("ses").Start(ctx, "SESEmailSender.Send")
	defer span.End()
	if err := ctx.Err(); err != nil {
		return err
	}

	message, err := email.RenderMessage(s.Config, to, tmpl, data)
	if err != nil {
		return err
	}

	input := sesv2.SendEmailInput{
		FromEmailAddress: &message.From,
		Destination: &types.Destination{
			ToAddresses: []string{
				message.To,
			},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data:    &message.Subject,
					Charset: stringPtr("UTF-8"),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data:    &message.HTML,
						Charset: stringPtr("UTF-8"),
					},
				},
			},
		},
	}

	_, err = s.sesClient.SendEmail(ctx, &input)
	if err != nil {
		// The caller records delivery failures without logging raw provider errors,
		// which may contain recipient addresses or one-time links.
		return fmt.Errorf("send email via SES: %w", err)
	}

	return nil
}

func stringPtr(value string) *string { return &value }
