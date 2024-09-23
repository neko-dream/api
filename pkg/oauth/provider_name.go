package oauth

import (
	"errors"
	"strings"
)

type AuthProviderName string

func (a AuthProviderName) String() string {
	return string(a)
}

const (
	ProviderGoogle AuthProviderName = "GOOGLE"
)

func NewAuthProviderName(provider string) (AuthProviderName, error) {
	switch strings.ToUpper(provider) {
	case ProviderGoogle.String():
		return ProviderGoogle, nil
	default:
		return "", errors.New("invalid auth provider")
	}
}
