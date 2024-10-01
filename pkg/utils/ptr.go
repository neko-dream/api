package utils

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
