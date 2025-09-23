package session

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type (
	sessionContextKey string

	TokenManager interface {
		Generate(ctx context.Context, user user.User, sessionID shared.UUID[Session]) (string, error)
		Parse(ctx context.Context, token string) (*Claim, error)
		SetSession(ctx context.Context) context.Context
	}

	Claim struct {
		Sub                    string     `json:"sub"` // subject (user)
		Iat                    int64      `json:"iat"` // issued at (seconds)
		Exp                    int64      `json:"exp"` // expiration time (seconds)
		Jti                    string     `json:"jti"` // JWT ID（SessionID）
		IconURL                *string    `json:"iconURL,omitempty"`
		DisplayName            *string    `json:"displayName,omitempty"`
		DisplayID              *string    `json:"displayID,omitempty"`
		IsRegistered           bool       `json:"isRegistered"`
		IsEmailVerified        bool       `json:"isEmailVerified"`
		RequiredPasswordChange bool       `json:"requiredPasswordChange"`
		OrgType                *int       `json:"orgType,omitempty"`          // 組織の種類
		OrganizationID         *string    `json:"organizationID,omitempty"`   // ログイン時に使用した組織ID
		OrganizationCode       *string    `json:"organizationCode,omitempty"` // ログイン時に使用した組織コード
		OrganizationRole       *string    `json:"organizationRole,omitempty"` // 組織でのロール名
		IsWithdrawn            bool       `json:"isWithdrawn,omitempty"`      // 退会ユーザーフラグ
		WithdrawalDate         *time.Time `json:"withdrawalDate,omitempty"`   // 退会日時
	}
)

func NewClaim(ctx context.Context, user user.User, sessionID shared.UUID[Session], requiredPasswordChange bool, orgType *int) Claim {
	ctx, span := otel.Tracer("session").Start(ctx, "NewClaim")
	defer span.End()
	now := clock.Now(ctx)
	return Claim{
		Sub:                    user.UserID().String(),
		Iat:                    now.Unix(),
		Exp:                    now.Add(time.Second * 60 * 60 * 24 * 7).Unix(),
		Jti:                    sessionID.String(),
		IconURL:                user.IconURL(),
		DisplayID:              user.DisplayID(),
		DisplayName:            user.DisplayName(),
		IsRegistered:           user.Verify(),
		IsEmailVerified:        user.IsEmailVerified(),
		RequiredPasswordChange: requiredPasswordChange,
		OrgType:                orgType,
	}
}

func NewClaimWithOrganization(ctx context.Context, user user.User, sessionID shared.UUID[Session], requiredPasswordChange bool, orgType *int, organizationID *string, organizationCode *string, organizationRole *string) Claim {
	ctx, span := otel.Tracer("session").Start(ctx, "NewClaimWithOrganization")
	defer span.End()
	now := clock.Now(ctx)
	return Claim{
		Sub:                    user.UserID().String(),
		Iat:                    now.Unix(),
		Exp:                    now.Add(time.Second * 60 * 60 * 24 * 7).Unix(),
		Jti:                    sessionID.String(),
		IconURL:                user.IconURL(),
		DisplayID:              user.DisplayID(),
		DisplayName:            user.DisplayName(),
		IsRegistered:           user.Verify(),
		IsEmailVerified:        user.IsEmailVerified(),
		RequiredPasswordChange: requiredPasswordChange,
		OrgType:                orgType,
		OrganizationID:         organizationID,
		OrganizationCode:       organizationCode,
		OrganizationRole:       organizationRole,
	}
}

