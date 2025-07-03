package user

import (
	"context"
	"time"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"go.opentelemetry.io/otel"
)

type (
	// UserAuth represents user authentication information
	UserAuth struct {
		userAuthID     shared.UUID[UserAuth]
		userID         shared.UUID[User]
		provider       string
		subject        string
		isVerified     bool
		withdrawalDate *time.Time
		createdAt      time.Time
	}

	// UserAuthRepository interface for UserAuth persistence
	UserAuthRepository interface {
		FindByUserID(context.Context, shared.UUID[User]) (*UserAuth, error)
		FindBySubject(context.Context, string) (*UserAuth, error)
		Update(context.Context, *UserAuth) error
		CheckReregistrationAllowed(context.Context, string) (bool, error)
	}
)

// NewUserAuth creates a new UserAuth
func NewUserAuth(
	userAuthID shared.UUID[UserAuth],
	userID shared.UUID[User],
	provider string,
	subject string,
	isVerified bool,
	createdAt time.Time,
) *UserAuth {
	return &UserAuth{
		userAuthID: userAuthID,
		userID:     userID,
		provider:   provider,
		subject:    subject,
		isVerified: isVerified,
		createdAt:  createdAt,
	}
}

// NewUserAuthWithWithdrawal creates a new UserAuth with withdrawal date
func NewUserAuthWithWithdrawal(
	userAuthID shared.UUID[UserAuth],
	userID shared.UUID[User],
	provider string,
	subject string,
	isVerified bool,
	createdAt time.Time,
	withdrawalDate *time.Time,
) *UserAuth {
	return &UserAuth{
		userAuthID:     userAuthID,
		userID:         userID,
		provider:       provider,
		subject:        subject,
		isVerified:     isVerified,
		createdAt:      createdAt,
		withdrawalDate: withdrawalDate,
	}
}

// UserAuthID returns the UserAuth ID
func (ua *UserAuth) UserAuthID() shared.UUID[UserAuth] {
	return ua.userAuthID
}

// UserID returns the User ID
func (ua *UserAuth) UserID() shared.UUID[User] {
	return ua.userID
}

// Provider returns the authentication provider
func (ua *UserAuth) Provider() string {
	return ua.provider
}

// Subject returns the authentication subject
func (ua *UserAuth) Subject() string {
	return ua.subject
}

// IsVerified returns whether the user auth is verified
func (ua *UserAuth) IsVerified() bool {
	return ua.isVerified
}

// WithdrawalDate returns the withdrawal date
func (ua *UserAuth) WithdrawalDate() *time.Time {
	return ua.withdrawalDate
}

// CreatedAt returns the creation time
func (ua *UserAuth) CreatedAt() time.Time {
	return ua.createdAt
}

// Withdraw sets the withdrawal date
func (ua *UserAuth) Withdraw(ctx context.Context) {
	ctx, span := otel.Tracer("user").Start(ctx, "UserAuth.Withdraw")
	defer span.End()

	now := time.Now()
	ua.withdrawalDate = &now
}

// IsWithdrawn checks if the user auth has been withdrawn
func (ua *UserAuth) IsWithdrawn() bool {
	return ua.withdrawalDate != nil
}

// CanReregister checks if enough time has passed since withdrawal for re-registration
func (ua *UserAuth) CanReregister() bool {
	if ua.withdrawalDate == nil {
		return true
	}
	
	// Allow re-registration after 30 days
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	return ua.withdrawalDate.Before(thirtyDaysAgo)
}