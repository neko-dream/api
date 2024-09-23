package user

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/pkg/oauth"
)

type UserName string

func (a UserName) String() string {
	return string(a)
}

type UserSubject string

func (a UserSubject) String() string {
	return string(a)
}

type (
	UserRepository interface {
		Create(context.Context, User) (User, error)
		FindByID(context.Context, shared.UUID[User]) (*User, error)
		FindBySubject(context.Context, UserSubject) (*User, error)
	}

	User struct {
		userID      shared.UUID[User]
		displayID   string
		displayName string
		subject     string
		provider    oauth.AuthProviderName
	}
)

func (u *User) ChangeName(name string) {
	u.displayName = name
}

func (u *User) UserID() shared.UUID[User] {
	return u.userID
}

func (u *User) DisplayName() string {
	return u.displayName
}

func (u *User) Provider() oauth.AuthProviderName {
	return u.provider
}

func NewUser(
	userID shared.UUID[User],
	displayID string,
	displayName string,
	subject string,
	provider oauth.AuthProviderName,
) User {
	return User{
		userID:      userID,
		displayID:   displayID,
		displayName: displayName,
		subject:     subject,
		provider:    provider,
	}
}
