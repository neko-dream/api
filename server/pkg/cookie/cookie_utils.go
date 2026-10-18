package cookie_utils

import (
	"net/http"
)

// EncodeCookies は ogen の Set-Cookie ヘッダー用に、Cookie を Set-Cookie の文字列にする。
// 標準ライブラリの http.Cookie.String() を使う。値は URL エンコードしない
// （以前は url.QueryEscape していたため、+ / = を含む値がエンコード後の値でブラウザに保存され、
// http.Request.Cookie で読んだ値と一致しなくなっていた）。
func EncodeCookies(cookies []*http.Cookie) []string {
	results := make([]string, 0, len(cookies))
	for _, cookie := range cookies {
		results = append(results, cookie.String())
	}
	return results
}
