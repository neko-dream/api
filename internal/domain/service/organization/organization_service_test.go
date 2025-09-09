package organization_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/neko-dream/server/internal/domain/messages"
	mock_image_model "github.com/neko-dream/server/internal/domain/model/mock/image"
	mock_organization_model "github.com/neko-dream/server/internal/domain/model/mock/organization"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	organization_service "github.com/neko-dream/server/internal/domain/service/organization"
	mock_organization "github.com/neko-dream/server/internal/domain/service/organization/mock/organization"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestOrganizationService_CreateOrganization(t *testing.T) {
	ctx := context.Background()
	ownerID := shared.NewUUID[user.User]()

	tests := []struct {
		name        string
		orgName     string
		code        string
		orgType     organization.OrganizationType
		env         config.ENV
		setupMocks  func(*mock_organization_model.MockOrganizationRepository, *mock_organization_model.MockOrganizationUserRepository, *mock_organization.MockOrganizationMemberManager)
		expectError error
	}{
		{
			name:    "ローカル環境で正常に組織を作成できる",
			orgName: "Test Organization",
			code:    "TEST123",
			orgType: organization.OrganizationTypeNormal,
			env:     config.LOCAL,
			setupMocks: func(orgRepo *mock_organization_model.MockOrganizationRepository, orgUserRepo *mock_organization_model.MockOrganizationUserRepository, memberManager *mock_organization.MockOrganizationMemberManager) {
				orgRepo.EXPECT().FindByName(gomock.Any(), "Test Organization").Return(nil, sql.ErrNoRows)
				orgRepo.EXPECT().FindByCode(gomock.Any(), "TEST123").Return(nil, sql.ErrNoRows)
				orgRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				orgUserRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectError: nil,
		},
		{
			name:    "開発環境で正常に組織を作成できる",
			orgName: "Dev Organization",
			code:    "DEV001",
			orgType: organization.OrganizationTypeNormal,
			env:     config.DEV,
			setupMocks: func(orgRepo *mock_organization_model.MockOrganizationRepository, orgUserRepo *mock_organization_model.MockOrganizationUserRepository, memberManager *mock_organization.MockOrganizationMemberManager) {
				orgRepo.EXPECT().FindByName(gomock.Any(), "Dev Organization").Return(nil, sql.ErrNoRows)
				orgRepo.EXPECT().FindByCode(gomock.Any(), "DEV001").Return(nil, sql.ErrNoRows)
				orgRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				orgUserRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectError: nil,
		},
		{
			name:    "本番環境でSuperAdmin権限を持つユーザーが組織を作成できる",
			orgName: "Prod Organization",
			code:    "PROD001",
			orgType: organization.OrganizationTypeGovernment,
			env:     config.PROD,
			setupMocks: func(orgRepo *mock_organization_model.MockOrganizationRepository, orgUserRepo *mock_organization_model.MockOrganizationUserRepository, memberManager *mock_organization.MockOrganizationMemberManager) {
				memberManager.EXPECT().IsSuperAdmin(gomock.Any(), ownerID).Return(true, nil)
				orgRepo.EXPECT().FindByName(gomock.Any(), "Prod Organization").Return(nil, sql.ErrNoRows)
				orgRepo.EXPECT().FindByCode(gomock.Any(), "PROD001").Return(nil, sql.ErrNoRows)
				orgRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				orgUserRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectError: nil,
		},
		{
			name:    "本番環境で一般ユーザーは組織を作成できない",
			orgName: "Prod Organization",
			code:    "PROD002",
			orgType: organization.OrganizationTypeNormal,
			env:     config.PROD,
			setupMocks: func(orgRepo *mock_organization_model.MockOrganizationRepository, orgUserRepo *mock_organization_model.MockOrganizationUserRepository, memberManager *mock_organization.MockOrganizationMemberManager) {
				memberManager.EXPECT().IsSuperAdmin(gomock.Any(), ownerID).Return(false, nil)
			},
			expectError: messages.OrganizationForbidden,
		},
		{
			name:    "既に存在する組織名では作成に失敗する",
			orgName: "Duplicate Organization",
			code:    "DUP001",
			orgType: organization.OrganizationTypeNormal,
			env:     config.LOCAL,
			setupMocks: func(orgRepo *mock_organization_model.MockOrganizationRepository, orgUserRepo *mock_organization_model.MockOrganizationUserRepository, memberManager *mock_organization.MockOrganizationMemberManager) {
				existingOrg := &organization.Organization{
					OrganizationID: shared.NewUUID[organization.Organization](),
					Name:           "Duplicate Organization",
				}
				orgRepo.EXPECT().FindByName(gomock.Any(), "Duplicate Organization").Return(existingOrg, nil)
			},
			expectError: messages.OrganizationAlreadyExists,
		},
		{
			name:    "短すぎるコードでは作成に失敗する",
			orgName: "Invalid Code Organization",
			code:    "123", // Too short
			orgType: organization.OrganizationTypeNormal,
			env:     config.LOCAL,
			setupMocks: func(orgRepo *mock_organization_model.MockOrganizationRepository, orgUserRepo *mock_organization_model.MockOrganizationUserRepository, memberManager *mock_organization.MockOrganizationMemberManager) {
				orgRepo.EXPECT().FindByName(gomock.Any(), "Invalid Code Organization").Return(nil, sql.ErrNoRows)
			},
			expectError: messages.OrganizationCodeTooShort,
		},
		{
			name:    "既に使用されているコードでは作成に失敗する",
			orgName: "Duplicate Code Organization",
			code:    "EXIST001",
			orgType: organization.OrganizationTypeNormal,
			env:     config.LOCAL,
			setupMocks: func(orgRepo *mock_organization_model.MockOrganizationRepository, orgUserRepo *mock_organization_model.MockOrganizationUserRepository, memberManager *mock_organization.MockOrganizationMemberManager) {
				orgRepo.EXPECT().FindByName(gomock.Any(), "Duplicate Code Organization").Return(nil, sql.ErrNoRows)
				existingOrg := &organization.Organization{
					OrganizationID: shared.NewUUID[organization.Organization](),
					Code:           "EXIST001",
				}
				orgRepo.EXPECT().FindByCode(gomock.Any(), "EXIST001").Return(existingOrg, nil)
			},
			expectError: messages.OrganizationCodeAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockOrgRepo := mock_organization_model.NewMockOrganizationRepository(ctrl)
			mockOrgUserRepo := mock_organization_model.NewMockOrganizationUserRepository(ctrl)
			mockMemberManager := mock_organization.NewMockOrganizationMemberManager(ctrl)
			mockImageStorage := mock_image_model.NewMockImageStorage(ctrl)

			tt.setupMocks(mockOrgRepo, mockOrgUserRepo, mockMemberManager)

			cfg := &config.Config{
				Env: tt.env,
			}

			service := organization_service.NewOrganizationService(mockOrgRepo, mockOrgUserRepo, mockMemberManager, mockImageStorage, cfg)

			result, err := service.CreateOrganization(ctx, tt.orgName, tt.code, nil, tt.orgType, ownerID)

			if tt.expectError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.orgName, result.Name)
				assert.Equal(t, tt.code, result.Code)
				// iconがnilなので、iconURLもnilになる
				assert.Nil(t, result.IconURL)
				assert.Equal(t, tt.orgType, result.OrganizationType)
				assert.Equal(t, ownerID, result.OwnerID)
			}
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
		setupMocks  func(*mock_organization_model.MockOrganizationRepository, *mock_organization_model.MockOrganizationUserRepository)
		expectOrgs  int
		expectError error
	}{
		{
			name: "ユーザーが複数の組織に所属している場合、全てを取得できる",
			setupMocks: func(orgRepo *mock_organization_model.MockOrganizationRepository, orgUserRepo *mock_organization_model.MockOrganizationUserRepository) {
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
				orgUserRepo.EXPECT().FindByUserID(gomock.Any(), userID).Return(orgUsers, nil)

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
				orgRepo.EXPECT().FindByIDs(gomock.Any(), []shared.UUID[organization.Organization]{orgID1, orgID2}).Return(orgs, nil)
			},
			expectOrgs:  2,
			expectError: nil,
		},
		{
			name: "ユーザーがどの組織にも所属していない場合、空のリストを返す",
			setupMocks: func(orgRepo *mock_organization_model.MockOrganizationRepository, orgUserRepo *mock_organization_model.MockOrganizationUserRepository) {
				orgUserRepo.EXPECT().FindByUserID(gomock.Any(), userID).Return([]*organization.OrganizationUser{}, nil)
				orgRepo.EXPECT().FindByIDs(gomock.Any(), []shared.UUID[organization.Organization]{}).Return([]*organization.Organization{}, nil)
			},
			expectOrgs:  0,
			expectError: nil,
		},
		{
			name: "データベースエラーが発生した場合、エラーを返す",
			setupMocks: func(orgRepo *mock_organization_model.MockOrganizationRepository, orgUserRepo *mock_organization_model.MockOrganizationUserRepository) {
				orgUserRepo.EXPECT().FindByUserID(gomock.Any(), userID).Return(nil, errors.New("database error"))
			},
			expectOrgs:  0,
			expectError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockOrgRepo := mock_organization_model.NewMockOrganizationRepository(ctrl)
			mockOrgUserRepo := mock_organization_model.NewMockOrganizationUserRepository(ctrl)
			mockMemberManager := mock_organization.NewMockOrganizationMemberManager(ctrl)
			mockImageStorage := mock_image_model.NewMockImageStorage(ctrl)

			cfg := &config.Config{}

			tt.setupMocks(mockOrgRepo, mockOrgUserRepo)

			service := organization_service.NewOrganizationService(mockOrgRepo, mockOrgUserRepo, mockMemberManager, mockImageStorage, cfg)

			result, err := service.GetUserOrganizations(ctx, userID)

			if tt.expectError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectOrgs)
			}
		})
	}
}

