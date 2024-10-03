package utils

import "github.com/neko-dream/server/internal/presentation/oas"

func ToPtrIfNotNullValue[T any](nullFlag bool, value T) *T {
	if nullFlag {
		return nil
	}
	return &value
}

func ToPtrIfNotNullFunc[T any](nullFlag bool, getValue func() T) *T {
	if nullFlag {
		return nil
	}
	val := getValue()
	return &val
}

// 三項演算子
func IfThenElse[T any](condition bool, thenValue T, elseValue T) T {
	if condition {
		return thenValue
	}
	return elseValue
}

// ogen用のユーティリティ関数

func StringToOptString(s *string) oas.OptString {
	if s == nil {
		return oas.OptString{Set: false}
	}
	return oas.OptString{
		Value: *s,
		Set:   true,
	}
}
