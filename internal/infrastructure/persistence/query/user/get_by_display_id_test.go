package user_query

import (
	"context"
	"testing"

	"github.com/neko-dream/server/internal/application/query/dto"
	user_query "github.com/neko-dream/server/internal/application/query/user"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock for user.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user user.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, userID shared.UUID[user.User]) (*user.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) FindBySubject(ctx context.Context, subject user.UserSubject) (*user.User, error) {
	args := m.Called(ctx, subject)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) FindByDisplayID(ctx context.Context, displayID string) (*user.User, error) {
	args := m.Called(ctx, displayID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user user.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) ChangeSubject(ctx context.Context, userID shared.UUID[user.User], newSubject string) error {
	args := m.Called(ctx, userID, newSubject)
	return args.Error(0)
}

func TestGetByDisplayIDHandler_Execute(t *testing.T) {
	tests := []struct {
		name           string
		input          user_query.GetByDisplayIDInput
		mockSetup      func(*MockUserRepository)
		expectedResult *user_query.GetByDisplayIDOutput
		expectedError  bool
	}{
		{
			name: "正常系: ユーザーが見つかる",
			input: user_query.GetByDisplayIDInput{
				DisplayID: "testuser123",
			},
			mockSetup: func(m *MockUserRepository) {
				userID := shared.NewUUID[user.User]()
				testUser := user.NewUser(
					userID,
					lo.ToPtr("testuser123"),
					lo.ToPtr("Test User"),
					"subject123",
					shared.AuthProviderName("google"),
					nil,
				)
				m.On("FindByDisplayID", mock.Anything, "testuser123").Return(&testUser, nil)
			},
			expectedResult: &user_query.GetByDisplayIDOutput{
				User: &dto.User{
					DisplayID:   "testuser123",
					DisplayName: "Test User",
				},
			},
			expectedError: false,
		},
		{
			name: "正常系: ユーザーが見つからない",
			input: user_query.GetByDisplayIDInput{
				DisplayID: "notfound",
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("FindByDisplayID", mock.Anything, "notfound").Return(nil, nil)
			},
			expectedResult: &user_query.GetByDisplayIDOutput{
				User: nil,
			},
			expectedError: false,
		},
		{
			name: "異常系: リポジトリエラー",
			input: user_query.GetByDisplayIDInput{
				DisplayID: "error",
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("FindByDisplayID", mock.Anything, "error").Return(nil, assert.AnError)
			},
			expectedResult: nil,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUserRepo := new(MockUserRepository)
			tt.mockSetup(mockUserRepo)

			handler := &GetByDisplayIDHandler{
				DBManager:      &db.DBManager{},
				userRepository: mockUserRepo,
			}

			// Execute
			result, err := handler.Execute(context.Background(), tt.input)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.expectedResult.User == nil {
					assert.Nil(t, result.User)
				} else {
					assert.NotNil(t, result.User)
					assert.Equal(t, tt.expectedResult.User.DisplayID, result.User.DisplayID)
					assert.Equal(t, tt.expectedResult.User.DisplayName, result.User.DisplayName)
				}
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}
