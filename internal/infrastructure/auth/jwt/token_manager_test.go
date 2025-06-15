package jwt_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/auth/jwt"
	"github.com/neko-dream/server/internal/infrastructure/di"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenManager_SecretValidation(t *testing.T) {
	t.Parallel()

	cont := di.BuildContainer()
	dbm := di.Invoke[*db.DBManager](cont)
	sessRepo := di.Invoke[session.SessionRepository](cont)
	orgUserRepo := di.Invoke[organization.OrganizationUserRepository](cont)
	orgRepo := di.Invoke[organization.OrganizationRepository](cont)

	tests := []struct {
		name         string
		firstSecret  string
		secondSecret string
		success      bool
		ctx          context.Context
	}{
		{
			name:         "matching_secrets",
			firstSecret:  "secret-key-123",
			secondSecret: "secret-key-123",
			success:      true,
			ctx:          context.Background(),
		},
		{
			name:         "different_secrets",
			firstSecret:  "secret-key-123",
			secondSecret: "different-secret-456",
			success:      false,
			ctx:          context.Background(),
		},
		{
			name:         "empty_first_secret",
			firstSecret:  "",
			secondSecret: "secret-key-123",
			success:      false,
			ctx:          context.Background(),
		},
		{
			name:         "empty_second_secret",
			firstSecret:  "secret-key-123",
			secondSecret: "",
			success:      false,
			ctx:          context.Background(),
		},
		{
			name:         "both_empty_secrets",
			firstSecret:  "",
			secondSecret: "",
			success:      false,
			ctx:          context.Background(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testUser := user.NewUser(
				shared.NewUUID[user.User](),
				lo.ToPtr("Secret"),
				lo.ToPtr("Test"),
				"secret.test@example.com",
				auth.ProviderGoogle,
				nil,
			)
			sessionID := shared.NewUUID[session.Session]()

			// Generate token with first secret
			tokenManager1 := jwt.NewTokenManagerWithSecret(tt.firstSecret, dbm, sessRepo, orgUserRepo, orgRepo)
			token, err := tokenManager1.Generate(tt.ctx, testUser, sessionID)

			if tt.firstSecret == "" {
				assert.Error(t, err, "Empty secret should cause error during generation")
				return
			}

			require.NoError(t, err, "Token generation should succeed with valid secret")
			assert.NotEmpty(t, token, "Generated token should not be empty")

			// Try to parse with second secret
			tokenManager2 := jwt.NewTokenManagerWithSecret(tt.secondSecret, dbm, sessRepo, orgUserRepo, orgRepo)
			parsedToken, err := tokenManager2.Parse(tt.ctx, token)

			if tt.success {
				assert.NoError(t, err, "Token parsing should succeed with matching secret")
				assert.NotNil(t, parsedToken, "Parsed token should not be nil")
			} else {
				assert.Error(t, err, "Token parsing should fail with mismatched or empty secret")
				assert.Nil(t, parsedToken, "Parsed token should be nil on error")
			}
		})
	}
}

func TestTokenManager_Generate(t *testing.T) {
	t.Parallel()

	cont := di.BuildContainer()
	dbm := di.Invoke[*db.DBManager](cont)
	sessRepo := di.Invoke[session.SessionRepository](cont)
	orgUserRepo := di.Invoke[organization.OrganizationUserRepository](cont)
	orgRepo := di.Invoke[organization.OrganizationRepository](cont)

	tests := []struct {
		name          string
		secret        string
		user          *user.User
		sessionID     shared.UUID[session.Session]
		expectedError bool
		errorContains string
	}{
		{
			name:   "valid_user_and_session",
			secret: "test-secret-key",
			user: user.NewUser(
				shared.NewUUID[user.User](),
				lo.ToPtr("John"),
				lo.ToPtr("Doe"),
				"john.doe@example.com",
				auth.ProviderGoogle,
				nil,
			),
			sessionID:     shared.NewUUID[session.Session](),
			expectedError: false,
		},
		{
			name:          "empty_secret",
			secret:        "",
			user: user.NewUser(
				shared.NewUUID[user.User](),
				lo.ToPtr("Jane"),
				lo.ToPtr("Smith"),
				"jane.smith@example.com",
				auth.ProviderGoogle,
				nil,
			),
			sessionID:     shared.NewUUID[session.Session](),
			expectedError: true,
		},
		{
			name:   "user_with_minimal_fields",
			secret: "test-secret-key",
			user: user.NewUser(
				shared.NewUUID[user.User](),
				nil,
				nil,
				"minimal@example.com",
				auth.ProviderGoogle,
				nil,
			),
			sessionID:     shared.NewUUID[session.Session](),
			expectedError: false,
		},
		{
			name:   "password_provider_user",
			secret: "test-secret-key",
			user: user.NewUser(
				shared.NewUUID[user.User](),
				lo.ToPtr("Password"),
				lo.ToPtr("User"),
				"password@example.com",
				auth.ProviderPassword,
				nil,
			),
			sessionID:     shared.NewUUID[session.Session](),
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenManager := jwt.NewTokenManagerWithSecret(tt.secret, dbm, sessRepo, orgUserRepo, orgRepo)

			token, err := tokenManager.Generate(context.Background(), tt.user, tt.sessionID)

			if tt.expectedError {
				assert.Error(t, err, "Expected error but got none")
				assert.Empty(t, token, "Token should be empty on error")
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err, "Unexpected error: %v", err)
				assert.NotEmpty(t, token, "Token should not be empty on success")
			}
		})
	}
}

