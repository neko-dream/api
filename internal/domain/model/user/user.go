package user

import (
	"context"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/pkg/oauth"
	"github.com/samber/lo"
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
		Create(context.Context, User) error
		FindByID(context.Context, shared.UUID[User]) (*User, error)
		FindBySubject(context.Context, UserSubject) (*User, error)
		FindByDisplayID(context.Context, string) (*User, error)
		Update(context.Context, User) error
	}

	// UserService ユーザードメインサービス
	UserService interface {
		// DisplayIDCheckDuplicate ユーザーの表示用IDが重複していないかチェック
		DisplayIDCheckDuplicate(context.Context, string) (bool, error)
	}

	User struct {
		userID       shared.UUID[User]
		displayID    *string
		displayName  *string
		picture      *string
		subject      string
		provider     oauth.AuthProviderName
		demographics *UserDemographics
	}
)

func (u *User) ChangeName(name string) {
	u.displayName = lo.ToPtr(name)
}

func (u *User) SetDisplayID(id string) error {
	if id == "" {
		return messages.UserDisplayIDInvalidError
	}
	if len([]rune(id)) > 30 {
		return messages.UserDisplayIDTooLong
	}
	if len([]rune(id)) < 4 {
		return messages.UserDisplayIDTooShort
	}
	u.displayID = lo.ToPtr(id)
	return nil
}

func (u *User) UserID() shared.UUID[User] {
	return u.userID
}

func (u *User) DisplayName() *string {
	return u.displayName
}

func (u *User) DisplayID() *string {
	return u.displayID
}

func (u *User) Subject() string {
	return u.subject
}

func (u *User) Picture() *string {
	return u.picture
}

func (u *User) Provider() oauth.AuthProviderName {
	return u.provider
}

func (u *User) Verify() bool {
	return u.displayID != nil && u.displayName != nil
}

func (u *User) SetDemographics(demographics UserDemographics) {
	u.demographics = lo.ToPtr(demographics)
}

func (u *User) Demographics() *UserDemographics {
	return u.demographics
}

func NewUser(
	userID shared.UUID[User],
	displayID *string,
	displayName *string,
	subject string,
	provider oauth.AuthProviderName,
	picture *string,
) User {
	return User{
		userID:      userID,
		displayID:   displayID,
		displayName: displayName,
		subject:     subject,
		provider:    provider,
		picture:     picture,
	}
}
