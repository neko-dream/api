package user

import (
	"context"
	"testing"
	"time"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/stretchr/testify/assert"
)

func TestUser_Withdraw(t *testing.T) {
	// Create a test user
	userID := shared.NewUUID[User]()
	user := NewUser(
		userID,
		stringPtr("test_user"),
		stringPtr("Test User"),
		"test_subject",
		shared.AuthProviderName("google"),
		stringPtr("https://example.com/icon.jpg"),
	)

	// Verify initial state
	assert.False(t, user.IsWithdrawn())
	assert.Equal(t, "test_user", *user.DisplayID())
	assert.Equal(t, "Test User", *user.DisplayName())
	assert.NotNil(t, user.IconURL())

	// Perform withdrawal
	ctx := context.Background()
	user.Withdraw(ctx)

	// Verify withdrawal state
	assert.True(t, user.IsWithdrawn())
	assert.Equal(t, "deleted_user", *user.DisplayID())
	assert.Equal(t, "削除されたユーザー", *user.DisplayName())
	assert.Nil(t, user.IconURL())
}

func TestUserAuth_Withdraw(t *testing.T) {
	// Create a test user auth
	userAuthID := shared.NewUUID[UserAuth]()
	userID := shared.NewUUID[User]()
	userAuth := NewUserAuth(
		userAuthID,
		userID,
		"google",
		"test_subject",
		true,
		time.Now(),
	)

	// Verify initial state
	assert.False(t, userAuth.IsWithdrawn())
	assert.Nil(t, userAuth.WithdrawalDate())
	assert.True(t, userAuth.CanReregister())

	// Perform withdrawal
	ctx := context.Background()
	userAuth.Withdraw(ctx)

	// Verify withdrawal state
	assert.True(t, userAuth.IsWithdrawn())
	assert.NotNil(t, userAuth.WithdrawalDate())
	assert.False(t, userAuth.CanReregister()) // Should not be able to re-register immediately
}

func TestUserAuth_CanReregister(t *testing.T) {
	userAuthID := shared.NewUUID[UserAuth]()
	userID := shared.NewUUID[User]()

	// Test with withdrawal date 31 days ago (should allow re-registration)
	thirtyOneDaysAgo := time.Now().AddDate(0, 0, -31)
	userAuth := NewUserAuthWithWithdrawal(
		userAuthID,
		userID,
		"google",
		"test_subject",
		true,
		time.Now().AddDate(0, 0, -60),
		&thirtyOneDaysAgo,
	)

	assert.True(t, userAuth.CanReregister())

	// Test with withdrawal date 29 days ago (should not allow re-registration)
	twentyNineDaysAgo := time.Now().AddDate(0, 0, -29)
	userAuth2 := NewUserAuthWithWithdrawal(
		userAuthID,
		userID,
		"google",
		"test_subject",
		true,
		time.Now().AddDate(0, 0, -60),
		&twentyNineDaysAgo,
	)

	assert.False(t, userAuth2.CanReregister())
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}