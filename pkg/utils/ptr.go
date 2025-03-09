package utils

import (
	"net/url"

	"github.com/neko-dream/server/internal/presentation/oas"
)

func ToPtrIfNotNullValue[T any](nullFlag bool, value T) *T {
	if nullFlag {
		return nil
	}
	return &value
}

func ToPtrIfNotNullFunc[T any](nullFlag bool, getValue func() *T) *T {
	if nullFlag {
		return nil
	}
	val := getValue()
	return val
}

func ToOpt[O any](v any) O {
	switch val := v.(type) {
	case *string:
		if val == nil {
			return any(oas.OptString{}).(O)
		} else {
			return any(oas.OptString{Value: *val, Set: true}).(O)
		}
	case string:
		return any(oas.OptString{Value: val, Set: true}).(O)
	case *int:
		if val == nil {
			return any(oas.OptInt{}).(O)
		} else {
			return any(oas.OptInt{Value: *val, Set: true}).(O)
		}
	case int:
		return any(oas.OptInt{Value: val, Set: true}).(O)
	case float64:
		return any(oas.OptFloat64{Value: val, Set: true}).(O)
	case *float64:
		if val == nil {
			return any(oas.OptFloat64{}).(O)
		} else {
			return any(oas.OptFloat64{Value: *val, Set: true}).(O)
		}
	case *bool:
		if val == nil {
			return any(oas.OptBool{}).(O)
		} else {
			return any(oas.OptBool{Value: *val, Set: true}).(O)
		}
	case bool:
		return any(oas.OptBool{Value: val, Set: true}).(O)
	case *url.URL:
		if val == nil {
			return any(oas.OptURI{}).(O)
		} else {
			return any(oas.OptURI{Value: *val, Set: true}).(O)
		}
	case url.URL:
		return any(oas.OptURI{Value: val, Set: true}).(O)
	default:
		var zero O
		return zero
	}
}

func ToOptNil[O any](v any) O {
	switch val := v.(type) {
	case *string:
		if val == nil {
			return any(oas.OptNilString{}).(O)
		} else {
			return any(oas.OptNilString{Value: *val, Set: true}).(O)
		}
	case *int:
		if val == nil {
			return any(oas.OptNilInt{}).(O)
		} else {
			return any(oas.OptNilInt{Value: *val, Set: true}).(O)
		}
	case *bool:
		if val == nil {
			return any(oas.OptNilBool{}).(O)
		} else {
			return any(oas.OptNilBool{Value: *val, Set: true}).(O)
		}
	default:
		var zero O
		return zero
	}
}
