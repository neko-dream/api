package auth_usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/neko-dream/api/internal/domain/messages"
	password_auth "github.com/neko-dream/api/internal/domain/model/auth/password"
	"github.com/neko-dream/api/internal/domain/model/session"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
)

type changePasswordTestManager struct {
	password_auth.PasswordAuthManager
	valid     bool
	updated   bool
	updateErr error
}

func (m *changePasswordTestManager) VerifyPassword(context.Context, shared.UUID[user.User], string) (bool, error) {
	return m.valid, nil
}

func (m *changePasswordTestManager) UpdatePassword(context.Context, shared.UUID[user.User], string) error {
	m.updated = m.updateErr == nil
	return m.updateErr
}

type changePasswordTestSessions struct {
	session.SessionRepository
	sessions map[shared.UUID[session.Session]]*session.Session
}

func (s *changePasswordTestSessions) FindByUserID(_ context.Context, userID shared.UUID[user.User]) ([]session.Session, error) {
	var out []session.Session
	for _, sess := range s.sessions {
		if sess.UserID() == userID {
			out = append(out, *sess)
		}
	}
	return out, nil
}

func (s *changePasswordTestSessions) Update(_ context.Context, sess session.Session) (*session.Session, error) {
	s.sessions[sess.SessionID()] = &sess
	return &sess, nil
}

func TestChangePasswordDeactivatesOtherSessions(t *testing.T) {
	ctx := context.Background()
	userID := shared.NewUUID[user.User]()
	otherUserID := shared.NewUUID[user.User]()
	expires := time.Now().Add(time.Hour)
	newSession := func(owner shared.UUID[user.User]) *session.Session {
		return session.NewSession(shared.NewUUID[session.Session](), owner, "PASSWORD", session.SESSION_ACTIVE, expires, time.Now())
	}
	current, other, otherUsers := newSession(userID), newSession(userID), newSession(otherUserID)
	newSessions := func() *changePasswordTestSessions {
		return &changePasswordTestSessions{sessions: map[shared.UUID[session.Session]]*session.Session{
			current.SessionID():    current,
			other.SessionID():      other,
			otherUsers.SessionID(): otherUsers,
		}}
	}
	input := ChangePasswordInput{UserID: userID, SessionID: current.SessionID(), OldPassword: "old-password", NewPassword: "new-password"}

	t.Run("現在のセッション以外を無効にする", func(t *testing.T) {
		sessions := newSessions()
		manager := &changePasswordTestManager{valid: true}
		uc := &changePasswordInteractor{passwordAuthManager: manager, SessionRepository: sessions, txRunner: &loginTestTx{}}

		if _, err := uc.Execute(ctx, input); err != nil {
			t.Fatal(err)
		}
		if !manager.updated {
			t.Fatal("password was not updated")
		}
		if !sessions.sessions[current.SessionID()].IsActive(ctx) {
			t.Fatal("current session was deactivated")
		}
		if sessions.sessions[other.SessionID()].IsActive(ctx) {
			t.Fatal("other session is still active")
		}
		if !sessions.sessions[otherUsers.SessionID()].IsActive(ctx) {
			t.Fatal("another user's session was deactivated")
		}
	})

	t.Run("旧パスワードが違えばセッションを変更しない", func(t *testing.T) {
		sessions := newSessions()
		uc := &changePasswordInteractor{passwordAuthManager: &changePasswordTestManager{valid: false}, SessionRepository: sessions, txRunner: &loginTestTx{}}

		if _, err := uc.Execute(ctx, input); !errors.Is(err, messages.InvalidPasswordError) {
			t.Fatalf("err = %v, want %v", err, messages.InvalidPasswordError)
		}
		if !sessions.sessions[other.SessionID()].IsActive(ctx) {
			t.Fatal("session was deactivated without password change")
		}
	})

	t.Run("パスワード更新に失敗すればセッションを変更しない", func(t *testing.T) {
		sessions := newSessions()
		uc := &changePasswordInteractor{passwordAuthManager: &changePasswordTestManager{valid: true, updateErr: messages.PasswordTooShort}, SessionRepository: sessions, txRunner: &loginTestTx{}}

		if _, err := uc.Execute(ctx, input); !errors.Is(err, messages.PasswordTooShort) {
			t.Fatalf("err = %v, want %v", err, messages.PasswordTooShort)
		}
		if !sessions.sessions[other.SessionID()].IsActive(ctx) {
			t.Fatal("session was deactivated without password change")
		}
	})
}
