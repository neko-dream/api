package oauth

import (
	"errors"
	"strings"

	"braces.dev/errtrace"
)

type AuthProviderName string

func (a AuthProviderName) String() string {
	return strings.ToUpper(string(a))
}

const (
	ProviderGoogle AuthProviderName = "GOOGLE"
)

func NewAuthProviderName(provider string) (AuthProviderName, error) {
	switch strings.ToUpper(provider) {
	case ProviderGoogle.String():
		return ProviderGoogle, nil
	default:
		return "", errtrace.Wrap(errors.New("invalid auth provider"))
	}
}
