package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	http_utils "github.com/neko-dream/api/pkg/http"
)

// ハンドラーは ctx からリクエスト（Cookie など）とレスポンス（ヘッダーの追加）を取り出すので、
// ルートの有無に関係なく必ず ctx に入れる
func TestWithHTTPContext(t *testing.T) {
	var gotReq *http.Request
	var gotRes http.ResponseWriter
	h := WithHTTPContext()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotReq = http_utils.GetHTTPRequest(r.Context())
		gotRes = http_utils.GetHTTPResponse(r.Context())
	}))

	req := httptest.NewRequest(http.MethodPost, "/v1/admin/refresh", nil)
	req.AddCookie(&http.Cookie{Name: "AdminRefreshToken", Value: "token"})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if gotReq == nil || gotRes == nil {
		t.Fatalf("request/response not in context: req=%v res=%v", gotReq, gotRes)
	}
	if c, err := gotReq.Cookie("AdminRefreshToken"); err != nil || c.Value != "token" {
		t.Fatalf("cookie from context request: %v %v", c, err)
	}
	if gotRes != rec {
		t.Fatal("response writer in context is not the one passed to the handler")
	}
}
