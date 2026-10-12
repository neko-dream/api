package auth_usecase

import (
	"context"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"testing"
)

func TestDevelopmentAuthenticationFailsClosed(t *testing.T) {
	for _, env := range []config.ENV{config.PROD, "", "prod", "dev", "test", "Production"} {
		t.Run(string(env), func(t *testing.T) {
			cfg := &config.Config{Env: env}
			// Nil dependencies ensure rejection happens before any user/DB/session access.
			login := &loginForDevInteractor{Config: cfg}
			if out, err := login.Execute(context.Background(), LoginForDevInput{Subject: "test"}); err == nil || out.Token != "" {
				t.Fatal("unsafe dev login allowed")
			}
			register := &passwordRegisterInteractor{cfg: cfg}
			if out, err := register.Execute(context.Background(), PasswordRegisterInput{Email: "test@example.com", Password: "test"}); err == nil || out != nil {
				t.Fatal("unsafe dev register allowed")
			}
		})
	}
}
