package session

import (
	"testing"

	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
)

func TestSignedToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		sessionID    string
		secret1      string
		secret2      string
		expectError  bool
		errorMessage string
	}{
		{
			name:         "正常に署名付きトークンを生成・検証できる",
			sessionID:    shared.NewUUID[any]().String(),
			secret1:      "test-secret",
			secret2:      "test-secret",
			expectError:  false,
			errorMessage: "",
		},
		{
			name:         "異なるシークレットでは検証に失敗する",
			sessionID:    shared.NewUUID[any]().String(),
			secret1:      "test-secret",
			secret2:      "wrong-secret",
			expectError:  true,
			errorMessage: "invalid token",
		},
		{
			name:         "不正な形式のトークンは検証に失敗する",
			sessionID:    "invalid-token-format",
			secret1:      "test-secret",
			secret2:      "test-secret",
			expectError:  true,
			errorMessage: "invalid token",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// トークン生成用のmanager
			cfg1 := &config.Config{TokenSecret: tt.secret1}
			manager1 := &sessionTokenManager{
				secret: cfg1.TokenSecret,
			}

			// トークンを生成
			var token string
			if tt.name == "不正な形式のトークンは検証に失敗する" {
				token = tt.sessionID
			} else {
				token = manager1.createSignedToken(tt.sessionID)
			}

			// トークン検証用のmanager
			cfg2 := &config.Config{TokenSecret: tt.secret2}
			manager2 := &sessionTokenManager{
				secret: cfg2.TokenSecret,
			}

			// トークンを検証
			result, err := manager2.verifySignedToken(token)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMessage, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.sessionID, result)
			}
		})
	}
}

func TestTokenFormat(t *testing.T) {
	cfg := &config.Config{TokenSecret: "test-secret"}
	manager := &sessionTokenManager{
		secret: cfg.TokenSecret,
	}

	sessionID := shared.NewUUID[any]().String()
	token := manager.createSignedToken(sessionID)

	// トークンが正しい形式であることを確認
	assert.Contains(t, token, ".")
	assert.True(t, len(token) > len(sessionID)+1)
}
