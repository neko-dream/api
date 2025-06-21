package shared

import (
	"errors"
	"strings"

	"braces.dev/errtrace"
)

type (
	IssuerURI        string
	AuthProviderName string
)

const (
	GoogleIssuerURI IssuerURI = "https://accounts.google.com"
	LineIssuerURI   IssuerURI = "https://access.line.me"
)

func (i IssuerURI) String() string {
	return string(i)
}

const (
	ProviderGoogle   AuthProviderName = "GOOGLE"
	ProviderLine     AuthProviderName = "LINE"
	ProviderDEV      AuthProviderName = "DEV"
	ProviderPassword AuthProviderName = "PASSWORD"
)

func NewAuthProviderName(provider string) (AuthProviderName, error) {
	switch strings.ToUpper(provider) {
	case ProviderGoogle.String():
		return ProviderGoogle, nil
	case ProviderLine.String():
		return ProviderLine, nil
	case ProviderDEV.String():
		return ProviderDEV, nil
	case ProviderPassword.String():
		return ProviderPassword, nil
	default:
		return "", errtrace.Wrap(errors.New("invalid auth provider"))
	}
}

func (a AuthProviderName) IssuerURI() IssuerURI {
	switch a {
	case ProviderGoogle:
		return GoogleIssuerURI
	case ProviderLine:
		return LineIssuerURI
	default:
		return ""
	}
}

func (a AuthProviderName) String() string {
	return string(a)
}

type OrganizationUserRole int

func NewOrganizationUserRole(role int) OrganizationUserRole {
	if role < int(OrganizationUserRoleSuperAdmin) || role > int(OrganizationUserRoleMember) {
		return OrganizationUserRoleMember
	}
	if role == 0 {
		return OrganizationUserRoleMember
	}

	return OrganizationUserRole(role)
}

func NameToRole(name string) OrganizationUserRole {
	switch name {
	case "メンバー":
		return OrganizationUserRoleMember
	case "管理者":
		return OrganizationUserRoleAdmin
	case "オーナー":
		return OrganizationUserRoleOwner
	case "運営":
		return OrganizationUserRoleSuperAdmin
	default:
		return OrganizationUserRoleMember
	}
}

func RoleToName(role OrganizationUserRole) string {
	switch role {
	case OrganizationUserRoleMember:
		return "メンバー"
	case OrganizationUserRoleAdmin:
		return "管理者"
	case OrganizationUserRoleOwner:
		return "オーナー"
	case OrganizationUserRoleSuperAdmin:
		return "運営"
	default:
		return "メンバー"
	}
}

const (
	OrganizationUserRoleSuperAdmin OrganizationUserRole = 10
	OrganizationUserRoleOwner      OrganizationUserRole = 20
	OrganizationUserRoleAdmin      OrganizationUserRole = 30
	OrganizationUserRoleMember     OrganizationUserRole = 40
)