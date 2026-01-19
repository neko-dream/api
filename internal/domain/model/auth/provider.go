package auth

//go:generate go tool mockgen -source=$GOFILE -package=mock_auth -destination=../../mock/auth/provider.go

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
)
