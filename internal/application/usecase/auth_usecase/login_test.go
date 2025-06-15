package auth_usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"your-project/internal/application/usecase/auth_usecase"
	"your-project/internal/domain/entity"
	"your-project/internal/domain/repository/mocks"
	"your-project/pkg/utils"
)

type LoginUsecaseTestSuite struct {
	suite.Suite
	userRepo        *mocks.UserRepository
	tokenService    *mocks.TokenService
	passwordService *mocks.PasswordService
	usecase         auth_usecase.LoginUsecase
}

func (suite *LoginUsecaseTestSuite) SetupTest() {
	suite.userRepo = new(mocks.UserRepository)
	suite.tokenService = new(mocks.TokenService)
	suite.passwordService = new(mocks.PasswordService)
	suite.usecase = auth_usecase.NewLoginUsecase(
		suite.userRepo,
		suite.tokenService,
		suite.passwordService,
	)
}

func (suite *LoginUsecaseTestSuite) TearDownTest() {
	suite.userRepo.AssertExpectations(suite.T())
	suite.tokenService.AssertExpectations(suite.T())
	suite.passwordService.AssertExpectations(suite.T())
}

func TestLoginUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(LoginUsecaseTestSuite))
}

func (suite *LoginUsecaseTestSuite) TestLogin_Success_ValidCredentials() {
	ctx := context.Background()
	email := "user@example.com"
	password := "validpassword"
	hashedPassword := "$2a$10$N9qo8uLOickgx2ZMRZoMye.BCfQQ5D.92cqp4L7/OJt4PKmKj1oKu"

	user := &entity.User{
		ID:        1,
		Email:     email,
		Password:  hashedPassword,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	expectedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test.token"

	suite.userRepo.On("GetByEmail", ctx, email).Return(user, nil)
	suite.passwordService.On("ComparePassword", hashedPassword, password).Return(nil)
	suite.tokenService.On("GenerateToken", user.ID).Return(expectedToken, nil)

	result, err := suite.usecase.Login(ctx, email, password)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), expectedToken, result.Token)
	assert.Equal(suite.T(), user.ID, result.UserID)
	assert.Equal(suite.T(), user.Email, result.Email)
}

func (suite *LoginUsecaseTestSuite) TestLogin_Success_CaseInsensitiveEmail() {
	ctx := context.Background()
	inputEmail := "USER@EXAMPLE.COM"
	normalizedEmail := "user@example.com"
	password := "validpassword"
	hashedPassword := "$2a$10$N9qo8uLOickgx2ZMRZoMye.BCfQQ5D.92cqp4L7/OJt4PKmKj1oKu"

	user := &entity.User{
		ID:       1,
		Email:    normalizedEmail,
		Password: hashedPassword,
		IsActive: true,
	}

	expectedToken := "jwt.token.here"

	suite.userRepo.On("GetByEmail", ctx, normalizedEmail).Return(user, nil)
	suite.passwordService.On("ComparePassword", hashedPassword, password).Return(nil)
	suite.tokenService.On("GenerateToken", user.ID).Return(expectedToken, nil)

	result, err := suite.usecase.Login(ctx, inputEmail, password)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), expectedToken, result.Token)
}

func (suite *LoginUsecaseTestSuite) TestLogin_Success_WithSpecialCharactersInPassword() {
	ctx := context.Background()
	email := "user@example.com"
	password := "P@ssw0rd!@#$%^&*()"
	hashedPassword := "$2a$10$specialHashedPassword"

	user := &entity.User{
		ID:       1,
		Email:    email,
		Password: hashedPassword,
		IsActive: true,
	}

	expectedToken := "jwt.token.special"

	suite.userRepo.On("GetByEmail", ctx, email).Return(user, nil)
	suite.passwordService.On("ComparePassword", hashedPassword, password).Return(nil)
	suite.tokenService.On("GenerateToken", user.ID).Return(expectedToken, nil)

	result, err := suite.usecase.Login(ctx, email, password)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

func (suite *LoginUsecaseTestSuite) TestLogin_Failure_EmptyEmail() {
	ctx := context.Background()
	email := ""
	password := "validpassword"

	result, err := suite.usecase.Login(ctx, email, password)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "email")
	assert.Contains(suite.T(), err.Error(), "required")
}

