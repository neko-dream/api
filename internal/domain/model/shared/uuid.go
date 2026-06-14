package shared

import (
	"braces.dev/errtrace"
	"github.com/google/uuid"
)

type UUID[T any] uuid.UUID

func NewUUID[T any]() UUID[T] {
	uid, _ := uuid.NewV7()
	return UUID[T](uid)
}

func (a UUID[T]) String() string {
	return uuid.UUID(a).String()
}

func (a UUID[T]) UUID() uuid.UUID {
	return uuid.UUID(a)
}

func (a UUID[T]) IsZero() bool {
	return uuid.UUID(a) == NilUUID
}

func MustParseUUID[T any](s string) UUID[T] {
	return UUID[T](uuid.MustParse(s))
}

func ParseUUID[T any](s string) (UUID[T], error) {
	u, err := uuid.Parse(s)
	return UUID[T](u), errtrace.Wrap(err)
}

// MarshalText は uuid.UUID の文字列形式（"xxxxxxxx-xxxx-..."）で出力する。
// 名前付きジェネリック型は uuid.UUID の MarshalText を継承しないため、明示的に実装する。
// encoding/json は TextMarshaler を引用符付きで使うため、これだけで JSONB へも正しい文字列で保存される。
func (a UUID[T]) MarshalText() ([]byte, error) {
	return uuid.UUID(a).MarshalText()
}

func (a *UUID[T]) UnmarshalText(data []byte) error {
	var u uuid.UUID
	if err := u.UnmarshalText(data); err != nil {
		return errtrace.Wrap(err)
	}
	*a = UUID[T](u)
	return nil
}

var NilUUID = uuid.MustParse("00000000-0000-0000-0000-000000000000")
