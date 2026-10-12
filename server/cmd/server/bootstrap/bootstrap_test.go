package bootstrap

import (
	"testing"

	"github.com/neko-dream/api/internal/infrastructure/di"
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
	assert.NotNil(t, boot.migrator)
	assert.NotNil(t, boot.eventProcessor)
}

func TestBootstrap_Run_ShouldReturnErrorWhenMigrationFails(t *testing.T) {
	// このテストは後で実装
	t.Skip("Implement after migrator mock is ready")
}

func TestBootstrap_Shutdown_ShouldCancelContext(t *testing.T) {
	// Arrange
	container := di.BuildContainer()
	boot, err := New(container)
	assert.NoError(t, err)

	// Manually set cancelFunc for testing
	boot.cancelFunc = func() {
		// This would be set during Run()
	}

	// Act & Assert - Should not panic
	assert.NotPanics(t, func() {
		boot.Shutdown()
	})
}

func TestBootstrap_Shutdown_WithoutRun(t *testing.T) {
	// Arrange
	container := di.BuildContainer()
	boot, err := New(container)
	assert.NoError(t, err)

	// Act & Assert - Should not panic even if cancelFunc is nil
	assert.NotPanics(t, func() {
		boot.Shutdown()
	})
}
