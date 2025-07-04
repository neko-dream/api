package oauth

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

// devAuthProvider 開発環境専用の認証プロバイダーです。本番環境では使用しないでください。
type (
	devAuthProvider struct {
		providerName shared.AuthProviderName
		conf         *config.Config
	}
)

func NewDevAuthProvider(
	ctx context.Context,
	providerName shared.AuthProviderName,
	conf *config.Config,
) (auth.AuthProvider, error) {
	_, span := otel.Tracer("oauth").Start(ctx, "NewDevAuthProvider")
	defer span.End()

	return &devAuthProvider{
		providerName: providerName,
		conf:         conf,
	}, nil
}

// GetAuthorizationURL implements auth.AuthProvider.
func (d *devAuthProvider) GetAuthorizationURL(ctx context.Context, state string) string {
	ctx, span := otel.Tracer("oauth").Start(ctx, "devAuthProvider.GetAuthorizationURL")
	defer span.End()

	_ = ctx

	panic("呼ばれないやつ")
}

// VerifyAndIdentify Auth処理を挟まないため、codeをそのままsubjectとして使用する
func (d *devAuthProvider) VerifyAndIdentify(ctx context.Context, code string) (*string, *string, error) {
	ctx, span := otel.Tracer("oauth").Start(ctx, "devAuthProvider.VerifyAndIdentify")
	defer span.End()

	_ = ctx

	// codeをそのままsubjectとして使用する
	return lo.ToPtr(code), nil, nil
}
