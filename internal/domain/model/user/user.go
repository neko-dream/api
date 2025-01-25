package user

import (
	"context"
	"mime/multipart"
	"unicode/utf8"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
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
		profileIcon  *ProfileIcon
		subject      string
		provider     auth.AuthProviderName
		demographics *UserDemographic
	}
)

func (u *User) ChangeName(ctx context.Context, name *string) {
	ctx, span := otel.Tracer("user").Start(ctx, "User.ChangeName")
	defer span.End()

	_ = ctx

	if name == nil || *name == "" {
		return
	}
	u.displayName = name
}

func (u *User) SetDisplayID(id string) error {
	if id == "" {
		return messages.UserDisplayIDInvalidError
	}
	if utf8.RuneCountInString(id) > 30 {
		return messages.UserDisplayIDTooLong
	}
	if utf8.RuneCountInString(id) < 4 {
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

func (u *User) ProfileIconURL() *string {
	if u.profileIcon == nil {
		return nil
	} else {
		return u.profileIcon.url
	}
}

func (u *User) ProfileIcon() *ProfileIcon {
	return u.profileIcon
}

func (u *User) IsIconUpdateRequired() bool {
	return u.profileIcon != nil && u.profileIcon.url == nil && u.profileIcon.ImageInfo() != nil
}

func (u *User) SetIconFile(ctx context.Context, file *multipart.FileHeader) error {
	ctx, span := otel.Tracer("user").Start(ctx, "User.SetIconFile")
	defer span.End()
	if file == nil {
		return nil
	}

	profileIcon := NewProfileIcon(nil)
	if err := profileIcon.SetProfileIconImage(ctx, file, *u); err != nil {
		return err
	}
	u.profileIcon = profileIcon
	u.profileIcon.url = nil
	return nil
}
func (u *User) DeleteIcon() {
	u.profileIcon = nil
}

func (u *User) Provider() auth.AuthProviderName {
	return u.provider
}

func (u *User) Verify() bool {
	return u.displayID != nil && u.displayName != nil
}

func (u *User) Demographics() *UserDemographic {
	return u.demographics
}

func (u *User) SetDemographics(demographics UserDemographic) {
	u.demographics = lo.ToPtr(demographics)
}

func NewUser(
	userID shared.UUID[User],
	displayID *string,
	displayName *string,
	subject string,
	provider auth.AuthProviderName,
	profileIcon *ProfileIcon,
) User {
	return User{
		userID:      userID,
		displayID:   displayID,
		displayName: displayName,
		subject:     subject,
		provider:    provider,
		profileIcon: profileIcon,
	}
}
