package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/di"
	"github.com/neko-dream/server/internal/infrastructure/http/middleware"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/internal/presentation/handler"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/rs/cors"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/swaggest/swgui/v5emb"
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

	conf := di.Invoke[*config.Config](container)
	migrator := di.Invoke[*db.Migrator](container)

	// migrator.Down()
	migrator.Up()
	// di.Invoke[*db.DummyInitializer](container).Initialize()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	routeFinder := middleware.MakeRouteFinder(srv)
	corsHandler := c.Handler(middleware.Wrap(
		srv,
		middleware.Instrument("kotohiro-api", routeFinder, di.Invoke[*sdktrace.TracerProvider](container)),
		middleware.Labeler(routeFinder),
	))

	mux := http.NewServeMux()
	mux.Handle("/", corsHandler)
	// if conf.Env != config.PROD {
	var domain string
	if conf.Env == config.DEV {
		domain = "https://api-dev.kotohiro.com/static/oas/openapi.yaml"
	} else if conf.Env == config.PROD {
		domain = "https://api.kotohiro.com/static/oas/openapi.yaml"
	} else {
		domain = "http://localhost:" + conf.PORT + "/static/oas/openapi.yaml"
	}
	mux.Handle("/static/", http.StripPrefix("/static/", handler.NewStaticHandler()))
	mux.Handle("/admin/", http.StripPrefix("/admin/", handler.NewAdminUIHandler()))
	mux.Handle("/assets/", http.StripPrefix("/assets/", handler.NewAdminUIAssetsHandler()))
	mux.Handle("/admin", http.RedirectHandler("/admin/index.html", http.StatusSeeOther))
	mux.Handle("/docs/", v5emb.New("kotohiro", domain, "/docs/"))
	// }

	if err := http.ListenAndServe(fmt.Sprintf(":%s", conf.PORT), mux); err != nil {
		log.Println("Error starting server:", err)
	}
}
