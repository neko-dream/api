package jwt_test

import (
	"context"
	"testing"

	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/auth/jwt"
	"github.com/neko-dream/server/internal/infrastructure/di"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
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
		},
		{
			name:         "fail",
			firstSecret:  "secret",
			secondSecret: "different_secret",
			success:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cont := di.BuildContainer()
			dbm := di.Invoke[*db.DBManager](cont)
			sessRepo := di.Invoke[session.SessionRepository](cont)
			orgUserRepo := di.Invoke[organization.OrganizationUserRepository](cont)
			orgRepo := di.Invoke[organization.OrganizationRepository](cont)

			util := jwt.NewTokenManagerWithSecret(tt.firstSecret, dbm, sessRepo, orgUserRepo, orgRepo)
			token, err := util.Generate(
				tt.ctx,
				user.NewUser(
					shared.NewUUID[user.User](),
					lo.ToPtr("test"),
					lo.ToPtr("test"),
					"test",
					auth.ProviderGoogle,
					nil,
				),
				shared.NewUUID[session.Session](),
			)

			if err != nil {
				t.Errorf("error: %v", err)
			}

			util = jwt.NewTokenManagerWithSecret(tt.secondSecret, dbm, sessRepo, orgUserRepo, orgRepo)
			_, err = util.Parse(tt.ctx, token)
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
