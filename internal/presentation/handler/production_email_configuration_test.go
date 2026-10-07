package handler

import (
	"encoding/json"
	"errors"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/presentation/oas"
	"os"
	"testing"
)

func TestProductionEmailConfiguration(t *testing.T) {
	raw, err := os.ReadFile("../../../.ecspresso/api/prd/ecs-task-def.json")
	if err != nil {
		t.Fatal(err)
	}
	var task struct {
		Containers []struct {
			Name        string                         `json:"name"`
			Environment []struct{ Name, Value string } `json:"environment"`
		} `json:"containerDefinitions"`
	}
	if err := json.Unmarshal(raw, &task); err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{}
	for _, container := range task.Containers {
		if container.Name != "api" {
			continue
		}
		for _, variable := range container.Environment {
			vars[variable.Name] = variable.Value
		}
	}
	if vars["ENV"] != "production" || vars["WEBSITE_URL"] != "https://www.kotohiro.com" ||
		vars["EMAIL_VERIFICATION_PAGE_URL"] != "https://www.kotohiro.com/auth/email/verify" ||
		vars["EMAIL_VERIFICATION_BASE_URL"] != "https://api.kotohiro.com" {
		t.Fatal("unsafe/inconsistent production task configuration")
	}
	a, fake := emailHandlerHarness(config.PROD)
	a.cfg.WEBSITE_URL = vars["WEBSITE_URL"]
	a.cfg.EMAIL_VERIFICATION_BASE_URL = vars["EMAIL_VERIFICATION_BASE_URL"]
	for _, method := range []string{"link", "PIN"} {
		ctx, _ := emailHandlerContext("POST", "https://www.kotohiro.com", "pending")
		var redirect string
		if method == "link" {
			res, err := a.VerifyPasswordEmail(ctx, &oas.VerifyPasswordEmailReq{Token: "test-token"})
			if err != nil {
				t.Fatal(err)
			}
			redirect = res.(*oas.VerifyPasswordEmailOKHeaders).Response.RedirectURL
		} else {
			res, err := a.VerifyPasswordEmailPIN(ctx, &oas.VerifyPasswordEmailPINReq{Pin: "000001"})
			if err != nil {
				t.Fatal(err)
			}
			redirect = res.(*oas.VerifyPasswordEmailPINOKHeaders).Response.RedirectURL
		}
		if redirect != "https://www.kotohiro.com/home" {
			t.Fatal("wrong production redirect", redirect)
		}
	}
	called := fake.verifies
	for _, origin := range []string{"https://localhost:3000", "https://develop.kotohiro.com", "https://evil.example"} {
		ctx, _ := emailHandlerContext("POST", origin, "pending")
		if _, err := a.VerifyPasswordEmailPIN(ctx, &oas.VerifyPasswordEmailPINReq{Pin: "000001"}); !errors.Is(err, messages.ForbiddenError) {
			t.Fatal("untrusted origin accepted", origin, err)
		}
	}
	if fake.verifies != called {
		t.Fatal("unsafe origin reached authentication")
	}
	ctx, _ := emailHandlerContext("POST", "https://kotohiro.com", "pending")
	if _, err := a.VerifyPasswordEmailPIN(ctx, &oas.VerifyPasswordEmailPINReq{Pin: "000001"}); err != nil {
		t.Fatal("production apex origin rejected", err)
	}
}
