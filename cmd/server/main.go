package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/neko-dream/server/internal/infrastructure/di"
	"github.com/neko-dream/server/internal/infrastructure/middleware"
	"github.com/neko-dream/server/internal/infrastructure/utils"
	"github.com/neko-dream/server/internal/presentation/handler"
	"github.com/neko-dream/server/internal/presentation/oas"
)

func main() {
	// .envを読み込む
	if err := utils.LoadEnv(); err != nil {
		panic(err)
	}

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
	port := os.Getenv("PORT")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), corsHandler); err != nil {
		panic(err)
	}
}
