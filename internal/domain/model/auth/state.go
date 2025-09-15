package auth

import (
	"context"
	"time"

	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/shared"
)

const (
	// StateExpirationDuration stateの有効期限（15分）
	StateExpirationDuration = 15 * time.Minute
)

var (
	ErrInvalidState = messages.InvalidStateError
	ErrExpiredState = messages.ExpiredStateError
)

type (
	// State OAuth認証のstate
	State struct {
		ID              int               // データベースの主キー
		State           string            // 認証stateの値
		Provider        string            // 認証プロバイダー名
		RedirectURL     string            // リダイレクトURL
		CreatedAt       time.Time         // 作成日時
		ExpiresAt       time.Time         // 有効期限
		RegistrationURL *string           // ログイン時に登録していない場合に飛ばすURL
		OrganizationID  *shared.UUID[any] // 組織ID（組織経由のログインの場合）
	}

	// StateRepository
	StateRepository interface {
		// Create 新しいstateを保存
		Create(ctx context.Context, state *State) error
		// Get 指定されたstateの値を取得
		Get(ctx context.Context, state string) (*State, error)
		// Delete 指定されたstateを削除
		Delete(ctx context.Context, state string) error
		// DeleteExpired 期限切れのstateを削除
		DeleteExpired(ctx context.Context) error
	}
)

// Validate 与えられたstateが有効かどうかを検証
// 1. stateの値が一致すること
// 2. stateの有効期限が切れていないこと
func (s *State) Validate(cookieState string) error {
	if s.State != cookieState {
		return ErrInvalidState
	}

	if time.Now().After(s.ExpiresAt) {
		return ErrExpiredState
	}

	return nil
}

// NewState
// state: 認証stateの値
// provider: 認証プロバイダー名
// redirectURL: リダイレクトURL
// expiresAt: 有効期限
// registrationURL: ログイン時に登録していない場合に飛ばすURL
// organizationID: 組織ID（組織経由のログインの場合）
func NewState(state string, provider string, redirectURL string, expiresAt time.Time, registrationURL *string, organizationID *shared.UUID[any]) *State {
	return &State{
		State:           state,
		Provider:        provider,
		RedirectURL:     redirectURL,
		ExpiresAt:       expiresAt,
		RegistrationURL: registrationURL,
		OrganizationID:  organizationID,
	}
}
