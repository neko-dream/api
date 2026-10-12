package handler

import (
	"context"
	"testing"

	"github.com/neko-dream/api/internal/application/query/dto"
	user_query "github.com/neko-dream/api/internal/application/query/user"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/presentation/oas"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGetByDisplayID is a mock for user_query.GetByDisplayID
type MockGetByDisplayID struct {
	mock.Mock
}

func (m *MockGetByDisplayID) Execute(ctx context.Context, input user_query.GetByDisplayIDInput) (*user_query.GetByDisplayIDOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user_query.GetByDisplayIDOutput), args.Error(1)
}

func TestUserHandler_GetUserByDisplayID(t *testing.T) {
	tests := []struct {
		name          string
		params        oas.GetUserByDisplayIDParams
		mockSetup     func(*MockGetByDisplayID)
		expectedType  interface{}
		expectedError bool
	}{
		{
			name: "正常系: ユーザーが見つかる",
			params: oas.GetUserByDisplayIDParams{
				DisplayID: "testuser123",
			},
			mockSetup: func(m *MockGetByDisplayID) {
				m.On("Execute", mock.Anything, user_query.GetByDisplayIDInput{
					DisplayID: "testuser123",
				}).Return(&user_query.GetByDisplayIDOutput{
					User: &dto.User{
						DisplayID:   "testuser123",
						DisplayName: "Test User",
					},
				}, nil)
			},
			expectedType:  &oas.User{},
			expectedError: false,
		},
		{
			name: "正常系: ユーザーが見つからない",
			params: oas.GetUserByDisplayIDParams{
				DisplayID: "notfound",
			},
			mockSetup: func(m *MockGetByDisplayID) {
				m.On("Execute", mock.Anything, user_query.GetByDisplayIDInput{
					DisplayID: "notfound",
				}).Return(&user_query.GetByDisplayIDOutput{
					User: nil,
				}, nil)
			},
			expectedType:  messages.UserNotFound,
			expectedError: true,
		},
		{
			name: "異常系: クエリ実行エラー",
			params: oas.GetUserByDisplayIDParams{
				DisplayID: "error",
			},
			mockSetup: func(m *MockGetByDisplayID) {
				m.On("Execute", mock.Anything, user_query.GetByDisplayIDInput{
					DisplayID: "error",
				}).Return(nil, assert.AnError)
			},
			expectedType:  &oas.GetUserByDisplayIDInternalServerError{},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockQuery := new(MockGetByDisplayID)
			tt.mockSetup(mockQuery)

			handler := &userHandler{
				getByDisplayID: mockQuery,
			}

			// Execute
			result, err := handler.GetUserByDisplayID(context.Background(), tt.params)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				assert.IsType(t, tt.expectedType, err)
			} else {
				assert.NoError(t, err)
				assert.IsType(t, tt.expectedType, result)

				if userResp, ok := result.(*oas.User); ok {
					assert.Equal(t, "testuser123", userResp.DisplayID)
					assert.Equal(t, "Test User", userResp.DisplayName)
				}
			}

			mockQuery.AssertExpectations(t)
		})
	}
}
