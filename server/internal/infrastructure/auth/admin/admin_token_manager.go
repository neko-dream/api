package admin

import (
	"context"
	"errors"
	"time"

	"braces.dev/errtrace"
	"github.com/golang-jwt/jwt/v5"
	adminmodel "github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

// admin トークンの発行者。user 用 (session.Issuer) とは別の値にして取り違えを防ぐ。
const adminIssuer = "https://api.kotohiro.com/admin"

type adminTokenManager struct {
	secret string
}

func NewAdminTokenManager(conf *config.Config) adminmodel.AdminTokenManager {
	return &adminTokenManager{secret: conf.TokenSecret}
}

func (m *adminTokenManager) Generate(ctx context.Context, adminUser adminmodel.AdminUser) (string, error) {
	ctx, span := otel.Tracer("adminTokenManager").Start(ctx, "Generate")
	defer span.End()

	now := clock.Now(ctx)
	claims := jwt.MapClaims{
		"iss":         adminIssuer,
		"sub":         adminUser.AdminUserID.String(),
		"email":       adminUser.Email,
		"displayName": adminUser.DisplayName,
		"role":        int(adminUser.Role),
		"isAdmin":     true,
		"iat":         now.Unix(),
		// access は短命。失効後は refresh トークンで再発行する。
		"exp": now.Add(time.Minute * 15).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return errtrace.Wrap2(token.SignedString([]byte(m.secret)))
}

func (m *adminTokenManager) Parse(ctx context.Context, token string) (*adminmodel.AdminClaim, error) {
	ctx, span := otel.Tracer("adminTokenManager").Start(ctx, "Parse")
	defer span.End()

	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			utils.HandleError(ctx, jwt.ErrInvalidKeyType, "InvalidKeyType")
			return nil, errtrace.Wrap(jwt.ErrInvalidKeyType)
		}
		return []byte(m.secret), nil
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	mapClaim, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errtrace.Wrap(jwt.ErrInvalidKeyType)
	}

	if iss, _ := mapClaim["iss"].(string); iss != adminIssuer {
		return nil, errtrace.Wrap(errors.New("invalid issuer"))
	}

	return adminClaimFromMap(mapClaim), nil
}

func adminClaimFromMap(claims jwt.MapClaims) *adminmodel.AdminClaim {
	claim := &adminmodel.AdminClaim{}
	if v, ok := claims["sub"].(string); ok {
		claim.Sub = v
	}
	if v, ok := claims["email"].(string); ok {
		claim.Email = v
	}
	if v, ok := claims["displayName"].(string); ok {
		claim.DisplayName = v
	}
	if v, ok := claims["role"].(float64); ok {
		claim.Role = int(v)
	}
	if v, ok := claims["iat"].(float64); ok {
		claim.Iat = int64(v)
	}
	if v, ok := claims["exp"].(float64); ok {
		claim.Exp = int64(v)
	}
	return claim
}
