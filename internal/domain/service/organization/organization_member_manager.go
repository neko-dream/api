package organization

import (
	"context"
	"errors"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/auth"
	password_auth "github.com/neko-dream/server/internal/domain/model/auth/password"
	"github.com/neko-dream/server/internal/domain/model/consent"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/email"
	email_template "github.com/neko-dream/server/internal/infrastructure/email/template"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/hash"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type InviteUserParams struct {
	OrganizationID shared.UUID[organization.Organization]
	Role           organization.OrganizationUserRole
	UserID         shared.UUID[user.User]
	Email          string
}

type OrganizationMemberManager interface {
	// ユーザーの発行
	InviteUser(ctx context.Context, params InviteUserParams) (*organization.OrganizationUser, error)

	IsSuperAdmin(ctx context.Context, userID shared.UUID[user.User]) (bool, error)
}

type organizationMemberManager struct {
	organizationRepo     organization.OrganizationRepository
	organizationUserRepo organization.OrganizationUserRepository
	cfg                  *config.Config
	userRep              user.UserRepository
	policyRep            consent.PolicyRepository
	consentService       consent.ConsentService
	passwordAuthManager  password_auth.PasswordAuthManager
	emailSender          email.EmailSender
	*db.DBManager
}

func NewOrganizationMemberManager(
	organizationRepo organization.OrganizationRepository,
	organizationUserRepo organization.OrganizationUserRepository,
	cfg *config.Config,
	userRep user.UserRepository,
	policyRep consent.PolicyRepository,
	consentService consent.ConsentService,
	passwordAuthManager password_auth.PasswordAuthManager,
	emailSender email.EmailSender,
	dbManager *db.DBManager,
) OrganizationMemberManager {
	return &organizationMemberManager{
		organizationRepo:     organizationRepo,
		organizationUserRepo: organizationUserRepo,
		cfg:                  cfg,
		userRep:              userRep,
		policyRep:            policyRep,
		consentService:       consentService,
		passwordAuthManager:  passwordAuthManager,
		emailSender:          emailSender,
		DBManager:            dbManager,
	}
}

// IsSuperAdmin implements OrganizationService.
func (s *organizationMemberManager) IsSuperAdmin(ctx context.Context, userID shared.UUID[user.User]) (bool, error) {
	ctx, span := otel.Tracer("organization").Start(ctx, "organizationMemberManager.IsSuperAdmin")
	defer span.End()

	// ユーザーの組織ユーザーを取得
	orgUsers, err := s.organizationUserRepo.FindByUserID(ctx, userID)
	if err != nil {
		return false, err
	}

	// スーパーユーザーかどうかをチェック
	for _, orgUser := range orgUsers {
		if orgUser.Role >= organization.OrganizationUserRoleSuperAdmin {
			return true, nil
		}
	}

	return false, nil
}

// CreateOrganizationUser implements OrganizationMemberManager.
func (s *organizationMemberManager) InviteUser(ctx context.Context, input InviteUserParams) (*organization.OrganizationUser, error) {
	ctx, span := otel.Tracer("organization").Start(ctx, "organizationMemberManager.CreateOrganizationUser")
	defer span.End()

	// SuperAdminのみがユーザーを作成できる
	isSuperAdmin, err := s.IsSuperAdmin(ctx, input.UserID)
	if err != nil {
		utils.HandleError(ctx, err, "OrganizationMemberManager.IsSuperAdmin")
		return nil, err
	}
	if !isSuperAdmin {
		return nil, messages.OrganizationForbidden
	}

	var orgUser *organization.OrganizationUser
	if err := s.DBManager.ExecTx(ctx, func(ctx context.Context) error {

		// 組織取得
		org, err := s.organizationRepo.FindByID(ctx, input.OrganizationID)
		if err != nil {
			utils.HandleError(ctx, err, "OrganizationRepository.FindByID")
			return messages.OrganizationForbidden
		}

		// 単純にユーザーを作成
		authProviderName, err := auth.NewAuthProviderName("password")
		if err != nil {
			utils.HandleError(ctx, err, "AuthProviderName")
			return errtrace.Wrap(err)
		}
		subject, err := hash.HashEmail(input.Email, s.cfg.HASH_PEPPER)
		if err != nil {
			utils.HandleError(ctx, err, "HashEmail")
			return messages.InvalidPasswordOrEmailError
		}
		existUser, err := s.userRep.FindBySubject(ctx, user.UserSubject(subject))
		if err != nil {
			utils.HandleError(ctx, err, "UserRepository.FindBySubject")
			return errtrace.Wrap(err)
		}
		if existUser != nil {
			return errors.New("既に登録済みです。")
		}

		newUser := user.NewUser(
			shared.NewUUID[user.User](),
			nil,
			nil,
			subject,
			authProviderName,
			nil,
		)
		newUser.ChangeEmail(input.Email)

		version, err := s.policyRep.FetchLatestPolicy(ctx)
		if err != nil {
			utils.HandleError(ctx, err, "PolicyRepository.GetLatestVersion")
			return errtrace.Wrap(err)
		}
		_, err = s.consentService.RecordConsent(
			ctx,
			newUser.UserID(),
			version.Version,
			"",
			"",
		)
		if err != nil {
			utils.HandleError(ctx, err, "ConsentService.RecordConsent")
			return errtrace.Wrap(err)
		}
		if err := s.userRep.Create(ctx, newUser); err != nil {
			utils.HandleError(ctx, err, "UserRepository.Create")
			return errtrace.Wrap(err)
		}
		pass := password_auth.GeneratePassword(16)
		if err := s.passwordAuthManager.RegisterPassword(ctx, newUser.UserID(), pass, true); err != nil {
			utils.HandleError(ctx, err, "PasswordAuthManager.RegisterPassword")
			return errtrace.Wrap(err)
		}
		// 組織アカウントを作成
		orgUsr := organization.OrganizationUser{
			OrganizationUserID: shared.NewUUID[organization.OrganizationUser](),
			OrganizationID:     input.OrganizationID,
			UserID:             newUser.UserID(),
			Role:               input.Role,
		}
		if err := s.organizationUserRepo.Create(ctx, orgUsr); err != nil {
			utils.HandleError(ctx, err, "OrganizationUserRepository.Create")
			return errtrace.Wrap(err)
		}
		orgUser = &orgUsr
		// メールにIDとパスワード、組織IDを送信
		if err := s.emailSender.Send(ctx, input.Email, email_template.OrganizationInvitationEmailTemplate, map[string]any{
			"Title":            "【ことひろ】招待が届いています",
			"CompanyLogo":      "https://github.com/neko-dream/api/raw/develop/docs/public/assets/icon.png",
			"AppName":          s.cfg.APP_NAME,
			"WebsiteURL":       s.cfg.WEBSITE_URL,
			"OrganizationName": org.Name,
			"Email":            input.Email,
			"Password":         pass,
			"InvitationURL":    s.cfg.WEBSITE_URL,
		}); err != nil {
			utils.HandleError(ctx, err, "EmailSender.Send")
			return errtrace.Wrap(err)
		}

		return nil
	}); err != nil {
		utils.HandleError(ctx, err, "DBManager.ExecTx")
		return nil, err
	}

	return orgUser, nil
}
