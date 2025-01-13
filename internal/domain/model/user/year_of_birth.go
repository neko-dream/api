package user

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type YearOfBirth int

func NewYearOfBirth(year *int) *YearOfBirth {
	if year == nil {
		return nil
	}
	if *year == 0 {
		return nil
	}
	if *year < 1900 {
		return nil
	}
	return lo.ToPtr(YearOfBirth(*year))

}

func (y YearOfBirth) Age(ctx context.Context) int {
	ctx, span := otel.Tracer("user").Start(ctx, "YearOfBirth.Age")
	defer span.End()

	if ctx == nil {
		ctx = context.Background()
	}
	return clock.Now(ctx).Year() - int(y)
}
