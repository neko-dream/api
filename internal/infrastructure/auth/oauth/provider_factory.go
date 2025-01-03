package oauth

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/infrastructure/config"
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
	authProviderName, err := auth.NewAuthProviderName(providerName)
	if err != nil {
		return nil, err
	}

	return NewAuthProvider(ctx, authProviderName, p.conf)
}
