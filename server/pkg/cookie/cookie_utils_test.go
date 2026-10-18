package cookie_utils

import (
	"net/http"
	"testing"
	"time"
)

func TestEncodeRevokeCookie(t *testing.T) {
	encoded := EncodeCookies([]*http.Cookie{{Name: "SessionId", Path: "/", MaxAge: -1, HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode}})
	if len(encoded) != 1 || encoded[0] != "SessionId=; Path=/; Max-Age=0; HttpOnly; Secure; SameSite=Lax" {
		t.Fatalf("cookie must expire immediately: %v", encoded)
	}
}

// 値を URL エンコードすると、ブラウザにはエンコード後の値が保存され、読み取り側（http.Request.Cookie）と
// 値が一致しなくなる。base64 の + / = は Cookie の値に使える文字なので、そのまま出す。
func TestEncodeCookiesKeepsValueAsIs(t *testing.T) {
	value := "Q4Whk+UbhG/pBen4H6gx5zGt+aNF="
	encoded := EncodeCookies([]*http.Cookie{{Name: "AdminRefreshToken", Value: value, Path: "/v1/admin", MaxAge: 60, HttpOnly: true, SameSite: http.SameSiteLaxMode}})
	if len(encoded) != 1 {
		t.Fatalf("encoded = %v", encoded)
	}
	req := &http.Request{Header: http.Header{"Cookie": {encoded[0]}}}
	c, err := req.Cookie("AdminRefreshToken")
	if err != nil || c.Value != value {
		t.Fatalf("value read back = %q (%v), Set-Cookie = %q", c.Value, err, encoded[0])
	}
}

func TestEncodeCookiesMatchesStandardLibrary(t *testing.T) {
	cookies := []*http.Cookie{
		{Name: "SessionId", Value: "a.b-c_d", Path: "/", Domain: "kotohiro.com", MaxAge: 3600, HttpOnly: true, Secure: true, SameSite: http.SameSiteNoneMode},
		{Name: "__Host-PasswordLoginPending", Value: "x", Path: "/", Expires: time.Date(2026, 10, 18, 12, 0, 0, 0, time.UTC), MaxAge: 900, HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode},
		{Name: "AdminSessionId", Path: "/v1", MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteLaxMode},
	}
	encoded := EncodeCookies(cookies)
	for i, c := range cookies {
		if encoded[i] != c.String() {
			t.Fatalf("encoded[%d] = %q, want %q", i, encoded[i], c.String())
		}
	}
}
