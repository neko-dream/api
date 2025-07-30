package middleware

import (
	"fmt"
	"net/http"

	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/infrastructure/di"
	http_utils "github.com/neko-dream/server/pkg/http"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/dig"
)

func Labeler(find RouteFinder) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := http_utils.WithHTTPResReqContext(r.Context(), r, w)
			route, ok := find(r.Method, r.URL)
			if !ok {
				h.ServeHTTP(w, r)
				return
			}

			attrs := []attribute.KeyValue{
				semconv.HTTPRouteKey.String(route.PathPattern()),
				semconv.HTTPMethodKey.String(r.Method),
				semconv.HTTPSchemeKey.String(r.URL.Scheme),
				semconv.HTTPHostKey.String(r.URL.Host),
				semconv.HTTPTargetKey.String(r.URL.Path),
				semconv.HTTPURLKey.String(r.URL.String()),
				semconv.HTTPUserAgentKey.String(r.UserAgent()),
			}

			span := trace.SpanFromContext(r.Context())
			span.SetAttributes(attrs...)

			labeler, _ := otelhttp.LabelerFromContext(r.Context())
			labeler.Add(attrs...)

			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}

func Instrument(serviceName string, find RouteFinder, traceProvider *sdktrace.TracerProvider) Middleware {
	return func(h http.Handler) http.Handler {
		return otelhttp.NewHandler(h, "",
			otelhttp.WithPropagators(otel.GetTextMapPropagator()),
			otelhttp.WithTracerProvider(traceProvider),
			otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
			otelhttp.WithServerName(serviceName),
			otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
				op, ok := find(r.Method, r.URL)
				if ok {
					return fmt.Sprintf("%s %s", r.Method, op.PathPattern())
				}
				return operation
			}),
		)
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

			// セッションの存在確認
			_, err = sessionRepository.FindBySessionID(ctx, sessionID)
			if err != nil {
				// セッションが見つからない場合もそのまま続行
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
