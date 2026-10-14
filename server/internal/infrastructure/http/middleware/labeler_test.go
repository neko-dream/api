package middleware

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type testRoute struct{}

func (testRoute) Name() string        { return "ChangePassword" }
func (testRoute) OperationID() string { return "changePassword" }
func (testRoute) PathPattern() string { return "/auth/password/change" }

func TestLabelerDoesNotRecordQueryString(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	find := func(string, *url.URL) (Route, bool) { return testRoute{}, true }

	labeler := &otelhttp.Labeler{}
	handler := Labeler(find)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))

	req := httptest.NewRequest(http.MethodPut, "/auth/password/change?oldPassword=old-secret&newPassword=new-secret", nil)
	ctx, span := tp.Tracer("test").Start(req.Context(), "request")
	ctx = otelhttp.ContextWithLabeler(ctx, labeler)
	handler.ServeHTTP(httptest.NewRecorder(), req.WithContext(ctx))
	span.End()

	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("spans = %d, want 1", len(spans))
	}
	attrs := append(spans[0].Attributes(), labeler.Get()...)
	if len(attrs) == 0 {
		t.Fatal("no attributes recorded")
	}
	for _, attr := range attrs {
		if v := attr.Value.Emit(); strings.Contains(v, "secret") {
			t.Fatalf("query string recorded in %s: %s", attr.Key, v)
		}
	}
}
