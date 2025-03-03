package middleware

import (
	"fmt"
	"net/http"

	http_utils "github.com/neko-dream/server/pkg/http"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
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
