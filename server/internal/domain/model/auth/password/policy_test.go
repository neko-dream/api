package password_auth

import (
	"errors"
	"strings"
	"testing"

	"github.com/neko-dream/api/internal/domain/messages"
)

func TestPasswordPolicy(t *testing.T) {
	for _, pepperBytes := range []int{0, 6, 17, 36, 60} {
		policy, err := NewPasswordPolicy(pepperBytes)
		if err != nil || policy.MinLength != 12 || policy.MaxBytes != 72-pepperBytes {
			t.Fatalf("incorrect public policy: %#v %v", policy, err)
		}
		cases := []struct {
			password string
			want     error
		}{
			{"", messages.PasswordTooShort},
			{strings.Repeat("a", 11), messages.PasswordTooShort},
			{strings.Repeat("a", 12), nil},
			{strings.Repeat("a", policy.MaxBytes), nil},
			{strings.Repeat("a", policy.MaxBytes+1), messages.PasswordTooLong},
			{" " + strings.Repeat("a", 10), messages.PasswordTooShort},
			{" " + strings.Repeat("a", 11), nil},
		}
		for _, tc := range cases {
			if err := policy.Validate(tc.password); !errors.Is(err, tc.want) {
				t.Fatalf("pepper bytes %d, password bytes %d: %v want %v", pepperBytes, len(tc.password), err, tc.want)
			}
		}
	}
	policy, _ := NewPasswordPolicy(17)
	for _, password := range []string{strings.Repeat("あ", 12), strings.Repeat("😀", 12), "あ" + strings.Repeat("a", 11)} {
		if err := policy.Validate(password); err != nil {
			t.Fatal("Unicode characters counted as bytes", err)
		}
	}
	if err := policy.Validate(strings.Repeat("😀", 11)); !errors.Is(err, messages.PasswordTooShort) {
		t.Fatal(err)
	}
	if err := policy.Validate(strings.Repeat("😀", 14)); !errors.Is(err, messages.PasswordTooLong) {
		t.Fatal(err)
	}
	for _, pepperBytes := range []int{-1, 61, 72, 100} {
		if _, err := NewPasswordPolicy(pepperBytes); err == nil {
			t.Fatal("invalid hashing configuration accepted")
		}
	}
}
