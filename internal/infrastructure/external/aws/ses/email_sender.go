package ses

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/email"
	email_template "github.com/neko-dream/server/internal/infrastructure/email/template"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type SESEmailSender struct {
	*config.Config
	sesClient *sesv2.Client
}

func NewSESEmailSender(
	cfg *config.Config,
	sesClient *sesv2.Client,
) email.EmailSender {
	return &SESEmailSender{
		Config:    cfg,
		sesClient: sesClient,
	}
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

	if to == "" {
		return fmt.Errorf("メール送信先アドレスが指定されていません")
	}

	dataWithCommon := email_template.DataWithCommonFields(
		s.Config,
		data,
	)
	t, err := email_template.LoadMailTemplate(tmpl)
	if err != nil {
		return err
	}

	resultBuf := new(bytes.Buffer)
	if err = t.ExecuteTemplate(resultBuf, string(tmpl), dataWithCommon); err != nil {
		return err
	}

	input := sesv2.SendEmailInput{
		FromEmailAddress: &s.Config.EMAIL_FROM,
		Destination: &types.Destination{
			ToAddresses: []string{
				to,
			},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: lo.ToPtr(data["Title"].(string)),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data: lo.ToPtr(resultBuf.String()),
					},
				},
			},
		},
	}

	_, err = s.sesClient.SendEmail(ctx, &input)
	if err != nil {
		return err
	}

	return nil
}
