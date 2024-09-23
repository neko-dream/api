package auth

import (
	"context"
	"net/url"

	"github.com/neko-dream/server/internal/domain/model/user"
)

type (
	AuthService interface {
		Authenticate(ctx context.Context, provider, code string) (*user.User, error)
		GetAuthURL(ctx context.Context, providerName string) (*url.URL, string, error)
	}
)
