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
)
