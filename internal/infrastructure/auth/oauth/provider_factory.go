package oauth

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"go.opentelemetry.io/otel"
)

type providerFactory struct {
	conf *config.Config
}

func NewProviderFactory(
	conf *config.Config,
) auth.AuthProviderFactory {
	return &providerFactory{
		conf: conf,
	}
}

func (p *providerFactory) NewAuthProvider(ctx context.Context, providerName string) (auth.AuthProvider, error) {
	ctx, span := otel.Tracer("oauth").Start(ctx, "providerFactory.NewAuthProvider")
	defer span.End()

	// 本番以外の場合のみDevAuthProviderを返す
	if p.conf.Env != config.PROD && providerName == "dev" {
		return NewDevAuthProvider(ctx, auth.AuthProviderName(providerName), p.conf)
	}

	authProviderName, err := auth.NewAuthProviderName(providerName)
	if err != nil {
		return nil, err
	}

	return NewAuthProvider(ctx, authProviderName, p.conf)
}
