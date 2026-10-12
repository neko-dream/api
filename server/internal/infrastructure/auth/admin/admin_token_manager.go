package admin

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	adminmodel "github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"time"
)

const adminIssuer = "kotohiro-admin"
const adminAudience = "kotohiro-admin-api"
const adminAccessTokenTTL = 15 * time.Minute

type adminTokenManager struct{ secret string }

func NewAdminTokenManager(conf *config.Config) adminmodel.AdminTokenManager {
	return &adminTokenManager{secret: conf.AdminJWTSecret}
}
func (m *adminTokenManager) Generate(ctx context.Context, user adminmodel.AdminUser) (string, error) {
	if len(m.secret) < 32 {
		return "", fmt.Errorf("ADMIN_JWT_SECRET must contain at least 32 bytes")
	}
	now := clock.Now(ctx)
	claims := jwt.RegisteredClaims{Issuer: adminIssuer, Audience: jwt.ClaimStrings{adminAudience}, Subject: user.AdminUserID.String(), IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(adminAccessTokenTTL))}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(m.secret))
}
func (m *adminTokenManager) Parse(ctx context.Context, raw string) (*adminmodel.AdminClaim, error) {
	if len(m.secret) < 32 {
		return nil, fmt.Errorf("ADMIN_JWT_SECRET must contain at least 32 bytes")
	}
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(raw, claims, func(*jwt.Token) (any, error) { return []byte(m.secret), nil }, jwt.WithValidMethods([]string{"HS256"}), jwt.WithIssuer(adminIssuer), jwt.WithAudience(adminAudience), jwt.WithExpirationRequired(), jwt.WithIssuedAt(), jwt.WithTimeFunc(func() time.Time { return clock.Now(ctx) }))
	if err != nil {
		return nil, err
	}
	claim := &adminmodel.AdminClaim{Sub: claims.Subject, Exp: claims.ExpiresAt.Unix()}
	if _, err := claim.AdminUserID(); err != nil {
		return nil, fmt.Errorf("invalid admin subject: %w", err)
	}
	if claims.IssuedAt != nil {
		claim.Iat = claims.IssuedAt.Unix()
	}
	return claim, nil
}
