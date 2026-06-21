package bootstrap

import (
	"context"
	"fmt"
	"log"

	"github.com/neko-dream/api/internal/application/event_processor"
	"github.com/neko-dream/api/internal/application/usecase/analysis_usecase"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/di"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"go.uber.org/dig"
)

// Bootstrap アプリケーションの起動を管理する
type Bootstrap struct {
	container          *dig.Container
	config             *config.Config
	migrator           *db.Migrator
	eventProcessor     *event_processor.EventProcessor
	analysisScheduler  analysis_usecase.AnalysisScheduler
	cancelFunc         context.CancelFunc
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

	analysisScheduler, err := di.InvokeWithError[analysis_usecase.AnalysisScheduler](container)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke analysis scheduler: %w", err)
	}

	return &Bootstrap{
		container:         container,
		config:            config,
		migrator:          migrator,
		eventProcessor:    eventProcessor,
		analysisScheduler: analysisScheduler,
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
	b.analysisScheduler.Start(ctx)

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
	log.Println("Shutting down analysis scheduler...")
	b.analysisScheduler.Stop()
}
