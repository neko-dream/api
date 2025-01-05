package mapper

import (
	"errors"
	"reflect"

	"github.com/jinzhu/copier"
)

// TypeConverterRegistry グローバルで利用するTypeConverterを登録するための構造体
type TypeConverterRegistry struct {
	converters []copier.TypeConverter
}

var globalRegistry = &TypeConverterRegistry{}

// RegisterConverter グローバルで利用するTypeConverterを登録する
func RegisterConverter[From, To any](fn func(From) (To, error)) {
	var fromZero From
	var toZero To
	converter := copier.TypeConverter{
		SrcType: getTypeInstance(fromZero),
		DstType: getTypeInstance(toZero),
		Fn: func(src interface{}) (dst interface{}, err error) {
			from, ok := src.(From)
			if !ok {
				return nil, errors.New("type assertion failed")
			}
			to, err := fn(from)
			return to, err
		},
	}
	globalRegistry.converters = append(globalRegistry.converters, converter)
}

// getTypeInstance は型のゼロ値からTypeConverterで使用する型情報を取得します
func getTypeInstance(v any) any {
	if reflect.TypeOf(v).Kind() == reflect.Ptr {
		return reflect.New(reflect.TypeOf(v).Elem()).Interface()
	}
	return v
}

func GetGlobalRegistry() *TypeConverterRegistry {
	return globalRegistry
}
