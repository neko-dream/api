package telemetry

import (
	"context"

	"github.com/getsentry/sentry-go"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func SentryProvider(conf *config.Config) *sdktrace.TracerProvider {
	if err := sentry.Init(sentry.ClientOptions{
		Environment:      conf.Env.String(),
		Dsn:              conf.SENTRY_DSN,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		Debug:            false,
		ServerName:       "kotohiro-server",
	}); err != nil {
		panic(err)
	}
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("kotohiro-server"),
			attribute.String("environment", conf.Env.String()),
		),
	)
	if err != nil {
		return nil
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor()),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(sentryotel.NewSentryPropagator())

	return tp
}
