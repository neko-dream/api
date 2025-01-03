package oauth

import (
	"errors"
	"strings"

	"braces.dev/errtrace"

	"github.com/neko-dream/server/internal/infrastructure/config"
)

type IssuerURI string

const (
	GoogleIssuerURI IssuerURI = "https://accounts.google.com"
	LineIssuerURI   IssuerURI = "https://access.line.me"
)

func (i IssuerURI) String() string {
	return string(i)
}

type AuthProviderName string

const (
	ProviderGoogle AuthProviderName = "GOOGLE"
	ProviderLine   AuthProviderName = "LINE"
)

func (a AuthProviderName) IssuerURI() IssuerURI {
	switch a {
	case ProviderGoogle:
		return GoogleIssuerURI
	case ProviderLine:
		return LineIssuerURI
	default:
		return ""
	}
}

func (a AuthProviderName) String() string {
	return string(a)
}

func NewAuthProviderName(provider string) (AuthProviderName, error) {
	switch strings.ToUpper(provider) {
	case ProviderGoogle.String():
		return ProviderGoogle, nil
	case ProviderLine.String():
		return ProviderLine, nil
	default:
		return "", errtrace.Wrap(errors.New("invalid auth provider"))
	}
}

type Provider struct {
	Issuer                IssuerURI
	AuthorizationEndpoint string
	TokenEndpoint         string
	UserInfoEndpoint      string
	JwksURI               string
	Algos                 []string
}

var (
	providers = map[IssuerURI]Provider{
		GoogleIssuerURI: {
			Issuer:                GoogleIssuerURI,
			AuthorizationEndpoint: "https://accounts.google.com/o/oauth2/v2/auth",
			TokenEndpoint:         "https://oauth2.googleapis.com/token",
			UserInfoEndpoint:      "https://openidconnect.googleapis.com/v1/userinfo",
			JwksURI:               "https://www.googleapis.com/oauth2/v3/certs",
			Algos:                 []string{"RS256"},
		},
		LineIssuerURI: {
			Issuer:                LineIssuerURI,
			AuthorizationEndpoint: "https://access.line.me/oauth2/v2.1/authorize",
			TokenEndpoint:         "https://api.line.me/oauth2/v2.1/token",
			UserInfoEndpoint:      "https://api.line.me/v2/profile",
			JwksURI:               "https://api.line.me/oauth2/v2.1/certs",
			Algos:                 []string{"ES256", "HS256"},
		},
	}
)

func GetProvider(issuerURI IssuerURI) (Provider, bool) {
	provider, ok := providers[issuerURI]
	return provider, ok
}

func (p Provider) ClientID(config *config.Config) string {
	switch p.Issuer {
	case LineIssuerURI:
		return config.LineClientID
	case GoogleIssuerURI:
		return config.GoogleClientID
	default:
		return ""
	}
}

func (p Provider) ClientSecret(config *config.Config) string {
	switch p.Issuer {
	case LineIssuerURI:
		return config.LineClientSecret
	case GoogleIssuerURI:
		return config.GoogleClientSecret
	default:
		return ""
	}
}

func (p Provider) Scopes() []string {
	switch p.Issuer {
	case LineIssuerURI:
		return []string{"openid", "profile"}
	case GoogleIssuerURI:
		return []string{"openid", "profile"}
	default:
		return []string{}
	}
}

func (p Provider) RedirectURL(config *config.Config) string {
	switch p.Issuer {
	case LineIssuerURI:
		return config.LineCallbackURL
	case GoogleIssuerURI:
		return config.GoogleCallbackURL
	default:
		return ""
	}
}