func TestTokenManager_Parse(t *testing.T) {
	t.Parallel()

	cont := di.BuildContainer()
	dbm := di.Invoke[*db.DBManager](cont)
	sessRepo := di.Invoke[session.SessionRepository](cont)
	orgUserRepo := di.Invoke[organization.OrganizationUserRepository](cont)
	orgRepo := di.Invoke[organization.OrganizationRepository](cont)

	secret := "test-secret-key"
	tokenManager := jwt.NewTokenManagerWithSecret(secret, dbm, sessRepo, orgUserRepo, orgRepo)

	// Create a valid user and session for testing
	testUser := user.NewUser(
		shared.NewUUID[user.User](),
		lo.ToPtr("Test"),
		lo.ToPtr("User"),
		"test@example.com",
		auth.ProviderGoogle,
		nil,
	)
	testSessionID := shared.NewUUID[session.Session]()

	// Generate a valid token for testing
	validToken, err := tokenManager.Generate(context.Background(), testUser, testSessionID)
	require.NoError(t, err, "Failed to generate test token")

	tests := []struct {
		name          string
		token         string
		secret        string
		expectedError bool
		errorContains string
	}{
		{
			name:          "valid_token",
			token:         validToken,
			secret:        secret,
			expectedError: false,
		},
		{
			name:          "invalid_token_format",
			token:         "invalid.token.format",
			secret:        secret,
			expectedError: true,
		},
		{
			name:          "empty_token",
			token:         "",
			secret:        secret,
			expectedError: true,
		},
		{
			name:          "wrong_secret",
			token:         validToken,
			secret:        "wrong-secret",
			expectedError: true,
		},
		{
			name:          "malformed_jwt",
			token:         "not.a.jwt",
			secret:        secret,
			expectedError: true,
		},
		{
			name:          "missing_signature",
			token:         "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
			secret:        secret,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tm session.TokenManager
			if tt.secret != secret {
				tm = jwt.NewTokenManagerWithSecret(tt.secret, dbm, sessRepo, orgUserRepo, orgRepo)
			} else {
				tm = tokenManager
			}

			parsedToken, err := tm.Parse(context.Background(), tt.token)

			if tt.expectedError {
				assert.Error(t, err, "Expected error but got none")
				assert.Nil(t, parsedToken, "Parsed token should be nil on error")
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err, "Unexpected error: %v", err)
				assert.NotNil(t, parsedToken, "Parsed token should not be nil on success")
			}
		})
	}
}

func TestTokenManager_DifferentProviders(t *testing.T) {
	t.Parallel()

	cont := di.BuildContainer()
	dbm := di.Invoke[*db.DBManager](cont)
	sessRepo := di.Invoke[session.SessionRepository](cont)
	orgUserRepo := di.Invoke[organization.OrganizationUserRepository](cont)
	orgRepo := di.Invoke[organization.OrganizationRepository](cont)

	tokenManager := jwt.NewTokenManagerWithSecret("provider-test-secret", dbm, sessRepo, orgUserRepo, orgRepo)

	providers := []auth.Provider{
		auth.ProviderGoogle,
		auth.ProviderGitHub,
		auth.ProviderPassword,
	}

	for _, provider := range providers {
		t.Run(string(provider), func(t *testing.T) {
			testUser := user.NewUser(
				shared.NewUUID[user.User](),
				lo.ToPtr("Provider"),
				lo.ToPtr("Test"),
				"provider.test@example.com",
				provider,
				nil,
			)
			sessionID := shared.NewUUID[session.Session]()

			// Generate token
			token, err := tokenManager.Generate(context.Background(), testUser, sessionID)
			require.NoError(t, err, "Failed to generate token for provider %v", provider)
			assert.NotEmpty(t, token, "Token should not be empty")

			// Parse token
			parsedToken, err := tokenManager.Parse(context.Background(), token)
			require.NoError(t, err, "Failed to parse token for provider %v", provider)
			assert.NotNil(t, parsedToken, "Parsed token should not be nil")
		})
	}
}

