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
	// ブラウザと同じく Set-Cookie を解析し、名前と値だけを Cookie ヘッダーで送り返す
	stored, err := http.ParseSetCookie(encoded[0])
	if err != nil {
		t.Fatalf("parse %q: %v", encoded[0], err)
	}
	if stored.Value != value {
		t.Fatalf("stored value = %q, want %q (Set-Cookie = %q)", stored.Value, value, encoded[0])
	}
	req := &http.Request{Header: http.Header{"Cookie": {stored.Name + "=" + stored.Value}}}
	c, err := req.Cookie("AdminRefreshToken")
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	if c.Value != value {
		t.Fatalf("value read back = %q, want %q", c.Value, value)
	}
}

func TestEncodeCookiesAttributes(t *testing.T) {
	for name, tc := range map[string]struct {
		cookie *http.Cookie
		want   string
	}{
		"ドメイン・SameSite=None": {
			&http.Cookie{Name: "SessionId", Value: "a.b-c_d", Path: "/", Domain: "kotohiro.com", MaxAge: 3600, HttpOnly: true, Secure: true, SameSite: http.SameSiteNoneMode},
			"SessionId=a.b-c_d; Path=/; Domain=kotohiro.com; Max-Age=3600; HttpOnly; Secure; SameSite=None",
		},
		"Expires・SameSite=Strict": {
			&http.Cookie{Name: "__Host-PasswordLoginPending", Value: "x", Path: "/", Expires: time.Date(2026, 10, 18, 12, 0, 0, 0, time.UTC), MaxAge: 900, HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode},
			"__Host-PasswordLoginPending=x; Path=/; Expires=Sun, 18 Oct 2026 12:00:00 GMT; Max-Age=900; HttpOnly; Secure; SameSite=Strict",
		},
		"取り消し（MaxAge<0）": {
			&http.Cookie{Name: "AdminSessionId", Path: "/v1", MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteLaxMode},
			"AdminSessionId=; Path=/v1; Max-Age=0; HttpOnly; SameSite=Lax",
		},
		"base64 の値": {
			&http.Cookie{Name: "AdminRefreshToken", Value: "a+b/c=", Path: "/v1/admin", MaxAge: 60},
			"AdminRefreshToken=a+b/c=; Path=/v1/admin; Max-Age=60",
		},
	} {
		t.Run(name, func(t *testing.T) {
			if got := EncodeCookies([]*http.Cookie{tc.cookie}); len(got) != 1 || got[0] != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}
