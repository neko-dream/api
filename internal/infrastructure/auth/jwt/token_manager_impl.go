package jwt

import (
	"context"
	"errors"
	"net/http"

	"braces.dev/errtrace"
	"github.com/golang-jwt/jwt/v5"
	"github.com/neko-dream/server/internal/domain/model/auth"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	http_utils "github.com/neko-dream/server/pkg/http"

	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type tokenManager struct {
	secret string
	*db.DBManager
	session.SessionRepository
	organization.OrganizationUserRepository
	organization.OrganizationRepository
}

// SetSession implements session.TokenManager.
func (j *tokenManager) SetSession(ctx context.Context) context.Context {
	ctx, span := otel.Tracer("jwt").Start(ctx, "tokenManager.SetSession")
	defer span.End()

	r := http_utils.GetHTTPRequest(ctx)
	if r == nil {
		return ctx
	}

	sessionCookie, err := r.Cookie("SessionId")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return ctx
		}
		// NoCookie以外の場合、エラーをハンドリング
		utils.HandleError(ctx, err, "CookieError")
		return ctx
	}
	// セッションIDを取得
	claim, err := j.Parse(ctx, sessionCookie.Value)
	if err != nil {
		return ctx
	}
	// トークンの有効性を確認 || // スキップするOperationの場合以外は、ユーザー登録済みか確認
	if claim.IsExpired(ctx) || !claim.IsRegistered {
		return ctx
	}

	sessID, err := claim.SessionID()
	if err != nil {
		return ctx
	}

	// サーバー側でセッションの有効性を確認
	sess, err := j.SessionRepository.FindBySessionID(ctx, sessID)
	if err != nil || sess == nil || !sess.IsActive(ctx) {
		return ctx
	}

	return session.SetSession(ctx, claim)
}

func (j *tokenManager) Generate(ctx context.Context, user user.User, sessionID shared.UUID[session.Session]) (string, error) {
	ctx, span := otel.Tracer("jwt").Start(ctx, "tokenManager.Generate")
	defer span.End()

	// セッション情報を取得
	sess, err := j.SessionRepository.FindBySessionID(ctx, sessionID)
	if err != nil {
		utils.HandleError(ctx, err, "FindBySessionID")
		return "", err
	}

	requiredPasswordChange := false
	if user.Provider() == auth.ProviderPassword {
		auths, err := j.DBManager.GetQueries(ctx).GetPasswordAuthByUserId(ctx, user.UserID().UUID())
		if err != nil {
			utils.HandleError(ctx, err, "GetPasswordAuthByUserId")
			return "", err
		}
		requiredPasswordChange = auths.PasswordAuth.RequiredPasswordChange
	}

	var orgType *int
	var organizationID *string
	var organizationCode *string
	var organizationRole *string

	// セッションに組織IDがある場合、その組織情報を優先的に使用
	if sess != nil && sess.OrganizationID() != nil && !sess.OrganizationID().IsZero() {
		orgID := shared.UUID[organization.Organization](sess.OrganizationID().UUID())
		org, err := j.OrganizationRepository.FindByID(ctx, orgID)
		if err == nil && org != nil {
			orgType = lo.ToPtr(int(org.OrganizationType))
			organizationID = lo.ToPtr(org.OrganizationID.String())
			organizationCode = lo.ToPtr(org.Code)
			// 組織でのユーザーのロールを取得
			orgUser, err := j.OrganizationUserRepository.FindByOrganizationIDAndUserID(ctx, orgID, user.UserID())
			if err == nil && orgUser != nil {
				organizationRole = lo.ToPtr(organization.RoleToName(orgUser.Role))
			}
		}
	}

	claim := session.NewClaimWithOrganization(ctx, user, sessionID, requiredPasswordChange, orgType, organizationID, organizationCode, organizationRole)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim.GenMapClaim())
	return errtrace.Wrap2(token.SignedString([]byte(j.secret)))
}

func (j *tokenManager) Parse(ctx context.Context, token string) (*session.Claim, error) {
	ctx, span := otel.Tracer("jwt").Start(ctx, "tokenManager.Parse")
	defer span.End()

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		// アルゴリズムの確認
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			utils.HandleError(ctx, jwt.ErrInvalidKeyType, "InvalidKeyType")
			return nil, errtrace.Wrap(jwt.ErrInvalidKeyType)
		}

		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	mapClaim, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errtrace.Wrap(jwt.ErrInvalidKeyType)
	}

	// issuerの検証
	if mapClaim["iss"].(string) != session.Issuer {
		return nil, errtrace.Wrap(errors.New("invalid issuer"))
	}

	claim := session.NewClaimFromMap(mapClaim)
	return &claim, nil
}

func NewTokenManager(
	sessRepo session.SessionRepository,
	conf *config.Config,
	dbm *db.DBManager,
	orgUserRep organization.OrganizationUserRepository,
	orgRep organization.OrganizationRepository,
) session.TokenManager {
	return &tokenManager{
		secret:                     conf.TokenSecret,
		SessionRepository:          sessRepo,
		DBManager:                  dbm,
		OrganizationUserRepository: orgUserRep,
		OrganizationRepository:     orgRep,
	}
}

func NewTokenManagerWithSecret(
	secret string,
	dbm *db.DBManager,
	sessRepo session.SessionRepository,
	orgUserRep organization.OrganizationUserRepository,
	orgRep organization.OrganizationRepository,
) session.TokenManager {
	return &tokenManager{
		secret:                     secret,
		SessionRepository:          sessRepo,
		DBManager:                  dbm,
		OrganizationUserRepository: orgUserRep,
		OrganizationRepository:     orgRep,
	}
}
