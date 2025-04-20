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

var NilUUID = uuid.MustParse("00000000-0000-0000-0000-000000000000")
