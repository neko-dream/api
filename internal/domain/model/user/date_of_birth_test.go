package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/stretchr/testify/assert"
)

// コンテキストにモッククロックを設定する関数
func withMockClock(t *testing.T, now time.Time) context.Context {
	ctx := context.Background()
    return clock.SetNow(ctx, now)
}

func TestIsValidDateFormat(t *testing.T) {
	tests := []struct {
		name     string
		date     int
		expected bool
	}{
		{"正常な日付 - 19900616", 19900616, true},
		{"正常な日付 - 20001231", 20001231, true},
		{"正常な日付 - 20000229(うるう年)", 20000229, true},
		{"不正な日付 - 0", 0, false},
		{"不正な日付 - 18991231(1900年未満)", 18991231, false},
		{"不正な日付 - 100000101(10000年以上)", 100000101, false},
		{"不正な日付 - 20001232(無効な日)", 20001232, false},
		{"不正な日付 - 20001301(無効な月)", 20001301, false},
		{"不正な日付 - 20010229(うるう年ではない)", 20010229, false},
		{"正常な日付 - 19000101(最小値)", 19000101, true},
		{"正常な日付 - 99991231(最大値)", 99991231, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := user.IsValidDateFormat(tt.date)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewDateOfBirth(t *testing.T) {
	t.Run("nilの場合", func(t *testing.T) {
		result := user.NewDateOfBirth(nil)
		assert.Nil(t, result)
	})

	t.Run("0の場合", func(t *testing.T) {
		zero := 0
		result := user.NewDateOfBirth(&zero)
		assert.Nil(t, result)
	})

	t.Run("不正な日付の場合", func(t *testing.T) {
		invalidDate := 20001301 // 13月は存在しない
		result := user.NewDateOfBirth(&invalidDate)
		assert.Nil(t, result)
	})

	t.Run("正常な日付の場合", func(t *testing.T) {
		validDate := 19900616
		result := user.NewDateOfBirth(&validDate)
		assert.NotNil(t, result)
		assert.Equal(t, user.DateOfBirth(validDate), *result)
	})
}

func TestDateOfBirth_Age(t *testing.T) {
	tests := []struct {
		name      string
		dob       user.DateOfBirth
		mockNow   time.Time
		expectedAge int
	}{
		{
			name:      "誕生日前",
			dob:       user.DateOfBirth(19900616),
			mockNow:   time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
			expectedAge: 32,
		},
		{
			name:      "誕生日当日",
			dob:       user.DateOfBirth(19900616),
			mockNow:   time.Date(2023, 6, 16, 0, 0, 0, 0, time.UTC),
			expectedAge: 33,
		},
		{
			name:      "誕生日後",
			dob:       user.DateOfBirth(19900616),
			mockNow:   time.Date(2023, 6, 17, 0, 0, 0, 0, time.UTC),
			expectedAge: 33,
		},
		{
			name:      "生まれ年の年末",
			dob:       user.DateOfBirth(19901231),
			mockNow:   time.Date(1990, 12, 31, 0, 0, 0, 0, time.UTC),
			expectedAge: 0,
		},
		{
			name:      "うるう年の誕生日 (2月29日)",
			dob:       user.DateOfBirth(20000229),
			mockNow:   time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC),
			expectedAge: 22, // うるう年でない場合、2/28を誕生日の前日と見なす
		},
		{
			name:      "うるう年の誕生日 (2月29日) - 翌日",
			dob:       user.DateOfBirth(20000229),
			mockNow:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			expectedAge: 23,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モッククロックでコンテキストを準備
			ctx := withMockClock(t, tt.mockNow)

			// 年齢計算
			age := tt.dob.Age(ctx)
			assert.Equal(t, tt.expectedAge, age)
		})
	}

	// コンテキストがnilの場合のテスト
	t.Run("コンテキストがnilの場合", func(t *testing.T) {
		dob := user.DateOfBirth(19900616)
		age := dob.Age(nil)
		// 実際の現在時刻に依存するので正確な値は予測できないが
		// エラーが発生しないことを確認するのじゃ
		assert.NotPanics(t, func() {
			_ = dob.Age(nil)
		})
		assert.GreaterOrEqual(t, age, 0)
	})
}
