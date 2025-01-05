package mapper

import (
	"errors"

	"github.com/jinzhu/copier"
)

// Map copier.Copyを型パラメータを利用してマップ
func Map[From, To any](from From) (To, error) {
	var to To
	err := copier.Copy(&to, from)
	return to, err
}

// 引数にコンバータを追加したMap
func MapWithConverters[From, To any](from From, converters ...copier.TypeConverter) (To, error) {
	var to To
	err := copier.CopyWithOption(&to, from, copier.Option{
		DeepCopy:   true,
		Converters: converters,
	})
	return to, err
}

func Converter[From, To any](fn func(From) (To, error)) copier.TypeConverter {
	var fromZero From
	var toZero To
	return copier.TypeConverter{
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
}
