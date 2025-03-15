package oauth

import (
	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/infrastructure/config"
)

type Provider struct {
	Issuer                auth.IssuerURI
	AuthorizationEndpoint string
	TokenEndpoint         string
	UserInfoEndpoint      string
	JwksURI               string
	Algos                 []string
	Config                *config.Config
}

var (
	providers = map[auth.IssuerURI]Provider{
		auth.GoogleIssuerURI: {
			Issuer:                auth.GoogleIssuerURI,
			AuthorizationEndpoint: "https://accounts.google.com/o/oauth2/v2/auth",
			TokenEndpoint:         "https://oauth2.googleapis.com/token",
			UserInfoEndpoint:      "https://openidconnect.googleapis.com/v1/userinfo",
			JwksURI:               "https://www.googleapis.com/oauth2/v3/certs",
			Algos:                 []string{"RS256"},
		},
		auth.LineIssuerURI: {
			Issuer:                auth.LineIssuerURI,
			AuthorizationEndpoint: "https://access.line.me/oauth2/v2.1/authorize",
			TokenEndpoint:         "https://api.line.me/oauth2/v2.1/token",
			UserInfoEndpoint:      "https://api.line.me/v2/profile",
			JwksURI:               "https://api.line.me/oauth2/v2.1/certs",
			Algos:                 []string{"ES256", "HS256"},
		},
	}
)

func GetProvider(issuerURI auth.IssuerURI, conf *config.Config) (Provider, bool) {
	provider, ok := providers[issuerURI]
	provider.Config = conf
	return provider, ok
}

func (p Provider) ClientID() string {
	switch p.Issuer {
	case auth.LineIssuerURI:
		return p.Config.LineClientID
	case auth.GoogleIssuerURI:
		return p.Config.GoogleClientID
	default:
		return ""
	}
}

func (p Provider) ClientSecret() string {
	switch p.Issuer {
	case auth.LineIssuerURI:
		return p.Config.LineClientSecret
	case auth.GoogleIssuerURI:
		return p.Config.GoogleClientSecret
	default:
		return ""
	}
}

func (p Provider) Scopes() []string {
	switch p.Issuer {
	case auth.LineIssuerURI:
		return []string{"openid", "email"}
	case auth.GoogleIssuerURI:
		return []string{"openid", "email"}
	default:
		return []string{}
	}
}

func (p Provider) RedirectURL() string {
	switch p.Issuer {
	case auth.LineIssuerURI:
		return p.Config.LineCallbackURL
	case auth.GoogleIssuerURI:
		return p.Config.GoogleCallbackURL
	default:
		return ""
	}
}
