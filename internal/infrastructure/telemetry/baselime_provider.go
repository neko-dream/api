package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/neko-dream/api/internal/infrastructure/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func BaselimeProvider(conf *config.Config) *sdktrace.TracerProvider {
	ctx := context.Background()
	var exporter sdktrace.SpanExporter
	var err error
	serviceName := fmt.Sprintf("kotohiro-api-%s", conf.Env)

	exporter, err = otlptrace.New(
		ctx,
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint("otel.baselime.io"),
			otlptracehttp.WithTimeout(time.Second),
			otlptracehttp.WithHeaders(map[string]string{
				"x-api-key": conf.BASELIME_API_KEY,
			}),
		),
	)
	if err != nil {
		panic(err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceNameKey.String(serviceName)),
	)
	if err != nil {
		panic(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter,
			sdktrace.WithBatchTimeout(5*time.Second),
			sdktrace.WithMaxExportBatchSize(512),
		),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)

	return tp
}