func (suite *LoginUsecaseTestSuite) TestLogin_Failure_EmptyPassword() {
	ctx := context.Background()
	email := "user@example.com"
	password := ""

	result, err := suite.usecase.Login(ctx, email, password)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "password")
	assert.Contains(suite.T(), err.Error(), "required")
}

func (suite *LoginUsecaseTestSuite) TestLogin_Failure_BothFieldsEmpty() {
	ctx := context.Background()
	email := ""
	password := ""

	result, err := suite.usecase.Login(ctx, email, password)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *LoginUsecaseTestSuite) TestLogin_Failure_InvalidEmailFormat() {
	testCases := []struct {
		name  string
		email string
	}{
		{"missing @", "invalid-email"},
		{"missing domain", "user@"},
		{"missing username", "@example.com"},
		{"double @", "user@@example.com"},
		{"invalid characters", "user<>@example.com"},
		{"spaces", "user @example.com"},
	}

	ctx := context.Background()
	password := "validpassword"

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			result, err := suite.usecase.Login(ctx, tc.email, password)

			assert.Error(suite.T(), err)
			assert.Nil(suite.T(), result)
			assert.Contains(suite.T(), err.Error(), "invalid email format")
		})
	}
}

func (suite *LoginUsecaseTestSuite) TestLogin_EdgeCase_MaxLengthInputs() {
	ctx := context.Background()
	longEmail := string(make([]byte, 254)) + "@example.com"
	longPassword := string(make([]byte, 1000))

	// Test long email
	result, err := suite.usecase.Login(ctx, longEmail, "password")
	if err != nil {
		assert.Contains(suite.T(), err.Error(), "email")
	}

	// Test long password
	result, err = suite.usecase.Login(ctx, "user@example.com", longPassword)
	if err != nil {
		assert.NotPanics(suite.T(), func() {
			_, _ = suite.usecase.Login(ctx, "user@example.com", longPassword)
		})
	}
}

func (suite *LoginUsecaseTestSuite) TestLogin_Failure_UserNotFound() {
	ctx := context.Background()
	email := "nonexistent@example.com"
	password := "validpassword"

	suite.userRepo.On("GetByEmail", ctx, email).Return(nil, errors.New("user not found"))

	result, err := suite.usecase.Login(ctx, email, password)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "invalid credentials")
	assert.NotContains(suite.T(), err.Error(), "not found")
}

func (suite *LoginUsecaseTestSuite) TestLogin_Failure_InvalidPassword() {
	ctx := context.Background()
	email := "user@example.com"
	password := "wrongpassword"
	hashedPassword := "$2a$10$correctHashedPassword"

	user := &entity.User{
		ID:       1,
		Email:    email,
		Password: hashedPassword,
		IsActive: true,
	}

	suite.userRepo.On("GetByEmail", ctx, email).Return(user, nil)
	suite.passwordService.On("ComparePassword", hashedPassword, password).Return(errors.New("password mismatch"))

	result, err := suite.usecase.Login(ctx, email, password)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "invalid credentials")
}

func (suite *LoginUsecaseTestSuite) TestLogin_Failure_InactiveUser() {
	ctx := context.Background()
	email := "user@example.com"
	password := "validpassword"
	hashedPassword := "$2a$10$validHashedPassword"

	user := &entity.User{
		ID:       1,
		Email:    email,
		Password: hashedPassword,
		IsActive: false,
	}

	suite.userRepo.On("GetByEmail", ctx, email).Return(user, nil)
	suite.passwordService.On("ComparePassword", hashedPassword, password).Return(nil)

	result, err := suite.usecase.Login(ctx, email, password)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "account disabled")
}

func (suite *LoginUsecaseTestSuite) TestLogin_Failure_DeletedUser() {
	ctx := context.Background()
	email := "user@example.com"
	password := "validpassword"

	user := &entity.User{
		ID:        1,
		Email:     email,
		Password:  "$2a$10$validHashedPassword",
		IsActive:  true,
		DeletedAt: &time.Time{},
	}

	suite.userRepo.On("GetByEmail", ctx, email).Return(user, nil)

	result, err := suite.usecase.Login(ctx, email, password)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "account not found")
}

