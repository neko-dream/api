package config

import "testing"
import "slices"

func TestEnvironmentValidationAndDevelopmentAuthentication(t *testing.T) {
	for _, tc := range []struct {
		env                ENV
		valid, development bool
	}{
		{PROD, true, false}, {DEV, true, true}, {LOCAL, true, true},
		{"", false, false}, {"prod", false, false}, {"dev", false, false},
		{"Production", false, false}, {"test", false, false},
	} {
		t.Run(string(tc.env), func(t *testing.T) {
			if (tc.env.Validate() == nil) != tc.valid {
				t.Fatal("invalid environment validation")
			}
			if tc.env.AllowsDevelopmentAuthentication() != tc.development {
				t.Fatal("unsafe development authentication policy")
			}
		})
	}
}

func TestBrowserOriginsProductionAliasAndDevelopmentIsolation(t *testing.T) {
	for _, website := range []string{"https://www.kotohiro.com", "https://www.kotohiro.com/"} {
		cfg := &Config{Env: PROD, WEBSITE_URL: website, EMAIL_VERIFICATION_BASE_URL: "https://api.kotohiro.com"}
		if !slices.Contains(cfg.BrowserOrigins(), "https://kotohiro.com") {
			t.Fatal("canonical origin rejected apex")
		}
		if slices.Contains(cfg.BrowserOrigins(), "https://api-dev.kotohiro.com") {
			t.Fatal("development alias leaked into production")
		}
	}
}
