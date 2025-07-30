package session

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/organization"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type sessionTokenManager struct {
	*db.DBManager
	secret string
	session.SessionRepository
	user.UserRepository
	organization.OrganizationUserRepository
	organization.OrganizationRepository
}

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidToken    = errors.New("invalid token")
)

// Generate implements session.TokenManager.
// SessionIDベースの実装では、署名付きトークンを生成する
func (s *sessionTokenManager) Generate(ctx context.Context, user user.User, sessionID shared.UUID[session.Session]) (string, error) {
	_, span := otel.Tracer("sessionTokenManager").Start(ctx, "Generate")
	defer span.End()

	// セッションIDと署名を組み合わせたトークンを生成
	token := s.createSignedToken(sessionID.String())
	return token, nil
}

// createSignedToken セッションIDに署名を付けたトークンを生成
func (s *sessionTokenManager) createSignedToken(sessionID string) string {
	// HMAC-SHA256で署名を生成
	h := hmac.New(sha256.New, []byte(s.secret))
	h.Write([]byte(sessionID))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	// セッションIDと署名を.で結合
	return fmt.Sprintf("%s.%s", sessionID, signature)
}

// verifySignedToken トークンの署名を検証し、セッションIDを返す
func (s *sessionTokenManager) verifySignedToken(token string) (string, error) {
	// トークンを分割
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return "", ErrInvalidToken
	}

	sessionID := parts[0]
	providedSignature := parts[1]
	decodedSignature, err := url.QueryUnescape(providedSignature)
	if err != nil {
		return "", ErrInvalidToken
	}

	// 署名を再計
	h := hmac.New(sha256.New, []byte(s.secret))
	h.Write([]byte(sessionID))
	expectedSignature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	// 署名を比較
	if !hmac.Equal([]byte(decodedSignature), []byte(expectedSignature)) {
		return "", ErrInvalidToken
	}

	return sessionID, nil
}

// Parse implements session.TokenManager.
// 署名付きトークンを検証し、セッション情報を取得してClaimオブジェクトを構築する
func (s *sessionTokenManager) Parse(ctx context.Context, token string) (*session.Claim, error) {
	ctx, span := otel.Tracer("sessionTokenManager").Start(ctx, "Parse")
	defer span.End()

	// トークンを検証してセッションIDを取得
	sessionIDStr, err := s.verifySignedToken(token)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	// セッションIDとしてパース
	sessionID, err := shared.ParseUUID[session.Session](sessionIDStr)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	// セッション情報を取得
	sess, err := s.SessionRepository.FindBySessionID(ctx, sessionID)
	if err != nil {
		utils.HandleError(ctx, err, "SessionRepository.FindBySessionID")
		return nil, errtrace.Wrap(ErrSessionNotFound)
	}

	// セッションが無効な場合はエラー
	if sess.Status() != session.SESSION_ACTIVE {
		return nil, errtrace.Wrap(ErrSessionNotFound)
	}

	// ユーザー情報を取得
	user, err := s.UserRepository.FindByID(ctx, sess.UserID())
	if err != nil {
		utils.HandleError(ctx, err, "UserRepository.FindByID")
		return nil, errtrace.Wrap(ErrUserNotFound)
	}

	// パスワード変更要求フラグの確認
	requiredPasswordChange := false
	if sess.Provider() == "password" {
		// パスワード認証の場合、必要に応じてパスワード変更要求フラグを設定
		// この実装は既存のJWT実装と同様の動作を維持
		requiredPasswordChange = false
	}

	// 組織情報の取得（セッションに組織IDが設定されている場合）
	var organizationID, organizationCode, organizationRole *string
	var orgType *int
	if sess.OrganizationID() != nil {
		orgUUID, err := shared.ParseUUID[organization.Organization](sess.OrganizationID().String())
		if err != nil {
			utils.HandleError(ctx, err, "ParseUUID")
			return nil, nil
		}
		org, err := s.OrganizationRepository.FindByID(ctx, orgUUID)
		if err == nil && org != nil {
			organizationID = lo.ToPtr(org.OrganizationID.String())
			organizationCode = lo.ToPtr(org.Code)
			orgType = lo.ToPtr(int(org.OrganizationType))

			// 組織でのロールを取得
			orgUser, err := s.OrganizationUserRepository.FindByOrganizationIDAndUserID(ctx, org.OrganizationID, sess.UserID())
			if err == nil && orgUser != nil {
				organizationRole = lo.ToPtr(organization.RoleToName(orgUser.Role))
			}
		}
	}

	// Claimオブジェクトを構築
	return &session.Claim{
		Sub:                    user.UserID().String(),
		Iat:                    sess.LastActivityAt().Unix(),
		Exp:                    sess.ExpiresAt().Unix(),
		Jti:                    sessionID.String(),
		IconURL:                user.IconURL(),
		DisplayName:            user.DisplayName(),
		DisplayID:              user.DisplayID(),
		IsRegistered:           user.Verify(),
		IsEmailVerified:        user.IsEmailVerified(),
		RequiredPasswordChange: requiredPasswordChange,
		OrgType:                orgType,
		OrganizationID:         organizationID,
		OrganizationCode:       organizationCode,
		OrganizationRole:       organizationRole,
	}, nil
}

// SetSession implements session.TokenManager.
// SessionIDベースの実装では、このメソッドは使用されません。
// SecurityHandlerでセッション情報がコンテキストに設定されるため、ここでは何もしません。
func (s *sessionTokenManager) SetSession(ctx context.Context) context.Context {
	ctx, span := otel.Tracer("sessionTokenManager").Start(ctx, "SetSession")
	defer span.End()

	// SessionIDベースの実装では、SecurityHandlerで既にセッション情報が設定されているため、
	// ここでは何もする必要がありません。
	return ctx
}

// NewSessionTokenManager creates a new session-based token manager
func NewSessionTokenManager(
	config *config.Config,
	dbManager *db.DBManager,
	sessionRepository session.SessionRepository,
	userRepository user.UserRepository,
	organizationUserRepository organization.OrganizationUserRepository,
	organizationRepository organization.OrganizationRepository,
) session.TokenManager {
	return &sessionTokenManager{
		secret:                     config.TokenSecret,
		DBManager:                  dbManager,
		SessionRepository:          sessionRepository,
		UserRepository:             userRepository,
		OrganizationUserRepository: organizationUserRepository,
		OrganizationRepository:     organizationRepository,
	}
}
