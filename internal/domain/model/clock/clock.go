package clock

import (
	"context"
	"go.opentelemetry.io/otel"
	"time"
)

var (
	// timeKey 現在時刻を取得するためのContextKey
	timeKey = &struct{}{}
)

// Now 現在時刻を取得する
func Now(ctx context.Context) time.Time {
	ctx, span := otel.Tracer("clock").Start(ctx, "Now")
	defer span.End()

	if t, ok := ctx.Value(timeKey).(time.Time); ok {
		if !t.IsZero() {
			return t
		}
	}
	return time.Now()
}

// SetNow 現在時刻を設定する
func SetNow(ctx context.Context, t time.Time) context.Context {
	ctx, span := otel.Tracer("clock").Start(ctx, "SetNow")
	defer span.End()

	if t.IsZero() {
		t = time.Now()
	}
	return context.WithValue(ctx, timeKey, t)
}
