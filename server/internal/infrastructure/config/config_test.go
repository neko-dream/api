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

func TestBrowserOriginsLocalFrontendAlias(t *testing.T) {
	local := (&Config{Env: LOCAL, WEBSITE_URL: "https://localhost:3000", EMAIL_VERIFICATION_BASE_URL: "https://localhost:3000/api"}).BrowserOrigins()
	if !slices.Contains(local, "https://local.kotohiro.com:3000") {
		t.Fatalf("local frontend alias not allowed: %v", local)
	}
	// ローカル開発の別名は development / production に漏らさない
	for _, env := range []ENV{DEV, PROD} {
		cfg := &Config{Env: env, WEBSITE_URL: "https://develop.kotohiro.com", EMAIL_VERIFICATION_BASE_URL: "https://api.dev.kotohiro.com"}
		if slices.Contains(cfg.BrowserOrigins(), "https://local.kotohiro.com:3000") {
			t.Fatalf("%s trusted local frontend alias", env)
		}
	}
}

func TestBrowserOriginsAdditionalOrigins(t *testing.T) {
	cfg := &Config{
		Env:                      PROD,
		WEBSITE_URL:              "https://www.kotohiro.com",
		AdditionalBrowserOrigins: []string{"https://staging.kotohiro.com", " https://staging2.kotohiro.com/ "},
	}
	origins := cfg.BrowserOrigins()
	for _, want := range []string{"https://www.kotohiro.com", "https://kotohiro.com", "https://staging.kotohiro.com", "https://staging2.kotohiro.com"} {
		if !slices.Contains(origins, want) {
			t.Fatalf("%s not allowed: %v", want, origins)
		}
	}

	// 明示的な HTTPS オリジン以外は、設定ミスでも信頼しない
	for _, invalid := range []string{
		"",
		"*",
		"https://*",
		"https://*.kotohiro.com",
		"http://staging.kotohiro.com",
		"https://staging.kotohiro.com/path",
		"https://staging.kotohiro.com?x=1",
		"https://staging.kotohiro.com#x",
		"https://user@staging.kotohiro.com",
		"staging.kotohiro.com",
	} {
		cfg := &Config{Env: PROD, AdditionalBrowserOrigins: []string{invalid}}
		if origins := cfg.BrowserOrigins(); len(origins) != 0 {
			t.Fatalf("%q accepted: %v", invalid, origins)
		}
	}

	if origins := (&Config{Env: "unknown", AdditionalBrowserOrigins: []string{"https://staging.kotohiro.com"}}).BrowserOrigins(); origins != nil {
		t.Fatalf("unknown environment trusted origins: %v", origins)
	}
}
