package config

import "testing"

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
