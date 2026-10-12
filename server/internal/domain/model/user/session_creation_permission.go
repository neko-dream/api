package user

import (
	"context"
	"errors"

	"github.com/neko-dream/api/internal/domain/model/shared"
)

var ErrSessionCreationPermissionIneligible = errors.New("session creation permission cannot be granted to an inactive or guest user")

// SessionCreationPermissionRepository stores the operator's grant, independently of organizations.
type SessionCreationPermissionRepository interface {
	Get(context.Context, shared.UUID[User]) (bool, error)
	Set(context.Context, shared.UUID[User], bool) error
}

func (u *User) CanCreateTalkSession(granted bool) bool {
	return CanCreateTalkSession(granted, u.IsGuest(), u.IsWithdrawn())
}

func CanCreateTalkSession(granted, guest, withdrawn bool) bool {
	return granted && !guest && !withdrawn
}
