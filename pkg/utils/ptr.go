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

// ToPtrIf conditionがtrueの場合にポインタを返す
// 変換元がsql.NullXXXのように、Validフィールドを持つような場合にの使用する。
// 単純なポインタ変換には、lo.ToPtrを使用すること。
// code:
//
//	// ok
//	num := sql.NullInt32{Int32: 29, Valid: true}
//	ptrNum := utils.ToPtrIf(num.Valid, num.Int32) // *int32 (29)
//	num = sql.NullInt32{Valid: false}
//	nilNum := utils.ToPtrIf(false, num.Int32) // nil
//	// bad
//	num := lo.ToPtr(10) // *int
//	nilNum := utils.ToPtrIf(num != nil, *num) // panic
func ToPtrIf[T any](condition bool, value T) *T {
	if condition {
		return &value
	}
	return nil
}

// MarshalToPtr conditionがtrueの場合に、valueをmarshalしてポインタを返す
func OptionalMarshalToPtr[T encoding.TextMarshaler, U any](
	isSet bool,
	value T,
	convert func(string) U,
) *U {
	if !isSet {
		return nil
	}
	if bytes, err := value.MarshalText(); err == nil {
		return lo.ToPtr(convert(string(bytes)))
	}
	return nil
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

// ToNullableSQL ポインタ型をsql.NullXXX型に変換する
// code:
//
//	str := lo.ToPtr("example")
//	nullStr := utils.ToNullableSQL[sql.NullString](str) // sql.NullString{String: "example", Valid: true}
//	str = nil
//	nullStr = utils.ToNullableSQL[sql.NullString](str) // sql.NullString{Valid: false}
//	num := 29
//	nullNum := utils.ToNullableSQL[sql.NullInt32](num) // sql.NullInt32{Int32: 29, Valid: true}
func ToNullableSQL[O any](v any) O {
	rv := reflect.ValueOf(v)

	if !rv.IsValid() {
		return newNullValue[O]()
	}

	// ポインタのチェーンを辿ってnilチェック
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return newNullValue[O]()
		}
		rv = rv.Elem()
	}

	if !rv.IsValid() {
		return newNullValue[O]()
	}

	switch rv.Kind() {
	case reflect.String:
		return any(sql.NullString{String: rv.String(), Valid: true}).(O)
	case reflect.Int8, reflect.Int16:
		return any(sql.NullInt16{Int16: int16(rv.Int()), Valid: true}).(O)
	case reflect.Int, reflect.Int32:
		return any(sql.NullInt32{Int32: int32(rv.Int()), Valid: true}).(O)
	case reflect.Int64:
		return any(sql.NullInt64{Int64: rv.Int(), Valid: true}).(O)
	case reflect.Uint8:
		return any(sql.NullInt16{Int16: int16(rv.Uint()), Valid: true}).(O)
	case reflect.Uint16:
		return any(sql.NullInt32{Int32: int32(rv.Uint()), Valid: true}).(O)
	case reflect.Uint, reflect.Uint32:
		return any(sql.NullInt64{Int64: int64(rv.Uint()), Valid: true}).(O)
	case reflect.Uint64:
		if rv.Uint() > math.MaxInt64 {
			return any(sql.NullInt64{Valid: false}).(O)
		}
		return any(sql.NullInt64{Int64: int64(rv.Uint()), Valid: true}).(O)
	case reflect.Float32, reflect.Float64:
		return any(sql.NullFloat64{Float64: rv.Float(), Valid: true}).(O)
	case reflect.Bool:
		return any(sql.NullBool{Bool: rv.Bool(), Valid: true}).(O)
	default:
		val := rv.Interface()
		switch v := val.(type) {
		case time.Time:
			return any(sql.NullTime{Time: v, Valid: true}).(O)
		case uuid.UUID:
			return any(uuid.NullUUID{UUID: v, Valid: true}).(O)
		default:
			return newNullValue[O]()
		}
	}
}

// newNullValue 指定されたsql.NullXXX型のゼロ値を返す
func newNullValue[O any]() O {
	var zero O
	t := reflect.TypeOf(zero)

	switch t {
	case reflect.TypeOf(sql.NullString{}):
		return any(sql.NullString{}).(O)
	case reflect.TypeOf(sql.NullInt16{}):
		return any(sql.NullInt16{}).(O)
	case reflect.TypeOf(sql.NullInt32{}):
		return any(sql.NullInt32{}).(O)
	case reflect.TypeOf(sql.NullInt64{}):
		return any(sql.NullInt64{}).(O)
	case reflect.TypeOf(sql.NullBool{}):
		return any(sql.NullBool{}).(O)
	case reflect.TypeOf(sql.NullFloat64{}):
		return any(sql.NullFloat64{}).(O)
	case reflect.TypeOf(sql.NullTime{}):
		return any(sql.NullTime{}).(O)
	case reflect.TypeOf(uuid.NullUUID{}):
		return any(uuid.NullUUID{}).(O)
	default:
		return zero
	}
}
