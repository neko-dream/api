package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/neko-dream/api/internal/domain/model/session"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"
)

// モック: TokenManager
type mockTokenManager struct {
	parseFunc func(ctx context.Context, token string) (*session.Claim, error)
}

func (m *mockTokenManager) Generate(ctx context.Context, user user.User, sessionID shared.UUID[session.Session]) (string, error) {
	return "", nil
}

func (m *mockTokenManager) Parse(ctx context.Context, token string) (*session.Claim, error) {
	return m.parseFunc(ctx, token)
}

func (m *mockTokenManager) SetSession(ctx context.Context) context.Context {
	return ctx
}

// モック: SessionRepository
type mockSessionRepository struct {
	findBySessionIDFunc func(ctx context.Context, sessionID shared.UUID[session.Session]) (*session.Session, error)
}

func (m *mockSessionRepository) Create(ctx context.Context, sess session.Session) (*session.Session, error) {
	return nil, nil
}

func (m *mockSessionRepository) Update(ctx context.Context, sess session.Session) (*session.Session, error) {
	return nil, nil
}

func (m *mockSessionRepository) DeactivateAllByUserID(ctx context.Context, userID shared.UUID[user.User]) error {
	return nil
}

func (m *mockSessionRepository) FindBySessionID(ctx context.Context, sessionID shared.UUID[session.Session]) (*session.Session, error) {
	return m.findBySessionIDFunc(ctx, sessionID)
}

func (m *mockSessionRepository) FindByUserID(ctx context.Context, userID shared.UUID[user.User]) ([]session.Session, error) {
	return nil, nil
}

// テスト用のDIコンテナを作成
func setupTestContainer(tokenManager session.TokenManager, sessionRepo session.SessionRepository) *dig.Container {
	container := dig.New()
	container.Provide(func() session.TokenManager { return tokenManager })
	container.Provide(func() session.SessionRepository { return sessionRepo })
	return container
}

// 有効なClaimを作成するヘルパー
func newValidClaim() *session.Claim {
	now := time.Now()
	return &session.Claim{
		Sub: "00000000-0000-0000-0000-000000000001",
		Iat: now.Unix(),
		Exp: now.Add(time.Hour).Unix(),
		Jti: "00000000-0000-0000-0000-000000000002",
	}
}

// 期限切れのClaimを作成するヘルパー
func newExpiredClaim() *session.Claim {
	now := time.Now()
	return &session.Claim{
		Sub: "00000000-0000-0000-0000-000000000001",
		Iat: now.Add(-2 * time.Hour).Unix(),
		Exp: now.Add(-1 * time.Hour).Unix(), // 1時間前に期限切れ
		Jti: "00000000-0000-0000-0000-000000000002",
	}
}

// アクティブなセッションを作成するヘルパー
func newActiveSession() *session.Session {
	sessionID := shared.MustParseUUID[session.Session]("00000000-0000-0000-0000-000000000002")
	userID := shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001")
	return session.NewSession(
		sessionID,
		userID,
		shared.ProviderGuest,
		*session.NewSessionStatus(0), // SESSION_ACTIVE
		time.Now().Add(time.Hour),
		time.Now(),
	)
}

// 無効（revoke済み）のセッションを作成するヘルパー
func newInactiveSession() *session.Session {
	sessionID := shared.MustParseUUID[session.Session]("00000000-0000-0000-0000-000000000002")
	userID := shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000001")
	return session.NewSession(
		sessionID,
		userID,
		shared.ProviderGuest,
		*session.NewSessionStatus(1), // SESSION_INACTIVE
		time.Now().Add(time.Hour),
		time.Now(),
	)
}

