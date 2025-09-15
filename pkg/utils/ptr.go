package utils

import (
	"database/sql"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/neko-dream/server/internal/presentation/oas"
)

// Deprecated: utils.ToPtrIf を使用すること
// ToPtrIfNotNullValue の場合は、nullFlagがtrueの場合にnilを返す
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

// ToOpt ポインタ型をoas.OptXXX型に変換する
// code:
//
//	str := lo.ToPtr("example")
//	optStr := utils.ToOpt[oas.OptString](str) // oas.OptString{Value: "example", Set: true}
//	str = nil
//	optStr = utils.ToOpt[oas.OptString](str) // oas.OptString{Set: false}
//	num := 29
//	optNum := utils.ToOpt[oas.OptInt](num) // oas.OptInt{Value: 29, Set: true}
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

// ToSQLNull ポインタ型をsql.NullXXX型に変換する
// code:
//
//	str := lo.ToPtr("example")
//	nullStr := utils.ToSQLNull[sql.NullString](str) // sql.NullString{String: "example", Valid: true}
//	str = nil
//	nullStr = utils.ToSQLNull[sql.NullString](str) // sql.NullString{Valid: false}
//	num := 29
//	nullNum := utils.ToSQLNull[sql.NullInt32](num) // sql.NullInt32{Int32: 29, Valid: true}
func ToSQLNull[O any](v any) O {
	switch val := v.(type) {
	case *string:
		if val == nil {
			return any(sql.NullString{}).(O)
		} else {
			return any(sql.NullString{String: *val, Valid: true}).(O)
		}
	case string:
		return any(sql.NullString{String: val, Valid: true}).(O)
	case *int:
		if val == nil {
			return any(sql.NullInt32{}).(O)
		} else {
			return any(sql.NullInt32{Int32: int32(*val), Valid: true}).(O)
		}
	case int:
		return any(sql.NullInt32{Int32: int32(val), Valid: true}).(O)
	case *int32:
		if val == nil {
			return any(sql.NullInt32{}).(O)
		} else {
			return any(sql.NullInt32{Int32: *val, Valid: true}).(O)
		}
	case int32:
		return any(sql.NullInt32{Int32: val, Valid: true}).(O)
	case *int64:
		if val == nil {
			return any(sql.NullInt64{}).(O)
		} else {
			return any(sql.NullInt64{Int64: *val, Valid: true}).(O)
		}
	case int64:
		return any(sql.NullInt64{Int64: val, Valid: true}).(O)
	case *bool:
		if val == nil {
			return any(sql.NullBool{}).(O)
		} else {
			return any(sql.NullBool{Bool: *val, Valid: true}).(O)
		}
	case bool:
		return any(sql.NullBool{Bool: val, Valid: true}).(O)
	case *float32:
		if val == nil {
			return any(sql.NullFloat64{}).(O)
		} else {
			return any(sql.NullFloat64{Float64: float64(*val), Valid: true}).(O)
		}
	case float32:
		return any(sql.NullFloat64{Float64: float64(val), Valid: true}).(O)
	case *float64:
		if val == nil {
			return any(sql.NullFloat64{}).(O)
		} else {
			return any(sql.NullFloat64{Float64: *val, Valid: true}).(O)
		}
	case float64:
		return any(sql.NullFloat64{Float64: val, Valid: true}).(O)
	case *time.Time:
		if val == nil {
			return any(sql.NullTime{}).(O)
		} else {
			return any(sql.NullTime{Time: *val, Valid: true}).(O)
		}
	case time.Time:
		return any(sql.NullTime{Time: val, Valid: true}).(O)
	case *uuid.UUID:
		if val == nil {
			return any(uuid.NullUUID{}).(O)
		} else {
			return any(uuid.NullUUID{UUID: *val, Valid: true}).(O)
		}
	case uuid.UUID:
		return any(uuid.NullUUID{UUID: val, Valid: true}).(O)
	default:
		var zero O
		return zero
	}
}
