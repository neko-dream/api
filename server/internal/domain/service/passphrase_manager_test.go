package service

import (
	"testing"

	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPassphraseManager_Hash(t *testing.T) {
	cfg := &config.Config{HASH_PEPPER: "test-pepper"}
	manager := NewPassphraseManager(cfg)

	hashed, err := manager.Hash("my-passphrase")
	require.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, "my-passphrase", hashed)
}

func TestPassphraseManager_Verify(t *testing.T) {
	cfg := &config.Config{HASH_PEPPER: "test-pepper"}
	manager := NewPassphraseManager(cfg)

	t.Run("正しい合言葉で検証成功", func(t *testing.T) {
		hashed, err := manager.Hash("correct")
		require.NoError(t, err)
		assert.True(t, manager.Verify("correct", hashed))
	})

	t.Run("誤った合言葉で検証失敗", func(t *testing.T) {
		hashed, err := manager.Hash("correct")
		require.NoError(t, err)
		assert.False(t, manager.Verify("wrong", hashed))
	})

	t.Run("異なるpepperでは検証失敗", func(t *testing.T) {
		hashed, err := manager.Hash("passphrase")
		require.NoError(t, err)

		otherManager := NewPassphraseManager(&config.Config{HASH_PEPPER: "other-pepper"})
		assert.False(t, otherManager.Verify("passphrase", hashed))
	})
}
