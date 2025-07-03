package user

import (
	"context"
	"unicode/utf8"

	"github.com/neko-dream/server/internal/domain/messages"
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
		userID        shared.UUID[User]
		displayID     *string
		displayName   *string
		iconURL       *string
		subject       string
		provider      shared.AuthProviderName
		demographics  *UserDemographic
		email         *string
		emailVerified bool
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

func (u *User) IconURL() *string {
	return u.iconURL
}
func (u *User) ChangeIconURL(url *string) {
	u.iconURL = url
}

func (u *User) DeleteIcon() {
	u.iconURL = nil
}

func (u *User) Provider() shared.AuthProviderName {
	return u.provider
}

func (u *User) Email() *string {
	return u.email
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

// ChangeEmail メールアドレスをセットする
func (u *User) ChangeEmail(email string) {
	u.email = lo.ToPtr(email)
}

func (u *User) SetEmailVerified(v bool) {
	u.emailVerified = v
}

func (u *User) IsEmailVerified() bool {
	return u.emailVerified
}

// Withdraw anonymizes the user data for withdrawal
func (u *User) Withdraw(ctx context.Context) {
	ctx, span := otel.Tracer("user").Start(ctx, "User.Withdraw")
	defer span.End()

	// Anonymize user data
	u.displayID = lo.ToPtr("deleted_user")
	u.displayName = lo.ToPtr("削除されたユーザー")
	u.iconURL = nil
}

// IsWithdrawn checks if the user has been withdrawn based on anonymized data
func (u *User) IsWithdrawn() bool {
	return u.displayID != nil && *u.displayID == "deleted_user"
}

func NewUser(
	userID shared.UUID[User],
	displayID *string,
	displayName *string,
	subject string,
	provider shared.AuthProviderName,
	iconURL *string,
) User {
	return User{
		userID:      userID,
		displayID:   displayID,
		displayName: displayName,
		subject:     subject,
		provider:    provider,
		iconURL:     iconURL,
	}
}
