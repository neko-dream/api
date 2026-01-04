package shared

import (
	"errors"
	"strings"

	"braces.dev/errtrace"
)

type (
	IssuerURI        string
	AuthProviderName string
)

const (
	GoogleIssuerURI IssuerURI = "https://accounts.google.com"
	LineIssuerURI   IssuerURI = "https://access.line.me"
)

func (i IssuerURI) String() string {
	return string(i)
}

const (
	ProviderGoogle   AuthProviderName = "GOOGLE"
	ProviderLine     AuthProviderName = "LINE"
	ProviderDEV      AuthProviderName = "DEV"
	ProviderPassword AuthProviderName = "PASSWORD"
	ProviderGuest    AuthProviderName = "GUEST"
)

func NewAuthProviderName(provider string) (AuthProviderName, error) {
	switch strings.ToUpper(provider) {
	case ProviderGoogle.String():
		return ProviderGoogle, nil
	case ProviderLine.String():
		return ProviderLine, nil
	case ProviderDEV.String():
		return ProviderDEV, nil
	case ProviderPassword.String():
		return ProviderPassword, nil
	case ProviderGuest.String():
		return ProviderGuest, nil
	default:
		return "", errtrace.Wrap(errors.New("invalid auth provider"))
	}
}

// IsGuest ゲストプロバイダーかどうかを判定
func (a AuthProviderName) IsGuest() bool {
	return a == ProviderGuest
}

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
