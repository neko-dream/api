package bootstrap

import (
	"context"
	"fmt"
	"log"

	"github.com/neko-dream/server/internal/application/event_processor"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/di"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"go.uber.org/dig"
)

// Bootstrap アプリケーションの起動を管理する
type Bootstrap struct {
	container      *dig.Container
	config         *config.Config
	migrator       *db.Migrator
	eventProcessor *event_processor.EventProcessor
	cancelFunc     context.CancelFunc
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

	eventProcessor, err := di.InvokeWithError[*event_processor.EventProcessor](container)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke event processor: %w", err)
	}

	return &Bootstrap{
		container:      container,
		config:         config,
		migrator:       migrator,
		eventProcessor: eventProcessor,
	}, nil
}

// Run アプリケーションを起動する
func (b *Bootstrap) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	b.cancelFunc = cancel

	if err := b.runMigrations(); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	b.startEventProcessor(ctx)

	return b.startHTTPServer()
}

func (b *Bootstrap) runMigrations() error {
	return b.migrator.Up()
}

// startEventProcessor イベントプロセッサーを起動する
func (b *Bootstrap) startEventProcessor(ctx context.Context) {
	go func() {
		log.Println("Starting event processor...")
		b.eventProcessor.Start(ctx)
	}()
}

// Shutdown アプリケーションを適切にシャットダウンする
func (b *Bootstrap) Shutdown() {
	if b.cancelFunc != nil {
		log.Println("Shutting down event processor...")
		b.cancelFunc()
	}
}