func TestOrganizationService_ResolveOrganizationIDFromCode(t *testing.T) {
	ctx := context.Background()
	orgID := shared.NewUUID[organization.Organization]()

	tests := []struct {
		name        string
		code        *string
		setupMocks  func(*mock_organization_model.MockOrganizationRepository)
		expectID    bool
		expectError error
	}{
		{
			name: "有効なコードで成功",
			code: lo.ToPtr("TEST123"),
			setupMocks: func(orgRepo *mock_organization_model.MockOrganizationRepository) {
				org := &organization.Organization{
					OrganizationID: orgID,
					Code:           "TEST123",
				}
				orgRepo.EXPECT().FindByCode(gomock.Any(), "TEST123").Return(org, nil)
			},
			expectID:    true,
			expectError: nil,
		},
		{
			name:        "無効なコードの場合nilを返す",
			code:        lo.ToPtr("123"), // Too short
			expectID:    false,
			expectError: nil,
		},
		{
			name:        "空のコードの場合nilを返す",
			code:        lo.ToPtr(""),
			expectID:    false,
			expectError: nil,
		},
		{
			name:        "nilのコードの場合nilを返す",
			code:        nil,
			expectID:    false,
			expectError: nil,
		},
		{
			name: "組織が見つからない場合nilを返す",
			code: lo.ToPtr("NOTFOUND"),
			setupMocks: func(orgRepo *mock_organization_model.MockOrganizationRepository) {
				orgRepo.EXPECT().FindByCode(gomock.Any(), "NOTFOUND").Return(nil, sql.ErrNoRows)
			},
			expectID:    false,
			expectError: nil,
		},
		{
			name: "データベースエラーの場合エラーを返す",
			code: lo.ToPtr("NOTFOUND"),
			setupMocks: func(orgRepo *mock_organization_model.MockOrganizationRepository) {
				orgRepo.EXPECT().FindByCode(gomock.Any(), "NOTFOUND").Return(nil, errors.New("database error"))
			},
			expectID:    false,
			expectError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockOrgRepo := mock_organization_model.NewMockOrganizationRepository(ctrl)
			mockOrgUserRepo := mock_organization_model.NewMockOrganizationUserRepository(ctrl)
			mockMemberManager := mock_organization.NewMockOrganizationMemberManager(ctrl)
			mockImageStorage := mock_image_model.NewMockImageStorage(ctrl)

			cfg := &config.Config{}

			if tt.setupMocks != nil {
				tt.setupMocks(mockOrgRepo)
			}

			service := organization_service.NewOrganizationService(mockOrgRepo, mockOrgUserRepo, mockMemberManager, mockImageStorage, cfg)

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
		})
	}
}
