package http_utils

import (
	"context"
	"go.opentelemetry.io/otel"
	"net/http"
)

type requestContextKey string

var (
	HTTPRequestContextKey  requestContextKey = "http_request"
	HTTPResponseContextKey requestContextKey = "http_response"
)

func WithHTTPResReqContext(ctx context.Context, req *http.Request, res http.ResponseWriter) context.Context {
	ctx, span := otel.Tracer("http_utils").Start(ctx, "WithHTTPResReqContext")
	defer span.End()

	ctx = context.WithValue(ctx, HTTPRequestContextKey, req)
	ctx = context.WithValue(ctx, HTTPResponseContextKey, res)

	return ctx
}

func GetHTTPRequest(ctx context.Context) *http.Request {
	ctx, span := otel.Tracer("http_utils").Start(ctx, "GetHTTPRequest")
	defer span.End()

	return ctx.Value(HTTPRequestContextKey).(*http.Request)
}

func GetHTTPResponse(ctx context.Context) http.ResponseWriter {
	ctx, span := otel.Tracer("http_utils").Start(ctx, "GetHTTPResponse")
	defer span.End()

	return ctx.Value(HTTPResponseContextKey).(http.ResponseWriter)
}
