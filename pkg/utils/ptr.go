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
