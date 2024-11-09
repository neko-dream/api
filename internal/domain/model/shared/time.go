package shared

import (
	"context"
	"time"
)

type Time struct {
	time.Time
}

func NewTime(ctx context.Context, t time.Time) Time {
	return Time{t}
}

func (t Time) FormatWithLayout(ctx context.Context, layout string) string {
	return t.Time.Format(layout)
}

func (t Time) Format(ctx context.Context) string {
	return t.FormatWithLayout(ctx, time.RFC3339)
}

func (t Time) Add(ctx context.Context, d time.Duration) Time {
	return Time{t.Time.Add(d)}
}

func Now(ctx context.Context) Time {
	return Time{time.Now()}
}

func Parse(ctx context.Context, str string) *Time {
	return ParseWithLayout(ctx, str, time.RFC3339)
}

func ParseWithLayout(ctx context.Context, str string, layout string) *Time {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return nil
	}
	return &Time{t}
}

const (
	Second = time.Second
	Minute = time.Minute
	Hour   = time.Hour
	Day    = 24 * time.Hour
	Week   = 7 * Day
	Month  = 30 * Day
	Year   = 365 * Day
)