func TestTokenManager_ConcurrentOperations(t *testing.T) {
	t.Parallel()

	cont := di.BuildContainer()
	dbm := di.Invoke[*db.DBManager](cont)
	sessRepo := di.Invoke[session.SessionRepository](cont)
	orgUserRepo := di.Invoke[organization.OrganizationUserRepository](cont)
	orgRepo := di.Invoke[organization.OrganizationRepository](cont)

	tokenManager := jwt.NewTokenManagerWithSecret("concurrent-test-secret", dbm, sessRepo, orgUserRepo, orgRepo)

	numGoroutines := 10
	results := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			testUser := user.NewUser(
				shared.NewUUID[user.User](),
				lo.ToPtr("Concurrent"),
				lo.ToPtr("User"),
				"concurrent@example.com",
				auth.ProviderGoogle,
				nil,
			)
			sessionID := shared.NewUUID[session.Session]()

			// Generate token
			token, err := tokenManager.Generate(context.Background(), testUser, sessionID)
			if err != nil {
				results <- err
				return
			}

			// Parse the same token
			_, err = tokenManager.Parse(context.Background(), token)
			results <- err
		}(i)
	}

	// Collect results
	for i := 0; i < numGoroutines; i++ {
		err := <-results
		assert.NoError(t, err, "Concurrent operation failed")
	}
}

func TestTokenManager_ContextHandling(t *testing.T) {
	t.Parallel()

	cont := di.BuildContainer()
	dbm := di.Invoke[*db.DBManager](cont)
	sessRepo := di.Invoke[session.SessionRepository](cont)
	orgUserRepo := di.Invoke[organization.OrganizationUserRepository](cont)
	orgRepo := di.Invoke[organization.OrganizationRepository](cont)

	tokenManager := jwt.NewTokenManagerWithSecret("context-test-secret", dbm, sessRepo, orgUserRepo, orgRepo)

	testUser := user.NewUser(
		shared.NewUUID[user.User](),
		lo.ToPtr("Context"),
		lo.ToPtr("Test"),
		"context.test@example.com",
		auth.ProviderGoogle,
		nil,
	)
	sessionID := shared.NewUUID[session.Session]()

	t.Run("cancelled_context_generate", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		token, err := tokenManager.Generate(ctx, testUser, sessionID)
		// Depending on implementation, this might or might not error
		// But the test should not hang
		if err != nil {
			assert.Contains(t, err.Error(), "context")
		}
		_ = token
	})

	t.Run("timeout_context_generate", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		// Small delay to ensure timeout
		time.Sleep(2 * time.Millisecond)

		token, err := tokenManager.Generate(ctx, testUser, sessionID)
		// Similar to above, test should not hang
		if err != nil {
			assert.Contains(t, err.Error(), "context")
		}
		_ = token
	})

	t.Run("valid_context_operations", func(t *testing.T) {
		ctx := context.Background()

		token, err := tokenManager.Generate(ctx, testUser, sessionID)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		parsedToken, err := tokenManager.Parse(ctx, token)
		assert.NoError(t, err)
		assert.NotNil(t, parsedToken)
	})
}

func TestTokenManager_EdgeCases(t *testing.T) {
	t.Parallel()

	cont := di.BuildContainer()
	dbm := di.Invoke[*db.DBManager](cont)
	sessRepo := di.Invoke[session.SessionRepository](cont)
	orgUserRepo := di.Invoke[organization.OrganizationUserRepository](cont)
	orgRepo := di.Invoke[organization.OrganizationRepository](cont)

	tokenManager := jwt.NewTokenManagerWithSecret("edge-case-secret", dbm, sessRepo, orgUserRepo, orgRepo)

	t.Run("very_long_secret", func(t *testing.T) {
		longSecret := string(make([]byte, 1000))
		for i := range longSecret {
			longSecret = longSecret[:i] + "a" + longSecret[i+1:]
		}

		tm := jwt.NewTokenManagerWithSecret(longSecret, dbm, sessRepo, orgUserRepo, orgRepo)

		testUser := user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("Long"),
			lo.ToPtr("Secret"),
			"long.secret@example.com",
			auth.ProviderGoogle,
			nil,
		)
		sessionID := shared.NewUUID[session.Session]()

		token, err := tm.Generate(context.Background(), testUser, sessionID)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		parsedToken, err := tm.Parse(context.Background(), token)
		assert.NoError(t, err)
		assert.NotNil(t, parsedToken)
	})

	t.Run("special_characters_in_user_data", func(t *testing.T) {
		testUser := user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("Special!@#$%^&*()"),
			lo.ToPtr("Characters<>?:{}[]"),
			"special.chars@example.com",
			auth.ProviderGoogle,
			nil,
		)
		sessionID := shared.NewUUID[session.Session]()

		token, err := tokenManager.Generate(context.Background(), testUser, sessionID)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		parsedToken, err := tokenManager.Parse(context.Background(), token)
		assert.NoError(t, err)
		assert.NotNil(t, parsedToken)
	})

	t.Run("unicode_characters_in_user_data", func(t *testing.T) {
		testUser := user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("Unicode日本語"),
			lo.ToPtr("Émoji🚀"),
			"unicode@example.com",
			auth.ProviderGoogle,
			nil,
		)
		sessionID := shared.NewUUID[session.Session]()

		token, err := tokenManager.Generate(context.Background(), testUser, sessionID)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		parsedToken, err := tokenManager.Parse(context.Background(), token)
		assert.NoError(t, err)
		assert.NotNil(t, parsedToken)
	})
}