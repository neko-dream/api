package main

import (
	"fmt"
	"net/http"

	swMiddleware "github.com/go-openapi/runtime/middleware"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/db"
	"github.com/neko-dream/server/internal/infrastructure/di"
	"github.com/neko-dream/server/internal/infrastructure/middleware"
	"github.com/neko-dream/server/internal/presentation/handler"
	"github.com/neko-dream/server/internal/presentation/oas"

	"github.com/rs/cors"
)

func main() {
	container := di.BuildContainer()
	srv, err := oas.NewServer(
		di.Invoke[oas.Handler](container),
		di.Invoke[oas.SecurityHandler](container),
		oas.WithErrorHandler(handler.CustomErrorHandler),
	)
	if err != nil {
		panic(err)
	}

	config := di.Invoke[*config.Config](container)
	migrator := di.Invoke[*db.Migrator](container)
	if config.Env != "production" {
		migrator.Down()
		migrator.Up()
		// dummyInitializer := di.Invoke[*db.DummyInitializer](container)
		// dummyInitializer.Initialize()
	}

	reqMiddleware := middleware.ReqMiddleware(srv)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*", "localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	corsHandler := c.Handler(reqMiddleware)
	mux := http.NewServeMux()

	mux.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static")),
		),
	)

	opts := swMiddleware.SwaggerUIOpts{SpecURL: "/static/openapi.yaml"}
	sh := swMiddleware.SwaggerUI(opts, nil)
	mux.Handle("/docs/", sh)
	mux.Handle("/", corsHandler)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.PORT), mux); err != nil {
		panic(err)
	}
}
