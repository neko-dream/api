package main

import (
	"fmt"
	"net/http"

	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/db"
	"github.com/neko-dream/server/internal/infrastructure/di"
	"github.com/neko-dream/server/internal/infrastructure/middleware"
	"github.com/neko-dream/server/internal/presentation/handler"
	"github.com/neko-dream/server/internal/presentation/oas"

	swMiddleware "github.com/go-openapi/runtime/middleware"
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
	// migrator.Down()
	migrator.Up()

	reqMiddleware := middleware.ReqMiddleware(srv)
	corsHandler := middleware.CORSMiddleware(reqMiddleware)
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
