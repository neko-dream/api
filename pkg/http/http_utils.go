package http_utils

import (
	"context"
	"net/http"
)

type requestContextKey string

var (
	HTTPRequestContextKey  requestContextKey = "http_request"
	HTTPResponseContextKey requestContextKey = "http_response"
)

func WithHTTPResReqContext(ctx context.Context, req *http.Request, res http.ResponseWriter) context.Context {
	ctx = context.WithValue(ctx, HTTPRequestContextKey, req)
	ctx = context.WithValue(ctx, HTTPResponseContextKey, res)

	return ctx
}

func GetHTTPRequest(ctx context.Context) *http.Request {
	return ctx.Value(HTTPRequestContextKey).(*http.Request)
}

func GetHTTPResponse(ctx context.Context) http.ResponseWriter {
	return ctx.Value(HTTPResponseContextKey).(http.ResponseWriter)
}
