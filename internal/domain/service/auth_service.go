package service

import (
	"context"
	"net/url"

	"braces.dev/errtrace"
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
		return nil, errtrace.Wrap(messages.InvalidProviderError)
	}

	provider, err := oauth.OIDCProviderFactory(
		ctx,
		authProviderName,
	)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	userInfo, err := provider.UserInfo(ctx, code)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	existUser, err := a.userRepository.FindBySubject(ctx, user.UserSubject(userInfo.Subject))
	if err != nil {
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
		&userInfo.Picture,
	)
	if err := a.userRepository.Create(ctx, newUser); err != nil {
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
		authProviderName,
	)
	if err != nil {
		return nil, "", errtrace.Wrap(err)
	}

	state := provider.GenerateState()
	authURL := provider.GetAuthURL(ctx, state)
	url, err := url.Parse(authURL)
	if err != nil {
		return nil, "", errtrace.Wrap(err)
	}

	return url, state, nil
}
