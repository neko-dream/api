package auth

import (
	"context"
)

type (
	AuthProvider interface {
		GetAuthorizationURL(ctx context.Context, state string) string
		// return subject, pictureURL, error
		VerifyAndIdentify(ctx context.Context, code string) (*string, *string, error)
	}

	// プロバイダ名よりAuthProviderを生成するファクトリ
	AuthProviderFactory interface {
		NewAuthProvider(ctx context.Context, providerName string) (AuthProvider, error)
	}

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