func TestSetContextCookieKey(t *testing.T) {
	t.Run("Cookieがない場合はセッションが設定されない", func(t *testing.T) {
		tokenManager := &mockTokenManager{}
		sessionRepo := &mockSessionRepository{}
		container := setupTestContainer(tokenManager, sessionRepo)

		middleware := SetContextCookieKey(container)

		var capturedCtx context.Context
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedCtx = r.Context()
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		middleware(handler).ServeHTTP(rec, req)

		assert.Nil(t, session.GetSession(capturedCtx))
	})

	t.Run("無効なトークンの場合はセッションが設定されない", func(t *testing.T) {
		tokenManager := &mockTokenManager{
			parseFunc: func(ctx context.Context, token string) (*session.Claim, error) {
				return nil, errors.New("invalid token")
			},
		}
		sessionRepo := &mockSessionRepository{}
		container := setupTestContainer(tokenManager, sessionRepo)

		middleware := SetContextCookieKey(container)

		var capturedCtx context.Context
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedCtx = r.Context()
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(&http.Cookie{Name: "SessionId", Value: "invalid-token"})
		rec := httptest.NewRecorder()

		middleware(handler).ServeHTTP(rec, req)

		assert.Nil(t, session.GetSession(capturedCtx))
	})

	t.Run("期限切れトークンの場合はセッションが設定されない", func(t *testing.T) {
		tokenManager := &mockTokenManager{
			parseFunc: func(ctx context.Context, token string) (*session.Claim, error) {
				return newExpiredClaim(), nil
			},
		}
		sessionRepo := &mockSessionRepository{}
		container := setupTestContainer(tokenManager, sessionRepo)

		middleware := SetContextCookieKey(container)

		var capturedCtx context.Context
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedCtx = r.Context()
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(&http.Cookie{Name: "SessionId", Value: "expired-token"})
		rec := httptest.NewRecorder()

		middleware(handler).ServeHTTP(rec, req)

		assert.Nil(t, session.GetSession(capturedCtx))
	})

	t.Run("セッションが存在しない場合はセッションが設定されない", func(t *testing.T) {
		tokenManager := &mockTokenManager{
			parseFunc: func(ctx context.Context, token string) (*session.Claim, error) {
				return newValidClaim(), nil
			},
		}
		sessionRepo := &mockSessionRepository{
			findBySessionIDFunc: func(ctx context.Context, sessionID shared.UUID[session.Session]) (*session.Session, error) {
				return nil, errors.New("session not found")
			},
		}
		container := setupTestContainer(tokenManager, sessionRepo)

		middleware := SetContextCookieKey(container)

		var capturedCtx context.Context
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedCtx = r.Context()
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(&http.Cookie{Name: "SessionId", Value: "valid-token"})
		rec := httptest.NewRecorder()

		middleware(handler).ServeHTTP(rec, req)

		assert.Nil(t, session.GetSession(capturedCtx))
	})

	t.Run("セッションがrevoke済み（IsActive=false）の場合はセッションが設定されない", func(t *testing.T) {
		tokenManager := &mockTokenManager{
			parseFunc: func(ctx context.Context, token string) (*session.Claim, error) {
				return newValidClaim(), nil
			},
		}
		sessionRepo := &mockSessionRepository{
			findBySessionIDFunc: func(ctx context.Context, sessionID shared.UUID[session.Session]) (*session.Session, error) {
				return newInactiveSession(), nil
			},
		}
		container := setupTestContainer(tokenManager, sessionRepo)

		middleware := SetContextCookieKey(container)

		var capturedCtx context.Context
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedCtx = r.Context()
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(&http.Cookie{Name: "SessionId", Value: "valid-token"})
		rec := httptest.NewRecorder()

		middleware(handler).ServeHTTP(rec, req)

		assert.Nil(t, session.GetSession(capturedCtx))
	})

	t.Run("有効なセッションの場合はセッションが設定される", func(t *testing.T) {
		claim := newValidClaim()
		tokenManager := &mockTokenManager{
			parseFunc: func(ctx context.Context, token string) (*session.Claim, error) {
				return claim, nil
			},
		}
		sessionRepo := &mockSessionRepository{
			findBySessionIDFunc: func(ctx context.Context, sessionID shared.UUID[session.Session]) (*session.Session, error) {
				return newActiveSession(), nil
			},
		}
		container := setupTestContainer(tokenManager, sessionRepo)

		middleware := SetContextCookieKey(container)

		var capturedCtx context.Context
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedCtx = r.Context()
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(&http.Cookie{Name: "SessionId", Value: "valid-token"})
		rec := httptest.NewRecorder()

		middleware(handler).ServeHTTP(rec, req)

		capturedSession := session.GetSession(capturedCtx)
		require.NotNil(t, capturedSession)
		assert.Equal(t, claim.Sub, capturedSession.Sub)
		assert.Equal(t, claim.Jti, capturedSession.Jti)
	})

	t.Run("セッションがnilの場合はセッションが設定されない", func(t *testing.T) {
		tokenManager := &mockTokenManager{
			parseFunc: func(ctx context.Context, token string) (*session.Claim, error) {
				return newValidClaim(), nil
			},
		}
		sessionRepo := &mockSessionRepository{
			findBySessionIDFunc: func(ctx context.Context, sessionID shared.UUID[session.Session]) (*session.Session, error) {
				return nil, nil // エラーなしでnilを返す
			},
		}
		container := setupTestContainer(tokenManager, sessionRepo)

		middleware := SetContextCookieKey(container)

		var capturedCtx context.Context
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedCtx = r.Context()
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(&http.Cookie{Name: "SessionId", Value: "valid-token"})
		rec := httptest.NewRecorder()

		middleware(handler).ServeHTTP(rec, req)

		assert.Nil(t, session.GetSession(capturedCtx))
	})
}
