package bootstrap

import (
	"testing"

	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
)

func TestBootstrap_getAdminRoutes(t *testing.T) {
	// Arrange
	boot := &Bootstrap{
		config: &config.Config{},
	}

	// Act
	routes := boot.getAdminRoutes()

	// Assert
	assert.Len(t, routes, 3)

	// /admin/ route
	assert.Equal(t, "/admin/", routes[0].Pattern)
	assert.Equal(t, "/admin/", routes[0].StripPrefix)
	assert.NotNil(t, routes[0].Handler)

	// /admin/assets/ route
	assert.Equal(t, "/admin/assets/", routes[1].Pattern)
	assert.Equal(t, "/admin/assets/", routes[1].StripPrefix)
	assert.NotNil(t, routes[1].Handler)

	// /admin redirect route
	assert.Equal(t, "/admin", routes[2].Pattern)
	assert.Equal(t, "", routes[2].StripPrefix)
	assert.NotNil(t, routes[2].Handler)
}

func TestBootstrap_getSwaggerRoutes(t *testing.T) {
	// Arrange
	boot := &Bootstrap{
		config: &config.Config{
			Env:  config.DEV,
			PORT: "8080",
		},
	}

	// Act
	routes := boot.getSwaggerRoutes()

	// Assert
	assert.Len(t, routes, 1)
	assert.Equal(t, "/docs/", routes[0].Pattern)
	assert.Equal(t, "", routes[0].StripPrefix)
	assert.NotNil(t, routes[0].Handler)
}

func TestBootstrap_setupRoutes_RegistersAllRoutes(t *testing.T) {
	// Arrange
	boot := &Bootstrap{
		container: nil, // テストでは使用しない
		config: &config.Config{
			Env:  config.DEV,
			PORT: "8080",
		},
	}

	// Act & Assert - ルート定義が取得できることを確認
	allRoutes := []Route{}
	allRoutes = append(allRoutes, boot.getAdminRoutes()...)
	allRoutes = append(allRoutes, boot.getSwaggerRoutes()...)

	// 期待されるルートパターン
	expectedPatterns := map[string]bool{
		"/admin/":        true,
		"/admin/assets/": true,
		"/admin":         true,
		"/docs/":         true,
	}

	for _, route := range allRoutes {
		_, ok := expectedPatterns[route.Pattern]
		assert.True(t, ok, "Unexpected route pattern: %s", route.Pattern)
		delete(expectedPatterns, route.Pattern)
	}

	assert.Empty(t, expectedPatterns, "Missing routes: %v", expectedPatterns)
}
