package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/neko-dream/server/internal/infrastructure/db"
	"github.com/neko-dream/server/internal/infrastructure/di"
	"github.com/neko-dream/server/internal/infrastructure/middleware"
	"github.com/neko-dream/server/internal/presentation/handler"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/neko-dream/server/pkg/utils"

	swMiddleware "github.com/go-openapi/runtime/middleware"
)

func main() {
	// .envを読み込む
	if err := utils.LoadEnv(); err != nil {
		panic(err)
	}

	db.Down()
	db.Migration()

	container := di.BuildContainer()
	srv, err := oas.NewServer(
		di.Invoke[oas.Handler](container),
		di.Invoke[oas.SecurityHandler](container),
		oas.WithErrorHandler(handler.CustomErrorHandler),
	)
	if err != nil {
		panic(err)
	}

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

	port := os.Getenv("PORT")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
		panic(err)
	}
}
