package organization

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories
type mockOrganizationRepository struct {
	mock.Mock
}

func (m *mockOrganizationRepository) FindByID(ctx context.Context, id shared.UUID[organization.Organization]) (*organization.Organization, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*organization.Organization), args.Error(1)
}

func (m *mockOrganizationRepository) FindByIDs(ctx context.Context, ids []shared.UUID[organization.Organization]) ([]*organization.Organization, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*organization.Organization), args.Error(1)
}

func (m *mockOrganizationRepository) FindByName(ctx context.Context, name string) (*organization.Organization, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*organization.Organization), args.Error(1)
}

func (m *mockOrganizationRepository) FindByCode(ctx context.Context, code string) (*organization.Organization, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*organization.Organization), args.Error(1)
}

func (m *mockOrganizationRepository) Create(ctx context.Context, org *organization.Organization) error {
	args := m.Called(ctx, org)
	return args.Error(0)
}

type mockOrganizationUserRepository struct {
	mock.Mock
}

func (m *mockOrganizationUserRepository) FindByOrganizationIDAndUserID(ctx context.Context, orgID shared.UUID[organization.Organization], userID shared.UUID[user.User]) (*organization.OrganizationUser, error) {
	args := m.Called(ctx, orgID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*organization.OrganizationUser), args.Error(1)
}

func (m *mockOrganizationUserRepository) FindByOrganizationID(ctx context.Context, orgID shared.UUID[organization.Organization]) ([]*organization.OrganizationUser, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*organization.OrganizationUser), args.Error(1)
}

func (m *mockOrganizationUserRepository) FindByUserID(ctx context.Context, userID shared.UUID[user.User]) ([]*organization.OrganizationUser, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*organization.OrganizationUser), args.Error(1)
}

func (m *mockOrganizationUserRepository) Create(ctx context.Context, orgUser organization.OrganizationUser) error {
	args := m.Called(ctx, orgUser)
	return args.Error(0)
}

type mockOrganizationMemberManager struct {
	mock.Mock
}

func (m *mockOrganizationMemberManager) IsSuperAdmin(ctx context.Context, userID shared.UUID[user.User]) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

func (m *mockOrganizationMemberManager) AddUser(ctx context.Context, params InviteUserParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *mockOrganizationMemberManager) InviteUser(ctx context.Context, params InviteUserParams) (*organization.OrganizationUser, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*organization.OrganizationUser), args.Error(1)
}

