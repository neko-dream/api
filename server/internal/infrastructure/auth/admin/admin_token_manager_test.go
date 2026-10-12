package admin

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	adminmodel "github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"testing"
	"time"
)

func TestAdminJWTValidation(t *testing.T) {
	key := "test-admin-only-secret-at-least-32-bytes"
	manager := NewAdminTokenManager(&config.Config{AdminJWTSecret: key})
	user := adminmodel.AdminUser{AdminUserID: shared.NewUUID[adminmodel.AdminUser](), Email: "private@example.com", DisplayName: "Private"}
	raw, err := manager.Generate(context.Background(), user)
	if err != nil {
		t.Fatal(err)
	}
	claim, err := manager.Parse(context.Background(), raw)
	if err != nil || claim.Sub != user.AdminUserID.String() {
		t.Fatalf("roundtrip: %v", err)
	}
	parsed, _, err := new(jwt.Parser).ParseUnverified(raw, jwt.MapClaims{})
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"email", "displayName", "isAdmin"} {
		if _, ok := parsed.Claims.(jwt.MapClaims)[field]; ok {
			t.Fatal("profile leaked into JWT", field)
		}
	}
	for _, name := range []string{"missing expiry", "expired", "wrong issuer", "wrong audience", "invalid subject", "future issued at", "wrong algorithm", "wrong key"} {
		t.Run(name, func(t *testing.T) {
			c := jwt.MapClaims{"iss": adminIssuer, "aud": adminAudience, "sub": user.AdminUserID.String(), "exp": time.Now().Add(time.Hour).Unix(), "iat": time.Now().Unix()}
			method := jwt.SigningMethodHS256
			signingKey := key
			switch name {
			case "missing expiry":
				delete(c, "exp")
			case "expired":
				c["exp"] = time.Now().Add(-time.Minute).Unix()
			case "wrong issuer":
				c["iss"] = "user"
			case "wrong audience":
				c["aud"] = "user-api"
			case "invalid subject":
				c["sub"] = "invalid"
			case "future issued at":
				c["iat"] = time.Now().Add(time.Hour).Unix()
			case "wrong algorithm":
				method = jwt.SigningMethodHS384
			case "wrong key":
				signingKey = "different-secret-at-least-32-bytes"
			}
			token, err := jwt.NewWithClaims(method, c).SignedString([]byte(signingKey))
			if err != nil {
				t.Fatal(err)
			}
			if _, err := manager.Parse(context.Background(), token); err == nil {
				t.Fatal("accepted invalid JWT")
			}
		})
	}
	if _, err := NewAdminTokenManager(&config.Config{TokenSecret: key}).Generate(context.Background(), user); err == nil {
		t.Fatal("fell back to user key")
	}
}
