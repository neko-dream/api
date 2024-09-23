package service

import (
	"context"
	"net/url"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/pkg/oauth"
)

type authService struct {
	userRepository user.UserRepository
}

func NewAuthService(
	userRepository user.UserRepository,
) auth.AuthService {
	return &authService{
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
		return nil, messages.InvalidProviderError
	}

	provider, err := oauth.OIDCProviderFactory(
		ctx,
		authProviderName,
	)
	if err != nil {
		return nil, err
	}

	subject, userName, err := provider.UserInfo(ctx, code)
	if err != nil {
		return nil, err
	}

	existUser, err := a.userRepository.FindBySubject(ctx, user.UserSubject(subject))
	if err != nil {
		return nil, err
	}
	if existUser != nil {
		return existUser, nil
	}

	newUser, err := a.userRepository.Create(ctx, user.NewUser(
		shared.NewUUID[user.User](),
		"",
		userName,
		subject,
		authProviderName,
	))
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (a *authService) GetAuthURL(
	ctx context.Context,
	providerName string,
) (*url.URL, string, error) {
	authProviderName, err := oauth.NewAuthProviderName(providerName)
	if err != nil {
		return nil, "", err
	}

	provider, err := oauth.OIDCProviderFactory(
		ctx,
		authProviderName,
	)
	if err != nil {
		return nil, "", err
	}

	state := provider.GenerateState()
	authURL := provider.GetAuthURL(ctx, state)
	url, err := url.Parse(authURL)
	if err != nil {
		return nil, "", err
	}

	return url, state, nil
}
