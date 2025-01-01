package service

import (
	"context"
	"net/url"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/pkg/oauth"
	"github.com/neko-dream/server/pkg/utils"
)

type authService struct {
	config         *config.Config
	userRepository user.UserRepository
}

func NewAuthService(
	config *config.Config,
	userRepository user.UserRepository,
) auth.AuthService {
	return &authService{
		config:         config,
		userRepository: userRepository,
	}
}

func (a *authService) Authenticate(
	ctx context.Context,
	providerName,
	code string,
) (*user.User, error) {
	authProviderName, err := oauth.NewAuthProviderName(providerName)
	if err != nil {
		utils.HandleError(ctx, err, "NewAuthProviderName")
		return nil, errtrace.Wrap(messages.InvalidProviderError)
	}

	provider, err := oauth.OIDCProviderFactory(
		ctx,
		a.config,
		authProviderName,
	)
	if err != nil {
		utils.HandleError(ctx, err, "OIDCProviderFactory")
		return nil, errtrace.Wrap(err)
	}

	userInfo, err := provider.UserInfo(ctx, code)
	if err != nil {
		utils.HandleError(ctx, err, "OIDCProvider.UserInfo")
		return nil, errtrace.Wrap(err)
	}
	existUser, err := a.userRepository.FindBySubject(ctx, user.UserSubject(userInfo.Subject))
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
		userInfo.Subject,
		authProviderName,
		nil,
	)
	if err := a.userRepository.Create(ctx, newUser); err != nil {
		utils.HandleError(ctx, err, "UserRepository.Create")
		return nil, errtrace.Wrap(err)
	}

	return &newUser, nil
}

func (a *authService) GetAuthURL(
	ctx context.Context,
	providerName string,
) (*url.URL, string, error) {
	authProviderName, err := oauth.NewAuthProviderName(providerName)
	if err != nil {
		return nil, "", errtrace.Wrap(err)
	}

	provider, err := oauth.OIDCProviderFactory(
		ctx,
		a.config,
		authProviderName,
	)
	if err != nil {
		utils.HandleError(ctx, err, "OIDCProviderFactory")
		return nil, "", errtrace.Wrap(err)
	}

	state, err := a.GenerateState(ctx)
	if err != nil {
		utils.HandleError(ctx, err, "GenerateState")
		return nil, "", errtrace.Wrap(err)
	}

	authURL := provider.GetAuthURL(ctx, state)
	url, err := url.Parse(authURL)
	if err != nil {
		utils.HandleError(ctx, err, "url.Parse")
		return nil, "", errtrace.Wrap(err)
	}

	return url, state, nil
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
