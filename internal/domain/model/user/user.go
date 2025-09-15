package user

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/pkg/random"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

const (
	// ReactivationPeriodDays 復活可能期間（日数）
	ReactivationPeriodDays = 30
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
		ChangeSubject(context.Context, shared.UUID[User], string) error
	}

	// UserService ユーザードメインサービス
	UserService interface {
		// DisplayIDCheckDuplicate ユーザーの表示用IDが重複していないかチェック
		DisplayIDCheckDuplicate(context.Context, string) (bool, error)
	}

	User struct {
		userID         shared.UUID[User]
		displayID      *string
		displayName    *string
		iconURL        *string
		subject        string
		provider       shared.AuthProviderName
		demographics   *UserDemographic
		email          *string
		emailVerified  bool
		withdrawalDate *time.Time
	}
)

func (u *User) ChangeName(ctx context.Context, name *string) {
	ctx, span := otel.Tracer("user").Start(ctx, "User.ChangeName")
	defer span.End()

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

// 退会関連メソッド

// Withdraw ユーザーを退会状態にする
func (u *User) Withdraw(now time.Time) error {
	if u.withdrawalDate != nil {
		return ErrAlreadyWithdrawn
	}
	u.withdrawalDate = &now
	return nil
}

// IsWithdrawn ユーザーが退会しているかどうか
func (u *User) IsWithdrawn() bool {
	return u.withdrawalDate != nil
}

// WithdrawalDate 退会日時を取得
func (u *User) WithdrawalDate() *time.Time {
	return u.withdrawalDate
}

// CanReactivate 復活可能かどうかを判定
func (u *User) CanReactivate(now time.Time) (bool, error) {
	if u.withdrawalDate == nil {
		return false, ErrNotWithdrawn
	}

	daysSinceWithdrawal := now.Sub(*u.withdrawalDate).Hours() / 24
	if daysSinceWithdrawal > ReactivationPeriodDays {
		return false, ErrReactivationPeriodExpired
	}

	return true, nil
}

// Reactivate ユーザーを復活させる
func (u *User) Reactivate(now time.Time) error {
	if _, err := u.CanReactivate(now); err != nil {
		return err
	}
	u.withdrawalDate = nil
	return nil
}

// IsReactivationPeriodExpired 復活可能期間が過ぎているかチェック
func (u *User) IsReactivationPeriodExpired(now time.Time) bool {
	if u.withdrawalDate == nil {
		return false
	}
	daysSinceWithdrawal := now.Sub(*u.withdrawalDate).Hours() / 24
	return daysSinceWithdrawal > ReactivationPeriodDays
}

// PrepareForDeleteUser 削除前の準備
func (u *User) PrepareForDeleteUser() (newSubject string) {
	if u.withdrawalDate == nil {
		return u.subject
	}

	withdrawnSuffix := fmt.Sprintf("_withdrawn_%d", u.withdrawalDate.Unix())
	newSubject = u.subject + withdrawnSuffix
	u.subject = newSubject
	u.email = nil           // 退会後はメールアドレスをクリア
	u.emailVerified = false // メールアドレスの検証状態もクリア
	u.iconURL = nil         // アイコンURLもクリア
	u.displayID = lo.ToPtr(fmt.Sprintf("%s_%s", *u.displayID, random.GenerateRandom()))
	u.displayName = lo.ToPtr("unknown")

	return newSubject
}

// ChangeSubject subjectを変更する（31日経過後の重複回避用）
func (u *User) ChangeSubject(newSubject string) {
	u.subject = newSubject
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

// NewUserWithWithdrawalDate リポジトリから読み込む際に使用
func NewUserWithWithdrawalDate(
	userID shared.UUID[User],
	displayID *string,
	displayName *string,
	subject string,
	provider shared.AuthProviderName,
	iconURL *string,
	withdrawalDate sql.NullTime,
) User {
	user := NewUser(userID, displayID, displayName, subject, provider, iconURL)
	if withdrawalDate.Valid {
		user.withdrawalDate = &withdrawalDate.Time
	}
	return user
}
