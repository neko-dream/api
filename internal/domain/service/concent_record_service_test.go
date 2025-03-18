package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/consent"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockConsentRecordRepository struct {
	mock.Mock
	consent.ConsentRecordRepository
}

func (m *mockConsentRecordRepository) FindByUserAndVersion(ctx context.Context, userID shared.UUID[user.User], version string) (*consent.ConsentRecord, error) {
	args := m.Called(mock.Anything, userID, version)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*consent.ConsentRecord), args.Error(1)
}

func (m *mockConsentRecordRepository) Create(ctx context.Context, record *consent.ConsentRecord) error {
	args := m.Called(mock.Anything, record)
	return args.Error(0)
}

type mockPolicyRepository struct {
	mock.Mock
	consent.PolicyRepository
}

func (m *mockPolicyRepository) FetchLatestPolicy(ctx context.Context) (*consent.Policy, error) {
	args := m.Called(mock.Anything)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*consent.Policy), args.Error(1)
}

func (m *mockPolicyRepository) FindByVersion(ctx context.Context, version string) (*consent.Policy, error) {
	args := m.Called(mock.Anything, version)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*consent.Policy), args.Error(1)
}

func TestConsentService_IsConsentValid(t *testing.T) {
	ctx := context.Background()
	userID := shared.NewUUID[user.User]()

	t.Run("最新のポリシーに同意済みの場合はtrueを返す", func(t *testing.T) {
		// Arrange
		mockCR := &mockConsentRecordRepository{}
		mockP := &mockPolicyRepository{}
		svc := NewConsentService(mockCR, mockP)

		mockP.On("FetchLatestPolicy", mock.Anything).Return(&consent.Policy{Version: "1.0.0"}, nil)
		mockCR.On("FindByUserAndVersion", mock.Anything, userID, "1.0.0").Return(&consent.ConsentRecord{}, nil)

		// Act
		result, err := svc.IsConsentValid(ctx, userID)

		// Assert
		assert.NoError(t, err)
		assert.True(t, result)
		mockCR.AssertExpectations(t)
		mockP.AssertExpectations(t)
	})

	t.Run("同意していない場合はfalseを返す", func(t *testing.T) {
		// Arrange
		mockCR := &mockConsentRecordRepository{}
		mockP := &mockPolicyRepository{}
		svc := NewConsentService(mockCR, mockP)

		mockP.On("FetchLatestPolicy", ctx).Return(&consent.Policy{Version: "1.0.0"}, nil)
		mockCR.On("FindByUserAndVersion", ctx, userID, "1.0.0").Return(nil, sql.ErrNoRows)

		// Act
		result, err := svc.IsConsentValid(ctx, userID)

		// Assert
		assert.NoError(t, err)
		assert.False(t, result)
		mockCR.AssertExpectations(t)
		mockP.AssertExpectations(t)
	})
}

func TestConsentService_RecordConsent(t *testing.T) {
	ctx := context.Background()
	userID := shared.NewUUID[user.User]()
	version := "1.0.0"
	ipAddress := "127.0.0.1"
	userAgent := "test-agent"

	t.Run("同意を記録できる", func(t *testing.T) {
		// Arrange
		mockCR := &mockConsentRecordRepository{}
		mockP := &mockPolicyRepository{}
		svc := NewConsentService(mockCR, mockP)

		mockP.On("FindByVersion", mock.Anything, version).Return(&consent.Policy{Version: version}, nil)
		mockCR.On("FindByUserAndVersion", mock.Anything, userID, version).Return(nil, sql.ErrNoRows)
		mockCR.On("Create", mock.Anything, mock.AnythingOfType("*consent.ConsentRecord")).Return(nil)

		// Act
		record, err := svc.RecordConsent(ctx, userID, version, ipAddress, userAgent)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, record)
		assert.Equal(t, userID, record.UserID)
		assert.Equal(t, version, record.Version)
		assert.Equal(t, ipAddress, record.IP)
		assert.Equal(t, userAgent, record.UA)
		mockCR.AssertExpectations(t)
		mockP.AssertExpectations(t)
	})

	t.Run("既に同意済みの場合はエラーを返す", func(t *testing.T) {
		// Arrange
		mockCR := &mockConsentRecordRepository{}
		mockP := &mockPolicyRepository{}
		svc := NewConsentService(mockCR, mockP)

		mockP.On("FindByVersion", ctx, version).Return(&consent.Policy{Version: version}, nil)
		mockCR.On("FindByUserAndVersion", ctx, userID, version).Return(&consent.ConsentRecord{}, nil)

		// Act
		record, err := svc.RecordConsent(ctx, userID, version, ipAddress, userAgent)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, record)
		assert.Equal(t, messages.PolicyAlreadyConsented, err)
		mockCR.AssertExpectations(t)
		mockP.AssertExpectations(t)
	})

	t.Run("存在しないバージョンの場合はエラーを返す", func(t *testing.T) {
		// Arrange
		mockCR := &mockConsentRecordRepository{}
		mockP := &mockPolicyRepository{}
		svc := NewConsentService(mockCR, mockP)

		mockP.On("FindByVersion", ctx, version).Return(nil, errors.New("not found"))

		// Act
		record, err := svc.RecordConsent(ctx, userID, version, ipAddress, userAgent)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, record)
		assert.Equal(t, messages.PolicyNotFound, err)
		mockCR.AssertExpectations(t)
		mockP.AssertExpectations(t)
	})
}
