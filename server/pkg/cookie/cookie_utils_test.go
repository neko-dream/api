package cookie_utils

import (
	"net/http"
	"testing"
)

func TestEncodeRevokeCookie(t *testing.T) {
	encoded := EncodeCookies([]*http.Cookie{{Name: "SessionId", Path: "/", MaxAge: -1, HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode}})
	if len(encoded) != 1 || encoded[0] != "SessionId=; Path=/; Max-Age=0; HttpOnly; Secure; SameSite=Lax" {
		t.Fatalf("cookie must expire immediately: %v", encoded)
	}
}
