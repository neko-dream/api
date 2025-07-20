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

	"github.com/swaggest/swgui"
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
	if err := migrator.Up(); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}
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
	mux.Handle("/admin/assets/", http.StripPrefix("/admin/assets/", handler.NewAdminUIAssetsHandler()))
	mux.Handle("/admin", http.RedirectHandler("/admin/", http.StatusSeeOther))
	tagsSorterFunc := "(a, b) => {" +
		"const priority = {\"auth\": 1, \"user\": 2, \"talk_session\": 3, \"opinion\": 4, \"organization\": 5, \"vote\": 6}; " +
		"const ap = priority[a.toLowerCase()]; " +
		"const bp = priority[b.toLowerCase()]; " +
		"if (ap && bp) return ap - bp; " +
		"if (ap) return -1; " +
		"if (bp) return 1; " +
		"return a.toLowerCase().localeCompare(b.toLowerCase());" +
		"}"
	swagger := v5emb.NewWithConfig(swgui.Config{
		Title:       "Kotohiro API",
		HideCurl:    true,
		SwaggerJSON: domain,
		BasePath:    "/docs/",
		JsonEditor:  true,
		SettingsUI: map[string]string{
			"deepLinking":              "true", // URLで各APIに直リンク可能
			"defaultModelsExpandDepth": "-1",
			"defaultModelExpandDepth":  "-1",
			"defaultModelRendering":    "\"model\"",
			"displayRequestDuration":   "true",
			"tryItOutEnabled":          "true",
			"layout":                   "\"BaseLayout\"",
			"showExtensions":           "true",
			"showCommonExtensions":     "true",
			"syntaxHighlight":          "{\"activate\": true,\"theme\": \"tomorrow-night\"}",
			"displayOperationId":       "true",
			"filter":                   "true",
			"operationsSorter":         "\"alpha\"",
			"tagsSorter":               tagsSorterFunc,
		},
	})
	mux.Handle("/docs/", swagger("Kotohiro API", domain, "/docs/"))
	// }

	if err := http.ListenAndServe(fmt.Sprintf(":%s", conf.PORT), mux); err != nil {
		log.Println("Error starting server:", err)
	}
}