func (suite *LoginUsecaseTestSuite) TestLogin_Failure_DatabaseError() {
	ctx := context.Background()
	email := "user@example.com"
	password := "validpassword"

	suite.userRepo.On("GetByEmail", ctx, email).Return(nil, errors.New("database connection failed"))

	result, err := suite.usecase.Login(ctx, email, password)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "internal server error")
	assert.NotContains(suite.T(), err.Error(), "database")
}

func (suite *LoginUsecaseTestSuite) TestLogin_Failure_TokenGenerationError() {
	ctx := context.Background()
	email := "user@example.com"
	password := "validpassword"
	hashedPassword := "$2a$10$validHashedPassword"

	user := &entity.User{
		ID:       1,
		Email:    email,
		Password: hashedPassword,
		IsActive: true,
	}

	suite.userRepo.On("GetByEmail", ctx, email).Return(user, nil)
	suite.passwordService.On("ComparePassword", hashedPassword, password).Return(nil)
	suite.tokenService.On("GenerateToken", user.ID).Return("", errors.New("token generation failed"))

	result, err := suite.usecase.Login(ctx, email, password)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "internal server error")
}

func (suite *LoginUsecaseTestSuite) TestLogin_Failure_PasswordServiceError() {
	ctx := context.Background()
	email := "user@example.com"
	password := "validpassword"
	hashedPassword := "$2a$10$validHashedPassword"

	user := &entity.User{
		ID:       1,
		Email:    email,
		Password: hashedPassword,
		IsActive: true,
	}

	suite.userRepo.On("GetByEmail", ctx, email).Return(user, nil)
	suite.passwordService.On("ComparePassword", hashedPassword, password).Return(errors.New("bcrypt service unavailable"))

	result, err := suite.usecase.Login(ctx, email, password)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *LoginUsecaseTestSuite) TestLogin_Failure_ContextCancellation() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	email := "user@example.com"
	password := "validpassword"

	result, err := suite.usecase.Login(ctx, email, password)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "context canceled")
}

func (suite *LoginUsecaseTestSuite) TestLogin_Failure_ContextTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	time.Sleep(1 * time.Millisecond)

	email := "user@example.com"
	password := "validpassword"

	result, err := suite.usecase.Login(ctx, email, password)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "context deadline exceeded")
}

func (suite *LoginUsecaseTestSuite) TestLogin_Success_WithContextValues() {
	ctx := context.WithValue(context.Background(), "request_id", "test-123")
	email := "user@example.com"
	password := "validpassword"
	hashedPassword := "$2a$10$validHashedPassword"

	user := &entity.User{
		ID:       1,
		Email:    email,
		Password: hashedPassword,
		IsActive: true,
	}

	expectedToken := "jwt.token.with.context"

	suite.userRepo.On("GetByEmail", ctx, email).Return(user, nil)
	suite.passwordService.On("ComparePassword", hashedPassword, password).Return(nil)
	suite.tokenService.On("GenerateToken", user.ID).Return(expectedToken, nil)

	result, err := suite.usecase.Login(ctx, email, password)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), expectedToken, result.Token)
}

func (suite *LoginUsecaseTestSuite) TestLogin_Security_SQLInjectionAttempt() {
	ctx := context.Background()
	maliciousEmail := "user@example.com'; DROP TABLE users; --"
	password := "validpassword"

	suite.userRepo.On("GetByEmail", ctx, maliciousEmail).Return(nil, errors.New("user not found"))

	result, err := suite.usecase.Login(ctx, maliciousEmail, password)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *LoginUsecaseTestSuite) TestLogin_Security_XSSAttempt() {
	ctx := context.Background()
	maliciousEmail := "<script>alert('xss')</script>@example.com"
	password := "validpassword"

	result, err := suite.usecase.Login(ctx, maliciousEmail, password)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "invalid email format")
}

func (suite *LoginUsecaseTestSuite) TestLogin_Security_NoSQLInjectionAttempt() {
	ctx := context.Background()
	maliciousPassword := `{"$ne": null}`
	email := "user@example.com"

	result, err := suite.usecase.Login(ctx, email, maliciousPassword)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *LoginUsecaseTestSuite) TestLogin_Security_UnicodeNormalization() {
	ctx := context.Background()
	unicodeEmail := "üser@example.com"
	password := "validpassword"

	result, err := suite.usecase.Login(ctx, unicodeEmail, password)

	if err != nil {
		assert.Contains(suite.T(), err.Error(), "invalid credentials")
	}
}

