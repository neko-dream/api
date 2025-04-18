package jwt

import (
	"context"
	"errors"
	"net/http"
	"sort"

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
	orgUsers, _ := j.OrganizationUserRepository.FindByUserID(ctx, user.UserID())
	// orgTypeでソート
	sort.SliceStable(orgUsers, func(i, j int) bool {
		return orgUsers[i].Role < orgUsers[j].Role
	})
	if len(orgUsers) > 0 {
		orgUser := orgUsers[0]
		// organizationをとる
		org, err := j.OrganizationRepository.FindByID(ctx, orgUser.OrganizationID)
		if err != nil {
			utils.HandleError(ctx, err, "GetOrganizationByID")
			return "", err
		}
		if org == nil {
			return "", errtrace.Wrap(errors.New("organization not found"))
		}
		// organizationのroleを取得
		orgType = lo.ToPtr(int(org.OrganizationType))
	}

	claim := session.NewClaim(ctx, user, sessionID, requiredPasswordChange, orgType)
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

func NewTokenManagerWithSecret(secret string, sessRepo session.SessionRepository) session.TokenManager {
	return &tokenManager{
		secret:            secret,
		SessionRepository: sessRepo,
	}
}
