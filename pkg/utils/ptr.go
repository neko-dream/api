package utils

func ToPtrIfNotNull[T any](nullFlag bool, value T) *T {
	if nullFlag {
		return nil
	}
	return &value
}
