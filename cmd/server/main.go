package main

import (
	"fmt"
	"net/http"

	swMiddleware "github.com/go-openapi/runtime/middleware"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/di"
	"github.com/neko-dream/server/internal/infrastructure/http/middleware"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/internal/presentation/handler"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/rs/cors"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {

	container := di.BuildContainer()
	srv, err := oas.NewServer(
		di.Invoke[oas.Handler](container),
		di.Invoke[oas.SecurityHandler](container),
		oas.WithErrorHandler(handler.CustomErrorHandler),
		oas.WithTracerProvider(di.Invoke[*sdktrace.TracerProvider](container)),
	)
	if err != nil {
		panic(err)
	}

	config := di.Invoke[*config.Config](container)
	migrator := di.Invoke[*db.Migrator](container)
	if config.Env != "production" {
		// migrator.Down()
		migrator.Up()
		// di.Invoke[*db.DummyInitializer](container).Initialize()
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
