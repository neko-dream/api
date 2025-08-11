package bootstrap

import (
	"testing"

	"github.com/neko-dream/server/internal/infrastructure/di"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_ShouldCreateBootstrapInstance(t *testing.T) {
	// Arrange
	container := di.BuildContainer()

	// Act
	boot, err := New(container)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, boot)
	assert.NotNil(t, boot.container)
	assert.NotNil(t, boot.config)
}

func TestBootstrap_Run_ShouldReturnErrorWhenMigrationFails(t *testing.T) {
	// このテストは後で実装
	t.Skip("Implement after migrator mock is ready")
}
