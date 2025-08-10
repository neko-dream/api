package main

import (
	"context"
	"log"

	"github.com/neko-dream/server/cmd/server/bootstrap"

	"github.com/neko-dream/server/internal/application/event_processor"
	"github.com/neko-dream/server/internal/infrastructure/di"
)

func main() {
	container := di.BuildContainer()
	boot, err := bootstrap.New(container)
	if err != nil {
		log.Fatal("Failed to initialize: ", err)
	}

	eventProcessor := di.Invoke[*event_processor.EventProcessor](container)
	ctx := context.Background()
	go func() {
		eventProcessor.Start(ctx)
		log.Printf("Event processor started")
	}()

	if err := boot.Run(); err != nil {
		log.Fatal("Server error: ", err)
	}
}
