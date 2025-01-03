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
	provider, err := a.authProviderFactory.NewAuthProvider(ctx, providerName)
	if err != nil {
		utils.HandleError(ctx, err, "AuthProviderFactory.NewAuthProvider")
		return nil, errtrace.Wrap(err)
	}

	subject, _, err := provider.VerifyAndIdentify(ctx, code)
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

	newUser := user.NewUser(
		shared.NewUUID[user.User](),
		nil,
		nil,
		*subject,
		auth.AuthProviderName(providerName),
		nil,
	)
	if err := a.userRepository.Create(ctx, newUser); err != nil {
		utils.HandleError(ctx, err, "UserRepository.Create")
		return nil, errtrace.Wrap(err)
	}

	return &newUser, nil
}

var randTable = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

func (a *authService) GenerateState(ctx context.Context) (string, error) {
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
