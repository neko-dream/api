package main

import (
	"log"

	"github.com/neko-dream/server/cmd/server/bootstrap"
	"github.com/neko-dream/server/internal/infrastructure/di"
)

func main() {
	container := di.BuildContainer()
	boot, err := bootstrap.New(container)
	if err != nil {
		log.Fatal("Failed to initialize: ", err)
	}

	if err := boot.Run(); err != nil {
		log.Fatal("Server error: ", err)
	}
}
