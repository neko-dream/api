package middleware

import (
	"net/http"

	"github.com/neko-dream/api/internal/domain/model/session"
	"github.com/neko-dream/api/internal/infrastructure/di"
	http_utils "github.com/neko-dream/api/pkg/http"
	"go.uber.org/dig"
)

// WithHTTPContext はリクエストとレスポンスを ctx に入れる。ハンドラーは ctx から Cookie を読んだり、
// ヘッダーを足したりする（http_utils.GetHTTPRequest / GetHTTPResponse）。
func WithHTTPContext() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r.WithContext(http_utils.WithHTTPResReqContext(r.Context(), r, w)))
		})
	}
}

// contextにCookieKeyをセットするミドルウェア
func SetContextCookieKey(cont *dig.Container) Middleware {
	tokenManager := di.Invoke[session.TokenManager](cont)
	sessionRepository := di.Invoke[session.SessionRepository](cont)
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			// headerからSessionIdを取得
			cookie, err := r.Cookie("SessionId")
			if err != nil {
				// Cookieが存在しない場合はそのまま続行
				h.ServeHTTP(w, r)
				return
			}

			// トークンをパース
			claim, err := tokenManager.Parse(ctx, cookie.Value)
			if err != nil {
				// パースエラーの場合もそのまま続行（オプショナル認証のため）
				h.ServeHTTP(w, r)
				return
			}

			// トークンの有効性を確認
			if claim.IsExpired(ctx) {
				// 期限切れの場合もそのまま続行（オプショナル認証のため）
				h.ServeHTTP(w, r)
				return
			}

			// セッションIDを取得して、セッションが存在するか確認
			sessionID, err := claim.SessionID()
			if err != nil {
				h.ServeHTTP(w, r)
				return
			}

			// セッションの存在確認とアクティブ状態を確認
			sess, err := sessionRepository.FindBySessionID(ctx, sessionID)
			if err != nil {
				// セッションが見つからない場合もそのまま続行
				h.ServeHTTP(w, r)
				return
			}
			// セッションが無効（revoke済み）の場合はセッションを設定しない
			if sess == nil || !sess.IsActive(ctx) {
				h.ServeHTTP(w, r)
				return
			}

			// contextにセッション情報を設定
			ctx = session.SetSession(ctx, claim)
			r = r.WithContext(ctx)

			h.ServeHTTP(w, r)
		})
	}
}

func Wrap(h http.Handler, middlewares ...Middleware) http.Handler {
	switch len(middlewares) {
	case 0:
		return h
	case 1:
		return middlewares[0](h)
	default:
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h)
		}
		return h
	}
}
