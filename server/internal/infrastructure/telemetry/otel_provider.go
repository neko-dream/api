package telemetry

import (
	"github.com/neko-dream/api/internal/infrastructure/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Otelのトレースプロバイダーを設定
// 基本的にプロバイダごとの設定はOpenTelemetry Collector側で行い、ここではOpenTelemetryCollectorへの通信のみを設定する
// トレースを無効化してnoopプロバイダーを使用
func SetupOtelProvider(conf *config.Config) trace.TracerProvider {
	noopProvider := trace.NewNoopTracerProvider()
	otel.SetTracerProvider(noopProvider)
	return noopProvider
}
