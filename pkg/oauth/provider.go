package oauth

import (
	"context"
	"errors"
	"fmt"

	"braces.dev/errtrace"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/pkg/utils"
	"golang.org/x/exp/slices"
	"golang.org/x/oauth2"
)

type (
	oidcProvider struct {
		providerName AuthProviderName
		oauthConf    oauth2.Config
		conf         *config.Config
		provider     *oidc.Provider
	}
	OIDCProvider interface {
		GetAuthURL(ctx context.Context, state string) string
		UserInfo(ctx context.Context, code string) (*UserInfo, error)
	}
)

func NewOIDCProvider(
	ctx context.Context,
	providerName AuthProviderName,
	conf *config.Config,
) (OIDCProvider, error) {
	issuer := providerName.IssuerURI()
	provider, ok := GetProvider(issuer)
	if !ok {
		return nil, errtrace.Wrap(errors.New("provider not found"))
	}

	config := oauth2.Config{
		ClientID:     provider.ClientID(conf),
		ClientSecret: provider.ClientSecret(conf),
		RedirectURL:  provider.RedirectURL(conf),
		Endpoint: oauth2.Endpoint{
			AuthURL:  provider.AuthorizationEndpoint,
			TokenURL: provider.TokenEndpoint,
		},
		Scopes: provider.Scopes(),
	}
	oidcProv, err := oidc.NewProvider(ctx, issuer.String())
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	return &oidcProvider{
		providerName: providerName,
		oauthConf:    config,
		conf:         conf,
		provider:     oidcProv,
	}, nil
}

func (p *oidcProvider) GetAuthURL(ctx context.Context, state string) string {
	return p.oauthConf.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (p *oidcProvider) GetProvider() *Provider {
	issuer := p.providerName.IssuerURI()
	provider, _ := GetProvider(issuer)
	return &provider
}

type UserInfo struct {
	Subject string `json:"sub"`
	Picture string `json:"picture"`
}

func (p *oidcProvider) UserInfo(ctx context.Context, code string) (*UserInfo, error) {
	token, err := p.oauthConf.Exchange(ctx, code)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	rawToken, ok := token.Extra("id_token").(string)
	if !ok {
		utils.HandleError(ctx, err, "id_token not found")
		return nil, errtrace.Wrap(errors.New("id_token not found"))
	}
	provider := p.GetProvider()
	// Google の公開鍵を取得
	keySet, err := jwk.Fetch(context.Background(), provider.JwksURI)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWK: %v", err)
	}

	jwtToken, err := jwt.Parse(rawToken, func(token *jwt.Token) (any, error) {
		// 許可されたアルゴリズムか確認
		if !slices.Contains(provider.Algos, token.Method.Alg()) {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		// algがES256, RS256の場合は、公開鍵を取得
		// HS256の場合は、シークレットを返す
		if token.Method.Alg() == jwt.SigningMethodHS256.Alg() {
			return []byte(provider.ClientSecret(p.conf)), nil
		}

		// キーIDを取得
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found")
		}

		// キーセットから対応する公開鍵を取得
		key, ok := keySet.LookupKeyID(kid)
		if !ok {
			return nil, fmt.Errorf("key %v not found", kid)
		}
		var rawKey any
		if err := jwk.Export(key, &rawKey); err != nil {
			utils.HandleError(ctx, err, "jwk.Export")
			return nil, fmt.Errorf("failed to export JWK: %v", err)
		}

		return rawKey, nil
	})
	if err != nil {
		utils.HandleError(ctx, err, "jwt.Parse")
		return nil, fmt.Errorf("failed to parse JWT: %v", err)
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		utils.HandleError(ctx, nil, "invalid token")
		return nil, fmt.Errorf("failed to get claims")
	}

	// audienceの検証
	aud, err := jwtToken.Claims.GetAudience()
	if err != nil {
		utils.HandleError(ctx, err, "GetAudience")
		return nil, fmt.Errorf("failed to get audience")
	}
	if !slices.Contains(aud, p.oauthConf.ClientID) {
		return nil, fmt.Errorf("invalid audience")
	}
	// issuerの検証
	iss, err := jwtToken.Claims.GetIssuer()
	if err != nil {
		utils.HandleError(ctx, err, "GetIssuer")
		return nil, fmt.Errorf("failed to get issuer")
	}
	if iss != provider.Issuer.String() {
		return nil, fmt.Errorf("invalid issuer")
	}

	userInfo := &UserInfo{
		Subject: claims["sub"].(string),
		Picture: claims["picture"].(string),
	}

	return userInfo, nil
}
