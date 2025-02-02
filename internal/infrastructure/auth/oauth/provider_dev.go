package oauth

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type (
	devAuthProvider struct {
		providerName auth.AuthProviderName
		conf         *config.Config
	}
)

func NewDevAuthProvider(
	ctx context.Context,
	providerName auth.AuthProviderName,
	conf *config.Config,
) (auth.AuthProvider, error) {
	ctx, span := otel.Tracer("oauth").Start(ctx, "NewDevAuthProvider")
	defer span.End()

	return &devAuthProvider{
		providerName: providerName,
		conf:         conf,
	}, nil
}

// GetAuthorizationURL implements auth.AuthProvider.
func (d *devAuthProvider) GetAuthorizationURL(ctx context.Context, state string) string {
	panic("呼ばれないやつ")
}

// VerifyAndIdentify Auth処理を挟まないため、codeをそのままsubjectとして使用する
func (d *devAuthProvider) VerifyAndIdentify(ctx context.Context, code string) (*string, *string, error) {
	// codeをそのままsubjectとして使用する
	return lo.ToPtr(code), nil, nil
}
