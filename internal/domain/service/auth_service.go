package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/consent"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

var (
	ErrNotAuthenticated        = errors.New("not authenticated")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrNotInOrganization       = errors.New("not in organization context")
)

type AuthenticationService interface {
	// 現在認証されているユーザーのコンテキストを取得
	GetCurrentUser(ctx context.Context) (*auth.AuthenticationContext, error)

	// ユーザーが認証されていることを確認
	RequireAuthentication(ctx context.Context) (*auth.AuthenticationContext, error)

	// 指定された役割以上の権限を持つことを確認
	RequireOrganizationRole(ctx context.Context, minRole organization.OrganizationUserRole) (*auth.AuthenticationContext, error)

	// スーパー管理者権限を持つことを確認
	RequireSuperAdmin(ctx context.Context) (*auth.AuthenticationContext, error)

	// オーナー権限を持つことを確認
	RequireOwner(ctx context.Context) (*auth.AuthenticationContext, error)

	// 管理者以上の権限を持つことを確認
	RequireAdmin(ctx context.Context) (*auth.AuthenticationContext, error)

	// 認証されているかをチェック
	IsAuthenticated(ctx context.Context) bool

	// 組織コンテキスト内かをチェック
	IsInOrganization(ctx context.Context) bool
}

type AuthService interface {
	Authenticate(ctx context.Context, provider, code string) (*user.User, error)
	GenerateState(ctx context.Context) (string, error)
}

type authService struct {
	config              *config.Config
	userRepository      user.UserRepository
	authProviderFactory auth.AuthProviderFactory
	consentService      consent.ConsentService
	policyRepository    consent.PolicyRepository
}

func NewAuthService(
	config *config.Config,
	userRepository user.UserRepository,
	authProviderFactory auth.AuthProviderFactory,
	consentService consent.ConsentService,
	policyRepository consent.PolicyRepository,
) AuthService {
	return &authService{
		config:              config,
		userRepository:      userRepository,
		authProviderFactory: authProviderFactory,
		consentService:      consentService,
		policyRepository:    policyRepository,
	}
}

func (a *authService) Authenticate(
	ctx context.Context,
	providerName,
	code string,
) (*user.User, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "authService.Authenticate")
	defer span.End()

	provider, err := a.authProviderFactory.NewAuthProvider(ctx, providerName)
	if err != nil {
		utils.HandleError(ctx, err, "AuthProviderFactory.NewAuthProvider")
		return nil, errtrace.Wrap(err)
	}

	subject, email, err := provider.VerifyAndIdentify(ctx, code)
	if err != nil {
		utils.HandleError(ctx, err, "OIDCProvider.UserInfo")
		return nil, errtrace.Wrap(err)
	}
	if subject == nil {
		return nil, ErrNotAuthenticated
	}

	existUser, err := a.userRepository.FindBySubject(ctx, user.UserSubject(*subject))
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			utils.HandleError(ctx, err, "UserRepository.FindBySubject")
			return nil, errtrace.Wrap(err)
		}
	}
	// 退会ユーザーで31日以上経過している場合の処理
	if existUser != nil && existUser.IsWithdrawn() && existUser.IsReactivationPeriodExpired(clock.Now(ctx)) {
		// 古いユーザーのsubjectとemailを変更して重複を回避
		newSubject := existUser.PrepareForDeleteUser()

		// 更新をDBに反映
		if err := a.userRepository.Update(ctx, *existUser); err != nil {
			utils.HandleError(ctx, err, "UserRepository.Update withdrawn user")
			return nil, errtrace.Wrap(err)
		}
		// ChangeSubjectクエリも実行
		if err := a.userRepository.ChangeSubject(ctx, existUser.UserID(), newSubject); err != nil {
			utils.HandleError(ctx, err, "UserRepository.ChangeSubject")
			return nil, errtrace.Wrap(err)
		}
		existUser = nil // 新規ユーザーとして処理を続行
	}
	if existUser != nil {
		return existUser, nil
	}

	authProviderName, err := shared.NewAuthProviderName(providerName)
	if err != nil {
		utils.HandleError(ctx, err, "AuthProviderName.NewAuthProviderName")
		return nil, errtrace.Wrap(err)
	}
	newUser := user.NewUser(
		shared.NewUUID[user.User](),
		nil,
		nil,
		*subject,
		authProviderName,
		nil,
	)
	if email != nil {
		newUser.ChangeEmail(*email)
		// Auth時点でemailが確認済みの場合はVerifyEmailを実行
		newUser.SetEmailVerified(true)
	}
	version, err := a.policyRepository.FetchLatestPolicy(ctx)
	if err != nil {
		utils.HandleError(ctx, err, "PolicyRepository.GetLatestVersion")
		return nil, errtrace.Wrap(err)
	}
	_, err = a.consentService.RecordConsent(
		ctx,
		newUser.UserID(),
		version.Version,
		"",
		"",
	)
	if err != nil {
		utils.HandleError(ctx, err, "ConsentService.RecordConsent")
		return nil, errtrace.Wrap(err)
	}

	if err := a.userRepository.Create(ctx, newUser); err != nil {
		utils.HandleError(ctx, err, "UserRepository.Create")
		return nil, errtrace.Wrap(err)
	}

	return &newUser, nil
}

var randTable = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

func (a *authService) GenerateState(ctx context.Context) (string, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "authService.GenerateState")
	defer span.End()

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		utils.HandleError(ctx, err, "rand.Read")
		return "", errtrace.Wrap(err)
	}

	for i, v := range b {
		b[i] = randTable[v%byte(len(randTable))]
	}

	return string(b), nil
}