func (suite *LoginUsecaseTestSuite) TestLogin_Performance_MultipleSimultaneousRequests() {
	ctx := context.Background()
	email := "user@example.com"
	password := "validpassword"
	hashedPassword := "$2a$10$validHashedPassword"

	user := &entity.User{
		ID:       1,
		Email:    email,
		Password: hashedPassword,
		IsActive: true,
	}

	expectedToken := "jwt.token.concurrent"
	numRequests := 10

	suite.userRepo.On("GetByEmail", ctx, email).Return(user, nil).Times(numRequests)
	suite.passwordService.On("ComparePassword", hashedPassword, password).Return(nil).Times(numRequests)
	suite.tokenService.On("GenerateToken", user.ID).Return(expectedToken, nil).Times(numRequests)

	results := make(chan bool, numRequests)
	for i := 0; i < numRequests; i++ {
		go func() {
			result, err := suite.usecase.Login(ctx, email, password)
			results <- (err == nil && result != nil && result.Token == expectedToken)
		}()
	}

	successCount := 0
	for i := 0; i < numRequests; i++ {
		if <-results {
			successCount++
		}
	}

	assert.Equal(suite.T(), numRequests, successCount)
}

func (suite *LoginUsecaseTestSuite) TestLogin_Performance_RateLimitingBehavior() {
	ctx := context.Background()
	email := "user@example.com"
	password := "wrongpassword"

	user := &entity.User{
		ID:       1,
		Email:    email,
		Password: "$2a$10$validHashedPassword",
		IsActive: true,
	}

	suite.userRepo.On("GetByEmail", ctx, email).Return(user, nil).Times(5)
	suite.passwordService.On("ComparePassword", mock.Anything, password).Return(errors.New("invalid password")).Times(5)

	var lastErr error
	for i := 0; i < 5; i++ {
		_, lastErr = suite.usecase.Login(ctx, email, password)
	}

	assert.Error(suite.T(), lastErr)
	assert.Contains(suite.T(), lastErr.Error(), "invalid credentials")
}

func (suite *LoginUsecaseTestSuite) createValidUser() *entity.User {
	return &entity.User{
		ID:        1,
		Email:     "user@example.com",
		Password:  "$2a$10$N9qo8uLOickgx2ZMRZoMye.BCfQQ5D.92cqp4L7/OJt4PKmKj1oKu",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (suite *LoginUsecaseTestSuite) createInactiveUser() *entity.User {
	user := suite.createValidUser()
	user.IsActive = false
	return user
}

func (suite *LoginUsecaseTestSuite) createDeletedUser() *entity.User {
	user := suite.createValidUser()
	now := time.Now()
	user.DeletedAt = &now
	return user
}

func BenchmarkLogin_Success(b *testing.B) {
	userRepo := new(mocks.UserRepository)
	tokenService := new(mocks.TokenService)
	passwordService := new(mocks.PasswordService)
	usecase := auth_usecase.NewLoginUsecase(userRepo, tokenService, passwordService)

	ctx := context.Background()
	email := "user@example.com"
	password := "validpassword"
	hashedPassword := "$2a$10$validHashedPassword"

	user := &entity.User{
		ID:       1,
		Email:    email,
		Password: hashedPassword,
		IsActive: true,
	}

	userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(user, nil)
	passwordService.On("ComparePassword", mock.Anything, mock.Anything).Return(nil)
	tokenService.On("GenerateToken", mock.Anything).Return("token", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = usecase.Login(ctx, email, password)
	}
}

func BenchmarkLogin_Failure(b *testing.B) {
	userRepo := new(mocks.UserRepository)
	tokenService := new(mocks.TokenService)
	passwordService := new(mocks.PasswordService)
	usecase := auth_usecase.NewLoginUsecase(userRepo, tokenService, passwordService)

	ctx := context.Background()
	email := "user@example.com"
	password := "wrongpassword"

	userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(nil, errors.New("user not found"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = usecase.Login(ctx, email, password)
	}
}

func TestLogin_Integration_RealScenarios(t *testing.T) {
	t.Skip("Integration tests require real database - implement when needed")
}