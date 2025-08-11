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
	return &Bootstrap{
		container: container,
		config:    di.Invoke[*config.Config](container),
		migrator:  di.Invoke[*db.Migrator](container),
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
