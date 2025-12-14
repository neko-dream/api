package telemetry

import (
	"context"
	"time"

	"github.com/neko-dream/api/internal/infrastructure/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Otelのトレースプロバイダーを設定
// 基本的にプロバイダごとの設定はOpenTelemetry Collector側で行い、ここではOpenTelemetryCollectorへの通信のみを設定する
func SetupOtelProvider(conf *config.Config) *sdktrace.TracerProvider {
	ctx := context.Background()
	var exporter sdktrace.SpanExporter
	var err error

	exporter, err = otlptrace.New(
		ctx,
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint("0.0.0.0:4318"),
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithTimeout(time.Second),
		),
	)
	if err != nil {
		panic(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)

	return tp
}
