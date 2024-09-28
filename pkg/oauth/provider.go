package oauth

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"braces.dev/errtrace"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type (
	oidcProvider struct {
		providerName AuthProviderName
		provider     *oidc.Provider
		config       oauth2.Config
	}
	OIDCProvider interface {
		GetAuthURL(ctx context.Context, state string) string
		Exchange(ctx context.Context, code string) (*oauth2.Token, error)
		Verify(ctx context.Context, token *oauth2.Token) (string, *oidc.IDToken, error)
		GenerateState() string
		Client(ctx context.Context, token *oauth2.Token) *http.Client
		UserInfo(ctx context.Context, code string) (*UserInfo, error)
	}
)

func OIDCProviderFactory(ctx context.Context, providerName AuthProviderName) (OIDCProvider, error) {
	var issuerURL, clientID, clientSecret, redirectURL string
	var scopes = []string{oidc.ScopeOpenID, "profile", "email"}

	switch providerName {
	case ProviderGoogle:
		issuerURL = os.Getenv("GOOGLE_ISSUER")
		clientID = os.Getenv("GOOGLE_CLIENT_ID")
		clientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
		redirectURL = os.Getenv("GOOGLE_CALLBACK_URL")
	default:
		return nil, errtrace.Wrap(errors.New("invalid provider"))
	}

	return errtrace.Wrap2(NewOIDCProvider(ctx, providerName, issuerURL, clientID, clientSecret, redirectURL, scopes))
}

func endpoint(providerName AuthProviderName, issuerURL string) oauth2.Endpoint {
	switch providerName {
	case ProviderGoogle:
		return google.Endpoint
	}
	return oauth2.Endpoint{}
}

func NewOIDCProvider(
	ctx context.Context,
	providerName AuthProviderName,
	issuerURL, clientID, clientSecret, redirectURL string,
	scopes []string,
) (OIDCProvider, error) {
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	end := endpoint(providerName, issuerURL)
	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     end,
		Scopes:       scopes,
	}

	return &oidcProvider{
		providerName: providerName,
		provider:     provider,
		config:       config,
	}, nil
}

func (p *oidcProvider) GetAuthURL(ctx context.Context, state string) string {
	return p.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (p *oidcProvider) Client(ctx context.Context, token *oauth2.Token) *http.Client {
	return p.config.Client(ctx, token)
}

func (p *oidcProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return errtrace.Wrap2(p.config.Exchange(ctx, code))
}

func (p *oidcProvider) Verify(ctx context.Context, token *oauth2.Token) (string, *oidc.IDToken, error) {

	rawToken, ok := token.Extra("id_token").(string)
	if !ok {
		return "", nil, errtrace.Wrap(errors.New("id_token not found"))
	}

	res, err := p.provider.Verifier(&oidc.Config{ClientID: p.config.ClientID}).Verify(ctx, rawToken)
	return rawToken, res, errtrace.Wrap(err)
}

func (p *oidcProvider) userInfo(ctx context.Context, token *oauth2.Token) (*oidc.UserInfo, error) {
	userInfo, err := p.provider.UserInfo(ctx, oauth2.StaticTokenSource(token))
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return userInfo, nil
}

func (p *oidcProvider) Refresh(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error) {
	newToken, err := p.config.TokenSource(ctx, token).Token()
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return newToken, nil
}

func (p *oidcProvider) GenerateState() string {
	return uuid.New().String()
}

type UserInfo struct {
	Subject string `json:"sub"`
	Name    string `json:"-"`
	Picture string `json:"picture"`
	Email   string `json:"email"`
}

func (p *oidcProvider) UserInfo(ctx context.Context, code string) (*UserInfo, error) {
	token, err := p.Exchange(ctx, code)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	if string(p.providerName) == string(ProviderGoogle) {
		_, idToken, err := p.Verify(ctx, token)
		if err != nil {
			return nil, errtrace.Wrap(err)
		}
		var userInfo UserInfo
		if err := idToken.Claims(&userInfo); err != nil {
			return nil, errtrace.Wrap(err)
		}

		return &UserInfo{
			Subject: userInfo.Subject,
			Picture: userInfo.Picture,
			Name:    strings.Split(userInfo.Email, "@")[0],
		}, nil
	} else {
		return nil, errtrace.Wrap(errors.New("invalid provider"))
	}

	panic("unreachable")
}
