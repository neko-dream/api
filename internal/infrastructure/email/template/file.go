package email_template

import (
	"embed"
	"html/template"

	"github.com/neko-dream/server/internal/infrastructure/config"
)

type EmailTemplateType string

//go:embed *.tpl
var templateFiles embed.FS

const (
	headerTemplate EmailTemplateType = "header.tpl"
	footerTemplate EmailTemplateType = "footer.tpl"
	// VerificationEmailTemplate
	VerificationEmailTemplate EmailTemplateType = "verification_email.tpl"
	// OrganizationInvitationEmailTemplate
	OrganizationInvitationEmailTemplate EmailTemplateType = "organization_invitation.tpl"
)

func LoadMailTemplate(templateType EmailTemplateType) (*template.Template, error) {
	t, err := template.ParseFS(
		templateFiles,
		string(templateType),
		string(headerTemplate),
		string(footerTemplate),
	)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func LoadMailTemplateWithFunc(templateType EmailTemplateType, funcMap template.FuncMap) (*template.Template, error) {
	t, err := template.New(string(templateType)).Funcs(funcMap).ParseFS(
		templateFiles,
		string(templateType),
		string(headerTemplate),
		string(footerTemplate),
	)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func DataWithCommonFields(cfg *config.Config, data map[string]any) map[string]any {
	commonFields := map[string]any{
		"AppName":    cfg.APP_NAME,
		"WebsiteURL": cfg.WEBSITE_URL,
	}

	for k, v := range data {
		commonFields[k] = v
	}

	return commonFields
}