func TestOrganizationService_CreateOrganization(t *testing.T) {
	ctx := context.Background()
	ownerID := shared.NewUUID[user.User]()

	tests := []struct {
		name        string
		orgName     string
		code        string
		orgType     organization.OrganizationType
		env         config.ENV
		setupMocks  func(*mockOrganizationRepository, *mockOrganizationUserRepository, *mockOrganizationMemberManager)
		expectError error
	}{
		{
			name:    "Success in local environment",
			orgName: "Test Organization",
			code:    "TEST123",
			orgType: organization.OrganizationTypeNormal,
			env:     config.LOCAL,
			setupMocks: func(orgRepo *mockOrganizationRepository, orgUserRepo *mockOrganizationUserRepository, memberManager *mockOrganizationMemberManager) {
				orgRepo.On("FindByName", mock.Anything, "Test Organization").Return(nil, sql.ErrNoRows)
				orgRepo.On("FindByCode", mock.Anything, "TEST123").Return(nil, sql.ErrNoRows)
				orgRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				orgUserRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: nil,
		},
		{
			name:    "Success in dev environment",
			orgName: "Dev Organization",
			code:    "DEV001",
			orgType: organization.OrganizationTypeNormal,
			env:     config.DEV,
			setupMocks: func(orgRepo *mockOrganizationRepository, orgUserRepo *mockOrganizationUserRepository, memberManager *mockOrganizationMemberManager) {
				orgRepo.On("FindByName", mock.Anything, "Dev Organization").Return(nil, sql.ErrNoRows)
				orgRepo.On("FindByCode", mock.Anything, "DEV001").Return(nil, sql.ErrNoRows)
				orgRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				orgUserRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: nil,
		},
		{
			name:    "Success with SuperAdmin in production",
			orgName: "Prod Organization",
			code:    "PROD001",
			orgType: organization.OrganizationTypeGovernment,
			env:     config.PROD,
			setupMocks: func(orgRepo *mockOrganizationRepository, orgUserRepo *mockOrganizationUserRepository, memberManager *mockOrganizationMemberManager) {
				memberManager.On("IsSuperAdmin", mock.Anything, ownerID).Return(true, nil)
				orgRepo.On("FindByName", mock.Anything, "Prod Organization").Return(nil, sql.ErrNoRows)
				orgRepo.On("FindByCode", mock.Anything, "PROD001").Return(nil, sql.ErrNoRows)
				orgRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				orgUserRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: nil,
		},
		{
			name:    "Fail without SuperAdmin in production",
			orgName: "Prod Organization",
			code:    "PROD002",
			orgType: organization.OrganizationTypeNormal,
			env:     config.PROD,
			setupMocks: func(orgRepo *mockOrganizationRepository, orgUserRepo *mockOrganizationUserRepository, memberManager *mockOrganizationMemberManager) {
				memberManager.On("IsSuperAdmin", mock.Anything, ownerID).Return(false, nil)
			},
			expectError: messages.OrganizationForbidden,
		},
		{
			name:    "Fail with duplicate name",
			orgName: "Duplicate Organization",
			code:    "DUP001",
			orgType: organization.OrganizationTypeNormal,
			env:     config.LOCAL,
			setupMocks: func(orgRepo *mockOrganizationRepository, orgUserRepo *mockOrganizationUserRepository, memberManager *mockOrganizationMemberManager) {
				existingOrg := &organization.Organization{
					OrganizationID: shared.NewUUID[organization.Organization](),
					Name:           "Duplicate Organization",
				}
				orgRepo.On("FindByName", mock.Anything, "Duplicate Organization").Return(existingOrg, nil)
			},
			expectError: messages.OrganizationAlreadyExists,
		},
		{
			name:    "Fail with invalid code",
			orgName: "Invalid Code Organization",
			code:    "123", // Too short
			orgType: organization.OrganizationTypeNormal,
			env:     config.LOCAL,
			setupMocks: func(orgRepo *mockOrganizationRepository, orgUserRepo *mockOrganizationUserRepository, memberManager *mockOrganizationMemberManager) {
				orgRepo.On("FindByName", mock.Anything, "Invalid Code Organization").Return(nil, sql.ErrNoRows)
			},
			expectError: messages.OrganizationCodeTooShort,
		},
		{
			name:    "Fail with duplicate code",
			orgName: "Duplicate Code Organization",
			code:    "EXIST001",
			orgType: organization.OrganizationTypeNormal,
			env:     config.LOCAL,
			setupMocks: func(orgRepo *mockOrganizationRepository, orgUserRepo *mockOrganizationUserRepository, memberManager *mockOrganizationMemberManager) {
				orgRepo.On("FindByName", mock.Anything, "Duplicate Code Organization").Return(nil, sql.ErrNoRows)
				existingOrg := &organization.Organization{
					OrganizationID: shared.NewUUID[organization.Organization](),
					Code:           "EXIST001",
				}
				orgRepo.On("FindByCode", mock.Anything, "EXIST001").Return(existingOrg, nil)
			},
			expectError: messages.OrganizationCodeAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgRepo := new(mockOrganizationRepository)
			orgUserRepo := new(mockOrganizationUserRepository)
			memberManager := new(mockOrganizationMemberManager)

			tt.setupMocks(orgRepo, orgUserRepo, memberManager)

			cfg := &config.Config{
				Env: tt.env,
			}

			service := NewOrganizationService(orgRepo, orgUserRepo, memberManager, cfg)

			result, err := service.CreateOrganization(ctx, tt.orgName, tt.code, tt.orgType, ownerID)

			if tt.expectError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.orgName, result.Name)
				assert.Equal(t, tt.code, result.Code)
				assert.Equal(t, tt.orgType, result.OrganizationType)
				assert.Equal(t, ownerID, result.OwnerID)
			}

			orgRepo.AssertExpectations(t)
			orgUserRepo.AssertExpectations(t)
			memberManager.AssertExpectations(t)
		})
	}
}

