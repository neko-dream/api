package bootstrap

import (
	"fmt"

	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/di"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"go.uber.org/dig"
)

// Bootstrap アプリケーションの起動を管理する
type Bootstrap struct {
	container *dig.Container
	config    *config.Config
	migrator  *db.Migrator
}

// New 新しいBootstrapインスタンスを作成する
func New(container *dig.Container) (*Bootstrap, error) {
	config, err := di.InvokeWithError[*config.Config](container)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke config: %w", err)
	}

	migrator, err := di.InvokeWithError[*db.Migrator](container)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke migrator: %w", err)
	}

	return &Bootstrap{
		container: container,
		config:    config,
		migrator:  migrator,
	}, nil
}

// Run アプリケーションを起動する
func (b *Bootstrap) Run() error {
	if err := b.runMigrations(); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	return b.startHTTPServer()
}

func (b *Bootstrap) runMigrations() error {
	return b.migrator.Up()
}
