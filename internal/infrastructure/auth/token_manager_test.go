package auth_test

import (
	"context"
	"testing"

	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/auth"
	"github.com/neko-dream/server/pkg/oauth"
	"github.com/samber/lo"
)

func TestNewTokenManagerTest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		firstSecret  string
		secondSecret string
		success      bool
		ctx          context.Context
	}{
		{
			name:         "success",
			firstSecret:  "secret",
			secondSecret: "secret",
			success:      true,
			ctx:          context.Background(),
		},
		{
			name:         "fail",
			firstSecret:  "secret",
			secondSecret: "different_secret",
			success:      false,
			ctx:          context.Background(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			util := auth.NewTokenManager(tt.firstSecret)
			token, err := util.GenerateToken(
				tt.ctx,
				user.NewUser(
					shared.NewUUID[user.User](),
					"test",
					"test",
					"test",
					oauth.ProviderGoogle,
					lo.ToPtr("https://example.com"),
				),
				shared.NewUUID[session.Session](),
			)

			if err != nil {
				t.Errorf("error: %v", err)
			}

			util = auth.NewTokenManager(tt.secondSecret)
			_, err = util.ParseToken(tt.ctx, token)
			if tt.success {
				if err != nil {
					t.Errorf("error: %v", err)
				}
				return
			}
			if err == nil {
				t.Errorf("error is nil")
			}
		})
	}
}
