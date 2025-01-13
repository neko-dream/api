package oauth

import (
	"context"
	"errors"
	"fmt"

	"braces.dev/errtrace"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
	"golang.org/x/exp/slices"
	"golang.org/x/oauth2"
)

type (
	authProvider struct {
		providerName auth.AuthProviderName
		oauthConf    oauth2.Config
		conf         *config.Config
		provider     *oidc.Provider
	}
)

func NewAuthProvider(
	ctx context.Context,
	providerName auth.AuthProviderName,
	conf *config.Config,
) (auth.AuthProvider, error) {
	ctx, span := otel.Tracer("oauth").Start(ctx, "NewAuthProvider")
	defer span.End()

	issuer := providerName.IssuerURI()
	provider, ok := GetProvider(issuer, conf)
	if !ok {
		return nil, errtrace.Wrap(errors.New("provider not found"))
	}

	config := oauth2.Config{
		ClientID:     provider.ClientID(),
		ClientSecret: provider.ClientSecret(),
		RedirectURL:  provider.RedirectURL(),
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

	return &authProvider{
		providerName: providerName,
		oauthConf:    config,
		conf:         conf,
		provider:     oidcProv,
	}, nil
}

func (p *authProvider) GetAuthorizationURL(ctx context.Context, state string) string {
	ctx, span := otel.Tracer("oauth").Start(ctx, "authProvider.GetAuthorizationURL")
	defer span.End()

	_ = ctx

	return p.oauthConf.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (p *authProvider) VerifyAndIdentify(ctx context.Context, code string) (*string, *string, error) {
	ctx, span := otel.Tracer("oauth").Start(ctx, "authProvider.VerifyAndIdentify")
	defer span.End()

	token, err := p.oauthConf.Exchange(ctx, code)
	if err != nil {
		return nil, nil, errtrace.Wrap(err)
	}
	rawToken, ok := token.Extra("id_token").(string)
	if !ok {
		utils.HandleError(ctx, err, "id_token not found")
		return nil, nil, errtrace.Wrap(errors.New("id_token not found"))
	}
	provider := p.getProvider()
	// Google の公開鍵を取得
	keySet, err := jwk.Fetch(context.Background(), provider.JwksURI)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch JWK: %v", err)
	}

	jwtToken, err := jwt.Parse(rawToken, func(token *jwt.Token) (any, error) {
		// 許可されたアルゴリズムか確認
		if !slices.Contains(provider.Algos, token.Method.Alg()) {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		// algがES256, RS256の場合は、公開鍵を取得
		// HS256の場合は、シークレットを返す
		if token.Method.Alg() == jwt.SigningMethodHS256.Alg() {
			return []byte(provider.ClientSecret()), nil
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
		return nil, nil, fmt.Errorf("failed to parse JWT: %v", err)
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		utils.HandleError(ctx, nil, "invalid token")
		return nil, nil, fmt.Errorf("failed to get claims")
	}

	// audienceの検証
	aud, err := jwtToken.Claims.GetAudience()
	if err != nil {
		utils.HandleError(ctx, err, "GetAudience")
		return nil, nil, fmt.Errorf("failed to get audience")
	}
	if !slices.Contains(aud, p.oauthConf.ClientID) {
		return nil, nil, fmt.Errorf("invalid audience")
	}
	// issuerの検証
	iss, err := jwtToken.Claims.GetIssuer()
	if err != nil {
		utils.HandleError(ctx, err, "GetIssuer")
		return nil, nil, fmt.Errorf("failed to get issuer")
	}
	if iss != provider.Issuer.String() {
		return nil, nil, fmt.Errorf("invalid issuer")
	}

	var picture, subject *string
	if v, ok := claims["picture"].(string); ok {
		picture = &v
	}
	if v, ok := claims["sub"].(string); ok {
		subject = &v
	}

	return subject, picture, nil
}

func (p *authProvider) getProvider() *Provider {
	issuer := p.providerName.IssuerURI()
	provider, _ := GetProvider(issuer, p.conf)
	return &provider
}
