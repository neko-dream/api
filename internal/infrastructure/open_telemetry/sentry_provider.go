package opentelemetry

import (
	"github.com/getsentry/sentry-go"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func SentryProvider(conf *config.Config) *sdktrace.TracerProvider {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:                conf.SENTRY_DSN,
		EnableTracing:      true,
		TracesSampleRate:   1.0,
		ProfilesSampleRate: 0.1,
		Debug:              false,
	}); err != nil {
		panic(err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(sentryotel.NewSentryPropagator())

	return tp
}