func TestOrganizationService_GetUserOrganizations(t *testing.T) {
	ctx := context.Background()
	userID := shared.NewUUID[user.User]()
	orgID1 := shared.NewUUID[organization.Organization]()
	orgID2 := shared.NewUUID[organization.Organization]()

	tests := []struct {
		name        string
		setupMocks  func(*mockOrganizationRepository, *mockOrganizationUserRepository)
		expectOrgs  int
		expectError error
	}{
		{
			name: "Success with multiple organizations",
			setupMocks: func(orgRepo *mockOrganizationRepository, orgUserRepo *mockOrganizationUserRepository) {
				orgUsers := []*organization.OrganizationUser{
					{
						OrganizationUserID: shared.NewUUID[organization.OrganizationUser](),
						OrganizationID:     orgID1,
						UserID:             userID,
						Role:               organization.OrganizationUserRoleMember,
					},
					{
						OrganizationUserID: shared.NewUUID[organization.OrganizationUser](),
						OrganizationID:     orgID2,
						UserID:             userID,
						Role:               organization.OrganizationUserRoleAdmin,
					},
				}
				orgUserRepo.On("FindByUserID", mock.Anything, userID).Return(orgUsers, nil)

				orgs := []*organization.Organization{
					{
						OrganizationID: orgID1,
						Name:           "Organization 1",
					},
					{
						OrganizationID: orgID2,
						Name:           "Organization 2",
					},
				}
				orgRepo.On("FindByIDs", mock.Anything, []shared.UUID[organization.Organization]{orgID1, orgID2}).Return(orgs, nil)
			},
			expectOrgs:  2,
			expectError: nil,
		},
		{
			name: "Success with no organizations",
			setupMocks: func(orgRepo *mockOrganizationRepository, orgUserRepo *mockOrganizationUserRepository) {
				orgUserRepo.On("FindByUserID", mock.Anything, userID).Return([]*organization.OrganizationUser{}, nil)
				orgRepo.On("FindByIDs", mock.Anything, []shared.UUID[organization.Organization]{}).Return([]*organization.Organization{}, nil)
			},
			expectOrgs:  0,
			expectError: nil,
		},
		{
			name: "Error from organization user repository",
			setupMocks: func(orgRepo *mockOrganizationRepository, orgUserRepo *mockOrganizationUserRepository) {
				orgUserRepo.On("FindByUserID", mock.Anything, userID).Return(nil, errors.New("database error"))
			},
			expectOrgs:  0,
			expectError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgRepo := new(mockOrganizationRepository)
			orgUserRepo := new(mockOrganizationUserRepository)
			memberManager := new(mockOrganizationMemberManager)
			cfg := &config.Config{}

			tt.setupMocks(orgRepo, orgUserRepo)

			service := NewOrganizationService(orgRepo, orgUserRepo, memberManager, cfg)

			result, err := service.GetUserOrganizations(ctx, userID)

			if tt.expectError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectOrgs)
			}

			orgRepo.AssertExpectations(t)
			orgUserRepo.AssertExpectations(t)
		})
	}
}

func TestOrganizationService_ResolveOrganizationIDFromCode(t *testing.T) {
	ctx := context.Background()
	orgID := shared.NewUUID[organization.Organization]()
	validCode := "TEST123"
	invalidCode := "123"
	emptyCode := ""

	tests := []struct {
		name        string
		code        *string
		setupMocks  func(*mockOrganizationRepository)
		expectID    bool
		expectError error
	}{
		{
			name: "Success with valid code",
			code: &validCode,
			setupMocks: func(orgRepo *mockOrganizationRepository) {
				org := &organization.Organization{
					OrganizationID: orgID,
					Code:           validCode,
				}
				orgRepo.On("FindByCode", mock.Anything, validCode).Return(org, nil)
			},
			expectID:    true,
			expectError: nil,
		},
		{
			name: "Return nil for invalid code",
			code: &invalidCode,
			setupMocks: func(orgRepo *mockOrganizationRepository) {
				// No mock calls expected - validation fails first
			},
			expectID:    false,
			expectError: nil,
		},
		{
			name: "Return nil for empty code",
			code: &emptyCode,
			setupMocks: func(orgRepo *mockOrganizationRepository) {
				// No mock calls expected - early return for empty code
			},
			expectID:    false,
			expectError: nil,
		},
		{
			name: "Return nil for nil code",
			code: nil,
			setupMocks: func(orgRepo *mockOrganizationRepository) {
				// No mock calls expected - early return for nil code
			},
			expectID:    false,
			expectError: nil,
		},
		{
			name: "Return nil when organization not found",
			code: &validCode,
			setupMocks: func(orgRepo *mockOrganizationRepository) {
				orgRepo.On("FindByCode", mock.Anything, validCode).Return(nil, sql.ErrNoRows)
			},
			expectID:    false,
			expectError: nil,
		},
		{
			name: "Return error for database error",
			code: &validCode,
			setupMocks: func(orgRepo *mockOrganizationRepository) {
				orgRepo.On("FindByCode", mock.Anything, validCode).Return(nil, errors.New("database error"))
			},
			expectID:    false,
			expectError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgRepo := new(mockOrganizationRepository)
			orgUserRepo := new(mockOrganizationUserRepository)
			memberManager := new(mockOrganizationMemberManager)
			cfg := &config.Config{}

			tt.setupMocks(orgRepo)

			service := NewOrganizationService(orgRepo, orgUserRepo, memberManager, cfg)

			result, err := service.ResolveOrganizationIDFromCode(ctx, tt.code)

			if tt.expectError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				if tt.expectID {
					assert.NotNil(t, result)
				} else {
					assert.Nil(t, result)
				}
			}

			orgRepo.AssertExpectations(t)
		})
	}
}
