package service

import (
	"context"
	"testing"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTalkSessionRepository struct {
	mock.Mock
	talksession.TalkSessionRepository
}

func (m *mockTalkSessionRepository) FindByID(ctx context.Context, id shared.UUID[talksession.TalkSession]) (*talksession.TalkSession, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*talksession.TalkSession), args.Error(1)
}

type mockUserRepository struct {
	mock.Mock
	user.UserRepository
}

func (m *mockUserRepository) FindByID(ctx context.Context, id shared.UUID[user.User]) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func TestCanUserJoin(t *testing.T) {
	ctx := context.Background()

	t.Run("制限なしの場合は参加可能", func(t *testing.T) {
		// Arrange
		mockTS := &mockTalkSessionRepository{}
		mockUser := &mockUserRepository{}
		talkSessionID := shared.NewUUID[talksession.TalkSession]()
		userID := shared.NewUUID[user.User]()
		svc := NewTalkSessionAccessControl(mockTS, mockUser)

		mockTS.On("FindByID", mock.Anything, talkSessionID).Return(&talksession.TalkSession{}, nil)
		mockUser.On("FindByID", mock.Anything, userID).Return(&user.User{}, nil)

		// Act
		result, err := svc.CanUserJoin(ctx, talkSessionID, &userID)

		// Assert
		assert.NoError(t, err)
		assert.True(t, result)
		mockTS.AssertExpectations(t)
		mockUser.AssertExpectations(t)
	})

	t.Run("トークセッションが存在しない場合はエラー", func(t *testing.T) {
		// Arrange
		mockTS := &mockTalkSessionRepository{}
		mockUser := &mockUserRepository{}
		talkSessionID := shared.NewUUID[talksession.TalkSession]()
		userID := shared.NewUUID[user.User]()
		svc := NewTalkSessionAccessControl(mockTS, mockUser)

		mockTS.On("FindByID", mock.Anything, talkSessionID).Return(nil, messages.TalkSessionNotFound)

		// Act
		result, err := svc.CanUserJoin(ctx, talkSessionID, &userID)

		// Assert
		assert.Error(t, err)
		assert.False(t, result)
		assert.Equal(t, messages.TalkSessionNotFound, err)
		mockTS.AssertExpectations(t)
	})

	t.Run("制限あり満たす場合は参加可能", func(t *testing.T) {
		// Arrange
		mockTS := &mockTalkSessionRepository{}
		mockUser := &mockUserRepository{}
		talkSessionID := shared.NewUUID[talksession.TalkSession]()
		userID := shared.NewUUID[user.User]()
		svc := NewTalkSessionAccessControl(mockTS, mockUser)

		ts := &talksession.TalkSession{}
		demographics := user.NewUserDemographic(
			ctx, shared.NewUUID[user.UserDemographic](),
			lo.ToPtr(2020),
			nil, nil, nil, nil, nil,
		)
		u := user.NewUser(
			userID, lo.ToPtr("u"), lo.ToPtr("u"), "u", auth.AuthProviderName("u"), nil,
		)
		u.SetDemographics(demographics)

		if err := ts.UpdateRestrictions(ctx, []string{string(talksession.DemographicsBirth)}); err != nil {
			t.Fatal("Failed to update restrictions:", err)
		}

		mockTS.On("FindByID", mock.Anything, talkSessionID).Return(ts, nil)
		mockUser.On("FindByID", mock.Anything, userID).Return(lo.ToPtr(u), nil)

		// Act
		result, err := svc.CanUserJoin(ctx, talkSessionID, &userID)

		// Assert
		assert.NoError(t, err)
		assert.True(t, result)
		mockTS.AssertExpectations(t)
		mockUser.AssertExpectations(t)
	})

	t.Run("制限あり満たさない場合は参加不可", func(t *testing.T) {
		// Arrange
		mockTS := &mockTalkSessionRepository{}
		mockUser := &mockUserRepository{}
		talkSessionID := shared.NewUUID[talksession.TalkSession]()
		userID := shared.NewUUID[user.User]()
		svc := NewTalkSessionAccessControl(mockTS, mockUser)

		ts := &talksession.TalkSession{}
		demographics := user.NewUserDemographic(
			ctx, shared.NewUUID[user.UserDemographic](),
			nil,
			nil, nil, nil, nil, nil,
		)
		u := user.NewUser(
			userID, lo.ToPtr("u"), lo.ToPtr("u"), "u", auth.AuthProviderName("u"), nil,
		)
		u.SetDemographics(demographics)

		if err := ts.UpdateRestrictions(ctx, []string{string(talksession.DemographicsBirth)}); err != nil {
			t.Fatal("Failed to update restrictions:", err)
		}

		mockTS.On("FindByID", mock.Anything, talkSessionID).Return(ts, nil)
		mockUser.On("FindByID", mock.Anything, userID).Return(lo.ToPtr(u), nil)

		// Act
		result, err := svc.CanUserJoin(ctx, talkSessionID, &userID)

		// Assert
		assert.Error(t, err)
		assert.False(t, result)
		assert.IsType(t, &messages.APIError{}, err)
		mockTS.AssertExpectations(t)
		mockUser.AssertExpectations(t)
	})
}
