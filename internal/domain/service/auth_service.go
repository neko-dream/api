package service

import (
	"context"
	"crypto/rand"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type AuthService interface {
	Authenticate(ctx context.Context, provider, code string) (*user.User, error)
	GenerateState(ctx context.Context) (string, error)
}

type authService struct {
	config              *config.Config
	userRepository      user.UserRepository
	authProviderFactory auth.AuthProviderFactory
}

func NewAuthService(
	config *config.Config,
	userRepository user.UserRepository,
	authProviderFactory auth.AuthProviderFactory,
) AuthService {
	return &authService{
		config:              config,
		userRepository:      userRepository,
		authProviderFactory: authProviderFactory,
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

	existUser, err := a.userRepository.FindBySubject(ctx, user.UserSubject(*subject))
	if err != nil {
		utils.HandleError(ctx, err, "UserRepository.FindBySubject")
		return nil, errtrace.Wrap(err)
	}
	if existUser != nil {
		return existUser, nil
	}

	authProviderName, err := auth.NewAuthProviderName(providerName)
	if err != nil {
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
		newUser.SetEmail(*email)
		// Auth時点でemailが確認済みの場合はVerifyEmailを実行
		newUser.VerifyEmail()
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
