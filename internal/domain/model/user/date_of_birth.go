package user

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type DateOfBirth int

func NewDateOfBirth(dateOfBirth *int) *DateOfBirth {
	if dateOfBirth == nil {
		return nil
	}
	if *dateOfBirth == 0 {
		return nil
	}

	// 1990616形式の日付をバリデーションする
	if !IsValidDateFormat(*dateOfBirth) {
		return nil
	}

	return lo.ToPtr(DateOfBirth(*dateOfBirth))
}

// 20001010形式の日付をバリデーションする
func IsValidDateFormat(date int) bool {
	if date < 19000101 || date > 99991231 {
		return false
	}

	// 年月日に分解
	year := date / 10000
	month := (date % 10000) / 100
	day := date % 100

	// 月のバリデーション (1-12)
	if month < 1 || month > 12 {
		return false
	}

	// 日のバリデーション (各月の最大日数)
	maxDays := 31
	if month == 4 || month == 6 || month == 9 || month == 11 {
		maxDays = 30
	} else if month == 2 {
		// うるう年の計算
		if year%400 == 0 || (year%4 == 0 && year%100 != 0) {
			maxDays = 29
		} else {
			maxDays = 28
		}
	}

	if day < 1 || day > maxDays {
		return false
	}

	return true
}

func (y DateOfBirth) Age(ctx context.Context) int {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, span := otel.Tracer("user").Start(ctx, "DateOfBirth.Age")
	defer span.End()
	// 年月日を分解
	birthYear := int(y) / 10000
	birthMonth := (int(y) % 10000) / 100
	birthDay := int(y) % 100

	now := clock.Now(ctx)
	currentYear := now.Year()
	currentMonth := int(now.Month())
	currentDay := now.Day()

	age := currentYear - birthYear

	// まだ誕生日が来ていなければ年齢を1つ減らす
	if currentMonth < birthMonth || (currentMonth == birthMonth && currentDay < birthDay) {
		age--
	}

	return age
}
}