func NewClaimFromMap(claims jwt.MapClaims) Claim {
	var picture, displayName, displayID *string

	if claims["iconURL"] != nil {
		picture = lo.ToPtr(claims["iconURL"].(string))
	}
	if claims["displayName"] != nil {
		displayName = lo.ToPtr(claims["displayName"].(string))
	}
	if claims["displayID"] != nil {
		displayID = lo.ToPtr(claims["displayID"].(string))
	}
	var isEmailVerified bool
	if claims["isEmailVerified"] != nil {
		isEmailVerified = claims["isEmailVerified"].(bool)
	}
	var requiredPasswordChange bool
	if claims["requiredPasswordChange"] != nil {
		requiredPasswordChange = claims["requiredPasswordChange"].(bool)
	}
	var orgType *int
	if claims["orgType"] != nil {
		orgType = lo.ToPtr(int(claims["orgType"].(float64)))
	}
	var organizationID *string
	if claims["organizationID"] != nil {
		organizationID = lo.ToPtr(claims["organizationID"].(string))
	}
	var organizationCode *string
	if claims["organizationCode"] != nil {
		organizationCode = lo.ToPtr(claims["organizationCode"].(string))
	}
	var organizationRole *string
	if claims["organizationRole"] != nil {
		organizationRole = lo.ToPtr(claims["organizationRole"].(string))
	}

	return Claim{
		Sub:                    claims["sub"].(string),
		Iat:                    int64(claims["iat"].(float64)),
		Exp:                    int64(claims["exp"].(float64)),
		Jti:                    claims["jti"].(string),
		IconURL:                picture,
		DisplayName:            displayName,
		DisplayID:              displayID,
		IsRegistered:           claims["isRegistered"].(bool),
		IsEmailVerified:        isEmailVerified,
		RequiredPasswordChange: requiredPasswordChange,
		OrgType:                orgType,
		OrganizationID:         organizationID,
		OrganizationCode:       organizationCode,
		OrganizationRole:       organizationRole,
	}
}

func (c *Claim) UserID() (shared.UUID[user.User], error) {
	return shared.ParseUUID[user.User](c.Sub)
}
func (c *Claim) SessionID() (shared.UUID[Session], error) {
	return shared.ParseUUID[Session](c.Jti)
}

func (c *Claim) Audience() string {
	return Audience
}
func (c *Claim) Issuer() string {
	return Issuer
}

func (c *Claim) IsExpired(ctx context.Context) bool {
	ctx, span := otel.Tracer("session").Start(ctx, "Claim.IsExpired")
	defer span.End()

	return clock.Now(ctx).Unix() > c.Exp
}
func (c *Claim) IssueAt() time.Time {
	return time.Unix(c.Iat, 0)
}
func (c *Claim) ExpiresAt() time.Time {
	return time.Unix(c.Exp, 0)
}

const (
	Issuer   = "https://api.kotohiro.com"
	Audience = "https://api.kotohiro.com"
)

func (c *Claim) GenMapClaim() *jwt.MapClaims {
	cl := &jwt.MapClaims{
		"exp":                    c.Exp,
		"iat":                    c.Iat,
		"jti":                    c.Jti,
		"sub":                    c.Sub,
		"iss":                    Issuer,
		"aud":                    Audience,
		"iconURL":                c.IconURL,
		"displayName":            c.DisplayName,
		"displayID":              c.DisplayID,
		"isRegistered":           c.IsRegistered,
		"isEmailVerified":        c.IsEmailVerified,
		"requiredPasswordChange": c.RequiredPasswordChange,
	}

	if c.OrgType != nil {
		(*cl)["orgType"] = *c.OrgType
	} else {
		(*cl)["orgType"] = nil
	}

	if c.OrganizationID != nil {
		(*cl)["organizationID"] = *c.OrganizationID
	}

	if c.OrganizationCode != nil {
		(*cl)["organizationCode"] = *c.OrganizationCode
	}

	if c.OrganizationRole != nil {
		(*cl)["organizationRole"] = *c.OrganizationRole
	}

	return cl
}

// SessionContextKey
var (
	sessKey sessionContextKey = "session"
)

func SetSession(ctx context.Context, claim *Claim) context.Context {
	ctx, span := otel.Tracer("session").Start(ctx, "SetSession")
	defer span.End()

	return context.WithValue(ctx, sessKey, claim)
}

func GetSession(ctx context.Context) *Claim {
	ctx, span := otel.Tracer("session").Start(ctx, "GetSession")
	defer span.End()

	if ctx == nil {
		return nil
	}

	value := ctx.Value(sessKey)
	if value == nil {
		return nil
	}

	claim, ok := value.(*Claim)
	if !ok {
		return nil
	}

	return claim
}
